package graphql

import (
	"boilerplate"
	"boilerplate/canlog"
	"boilerplate/db"
	"context"
	"errors"
	"sync"

	"github.com/gofrs/uuid"
)

func insertIntoString(s, is string, index int) string {
	return s[:index] + is + s[index:]
}

func (r *queryResolver) Notes(ctx context.Context) ([]*db.Note, error) {
	result, err := boilerplate.Notes(ctx, r.NoteStorer)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}
func (r *queryResolver) NoteByID(ctx context.Context, noteID string) (*db.Note, error) {
	if noteID == "" {
		user, err := boilerplate.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
		if err != nil {
			return nil, err
		}
		return r.newNote(user.ID)
	}
	noteUUID, err := uuid.FromString(noteID)
	if err != nil {
		canlog.AppendErr(ctx, "6d525675-3336-4f41-aa63-40889e63c96b")
		return nil, boilerplate.ErrParse
	}
	note, err := r.NoteStorer.Get(noteUUID)
	return note, nil
}

func (r *mutationResolver) NoteCreate(ctx context.Context, input CreateNote) (*db.Note, error) {
	return boilerplate.NoteCreate(ctx, r.NoteStorer, input.Name, input.Body)
}
func (r *mutationResolver) NoteUpdate(ctx context.Context, input UpdateNote) (*db.Note, error) {
	id, err := uuid.FromString(input.ID)
	if err != nil {
		return nil, boilerplate.ErrParse
	}
	return boilerplate.NoteUpdate(ctx, r.NoteStorer, id, input.Text)
}

type noteResolver struct{ *Resolver }

func (r *noteResolver) Owner(ctx context.Context, obj *db.Note) (*db.User, error) {
	id, err := uuid.FromString(obj.OwnerID)
	if err != nil {
		return nil, boilerplate.ErrParse
	}
	result, err := boilerplate.User(ctx, id)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}

type Noteroom struct {
	Note      *db.Note
	Observers map[string]struct {
		UserID string
		Chan   chan *NoteEvent
	}
	*sync.RWMutex
}

func (r *mutationResolver) NoteChange(ctx context.Context, input NoteChange) (*NoteEventResult, error) {
	result := &NoteEventResult{
		Success: false,
	}

	user, err := boilerplate.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
	if err != nil {
		return result, err
	}

	noteID := input.NoteID
	eventID := input.EventID

	event := &NoteEvent{
		NoteID:   noteID,
		EventID:  eventID,
		UserID:   user.ID,
		UserName: user.Name,
	}

	r.RLock()
	noteRoom := r.noterooms[noteID]
	noteRoom.Lock()

	if input.Insert != nil {
		text := input.Insert.Text
		index := input.Insert.Index
		noteRoom.Note.Body = insertIntoString(noteRoom.Note.Body, text, index)
		event.Insert = &TextInsert{
			Text:  text,
			Index: index,
		}
	}

	if input.Cursor != nil {
		cursor := input.Cursor
		event.Cursor = &CursorPlacement{
			LineNumber: cursor.LineNumber,
			Column:     cursor.Column,
		}
	}

	for _, observer := range noteRoom.Observers {
		observer.Chan <- event
	}

	result.Success = true

	noteRoom.Unlock()

	r.RUnlock()
	return result, nil
}

func (r *Resolver) noteRoom(note *db.Note) (*Noteroom, error) {
	r.Lock()
	defer r.Unlock()
	if note.ID == "" {
		return nil, errors.New("note does not exist")
	}

	room := r.noterooms[note.ID]
	if room == nil {
		room = &Noteroom{
			Note: note,
			Observers: map[string]struct {
				UserID string
				Chan   chan *NoteEvent
			}{},
			RWMutex: &sync.RWMutex{},
		}
		r.noterooms[note.ID] = room
	}
	return room, nil
}

func (r *Resolver) newNote(userID string) (*db.Note, error) {
	note := &db.Note{
		OwnerID: userID,
		Name:    "",
		Body:    "",
	}

	return r.NoteStorer.Insert(note)
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) NoteEvent(ctx context.Context, noteID string) (<-chan *NoteEvent, error) {
	var note *db.Note

	user, err := boilerplate.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)

	uid, err := uuid.FromString(noteID)
	if err != nil {
		note, err = r.newNote(user.ID)
	} else {
		note, err = r.NoteStorer.Get(uid)
		if err != nil {
			note, err = r.newNote(user.ID)
		}
	}

	room, err := r.noteRoom(note)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	idStr := id.String()
	if err != nil {
		//TODO: find out why would newv4 ever return an error?
		panic(err)
	}

	events := make(chan *NoteEvent, 1)

	go func() {
		<-ctx.Done()
		r.Lock()
		delete(room.Observers, idStr)
		r.Unlock()
	}()

	r.Lock()
	defer r.Unlock()
	room.Observers[idStr] = struct {
		UserID string
		Chan   chan *NoteEvent
	}{UserID: user.ID, Chan: events}

	return events, nil
}
