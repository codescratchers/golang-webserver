package user

import (
	"database/sql"
)

type IRoleRepository interface {
	Save(tx *sql.Tx, role Role) error
}

type roleRepository struct{}

func NewRoleRepository() IRoleRepository {
	return roleRepository{}
}

func (r roleRepository) Save(tx *sql.Tx, role Role) error {
	_, err := tx.Exec("INSERT INTO role (role, user_id) VALUES (?, ?)", role.Role, role.UserId)
	return err
}
