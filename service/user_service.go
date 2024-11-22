package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/user"
)

type UserService interface {
	Create(ctx context.Context, request user.UserCreateRequest) (user.UserResponse, error)
	Update(ctx context.Context, request user.UserUpdateRequest) (user.UserResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]user.UserResponse, error)
	FindById(ctx context.Context, id int) (user.UserResponse, error)
}
