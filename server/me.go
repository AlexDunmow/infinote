package infinote

import (
	"boilerplate/crypto"
	"boilerplate/db"
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	gocheckpasswd "github.com/ninja-software/go-check-passwd"
)

// ChangePassword of a user
func ChangePassword(ctx context.Context, bs BlacklistProvider, us UserStorer, id uuid.UUID, oldPassword, newPassword string) error {
	// make sure password met minimum length
	if len(newPassword) < 8 {
		return Err("5d46705b-505c-4976-bcc3-ccd324095e20", "New password is too short", errors.New("password too short"), KindInput)
	}

	// make sure user not using common password
	if gocheckpasswd.IsCommon(newPassword) {
		return Err("5dbb7fb4-b39f-4c5e-8607-e07a714daccc", "New password is too common", errors.New("password too common"), KindInput)
	}
	user, err := us.Get(id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	// When user first verifies email, they are sent to set their password for the first time, does not need old password
	// After this, old password should be required to set new password.
	if user.RequireOldPassword {
		err := ValidatePassword(ctx, us, user.Email, oldPassword)
		if err != nil {
			return fmt.Errorf("validate password: %w", err)
		}
	}

	// Set value to true, as initial password set without old password will only be allowed once.
	if !user.RequireOldPassword {
		user.RequireOldPassword = true
	}

	hashed := crypto.HashPassword(newPassword)

	user.PasswordHash = hashed

	_, err = us.Update(user)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	// Blacklist all old tokens
	err = bs.BlacklistAll(user.ID)
	if err != nil {
		return fmt.Errorf("get blacklist: %w", err)
	}

	return nil
}

// ChangeDetails of a user
func ChangeDetails(us UserStorer, id uuid.UUID, name string) (*db.User, error) {
	if name == "" {
		err := fmt.Errorf("found blank name field when trying to update details")
		return nil, Err("30c82163-d691-4ce0-9801-95c7999d9b7a", "Your name must not be blank", err, KindInput, "userID", id.String())
	}
	user, err := us.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	user.Name = name
	_, err = us.Update(user)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}
