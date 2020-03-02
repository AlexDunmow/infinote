package infinote

import (
	"infinote/canlog"
	"infinote/db"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/jwtauth"
)

// Auther to handle JWT authentication
type Auther struct {
	TokenExpirationDays int
	TokenAuth           *jwtauth.JWTAuth
	Blacklister         BlacklistProvider
	TokenStore          TokenStorer
	UserStore           UserStorer
}

// NewAuther for JWT and blacklisting
func NewAuther(jwtsecret string, userStore UserStorer, blacklister BlacklistProvider, tokenStorer TokenStorer, tokenExpirationDays int) *Auther {
	result := &Auther{
		TokenAuth:           jwtauth.New("HS256", []byte(jwtsecret), []byte(jwtsecret)),
		Blacklister:         blacklister,
		TokenStore:          tokenStorer,
		UserStore:           userStore,
		TokenExpirationDays: tokenExpirationDays,
	}
	return result
}

// ClaimsFromContext a map of all claims in JWT
func ClaimsFromContext(ctx context.Context) (map[string]string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	result := map[string]string{}
	for k, v := range claims {
		val, ok := v.(string)
		if !ok {
			continue
		}
		result[k] = val
	}
	return result, err
}

// ClaimKey is a type used to set values in the JWT
type ClaimKey string

// ClaimUserName JWT key value
const ClaimUserName ClaimKey = "username"

// ClaimUserID JWT key value
const ClaimUserID ClaimKey = "uid"

// ClaimOrgID JWT key value
const ClaimOrgID ClaimKey = "oid"

// ClaimRoles JWT key value
const ClaimRoles ClaimKey = "roles"

// ClaimTokenID JWT key value
const ClaimTokenID ClaimKey = "tokenId"

// ClaimExistsInContext returns a specific claim value from the JWT
func ClaimExistsInContext(ctx context.Context, key ClaimKey) bool {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return false
	}
	_, ok := claims.Get(string(key))
	if !ok {
		return false
	}
	return true
}

// ClaimValueFromContext returns a specific claim value from the JWT
func ClaimValueFromContext(ctx context.Context, key ClaimKey) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		canlog.AppendErr(ctx, "58a90dfe-5248-41d0-ab9b-aaf98678f298")
		return "", ErrBadContext
	}
	valI, ok := claims.Get(string(key))
	if !ok {
		canlog.AppendErr(ctx, "573c3ae7-4182-44fa-8888-a356495b8e93")
		canlog.Set(ctx, "claimkey", key)
		return "", ErrBadClaims
	}
	val, ok := valI.(string)
	if !ok {
		canlog.AppendErr(ctx, "e22d2fd6-0eff-4663-a12b-190efea26c9f")
		return "", ErrTypeCast
	}
	return val, nil
}

// UserFromContext grabs the user from the context if a JWT is inside
// This is a heavy func, don't use it too often, and if you do, bend your knees
func UserFromContext(ctx context.Context, us UserStorer, bs BlacklistProvider) (*db.User, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		canlog.AppendErr(ctx, "9938abc7-dea1-4eee-a1a1-b3218f25fdf8")
		return nil, ErrBadContext
	}

	userIDI, ok := claims.Get(string(ClaimUserID))
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimUserID)
		canlog.AppendErr(ctx, "64a7ffe7-f615-4609-ba80-2e38bedd23d5")
		return nil, ErrBadClaims
	}

	tokenID, ok := claims.Get(string(ClaimTokenID))
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimTokenID)
		canlog.AppendErr(ctx, "2907cf77-a6c6-40b5-9604-daff9e49f640")
		return nil, ErrBadClaims
	}

	blacklisted := bs.OnList(tokenID.(string))
	if blacklisted {
		canlog.AppendErr(ctx, "66d879f8-84b4-4b97-84b9-38c4bd079901")
		return nil, ErrBlacklisted
	}

	userIDStr, ok := userIDI.(string)
	if !ok {
		canlog.AppendErr(ctx, "2b3391f7-f96d-41a7-b89d-a123932e267b")
		return nil, ErrTypeCast
	}
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		canlog.AppendErr(ctx, "bef39ebe-62fa-4051-8928-27f7f01824ad")
		return nil, ErrParse
	}
	u, err := us.Get(userID)
	if err != nil {
		canlog.AppendErr(ctx, "545e072f-6f67-4e93-895e-1ebeeaf4b59e")
		return nil, fmt.Errorf("get user: %w", err)
	}

	return u, nil

}

// GenerateJWT returns the token for client side persistence
func (a *Auther) GenerateJWT(ctx context.Context, user *db.User, userAgent string) (string, error) {
	// Record token in issued token records

	//TODO: device currently being set by request.UserAgent() ... might need to parse to get a better device name
	newToken := &db.IssuedToken{
		UserID:       user.ID,
		CompanyID:    user.CompanyID,
		Device:       userAgent,
		TokenExpires: time.Now().Add(time.Hour * time.Duration(24) * time.Duration(a.TokenExpirationDays)),
	}
	token, err := a.TokenStore.Insert(newToken)
	if err != nil {
		canlog.AppendErr(ctx, "107be839-712e-4c3a-837b-471d647a55ca")
		return "", fmt.Errorf("insert token: %w", err)
	}
	_, tokenString, err := a.TokenAuth.Encode(jwtauth.Claims{
		string(ClaimUserName): user.Email,
		string(ClaimUserID):   user.ID,
		string(ClaimTokenID):  token.ID,
	})
	if err != nil {
		canlog.AppendErr(ctx, "e61e7b60-b795-41f7-8fdb-0736c4e472d2")
		return "", fmt.Errorf("encode token: %w", err)
	}
	return tokenString, nil
}

// VerifyMiddleware for authentication adds JWT to context down the HTTP chain
func (a *Auther) VerifyMiddleware() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.TokenAuth)
}

// ValidatePassword will check the login details
func ValidatePassword(ctx context.Context, us UserStorer, email string, password string) error {
	user, err := us.GetByEmail(email)
	if user == nil {
		canlog.AppendErr(ctx, "750b8c66-0c19-43db-8151-535af516c74b")
		return fmt.Errorf("get user: %w", ErrAuthNoEmail)
	}
	if err != nil {
		canlog.AppendErr(ctx, "e9e43a1c-a86d-4428-a8d3-44c2a8d988ef")
		return fmt.Errorf("get user: %w", err)
	}

	storedHash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		canlog.AppendErr(ctx, "b28c8525-c68f-4bc0-ab72-4bc7caa57e0a")
		return fmt.Errorf("decode hash: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		canlog.AppendErr(ctx, "3bd139d1-f4d9-46a4-8454-a913528c1b2a")
		return fmt.Errorf("get user: %w", ErrAuthWrongPassword)
	}
	if err != nil {
		canlog.AppendErr(ctx, "740f2709-e579-435b-8f60-8a7b9e0a8e65")
		return fmt.Errorf("compare hash: %w", err)
	}

	return nil
}
