package store

import (
	"boilerplate/db"
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"syreclabs.com/go/faker"
)

// UserFactory creates users
func UserFactory() *db.User {
	u := &db.User{
		ID:           uuid.Must(uuid.NewV4()).String(),
		Name:         faker.Name().Name(),
		Email:        faker.Internet().Email(),
		PasswordHash: faker.Internet().Password(8, 20),
	}
	return u
}

// NewUserStore returns a new user repo that implements UserMutator, UserArchiver and UserQueryer
func NewUserStore(conn *sqlx.DB) *Users {
	r := &Users{conn}
	return r
}

// Users for persistence
type Users struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *Users) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// GetByVerifyToken returns a user with the matching verify token
func (s *Users) GetByVerifyToken(token string, txes ...*sql.Tx) (*db.User, error) {
	return db.Users(db.UserWhere.VerifyToken.EQ(token)).One(s.Conn)
}

// All users
func (s *Users) All(txes ...*sql.Tx) (db.UserSlice, error) {
	return db.Users().All(s.Conn)
}

// GetByCompany users by org
func (s *Users) GetByCompany(orgID uuid.UUID, txes ...*sql.Tx) (db.UserSlice, error) {
	return db.Users(db.UserWhere.CompanyID.EQ(orgID.String())).All(s.Conn)
}

// GetManyByIDs users given a list of IDs
func (s *Users) GetManyByIDs(keys []string, txes ...*sql.Tx) (db.UserSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	args := []interface{}{}
	for _, key := range keys {
		args = append(args, key)
	}
	records, err := db.Users(qm.WhereIn("id in ?", args...)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.User{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.User{}
	for _, key := range keys {
		for _, record := range records {
			if record.ID == key {
				result = append(result, record)
				break
			}
		}
	}
	return result, nil
}

// Get a user given their ID
func (s *Users) Get(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	return db.FindUser(s.Conn, id.String())
}

// Insert a user
func (s *Users) Insert(u *db.User, txes ...*sql.Tx) (*db.User, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return u.Insert(tx, boil.Infer())
	}, txes...)

	err = u.Reload(s.Conn)
	if err != nil {
		return nil, err
	}
	return u, err
}

// Update a user
func (s *Users) Update(u *db.User, txes ...*sql.Tx) (*db.User, error) {
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail returns a user given an email
func (s *Users) GetByEmail(email string, txes ...*sql.Tx) (*db.User, error) {
	return db.Users(db.UserWhere.Email.EQ(email)).One(s.Conn)
}

// Create a user
func (s *Users) Create(input *db.User, txes ...*sql.Tx) (*db.User, error) {
	err := input.Insert(s.Conn, boil.Infer())
	return input, err
}

// Archive will archive users
func (s *Users) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.FindUser(s.Conn, id.String())
	if err != nil {
		return nil, err
	}
	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.UserColumns.Archived, db.UserColumns.ArchivedAt))
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Unarchive will unarchive users
func (s *Users) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.FindUser(s.Conn, id.String())
	if err != nil {
		return nil, err
	}
	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.UserColumns.Archived, db.UserColumns.ArchivedAt))
	if err != nil {
		return nil, err
	}
	return u, nil
}
