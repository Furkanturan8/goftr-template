package service

import (
	"context"
	"goftr-v1/backend/internal/dto"
	"goftr-v1/backend/internal/repository"
	"goftr-v1/backend/pkg/errorx"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) List(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			Status:    string(user.Status),
		}
	}

	return &dto.UsersResponse{
		Users: userResponses,
		Total: int64(len(users)),
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.ErrNotFound
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		Status:    string(user.Status),
	}, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req *dto.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.ErrNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		// Email değişiyorsa, yeni email'in başka bir kullanıcıda olmadığından emin ol
		if req.Email != user.Email {
			exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
			if err != nil {
				return errorx.ErrDatabaseOperation
			}
			if exists {
				return errorx.ErrDuplicate
			}
		}
		user.Email = req.Email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.ErrDatabaseOperation
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return errorx.ErrDatabaseOperation
	}
	return nil
}
