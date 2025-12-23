package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mohod24/go-project-management/models"
	"github.com/mohod24/go-project-management/repositories"
)

// BoardService defines the interface for board-related business logic.
type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	GetByPublicID(publicID string) (*models.Board, error)
	AddMember(boardPublicID string, userPublicIDs []string) error
}

// boardService implements the BoardService interface.
type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

// NewBoardService creates a new instance of BoardService.
func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
}

// Create creates a new board.
func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

// Update updates an existing board.
func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

// GetByPublicID retrieves a board by its public ID.
func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

// AddMember adds members to a board.
func (s *boardService) AddMember(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}
	// Cek keanggotaaan sebelum ditambahkan
	existingMembers, err := s.boardMemberRepo.GetMembers(string(boardPublicID))
	if err != nil {
		return errors.New("failed to check existing members")
	}

	// cek cepat pakai map
	memberMap := make(map[uint]bool)
	// isi memberMap dengan existingMembers
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true //memberMap[1] = true
	}
	// filter userInternalIDs yang sudah menjadi member
	// misal userInternalIDs = [1,2,3,4], memberMap[1]=true, memberMap[3]=true
	// maka newMemberIDs = [2,4]
	var newMemberIDs []uint
	for _, userID := range userInternalIDs {
		// jika userID tidak ada di memberMap, berarti bukan member
		if !memberMap[userID] {
			newMemberIDs = append(newMemberIDs, userID)
		}
	}
	if len(newMemberIDs) == 0 {
		return nil // tidak ada member baru untuk ditambahkan
	}
	
	// tambahkan member baru
	return s.boardRepo.AddMember(uint(board.InternalID), newMemberIDs)
}