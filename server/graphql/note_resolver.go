package graphql

import (
	"context"
	"errors"
	"fmt"
	"infinote"
	"infinote/canlog"
	"infinote/db"
	"log"
	"sync"

	"github.com/gofrs/uuid"
)

func insertIntoString(s, is string, index int) string {
	runes := []rune(s)
	return string(runes[:index]) + is + string(runes[index:])
}

func replaceIntoString(s, rs string, index, length int) string {
	runes := []rune(s)

	news := string(runes[:index]) + rs + string(runes[index+length:])

	return news
}

func (r *queryResolver) Notes(ctx context.Context) ([]*db.Note, error) {
	user, err := infinote.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.FromString(user.ID)
	if err != nil {
		return nil, err
	}
	result, err := infinote.Notes(ctx, r.NoteStorer, uid)

	if errors.Is(err, infinote.ErrUnauthorized) {
		return nil, nil
	}
	return result, err
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

type Observer struct {
	UserID         string
	CursorPosition CursorPlacement
	Chan           chan *NoteEvent
}

type Noteroom struct {
	Note      *db.Note
	Observers map[string]*Observer
	*sync.RWMutex
}

var ErrNoNoteroom = errors.New("no noteroom created")

func (r *mutationResolver) NoteChange(ctx context.Context, input NoteChange) (*NoteEventResult, error) {
	fmt.Println("CHANGES", input.Cursor)
	result := &NoteEventResult{
		Success: false,
	}

	user, err := infinote.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
	if err != nil {
		return result, err
	}

	noteID := input.NoteID
	eventID := input.EventID
	sessionID := input.SessionID

	log.Println("SessionID", input.SessionID)

	event := &NoteEvent{
		NoteID:    noteID,
		EventID:   eventID,
		UserID:    user.ID,
		UserName:  user.Name,
		SessionID: input.SessionID,
	}

	fmt.Println("Locking resolver. Changing.", input.Cursor)
	r.RLock()
	defer func() {
		fmt.Println("Unlocking resolver. Changing.")
		r.RUnlock()
	}()
	noteRoom, ok := r.noterooms[noteID]
	if !ok {
		fmt.Println("No room.", noteID, len(r.noterooms))
		return nil, ErrNoNoteroom
	}

	if input.Cursor != nil {
		obv, ok := noteRoom.Observers[sessionID]
		if !ok {
			panic(err)
		}
		obv.CursorPosition.LineNumber = input.Cursor.LineNumber
		obv.CursorPosition.Column = input.Cursor.Column
		event.Cursor = &CursorPlacement{
			LineNumber: obv.CursorPosition.LineNumber,
			Column:     obv.CursorPosition.Column,
		}
	}

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
		length := input.Replace.Length
		noteRoom.Note.Body = replaceIntoString(noteRoom.Note.Body, text, index, length)
		event.Replace = &ReplaceTextNote{
			Text:   text,
			Index:  index,
			Length: length,
		}
	}

	if input.Remove != nil {
		index := input.Remove.Index
		length := input.Remove.Length
		noteRoom.Note.Body = replaceIntoString(noteRoom.Note.Body, "", index, length)
		event.Remove = &DeleteTextNote{
			Index:  index,
			Length: length,
		}
	}

	_, err = r.NoteStorer.Update(noteRoom.Note)
	if err != nil {
		panic(err)
	}

	// no more actual text changes

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

func (r *Resolver) newNote(userID string) (*db.Note, error) {
	note := &db.Note{
		OwnerID: userID,
		Name:    "",
		Body:    "",
	}

	return r.NoteStorer.Insert(note)
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) NoteEvent(ctx context.Context, noteID, sessionID string) (<-chan *NoteEvent, error) {
	var note *db.Note

	fmt.Println("Subscribing...", sessionID)

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

	fmt.Println("Checking for note", note.ID)

	if note == nil {
		fmt.Println("No note found")
		return nil, errors.New("no note found")
	}

	events := make(chan *NoteEvent, 1)

	go func() {
		<-ctx.Done()
		fmt.Println("Locking. Leaving.")
		r.Lock()
		room, _ := r.noterooms[noteID]
		defer func() {
			fmt.Println("Unlocking. Leaving")
			r.Unlock()
		}()
		delete(room.Observers, sessionID)
	}()

	fmt.Println("Locking. Entering.")
	r.Lock()
	defer func() {
		fmt.Println("Unlocking. Entering")
		r.Unlock()
	}()
	room, ok := r.noterooms[noteID]
	if !ok {
		room = &Noteroom{
			Note:      note,
			Observers: map[string]*Observer{},
			RWMutex:   &sync.RWMutex{},
		}
		fmt.Println("Noteroom created")
		r.noterooms[note.ID] = room
	}

	room.Observers[sessionID] = &Observer{UserID: user.ID, Chan: events, CursorPosition: CursorPlacement{
		LineNumber: 0,
		Column:     0,
	}}

	return events, nil
}
