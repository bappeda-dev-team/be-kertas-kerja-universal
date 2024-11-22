package user

type UserResponse struct {
	Id       int
	Nip      string
	Email    string
	IsActive bool
	Role     []RoleResponse
}
