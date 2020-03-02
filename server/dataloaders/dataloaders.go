//go:generate go run github.com/vektah/dataloaden UserLoader string *infinote/db.User
//go:generate go run github.com/vektah/dataloaden NoteLoader string *infinote/db.Note
//go:generate go run github.com/vektah/dataloaden CompanyLoader string *infinote/db.Company
//go:generate go run github.com/vektah/dataloaden CompanyUsersLoader string []*infinote/db.User
//go:generate go run github.com/vektah/dataloaden UserNotesLoader string []*infinote/db.Note

package dataloaders
