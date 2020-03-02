package store

import (
	"infinote/db"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Note for persistence
type Note struct {
	Conn *sqlx.DB
}

// NewNoteStore handle Note methods
func NewNoteStore(conn *sqlx.DB) *Note {
	ts := &Note{conn}
	return ts
}

// All Notes
func (s *Note) All() (db.NoteSlice, error) {
	return db.Notes().All(s.Conn)
}

// GetMany Notes
func (s *Note) GetMany(keys []string) (db.NoteSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	args := []interface{}{}
	for _, key := range keys {
		args = append(args, key)
	}
	records, err := db.Notes(qm.WhereIn("id in ?", args...)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Note{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Note{}
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

// GetByUser Notes by user
func (s *Note) GetByUser(userID uuid.UUID, txes ...*sql.Tx) ([]*db.Note, error) {
	return db.Notes(db.NoteWhere.OwnerID.EQ(userID.String())).All(s.Conn)
}

// Select Notes
func (s *Note) Select(userID uuid.UUID, limit int, offset int) ([]*db.Note, error) {
	return db.Notes(db.NoteWhere.OwnerID.EQ(userID.String()),
		qm.Limit(limit),
		qm.Offset(offset),
		qm.OrderBy(db.NoteColumns.CreatedAt),
	).All(s.Conn)
}

// Get Notes
func (s *Note) Get(id uuid.UUID) (*db.Note, error) {
	return db.FindNote(s.Conn, id.String())
}

// Insert Notes
func (s *Note) Insert(record *db.Note, txes ...*sql.Tx) (*db.Note, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, err
	}
	return record, nil
}

// Update Notes
func (s *Note) Update(record *db.Note, txes ...*sql.Tx) (*db.Note, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, err
	}
	return record, nil
}
