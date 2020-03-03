package graphql

import (
	"context"
	"errors"
	"infinote"
	"infinote/canlog"
	"infinote/db"
	"log"
	"sync"

	"github.com/gofrs/uuid"
)

func insertIntoString(s, is string, index int) string {
	return s[:index] + is + s[index:]
}

func (r *queryResolver) Notes(ctx context.Context) ([]*db.Note, error) {
	result, err := infinote.Notes(ctx, r.NoteStorer)
	if errors.Is(err, infinote.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}
func (r *queryResolver) NoteByID(ctx context.Context, noteID string) (*db.Note, error) {
	if noteID == "" {
		user, err := infinote.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
		if err != nil {
			return nil, err
		}
		return r.newNote(user.ID)
	}

	r.RLock()
	defer r.RUnlock()
	noteRoom, ok := r.noterooms[noteID]
	if ok {
		return noteRoom.Note, nil
	}

	noteUUID, err := uuid.FromString(noteID)
	if err != nil {
		canlog.AppendErr(ctx, "6d525675-3336-4f41-aa63-40889e63c96b")
		return nil, infinote.ErrParse
	}
	note, err := r.NoteStorer.Get(noteUUID)
	return note, nil
}

func (r *mutationResolver) NoteCreate(ctx context.Context, input CreateNote) (*db.Note, error) {
	return infinote.NoteCreate(ctx, r.NoteStorer, input.Name, input.Body)
}
func (r *mutationResolver) NoteUpdate(ctx context.Context, input UpdateNote) (*db.Note, error) {
	id, err := uuid.FromString(input.ID)
	if err != nil {
		return nil, infinote.ErrParse
	}
	return infinote.NoteUpdate(ctx, r.NoteStorer, id, input.Text)
}

type noteResolver struct{ *Resolver }

func (r *noteResolver) Owner(ctx context.Context, obj *db.Note) (*db.User, error) {
	id, err := uuid.FromString(obj.OwnerID)
	if err != nil {
		return nil, infinote.ErrParse
	}
	result, err := infinote.User(ctx, id)
	if errors.Is(err, infinote.ErrUnauthorized) {
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

	user, err := infinote.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
	if err != nil {
		return result, err
	}

	noteID := input.NoteID
	eventID := input.EventID

	log.Println("SESSOIN ID", input.SessionID)

	event := &NoteEvent{
		NoteID:    noteID,
		EventID:   eventID,
		UserID:    user.ID,
		UserName:  user.Name,
		SessionID: input.SessionID,
	}

	r.RLock()
	defer r.RUnlock()
	noteRoom, ok := r.noterooms[noteID]
	if !ok {
		uid, err := uuid.FromString(noteID)
		if err != nil {
			return nil, err
		}
		note, err := r.NoteStorer.Get(uid)
		if err != nil {
			note, err = r.newNote(user.ID)
		}
		noteRoom, err = r.noteRoom(note)
		if err != nil {
			return nil, err
		}
	}
	noteRoom.Lock()
	defer noteRoom.Unlock()

	if input.Insert != nil {
		text := input.Insert.Text
		index := input.Insert.Index
		noteRoom.Note.Body = insertIntoString(noteRoom.Note.Body, text, index)
		event.Insert = &TextInsert{
			Text:  text,
			Index: index,
		}
	}

	if input.Replace != nil {
		text := input.Replace.Text
		index := input.Replace.Index
		end := input.Replace.End
		noteRoom.Note.Body = insertIntoString(noteRoom.Note.Body, text, index)
		event.Replace = &ReplaceTextNote{
			Text:  text,
			Index: index,
			End:   end,
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

	user, err := infinote.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)

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
		defer r.Unlock()
		delete(room.Observers, idStr)
	}()

	r.Lock()
	defer r.Unlock()
	room.Observers[idStr] = struct {
		UserID string
		Chan   chan *NoteEvent
	}{UserID: user.ID, Chan: events}

	return events, nil
}
