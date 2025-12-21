package repositories

import (
	"strings"

	"github.com/mohod24/go-project-management/config"
	"github.com/mohod24/go-project-management/models"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindAllPagination(filter, sort string, limit, ofset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct{}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Create adds a new user to the database.
func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// FindByEmail retrieves a user by their email.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByID retrieves a user by their internal ID.
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

// FindByPublicID retrieves a user by their public ID.
func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

// FindAllPagination retrieves users with pagination, filtering, and sorting.
func (r *userRepository) FindAllPagination(filter, sort string, limit, ofset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := config.DB.Model(&models.User{})

	//filtering
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("name Ilike ? OR email Ilike ?", filterPattern, filterPattern)
	}
	//count total data
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//sorting
	if sort != "" {
		//Misalnya sort=name (ASC ascending ) sort =-name (DESC descending)
		if sort == "-id" {
			sort = "-internal_id"
		} else if sort == "id" {
			sort = "internal_id"
		}

		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		} else {
			sort += " ASC"
		}

		db = db.Order(sort)
	}

	err := db.Limit(limit).Offset(ofset).Find(&users).Error
	return users, total, err

}

// Update modifies an existing user's information.
func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).
		Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
	}).Error
}

// Delete removes a user from the database by their internal ID.
func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
