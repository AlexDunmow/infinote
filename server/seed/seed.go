package seed

import (
	"boilerplate/crypto"
	"boilerplate/db"
	"boilerplate/graphql"
	"boilerplate/store"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"
)

// Run for database spinup
func Run(conn *sqlx.DB) error {
	var err error
	fmt.Println("Seeding roles")
	err = Roles(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println("Seeding Companys")
	err = Companys(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println("Seeding users")
	err = Users(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println("Seed complete")
	return nil
}

// Roles for database spinup
func Roles(conn *sqlx.DB) error {
	allPerms := types.StringArray{}
	for _, perm := range graphql.AllPerm {
		allPerms = append(allPerms, string(perm))
	}
	r := &db.Role{Name: store.RoleSuperAdmin, Permissions: allPerms}
	err := r.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	r2 := &db.Role{
		Name: store.RoleCompanyAdmin,
		Permissions: types.StringArray{
			string(graphql.PermNoteList),
			string(graphql.PermNoteCreate),
			string(graphql.PermNoteRead),
			string(graphql.PermNoteUpdate),
			string(graphql.PermNoteArchive),
			string(graphql.PermNoteUnarchive),
			string(graphql.PermUserList),
			string(graphql.PermUserCreate),
			string(graphql.PermUserRead),
			string(graphql.PermUserUpdate),
			string(graphql.PermUserArchive),
			string(graphql.PermUserUnarchive),
			string(graphql.PermCompanyRead),
		},
	}
	err = r2.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	r3 := &db.Role{
		Name: store.RoleMember,
		Permissions: types.StringArray{
			string(graphql.PermNoteList),
			string(graphql.PermNoteCreate),
			string(graphql.PermNoteRead),
			string(graphql.PermNoteUpdate),
			string(graphql.PermNoteArchive),
			string(graphql.PermNoteUnarchive),
			string(graphql.PermUserRead),
			string(graphql.PermUserUpdate),
		},
	}
	err = r3.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Companys for database spinup
func Companys(conn *sqlx.DB) error {
	for i := 0; i < 5; i++ {
		o := store.CompanyFactory()
		err := o.Insert(conn, boil.Infer())
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// Users for database spinup
func Users(conn *sqlx.DB) error {
	orgs, err := db.Companies().All(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	rs := store.NewRoleStore(conn)
	superAdmin, err := rs.RoleByName(store.RoleSuperAdmin)
	if err != nil {
		return fmt.Errorf("get user role: %w", err)
	}
	orgAdmin, err := rs.RoleByName(store.RoleCompanyAdmin)
	if err != nil {
		return fmt.Errorf("get user role: %w", err)
	}
	member, err := rs.RoleByName(store.RoleMember)
	if err != nil {
		return fmt.Errorf("get user role: %w", err)
	}
	for i, org := range orgs {
		if i == 0 {
			fmt.Println("insert superadmin user")
			u := store.UserFactory()
			u.Email = "alex@ninjasoftware.com.au"
			u.PasswordHash = crypto.HashPassword("devdev!")
			u.CompanyID = org.ID
			u.RoleID = superAdmin.ID
			u.Name = "Alex Dunmow"
			err := u.Insert(conn, boil.Infer())
			if err != nil {
				return fmt.Errorf("insert user: %w", err)
			}

			fmt.Println("insert orgadmin user")
			u2 := store.UserFactory()
			u2.Email = "john@ninjasoftware.com.au"
			u2.Name = "John Nguyen"
			u2.PasswordHash = crypto.HashPassword("devdev!")
			u2.CompanyID = org.ID
			u2.RoleID = orgAdmin.ID
			err = u2.Insert(conn, boil.Infer())
			if err != nil {
				return fmt.Errorf("insert user: %w", err)
			}

			fmt.Println("insert member user")
			u3 := store.UserFactory()
			u3.Email = "member@example.com"
			u3.Name = "Lan Tran"
			u3.PasswordHash = crypto.HashPassword("devdev!")
			u3.CompanyID = org.ID
			u3.RoleID = member.ID
			err = u3.Insert(conn, boil.Infer())
			if err != nil {
				return fmt.Errorf("insert user: %w", err)
			}

		}

		for i := 0; i < 5; i++ {
			u := store.UserFactory()
			u.CompanyID = org.ID
			u.RoleID = member.ID
			err := u.Insert(conn, boil.Infer())
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}
