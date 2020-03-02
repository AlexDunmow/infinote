package infinote

import (
	"boilerplate/canlog"
	"boilerplate/db"
	"context"
	"fmt"

	"github.com/gofrs/uuid"
)

// Users resolver
func Users(ctx context.Context, us UserStorer) ([]*db.User, error) {
	result, err := us.All()
	if err != nil {
		canlog.AppendErr(ctx, "bf19a6af-aa27-46c7-9102-c2b8abeb9b4d")
		return nil, fmt.Errorf("list user: %w", err)
	}
	return result, nil
}

// User resolver
func User(ctx context.Context, id uuid.UUID) (*db.User, error) {
	result, err := UserLoaderFromContext(ctx, id)
	if err != nil {
		canlog.AppendErr(ctx, "dcfd7cf4-4d8b-4a1b-8d22-0bffe1d44490")
		return nil, fmt.Errorf("get user: %w", err)
	}
	return result, nil
}

// UsersByCompanyID gets users for a given Company ID
func UsersByCompanyID(ctx context.Context, orgID uuid.UUID) ([]*db.User, error) {
	result, err := CompanyUsersLoaderFromContext(ctx, orgID)
	if err != nil {
		canlog.AppendErr(ctx, "218e064a-e5c3-4b18-b570-4b1ea9425e8e")
		return nil, fmt.Errorf("list Company user: %w", err)
	}
	return result, nil
}

// UserUpdateDetails to change their name etc
func UserUpdateDetails(ctx context.Context, us UserStorer, id uuid.UUID, name string) (*db.User, error) {
	u, err := us.Get(id)
	if err != nil {
		canlog.AppendErr(ctx, "fef4879e-bc40-45ae-b7a8-2ab61fba5401")
		return nil, fmt.Errorf("get user: %w", err)
	}

	u.Name = name

	result, err := us.Update(u)
	if err != nil {
		canlog.AppendErr(ctx, "f7c5b7ff-73ec-49d8-84c1-a60564dfd1ad")
		return nil, fmt.Errorf("update user: %w", err)
	}
	return result, nil
}
