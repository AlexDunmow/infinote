package store

import (
	"infinote/db"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

// RoleSuperAdmin role for authorization
var RoleSuperAdmin = "SUPERADMIN"

// RoleCompanyAdmin role for authorization
var RoleCompanyAdmin = "CompanyADMIN"

// RoleMember role for authorization
var RoleMember = "MEMBER"

// Role for persistence
type Role struct {
	Conn *sqlx.DB
}

// NewRoleStore returns a new store
func NewRoleStore(conn *sqlx.DB) *Role {
	rs := &Role{conn}
	return rs
}

// RoleByName returns the role given its name
func (s *Role) RoleByName(name string) (*db.Role, error) {
	return db.Roles(db.RoleWhere.Name.EQ(name)).One(s.Conn)
}

// ByUser roles
func (s *Role) ByUser(id uuid.UUID, tx ...*sql.Tx) (*db.Role, error) {
	u, err := db.FindUser(s.Conn, id.String())
	if err != nil {
		return nil, err
	}
	return u.Role().One(s.Conn)
}
