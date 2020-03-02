package boilerplate

import (
	"boilerplate/dataloaders"
	"boilerplate/db"
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// ContextKey holds a custom String func for uniqueness
type ContextKey string

func (k ContextKey) String() string {
	return "dataloader_" + string(k)
}

// UserLoaderKey declares a statically typed key for context reference in other packages
const UserLoaderKey ContextKey = "user_loader"

// CompanyLoaderKey declares a statically typed key for context reference in other packages
const CompanyLoaderKey ContextKey = "Company_loader"

// CompanyUsersLoaderKey declares a statically typed key for context reference in other packages
const CompanyUsersLoaderKey ContextKey = "Company_users_loader"

// UserNotesLoaderKey declares a statically typed key for context reference in other packages
const UserNotesLoaderKey ContextKey = "user_Notes_loader"

// NoteLoaderKey declares a statically typed key for context reference in other packages
const NoteLoaderKey ContextKey = "Note_loader"

// CompanyUsersLoaderFromContext runs the dataloader inside the context
func CompanyUsersLoaderFromContext(ctx context.Context, id uuid.UUID) ([]*db.User, error) {
	return ctx.Value(CompanyUsersLoaderKey).(*dataloaders.CompanyUsersLoader).Load(id.String())
}

// CompanyLoaderFromContext runs the dataloader inside the context
func CompanyLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Company, error) {
	return ctx.Value(CompanyLoaderKey).(*dataloaders.CompanyLoader).Load(id.String())
}

// UserLoaderFromContext runs the dataloader inside the context
func UserLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return ctx.Value(UserLoaderKey).(*dataloaders.UserLoader).Load(id.String())
}

// UserNotesLoaderFromContext runs the dataloader inside the context
func UserNotesLoaderFromContext(ctx context.Context, id uuid.UUID) ([]*db.Note, error) {
	return ctx.Value(UserNotesLoaderKey).(*dataloaders.UserNotesLoader).Load(id.String())
}

// NoteLoaderFromContext runs the dataloader inside the context
func NoteLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Note, error) {
	return ctx.Value(NoteLoaderKey).(*dataloaders.NoteLoader).Load(id.String())
}

// WithDataloaders returns a new context that contains dataloaders
func WithDataloaders(
	ctx context.Context,
	NoteStorer NoteStorer,
	CompanyStorer CompanyStorer,
	UserStorer UserStorer,
) context.Context {
	userloader := dataloaders.NewUserLoader(
		dataloaders.UserLoaderConfig{
			Fetch: func(ids []string) ([]*db.User, []error) {
				return UserStorer.GetManyByIDs(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	CompanyLoader := dataloaders.NewCompanyLoader(
		dataloaders.CompanyLoaderConfig{
			Fetch: func(ids []string) ([]*db.Company, []error) {
				return CompanyStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	NoteLoader := dataloaders.NewNoteLoader(
		dataloaders.NoteLoaderConfig{
			Fetch: func(ids []string) ([]*db.Note, []error) {
				return NoteStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)
	userNoteLoader := dataloaders.NewUserNotesLoader(
		dataloaders.UserNotesLoaderConfig{
			Fetch: func(ids []string) ([][]*db.Note, []error) {
				result := [][]*db.Note{}
				for _, id := range ids {
					userID, err := uuid.FromString(id)
					if err != nil {
						return nil, []error{ErrParse}
					}
					records, err := NoteStorer.GetByUser(userID)
					if err != nil {
						return nil, []error{ErrDataloader}
					}
					result = append(result, records)
				}
				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)
	CompanyUsersLoader := dataloaders.NewCompanyUsersLoader(
		dataloaders.CompanyUsersLoaderConfig{
			Fetch: func(ids []string) ([][]*db.User, []error) {
				result := [][]*db.User{}
				for _, id := range ids {
					userID, err := uuid.FromString(id)
					if err != nil {
						return nil, []error{ErrParse}
					}
					records, err := UserStorer.GetByCompany(userID)
					if err != nil {
						return nil, []error{ErrDataloader}
					}
					result = append(result, records)
				}
				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)
	ctx = context.WithValue(ctx, UserLoaderKey, userloader)
	ctx = context.WithValue(ctx, CompanyLoaderKey, CompanyLoader)
	ctx = context.WithValue(ctx, CompanyUsersLoaderKey, CompanyUsersLoader)
	ctx = context.WithValue(ctx, NoteLoaderKey, NoteLoader)
	ctx = context.WithValue(ctx, UserNotesLoaderKey, userNoteLoader)
	return ctx
}

// DataloaderMiddleware runs before each API call and loads the dataloaders into context
func DataloaderMiddleware(
	ts NoteStorer,
	os CompanyStorer,
	us UserStorer,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithDataloaders(r.Context(), ts, os, us))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
