package infinote

import (
	"infinote/db"
	"infinote/email"
	"infinote/report"
	"infinote/store"
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

var _ UserStorer = &store.Users{}
var _ Notifier = &email.Console{}
var _ Notifier = &email.Mailer{}
var _ Reporter = &report.Sentry{}
var _ Reporter = &report.Console{}
var _ TokenStorer = &store.Tokens{}
var _ AuthProvider = &Auther{}

// AuthProvider contains methods for authentication
type AuthProvider interface {
	GenerateJWT(ctx context.Context, user *db.User, userAgent string) (string, error)
}

// TokenStorer collects all token methods
type TokenStorer interface {
	Get(id string) (*db.IssuedToken, error)
	Insert(t *db.IssuedToken) (*db.IssuedToken, error)
	Update(t *db.IssuedToken) (*db.IssuedToken, error)
	GetAllByUser(userID string) ([]*db.IssuedToken, error)
	GetAllExpired() ([]*db.IssuedToken, error)
	Delete(t *db.IssuedToken) error
	Blacklist() (store.Blacklist, error)
}

// Reporter is used for notifying the admin of errors
type Reporter interface {
	LogExternal(err error)
}

// Notifier is used for notifying the user of things
type Notifier interface {
	ReceivedSignup(email string) error
	ForgotPassword(email string) error
}

// CompanyStorer collects all Note methods
type CompanyStorer interface {
	All() (db.CompanySlice, error)
	Get(id uuid.UUID) (*db.Company, error)
	GetMany(keys []string) (db.CompanySlice, []error)
	Insert(record *db.Company, txes ...*sql.Tx) (*db.Company, error)
	Update(record *db.Company, txes ...*sql.Tx) (*db.Company, error)
}

// NoteStorer collects all Note methods
type NoteStorer interface {
	All() (db.NoteSlice, error)
	Select(userID uuid.UUID, limit, offset int) ([]*db.Note, error)
	GetMany(keys []string) (db.NoteSlice, []error)
	Get(id uuid.UUID) (*db.Note, error)
	Insert(record *db.Note, txes ...*sql.Tx) (*db.Note, error)
	Update(t *db.Note, txes ...*sql.Tx) (*db.Note, error)
	GetByUser(userID uuid.UUID, txes ...*sql.Tx) ([]*db.Note, error)
}

// UserStorer collects all user methods
type UserStorer interface {
	BeginTransaction() (*sql.Tx, error)
	GetByVerifyToken(token string, txes ...*sql.Tx) (*db.User, error)
	GetByEmail(email string, txes ...*sql.Tx) (*db.User, error)
	Get(id uuid.UUID, txes ...*sql.Tx) (*db.User, error)
	GetManyByIDs(keys []string, txes ...*sql.Tx) (db.UserSlice, []error)
	All(txes ...*sql.Tx) (db.UserSlice, error)
	GetByCompany(orgID uuid.UUID, txes ...*sql.Tx) (db.UserSlice, error)
	Insert(u *db.User, tx ...*sql.Tx) (*db.User, error)
	Update(u *db.User, tx ...*sql.Tx) (*db.User, error)
}

type RoleStorer interface {
	ByUser(id uuid.UUID, tx ...*sql.Tx) (*db.Role, error)
}
