package store_test

import (
	"infinote/store"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"

	"github.com/gofrs/uuid"
)

func TestUserRepo_GetMany(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)

	t.Run("wrong ID", func(t *testing.T) {
		drop()
		migrate()
		randomID := uuid.Must(uuid.NewV4()).String()
		result, errs := repo.GetManyByIDs([]string{randomID})
		if errs != nil {
			t.Errorf("errs: got %v, expected %v", len(errs), 0)
		}
		if len(result) != 0 {
			t.Errorf("result count: got %v, expected %v", len(result), 0)
		}
	})
	t.Run("one record", func(t *testing.T) {
		drop()
		migrate()

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
			t.Error(err)
		}
		result, errs := repo.GetManyByIDs([]string{u.ID})
		if len(errs) != 0 {
			t.Errorf("errs: got %v, expected %v", len(errs), 0)
		}
		if len(result) != 1 {
			t.Errorf("result count: got %v, expected %v", len(result), 1)
		}
	})

	t.Run("another record", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u1 := store.UserFactory()
		u1.CompanyID = o.ID
		err = u1.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}
		u2 := store.UserFactory()
		u2.CompanyID = o.ID
		err = u2.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		result, errs := repo.GetManyByIDs([]string{u1.ID, u2.ID})
		if len(errs) != 0 {
			t.Errorf("errs: got %v, expected %v", len(errs), 0)
		}
		if len(result) != 2 {
			t.Errorf("result count: got %v, expected %v", len(result), 1)
		}
	})
	t.Run("nil keys", func(t *testing.T) {
		drop()
		migrate()
		result, errs := repo.GetManyByIDs(nil)
		if len(errs) == 0 {
			t.Errorf("errs: got %v, expected %v", len(errs), 0)
		}
		if result != nil {
			t.Errorf("result count: got %v, expected %v", len(result), 1)
		}
	})

}

func TestUserRepo_UserGetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)

	t.Run("no records", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		result, err := repo.GetByEmail("random@email.com")
		if err == nil {
			t.Error("expected error, got none")
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
	t.Run("existing records", func(t *testing.T) {
		drop()
		migrate()

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
			t.Error(err)
		}

		result, err := repo.GetByEmail(u.Email)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected user, got %v", nil)
		}

	})
	t.Run("bad email", func(t *testing.T) {
		drop()
		migrate()

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
			t.Error(err)
		}

		result, err := repo.GetByEmail("incorrect@email.com")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil user, got %v", result)
		}
	})
}

func TestUserRepo_UserList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("no records", func(t *testing.T) {
		drop()
		migrate()

		result, err := repo.All()
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected nil records, got %v", result)
		}
	})
	t.Run("existing records", func(t *testing.T) {
		drop()
		migrate()

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
			t.Error(err)
		}
		result, err := repo.All()
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected records, got %v", result)
		}
	})
}

func TestUserRepo_UserGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("no records", func(t *testing.T) {
		drop()
		migrate()

		randomID := uuid.Must(uuid.NewV4())
		result, err := repo.Get(randomID)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected nil record, got %v", result)
		}
	})
	t.Run("existing records", func(t *testing.T) {
		drop()
		migrate()

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
			t.Error(err)
		}
		for i := 0; i < 5; i++ {
			u := store.UserFactory()
			u.CompanyID = o.ID
			err = u.Insert(conn, boil.Infer())
			if err != nil {
				t.Error(err)
			}
		}
		result, err := repo.Get(uuid.Must(uuid.FromString(u.ID)))
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected records, got %v", result)
		}
	})
	t.Run("bad ID", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < 5; i++ {
			u := store.UserFactory()
			u.CompanyID = o.ID

			err = u.Insert(conn, boil.Infer())
			if err != nil {
				t.Error(err)
			}
		}
		randomID := uuid.Must(uuid.NewV4())
		result, err := repo.Get(randomID)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected nil record, got %v", result)
		}
	})
}

func TestUserRepo_UserCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("insert one", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID

		result, err := repo.Create(u)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result == nil {
			t.Errorf("expected record, got %v", result)
		}
	})
	t.Run("existing email", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID

		_, err = repo.Create(u)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		u2 := store.UserFactory()
		u2.CompanyID = o.ID
		u2.Email = u.Email
		_, err = repo.Create(u2)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserRepo_UserUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("update existing", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID
		created, err := repo.Create(u)
		created.Email = "changed@email.com"
		result, err := repo.Update(uuid.Must(uuid.FromString(u.ID)), created)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result.Email != "changed@email.com" {
			t.Errorf("expected changed@email.com, got %v", result.Email)
		}
	})
	t.Run("update non existing", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID
		_, err = repo.Update(uuid.Must(uuid.FromString(u.ID)), u)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
	t.Run("update email to already existing email", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID
		created1, err := repo.Create(u)
		if err != nil {
			t.Error(err)
		}

		u2 := store.UserFactory()
		u2.CompanyID = o.ID
		created2, err := repo.Create(u2)
		u2.Email = created1.Email
		_, err = repo.Update(uuid.Must(uuid.FromString(created2.ID)), created2)
		if err == nil {
			t.Errorf("expected non-nil error, got %v", err)
		}
	})
}

func TestUserRepo_UserArchive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("archive existing", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID

		created, err := repo.Create(u)
		result, err := repo.Archive(uuid.FromStringOrNil(created.ID))
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result.Archived != true {
			t.Errorf("expected result.Archived == true, got %v", result.Archived)
		}
		if !result.ArchivedAt.Valid {
			t.Errorf("expected result.ArchivedAt.Valid == true, got %v", result.ArchivedAt.Valid)
		}
	})
	t.Run("archive non existing", func(t *testing.T) {
		drop()
		migrate()

		_, err := repo.Archive(uuid.Must(uuid.NewV4()))
		if err == nil {
			t.Errorf("expected non-nil error, got %v", err)
		}
	})
}

func TestUserRepo_UserUnarchive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := store.NewUserStore(conn)
	t.Run("unarchive existing", func(t *testing.T) {
		drop()
		migrate()

		o := store.CompanyFactory()
		o.ID = uuid.Must(uuid.NewV4()).String()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			t.Error(err)
		}

		u := store.UserFactory()
		u.CompanyID = o.ID
		created, err := repo.Create(u)
		_, err = repo.Archive(uuid.FromStringOrNil(created.ID))
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		_, err = repo.Unarchive(uuid.FromStringOrNil(created.ID))
		result, err := repo.Get(uuid.FromStringOrNil(created.ID))
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if result.Archived != false {
			t.Errorf("expected result.Archived == false, got %v", result.Archived)
		}
		if result.ArchivedAt.Valid {
			t.Errorf("expected result.ArchivedAt.Valid == true, got %v", result.ArchivedAt.Valid)
		}
	})
	t.Run("unarchive non existing", func(t *testing.T) {
		drop()
		migrate()

		_, err := repo.Unarchive(uuid.Must(uuid.NewV4()))
		if err == nil {
			t.Errorf("expected non-nil error, got %v", err)
		}
	})
}
