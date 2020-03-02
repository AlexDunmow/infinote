//go:generate go run github.com/vektah/dataloaden UserLoader string *boilerplate/db.User
//go:generate go run github.com/vektah/dataloaden NoteLoader string *boilerplate/db.Note
//go:generate go run github.com/vektah/dataloaden CompanyLoader string *boilerplate/db.Company
//go:generate go run github.com/vektah/dataloaden CompanyUsersLoader string []*boilerplate/db.User
//go:generate go run github.com/vektah/dataloaden UserNotesLoader string []*boilerplate/db.Note

package dataloaders
