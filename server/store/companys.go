package store

import (
	"infinote/db"
	"database/sql"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"syreclabs.com/go/faker"
)

// CompanyFactory creates orgs
func CompanyFactory() *db.Company {
	name := faker.Company().Name()
	u := &db.Company{
		ID:   uuid.Must(uuid.NewV4()).String(),
		Key:  strings.Replace(name, " ", "", -1),
		Name: name,
	}
	return u
}

// Company for persistence
type Company struct {
	Conn *sqlx.DB
}

// NewCompanyStore returns a new store
func NewCompanyStore(conn *sqlx.DB) *Company {
	os := &Company{conn}
	return os
}

// All Companyss
func (s *Company) All() (db.CompanySlice, error) {
	return db.Companies().All(s.Conn)
}

// Get Companys
func (s *Company) Get(id uuid.UUID) (*db.Company, error) {
	return db.FindCompany(s.Conn, id.String())
}

// GetMany Companys
func (s *Company) GetMany(keys []string) (db.CompanySlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	args := []interface{}{}
	for _, key := range keys {
		args = append(args, key)
	}
	records, err := db.Companies(qm.WhereIn("id in ?", args...)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Company{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Company{}
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

// Insert Companys
func (s *Company) Insert(record *db.Company, txes ...*sql.Tx) (*db.Company, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)

	err = record.Reload(s.Conn)
	if err != nil {
		return nil, err
	}
	return record, err
}

// Update Companys
func (s *Company) Update(record *db.Company, txes ...*sql.Tx) (*db.Company, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, err
	}
	return record, nil
}
