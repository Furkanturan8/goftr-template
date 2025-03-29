package dto

import (
	"goftr-v1/backend/internal/model"
)

type UserCreateDTO struct {
	Email     string       `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	FirstName string       `json:"first_name" validate:"required,max=100"`
	LastName  string       `json:"last_name" validate:"required,max=100"`
	Password  string       `json:"password" validate:"required,min=3,max=100"`
	Status    model.Status `json:"status" validate:"omitempty,oneof=active inactive"`
	Role      model.Role   `json:"role" validate:"required"`
}

func (vm UserCreateDTO) ToDBModel(m model.User) model.User {
	m.Email = vm.Email
	m.FirstName = vm.FirstName
	m.LastName = vm.LastName
	if vm.Password != "" {
		_ = m.SetPassword(vm.Password)
	}
	m.Role = vm.Role
	if vm.Status == "" {
		m.Status = model.StatusActive
	} else {
		m.Status = vm.Status
	}

	return m
}

type UserResponseDTO struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"active"`
}

func (vm UserResponseDTO) ToResponseModel(m model.User) UserResponseDTO {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.FirstName = m.FirstName
	vm.LastName = m.LastName
	vm.Role = string(m.Role)
	vm.Status = string(m.Status)

	return vm
}
