package repositories

import (
	"time"

	"github.com/mohod24/go-project-management/config"
	"github.com/mohod24/go-project-management/models"
)

// BoardRepository defines the interface for board-related database operations.
type BoardRepository interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMember(boardID uint, userIDs []uint) error
}

// boardRepository implements the BoardRepository interface.
type boardRepository struct {
}

// NewBoardRepository creates a new instance of BoardRepository.
func NewBoardRepository() BoardRepository {
	return &boardRepository{}

}

// Create saves a new board to the database.
func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}

// Update modifies an existing board in the database.
func (r *boardRepository) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":        board.Title,
		"description":  board.Description,
		"due_date":     board.DueDate,
	}).Error
}

// FindByPublicID retrieves a board by its public ID.
func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board
	err := config.DB.Where("public_id = ?", publicID).First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// AddMember adds members to a board.
func (r *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	// Implementation for adding members to a board
	if len(userIDs) == 0 {
		return nil
	}
	// Create BoardMember entries
	now := time.Now()
	var members []models.BoardMember

	// Loop through userIDs and create BoardMember structs
	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardID: int64(boardID),
			UserID:  int64(userID),
			JoinedAt: now,
		})
	}
	// Bulk insert members
	return config.DB.Create(&members).Error
}