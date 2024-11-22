package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error) {
	script := "INSERT INTO tb_users(nip, email, password, is_active) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, users.Nip, users.Email, users.Password, users.IsActive)
	if err != nil {
		return users, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return users, err
	}
	users.Id = int(id)

	userRole := domain.UsersRole{
		UserId: users.Id,
		RoleId: 2,
	}

	scriptRole := "INSERT INTO tb_user_role(user_id, role_id) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, scriptRole, userRole.UserId, userRole.RoleId)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, users domain.Users) (domain.Users, error) {
	script := "UPDATE tb_users SET nip = ?, email = ?, password = ?, is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, users.Nip, users.Email, users.Password, users.IsActive, users.Id)
	if err != nil {
		return users, err
	}

	userRole := domain.UsersRole{
		UserId: users.Id,
		RoleId: 2,
	}

	scriptRole := "UPDATE tb_user_role SET role_id = ? WHERE user_id = ?"
	_, err = tx.ExecContext(ctx, scriptRole, userRole.RoleId, userRole.UserId)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Users, error) {
	script := `
		SELECT u.id, u.nip, u.email, u.is_active, ur.role_id, r.role 
		FROM tb_users u
		LEFT JOIN tb_user_role ur ON u.id = ur.user_id
		LEFT JOIN tb_role r ON ur.role_id = r.id
	`
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domain.Users{}, err
	}
	defer rows.Close()

	var users []domain.Users
	for rows.Next() {
		user := domain.Users{}
		userRole := domain.UsersRole{}
		var roleName string

		err := rows.Scan(
			&user.Id,
			&user.Nip,
			&user.Email,
			&user.IsActive,
			&userRole.RoleId,
			&roleName,
		)
		if err != nil {
			return []domain.Users{}, err
		}

		user.Role = domain.Roles{
			Id:   userRole.RoleId,
			Role: roleName,
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Users, error) {
	script := `
		SELECT u.id, u.nip, u.email, u.is_active, ur.role_id, r.role 
		FROM tb_users u
		LEFT JOIN tb_user_role ur ON u.id = ur.user_id
		LEFT JOIN tb_role r ON ur.role_id = r.id
		WHERE u.id = ?
	`
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.Users{}, err
	}
	defer rows.Close()

	user := domain.Users{}
	if rows.Next() {
		userRole := domain.UsersRole{}
		var roleName string

		err := rows.Scan(
			&user.Id,
			&user.Nip,
			&user.Email,
			&user.IsActive,
			&userRole.RoleId,
			&roleName,
		)
		if err != nil {
			return domain.Users{}, err
		}

		user.Role = domain.Roles{
			Id:   userRole.RoleId,
			Role: roleName,
		}
	}

	return user, nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	scriptRole := "DELETE FROM tb_user_role WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, scriptRole, id)
	if err != nil {
		return err
	}

	scriptUser := "DELETE FROM tb_users WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptUser, id)
	if err != nil {
		return err
	}

	return nil
}
