package infinote

import (
	"context"
	"fmt"
	"infinote/canlog"
	"infinote/db"

	"github.com/gofrs/uuid"
)

// UserNotesSelect gets some Notes
func UserNotesSelect(ctx context.Context, ts NoteStorer, ownerID uuid.UUID, limit, offset int) ([]*db.Note, error) {
	result, err := ts.Select(ownerID, limit, offset)
	if err != nil {
		canlog.AppendErr(ctx, "1202481c-9040-49a1-bde9-d19e1a74febe")
		return nil, fmt.Errorf("user Note select: %w", err)
	}
	return result, nil
}

// Notes gets all Notes
func Notes(ctx context.Context, ts NoteStorer, ownerID uuid.UUID) ([]*db.Note, error) {
	result, err := ts.GetByUser(ownerID)

	if err != nil {
		canlog.AppendErr(ctx, "dae3d21f-921a-4232-89e3-ff00b94d1230")
		return nil, fmt.Errorf("list Note: %w", err)
	}
	return result, nil
}

// UserNotes gets all Notes for a user
func UserNotes(ctx context.Context, ownerID uuid.UUID) ([]*db.Note, error) {
	result, err := UserNotesLoaderFromContext(ctx, ownerID)
	if err != nil {
		canlog.AppendErr(ctx, "84fd8b01-9af9-419f-b72a-bbb50fa79fab")
		return nil, fmt.Errorf("list user Note: %w", err)
	}
	return result, nil
}

// NoteCreate inserts the Note item in the DB
func NoteCreate(ctx context.Context, ts NoteStorer, name, body string) (*db.Note, error) {

	ownerID, err := ClaimValueFromContext(ctx, ClaimUserID)
	if err != nil {
		canlog.AppendErr(ctx, "9dd9be6e-ec54-4e43-ab09-90efe7061b67")
		return nil, fmt.Errorf("claim value: %w", err)
	}
	t := &db.Note{
		OwnerID: ownerID,
		Name:    name,
		Body:    body,
	}
	created, err := ts.Insert(t)
	if err != nil {
		canlog.AppendErr(ctx, "bcde6718-03f5-46bf-b938-2df5ed70c080")
		return nil, fmt.Errorf("create Note: %w", err)
	}
	return created, nil
}

// NoteUpdate updates the Note item in the DB
func NoteUpdate(ctx context.Context, ts NoteStorer, id uuid.UUID, name string) (*db.Note, error) {
	t, err := ts.Get(id)
	if err != nil {
		canlog.AppendErr(ctx, "8cb2def3-490b-4c57-8dc9-fb720533e957")
		return nil, fmt.Errorf("update Note: %w", err)
	}

	t.Name = name
	updated, err := ts.Update(t)
	if err != nil {
		canlog.AppendErr(ctx, "c558de61-6d5e-4519-bc79-a93737d92afa")
		return nil, fmt.Errorf("update Note: %w", err)
	}
	return updated, nil
}
