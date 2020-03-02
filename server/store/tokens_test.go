package store_test

import (
	"infinote/store"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestTokenStore_Blacklist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()
		result, err := repo.Blacklist()
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestTokenStore_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		randomID := uuid.Must(uuid.NewV4())
		result, err := repo.Get(randomID.String())
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected nil record, got %v", result)
		}
	})
}

func TestTokenStore_GetAllByUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		randomID := uuid.Must(uuid.NewV4())
		result, err := repo.GetAllByUser(randomID.String())
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})
}

func TestTokenStore_GetAllExpired(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		result, err := repo.GetAllExpired()
		if err != nil {
			t.Errorf("expected error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})
}

func TestTokenStore_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		record := store.TokenFactory()
		u := store.UserFactory()
		u.CompanyID = o.ID
		err = u.Insert(conn, boil.Infer())
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		record.UserID = u.ID
		result, err := repo.Insert(record)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected nil record, got %v", result)
		}
	})
}

func TestTokenStore_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		record := store.TokenFactory()
		u := store.UserFactory()
		u.CompanyID = o.ID
		err = u.Insert(conn, boil.Infer())
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		record.UserID = u.ID
		record.CompanyID = u.CompanyID
		result, err := repo.Insert(record)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected record, got %v", result)
		}
		_, err = repo.Update(record)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}

func TestTokenStore_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.Tokens{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()
		record := store.TokenFactory()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID
		err = u.Insert(conn, boil.Infer())
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		record.UserID = u.ID
		record.CompanyID = u.CompanyID
		result, err := repo.Insert(record)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected record, got %v", result)
		}
		err = repo.Delete(record)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}
