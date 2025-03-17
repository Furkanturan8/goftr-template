package dto

type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"active"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}
