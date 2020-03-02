package graphql

import (
	"boilerplate"
	"boilerplate/canlog"
	"boilerplate/db"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
)

type Resolver struct {
	Auther               boilerplate.AuthProvider
	NoteStorer           boilerplate.NoteStorer
	RoleStorer           boilerplate.RoleStorer
	CompanyStorer        boilerplate.CompanyStorer
	UserStorer           boilerplate.UserStorer
	BlacklistProvider    boilerplate.BlacklistProvider
	SubscriptionResolver subscriptionResolver
	noterooms            map[string]*Noteroom
	*sync.RWMutex
}

type ResolverOpts struct {
	Auther               boilerplate.AuthProvider
	NoteStorer           boilerplate.NoteStorer
	RoleStorer           boilerplate.RoleStorer
	CompanyStorer        boilerplate.CompanyStorer
	UserStorer           boilerplate.UserStorer
	BlacklistProvider    boilerplate.BlacklistProvider
	SubscriptionResolver subscriptionResolver
}

func NewResolver(opts *ResolverOpts) *Resolver {
	return &Resolver{
		RoleStorer:        opts.RoleStorer,
		Auther:            opts.Auther,
		NoteStorer:        opts.NoteStorer,
		CompanyStorer:     opts.CompanyStorer,
		UserStorer:        opts.UserStorer,
		BlacklistProvider: opts.BlacklistProvider,
		RWMutex:           &sync.RWMutex{},
		noterooms:         map[string]*Noteroom{},
	}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Company() CompanyResolver {
	return &organisationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Note() NoteResolver {
	return &noteResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RequestToken(ctx context.Context, input *RequestToken) (string, error) {
	canlog.Set(ctx, "email", input.Email)
	err := boilerplate.ValidatePassword(ctx, r.UserStorer, input.Email, input.Password)
	if err != nil {
		return "", boilerplate.ErrBadCredentials
	}
	u, err := r.UserStorer.GetByEmail(input.Email)
	if err != nil {
		return "", boilerplate.ErrBadCredentials
	}
	//userID, err := uuid.FromString(u.ID)
	//if err != nil {
	//return "", boilerplate.ErrParse
	//}
	//role, err := r.RoleStorer.ByUser(userID)
	//if err != nil {
	//	return "", fmt.Errorf("get roles: %w", err)
	//}
	token, err := r.Auther.GenerateJWT(ctx, u, "")
	if err != nil {
		return "", boilerplate.ErrBadCredentials
	}
	return token, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	u, err := boilerplate.UserFromContext(ctx, r.UserStorer, r.BlacklistProvider)
	if err != nil {
		return nil, fmt.Errorf("user from context: %w", err)
	}
	return u, nil
}
func (r *queryResolver) Companys(ctx context.Context) ([]*db.Company, error) {
	result, err := boilerplate.Companys(ctx, r.CompanyStorer)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}
func (r *queryResolver) Users(ctx context.Context) ([]*db.User, error) {
	result, err := boilerplate.Users(ctx, r.UserStorer)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}

type organisationResolver struct{ *Resolver }

func (r *organisationResolver) Users(ctx context.Context, obj *db.Company) ([]*db.User, error) {
	orgID, err := uuid.FromString(obj.ID)
	if err != nil {
		return nil, err
	}

	result, err := boilerplate.UsersByCompanyID(ctx, orgID)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) NotesConnection(ctx context.Context, obj *db.User, limit int, offset int) (*NotesConnection, error) {
	userUUID, err := uuid.FromString(obj.ID)
	if err != nil {
		canlog.AppendErr(ctx, "74dd216c-220f-4247-8871-e82f5a80a8ec")
		return nil, boilerplate.ErrParse
	}
	result, err := boilerplate.UserNotesSelect(ctx, r.NoteStorer, userUUID, limit, offset)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		canlog.AppendErr(ctx, "79c6f9dc-56b3-4407-aeac-8bbd015087c4")
		return nil, boilerplate.ErrUnauthorized
	}
	totalCount := len(result)
	pageInfo := &PageInfo{
		StartCursor: result[0].ID,
		EndCursor:   result[len(result)-1].ID,
	}
	edges := []*NotesEdge{}
	for _, node := range result {
		edges = append(edges, &NotesEdge{
			Cursor: node.ID,
			Node:   node,
		})
	}
	conn := &NotesConnection{
		TotalCount: totalCount,
		PageInfo:   pageInfo,
		Edges:      edges,
	}
	return conn, nil
}
func (r *userResolver) Company(ctx context.Context, obj *db.User) (*db.Company, error) {
	id, err := uuid.FromString(obj.CompanyID)
	if err != nil {
		return nil, boilerplate.ErrParse
	}
	result, err := boilerplate.Company(ctx, r.UserStorer, r.BlacklistProvider, id)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}
func (r *userResolver) Notes(ctx context.Context, obj *db.User) ([]*db.Note, error) {
	id, err := uuid.FromString(obj.ID)
	if err != nil {
		return nil, err
	}
	result, err := boilerplate.UserNotes(ctx, id)
	if errors.Is(err, boilerplate.ErrUnauthorized) {
		return nil, nil
	}
	return result, nil
}
