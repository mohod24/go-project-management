package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/mohod24/go-project-management/models"
	"github.com/mohod24/go-project-management/services"
	"github.com/mohod24/go-project-management/utils"
)

// BoardController handles HTTP requests related to boards.
type BoardController struct {
	service services.BoardService
}

// NewBoardController creates a new instance of BoardController.
func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

// CreateBoard handles the creation of a new board.
func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	Claims := user.Claims.(jwt.MapClaims)
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing permintaan", err.Error())
	}

	userID, err = uuid.Parse(Claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing permintaan", err.Error())
	}
	board.OwnerPublicID = userID

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal menyimpan data", err.Error())
	}
	return utils.Success(ctx, "Berhasil membuat board", board)
}

// UpdateBoard handles the updating of an existing board.
func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	// Get the board public ID from the URL parameters
	publicID := ctx.Params("id")
	board := new(models.Board)
	// Parse the request body into the board struct
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing permintaan", err.Error())
	}
	// Validate the public ID
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Public ID tidak valid", err.Error())
	}
	// Retrieve the existing board to ensure it exists
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	}
	// Set the public ID and owner public ID to ensure they are not changed
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt
	board.OwnerID = existingBoard.OwnerID
	// Proceed to update the board
	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal update board", err.Error())
	}
	return utils.Success(ctx, "Berhasil update board", board)
}

func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	
	var userIDs []string
	// Parse the request body to get user IDs
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing permintaan", err.Error())
	}
	// Add members to the board
	if err := c.service.AddMember(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan anggota", err.Error())
	}
	return utils.Success(ctx, "Berhasil menambahkan anggota", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIDs []string
	// Parse the request body to get user IDs
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing permintaan", err.Error())
	}
	// Remove members from the board
	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus anggota", err.Error())
	}
	return utils.Success(ctx, "Berhasil menghapus anggota", nil)
}