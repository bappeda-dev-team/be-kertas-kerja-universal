package dataseeder

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/repository"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserSeederImpl struct {
	UserRepository repository.UserRepository
	RoleRepository repository.RoleRepository
}

func NewUserSeederImpl(userRepository repository.UserRepository, roleRepository repository.RoleRepository) *UserSeederImpl {
	return &UserSeederImpl{
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func (seeder *UserSeederImpl) Seed(ctx context.Context, tx *sql.Tx) error {
	users, err := seeder.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		log.Println("Users sudah ada, skip seeding users")
		return nil
	}

	roles, err := seeder.RoleRepository.FindAll(ctx, tx)
	if err != nil {
		return err
	}

	roleMap := make(map[string]domain.Roles)
	for _, role := range roles {
		roleMap[role.Role] = role
	}

	defaultUsers := []struct {
		user     domain.Users
		roleKeys []string
	}{
		{
			user: domain.Users{
				Nip:      "admin1",
				Email:    "admin@madiunkabtest.com",
				Password: "KabKabMadiun2024",
				IsActive: true,
			},
			roleKeys: []string{"super_admin"},
		},
		{
			user: domain.Users{
				Nip:      "admin2",
				Email:    "admin2@madiunkabtest.com",
				Password: "KabKabMadiun2024",
				IsActive: true,
			},
			roleKeys: []string{"super_admin"},
		},
	}

	for _, userData := range defaultUsers {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userData.user.Password = string(hashedPassword)

		for _, roleKey := range userData.roleKeys {
			if role, exists := roleMap[roleKey]; exists {
				userData.user.Role = append(userData.user.Role, role)
			}
		}

		_, err = seeder.UserRepository.Create(ctx, tx, userData.user)
		if err != nil {
			return err
		}
	}
	log.Println("Users berhasil di-seed")
	return nil
}
