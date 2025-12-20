package services

//go:generate mockgen -source=user_service.go -destination=../mocks/user_service_mock.go -package=mocks
import (
	"errors"

	"github.com/google/uuid"
	"github.com/mohod24/go-project-management/models"
	"github.com/mohod24/go-project-management/repositories"
	"github.com/mohod24/go-project-management/utils"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(id string) (*models.User, error)
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// userService is the concrete implementation of UserService.
type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// Register registers a new user.
func (s *userService) Register(user *models.User) error {
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased
	user.Role = "user"
	user.PublicID = uuid.New()

	return s.repo.Create(user)
}

// Login authenticates a user with email and password.
func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credential")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credential")
	}
	return user, nil

}

// GetByID retrieves a user by their internal ID.
func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

// GetByPublicID retrieves a user by their public ID.
func (s *userService) GetByPublicID(id string) (*models.User, error) {
	return s.repo.FindByPublicID(id)
}

// GetAllPagination retrieves users with pagination, filtering, and sorting.
func (s *userService) GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.FindAllPagination(filter, sort, limit, offset)
}

// Update updates user information.
func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

// Delete removes a user by their internal ID.
func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}