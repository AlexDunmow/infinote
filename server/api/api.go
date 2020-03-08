package api

import (
	"encoding/json"
	"errors"
	"infinote"
	"strings"

	"go.uber.org/zap"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"

	"context"
	"fmt"
	"infinote/canlog"
	"infinote/graphql"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/gqlerror"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

// BasicResponse just contain success flag
type BasicResponse struct {
	Errors  []*error `json:"errors,omitempty"`
	Success bool     `json:"success"`
}

type ControllerOpts struct {
	NoteStorer        infinote.NoteStorer
	UserStorer        infinote.UserStorer
	CompanyStorer     infinote.CompanyStorer
	BlacklistProvider infinote.BlacklistProvider
	TokenStorer       infinote.TokenStorer
	RoleStorer        infinote.RoleStorer
	//SubscriptionResolver infinote.SubscriptionProvider
	JWTSecret string
	Auther    *infinote.Auther
	Logger    *zap.SugaredLogger
}

// log *zap.SugaredLogger, ts infinote.NoteStorer, os , us infinote.UserStorer,
// bs infinote.BlacklistProvider, tks infinote.TokenStorer, rs infinote.RoleStorer, jwtSecret string, auther *infinote.Auther

// NewAPIController returns the public router layer
func NewAPIController(opts *ControllerOpts) http.Handler {
	authentication := NewAuthRoutes(opts.UserStorer, opts.RoleStorer, opts.Auther)
	// Websocket Upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(opts.Auther.VerifyMiddleware())
	r.Use(canonicalLogger(opts.Logger.Desugar()))
	r.Use(infinote.DataloaderMiddleware(opts.NoteStorer, opts.CompanyStorer, opts.UserStorer))
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Mount("/api/auth", authentication)
	r.Route("/api/gql", func(r chi.Router) {
		r.Handle("/", handler.Playground("GraphQL playground", "/api/gql/query"))
		r.Handle("/query", handler.GraphQL(
			graphql.NewExecutableSchema(
				graphql.Config{
					Resolvers: graphql.NewResolver(&graphql.ResolverOpts{
						RoleStorer:        opts.RoleStorer,
						Auther:            opts.Auther,
						NoteStorer:        opts.NoteStorer,
						CompanyStorer:     opts.CompanyStorer,
						UserStorer:        opts.UserStorer,
						BlacklistProvider: opts.BlacklistProvider,
					}),
				},
			),
			handler.ErrorPresenter(
				func(ctx context.Context, e error) *gqlerror.Error {
					var bErr *infinote.Error
					canlog.SetErr(ctx, e)
					if errors.As(e, &bErr) {
						canlog.SetErr(ctx, errors.New(bErr.Error()))
						canlog.AppendErr(ctx, bErr.ID)
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New(bErr.Message))
					}
					if errors.Is(e, infinote.ErrBadContext) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("There was a problem reading your credentials. Please sign in and try again."))
					}
					if errors.Is(e, infinote.ErrBadClaims) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("There was a problem reading your credentials. Please sign in and try again."))
					}
					if errors.Is(e, infinote.ErrBlacklisted) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("Login token is no longer valid. Please login again."))
					}
					if errors.Is(e, infinote.ErrUnauthorized) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("You are not authorized to do this action."))
					}
					if errors.Is(e, infinote.ErrBadCredentials) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("Please check your username or password and try again."))
					}
					if errors.Is(e, infinote.ErrNotImplemented) {
						return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("This functionality is not yet available."))
					}
					return gqlgraphql.DefaultErrorPresenter(ctx, errors.New("There was a problem with the server. Please try again later."))
				},
			),
			handler.WebsocketKeepAliveDuration(5*time.Second),
			// hack to make firecamp ws extension work because of origin security check. remove when done
			handler.WebsocketUpgrader(upgrader),
		),
		)
	})

	return r
}

// AuthController contains handlers involving authentication
type AuthController struct {
	cookieDefaults CookieSettings
	userStorer     infinote.UserStorer
	roleStorer     infinote.RoleStorer
	auther         *infinote.Auther
}

// CookieSettings are the default values used to set cookies
type CookieSettings struct {
	SameSite http.SameSite
	HTTPOnly bool
	Secure   bool
	Path     string
}

// NewAuthRoutes returns a router for use in authentication
func NewAuthRoutes(userStorer infinote.UserStorer, roleStorer infinote.RoleStorer, auther *infinote.Auther) chi.Router {
	cookieDefaults := CookieSettings{
		SameSite: http.SameSiteDefaultMode,
		HTTPOnly: true,
		Secure:   false, // TODO: set to true when https available. currently not enabled for dev
		Path:     "/",
	}

	c := &AuthController{
		cookieDefaults: cookieDefaults,
		userStorer:     userStorer,
		roleStorer:     roleStorer,
		auther:         auther,
	}

	r := chi.NewRouter()
	r.Post("/login", c.login())
	r.Post("/logout", c.logout())
	r.Post("/verify_account", c.verifyAccount())
	return r
}

// LoginRequest structs for the HTTP request/response cycle
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest structs for the HTTP request/response cycle
type VerifyRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

// LoginResponse structs for the HTTP request/response cycle
type LoginResponse struct {
	Verified bool `json:"verified"`
	Success  bool `json:"success"`
}

func httpWriteError(w http.ResponseWriter, err error, message string, code int) {
	// TODO, figure out context
	ctx := context.Background()

	// TODO, figure out log external
	infinote.LogExternal(ctx, nil, err)

	// message is the friendly error text
	// do not send user internal error log
	http.Error(w, fmt.Sprintf(`{"error":"%s"}`, message), code)
}

func httpWriteJSON(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		httpWriteError(w, err, "fail to prepare data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(b)
}

// login logs a user in
func (c *AuthController) login() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			httpWriteError(w, err, "invalid user input", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// email must be lower case
		email := strings.ToLower(req.Email)

		err = infinote.ValidatePassword(ctx, c.userStorer, email, req.Password)
		if err != nil {
			httpWriteError(w, err, "email or password fail", http.StatusUnauthorized)
			return
		}

		// load user details
		user, err := c.userStorer.GetByEmail(email)
		if err != nil {
			httpWriteError(w, err, "fail to load user", http.StatusInternalServerError)
			return
		}

		// save user detail in encrypted cookie and make it persist
		expiration := time.Now().Add(time.Duration(c.auther.TokenExpirationDays) * time.Hour * 24)
		jwt, err := c.auther.GenerateJWT(ctx, user, r.UserAgent())
		if err != nil {
			httpWriteError(w, err, "cookie fail", http.StatusBadRequest)
			return
		}

		// push cookie to browser
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    jwt,
			Expires:  expiration,
			HttpOnly: c.cookieDefaults.HTTPOnly,
			Path:     c.cookieDefaults.Path,
			SameSite: c.cookieDefaults.SameSite,
			Secure:   c.cookieDefaults.Secure,
		}
		http.SetCookie(w, &cookie)

		lr := &LoginResponse{
			Verified: user.Verified,
			Success:  true,
		}
		httpWriteJSON(w, lr)
	}
	return fn
}

// logout
func (c *AuthController) logout() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// clear and expire cookie and push to browser
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().AddDate(-1, 0, 0),
			HttpOnly: c.cookieDefaults.HTTPOnly,
			Path:     c.cookieDefaults.Path,
			SameSite: c.cookieDefaults.SameSite,
			Secure:   c.cookieDefaults.Secure,
		}
		http.SetCookie(w, &cookie)

		resp := &BasicResponse{
			Success: true,
		}

		httpWriteJSON(w, resp)
	}
	return fn
}

// verifyAccount verifies an account and logs the user in
func (c *AuthController) verifyAccount() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		req := &VerifyRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			httpWriteError(w, err, "invalid user input", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// email and token must be lower case
		email := strings.ToLower(req.Email)
		token := strings.ToLower(req.Token)

		// load user details
		user, err := c.userStorer.GetByEmail(email)
		if err != nil {
			httpWriteError(w, err, "fail to load user", http.StatusInternalServerError)
			return
		}

		if user.VerifyToken != token {
			httpWriteError(w, err, "fail to validate token", http.StatusInternalServerError)
			return
		}

		user.Verified = true
		user, err = c.userStorer.Update(user)
		if err != nil {
			panic(err)
		}

		resp := &BasicResponse{
			Success: true,
		}

		httpWriteJSON(w, resp)
	}
	return fn
}
