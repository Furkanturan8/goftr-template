package service

import (
	"context"
	"goftr-v1/backend/internal/dto"
	"goftr-v1/backend/internal/model"
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

func (s *UserService) Create(ctx context.Context, req *dto.UserCreateDTO) error {
	// Email kontrolü
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return errorx.ErrDatabaseOperation
	}
	if exists {
		return errorx.ErrDuplicate
	}

	user := req.ToDBModel(model.User{})

	if err = s.userRepo.Create(ctx, &user); err != nil {
		return errorx.ErrDatabaseOperation
	}

	return nil
}

func (s *UserService) List(ctx context.Context) (*[]dto.UserResponseDTO, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, errorx.ErrDatabaseOperation
	}

	userResponses := make([]dto.UserResponseDTO, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponseDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			Status:    string(user.Status),
		}
	}

	return &userResponses, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*dto.UserResponseDTO, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.ErrNotFound
	}

	return &dto.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		Status:    string(user.Status),
	}, nil
}

func (s *UserService) Update(ctx context.Context, id int64, updatedUser model.User) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.ErrNotFound
	}

	if updatedUser.Email != "" {
		// Email değişiyorsa, yeni email'in başka bir kullanıcıda olmadığından emin ol
		if updatedUser.Email != user.Email {
			exists, err := s.userRepo.ExistsByEmail(ctx, updatedUser.Email)
			if err != nil {
				return errorx.ErrDatabaseOperation
			}
			if exists {
				return errorx.ErrDuplicate
			}
		}
	}

	if err = s.userRepo.Update(ctx, &updatedUser); err != nil {
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
