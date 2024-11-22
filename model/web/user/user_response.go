package user

type UserResponse struct {
	Id       int            `json:"id"`
	Nip      string         `json:"nip"`
	Email    string         `json:"email"`
	IsActive bool           `json:"is_active"`
	Role     []RoleResponse `json:"role"`
}
