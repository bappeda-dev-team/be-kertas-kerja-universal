package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/user"
	"ekak_kabupaten_madiun/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository    repository.UserRepository
	RoleRepository    repository.RoleRepository
	PegawaiRepository repository.PegawaiRepository
	DB                *sql.DB
}

func NewUserServiceImpl(userRepository repository.UserRepository, roleRepository repository.RoleRepository, pegawaiRepository repository.PegawaiRepository, db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository:    userRepository,
		RoleRepository:    roleRepository,
		PegawaiRepository: pegawaiRepository,
		DB:                db,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request user.UserCreateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi input dasar
	if request.Nip == "" {
		return user.UserResponse{}, errors.New("nip harus diisi")
	}
	if request.Email == "" {
		return user.UserResponse{}, errors.New("email harus diisi")
	}
	if request.Password == "" {
		return user.UserResponse{}, errors.New("password harus diisi")
	}
	if len(request.Role) == 0 {
		return user.UserResponse{}, errors.New("role harus diisi")
	}

	// Validasi NIP dengan data pegawai
	_, err = service.PegawaiRepository.FindByNip(ctx, tx, request.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.UserResponse{}, errors.New("nip tidak terdaftar di data pegawai")
		}
		return user.UserResponse{}, err
	}

	// Ambil data role pertama
	role, err := service.RoleRepository.FindById(ctx, tx, request.Role[0].RoleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return user.UserResponse{}, errors.New("role tidak ditemukan")
		}
		return user.UserResponse{}, err
	}

	userDomain := domain.Users{
		Nip:      request.Nip,
		Email:    request.Email,
		Password: request.Password,
		IsActive: request.IsActive,
		Role:     role, // Menggunakan struct Roles langsung
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.UserResponse{}, err
	}
	userDomain.Password = string(hashedPassword)

	createdUser, err := service.UserRepository.Create(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	roleResponse := user.RoleResponse{
		Id:   createdUser.Role.Id,
		Role: createdUser.Role.Role,
	}

	response := user.UserResponse{
		Id:       createdUser.Id,
		Nip:      createdUser.Nip,
		Email:    createdUser.Email,
		IsActive: createdUser.IsActive,
		Role:     []user.RoleResponse{roleResponse},
	}

	return response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request user.UserUpdateRequest) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	existingUser, err := service.UserRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return user.UserResponse{}, err
	}
	if existingUser.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	userDomain := domain.Users{
		Id:       existingUser.Id,
		Nip:      request.Nip,
		Email:    request.Email,
		Password: request.Password,
		IsActive: request.IsActive,
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return user.UserResponse{}, err
		}
		userDomain.Password = string(hashedPassword)
	} else {
		userDomain.Password = existingUser.Password
	}

	updatedUser, err := service.UserRepository.Update(ctx, tx, userDomain)
	if err != nil {
		return user.UserResponse{}, err
	}

	var roles []user.RoleResponse
	// TODO: Implement role update and conversion

	response := user.UserResponse{
		Id:       updatedUser.Id,
		Nip:      updatedUser.Nip,
		Email:    updatedUser.Email,
		IsActive: updatedUser.IsActive,
		Role:     roles,
	}

	return response, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	existingUser, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if existingUser.Id == 0 {
		return errors.New("user tidak ditemukan")
	}

	err = service.UserRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var userResponses []user.UserResponse
	for _, u := range users {
		var roles []user.RoleResponse

		userResponse := user.UserResponse{
			Id:       u.Id,
			Nip:      u.Nip,
			Email:    u.Email,
			IsActive: u.IsActive,
			Role:     roles,
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

func (service *UserServiceImpl) FindById(ctx context.Context, id int) (user.UserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return user.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cari user berdasarkan ID
	userDomain, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return user.UserResponse{}, err
	}

	// Cek apakah user ditemukan
	if userDomain.Id == 0 {
		return user.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// TODO: Ambil role untuk user
	var roles []user.RoleResponse

	// Convert ke response
	response := user.UserResponse{
		Id:       userDomain.Id,
		Nip:      userDomain.Nip,
		Email:    userDomain.Email,
		IsActive: userDomain.IsActive,
		Role:     roles,
	}

	return response, nil
}
