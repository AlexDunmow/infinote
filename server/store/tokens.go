package store

import (
	"boilerplate/db"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
)

// TokenFactory creates tokens
func TokenFactory() *db.IssuedToken {
	u := &db.IssuedToken{
		ID:           uuid.Must(uuid.NewV4()).String(),
		UserID:       "",
		Device:       "",
		TokenCreated: time.Now(),
		TokenExpires: time.Now(),
		Blacklisted:  false,
	}
	return u
}

// Blacklist type for recording blacklisted tokens
type Blacklist map[string]struct{}

// Token for persistence
type Tokens struct {
	Conn boil.Executor
}

func NewTokenStore(conn *sqlx.DB) *Tokens {
	return &Tokens{conn}
}

// Blacklist returns the token blacklist
func (s *Tokens) Blacklist() (Blacklist, error) {
	tokens, err := db.IssuedTokens(db.IssuedTokenWhere.Blacklisted.EQ(true)).All(s.Conn)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	list := Blacklist{}
	for _, token := range tokens {
		list[token.ID] = struct{}{}
	}
	return list, nil
}

// Get returns the token given its ID
func (s *Tokens) Get(id string) (*db.IssuedToken, error) {
	return db.FindIssuedToken(s.Conn, id)
}

// GetAllByUser returns the GetAllByUser
func (s *Tokens) GetAllByUser(userID string) ([]*db.IssuedToken, error) {
	return db.IssuedTokens(db.IssuedTokenWhere.UserID.EQ(userID)).All(s.Conn)
}

// GetAllExpired returns the GetAllExpired
func (s *Tokens) GetAllExpired() ([]*db.IssuedToken, error) {
	return db.IssuedTokens(db.IssuedTokenWhere.TokenExpires.LT(time.Now())).All(s.Conn)
}

// Insert returns the Insert
func (s *Tokens) Insert(t *db.IssuedToken) (*db.IssuedToken, error) {
	err := t.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Update returns the Update
func (s *Tokens) Update(t *db.IssuedToken) (*db.IssuedToken, error) {
	_, err := t.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Delete returns the Delete
func (s *Tokens) Delete(t *db.IssuedToken) error {
	_, err := t.Delete(s.Conn)
	if err != nil {
		return err
	}
	return nil
}
