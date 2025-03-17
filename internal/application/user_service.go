package application

import (
	"auth0_test/internal/domain"
	"auth0_test/internal/infrastructure/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email string) (*domain.User, error) {
	user := domain.NewUser(name, email)
	err := s.repo.Create(user)
	return user, err
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(uid)
}

func (s *UserService) UpdateUser(id string, name, email string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	user, err := s.repo.GetByID(uid)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	user.Name = name
	user.Email = email
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(uid)
}
