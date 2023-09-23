package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "hackathon/db/sqlc"
	"hackathon/utils"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type userResponse struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.UserID,
		FirstName: user.FirstName,
		Email:     user.Email,
	}
}

// CreateUserRequest represents the request payload needed to create a new user.
type createUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid username or password: %s", err.Error()))
		return
	}
	//password, err := utils.HashPassword(req.Password)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}

	args := db.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, err)
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

// ProcessReceiptRequest represents the request payload needed to process a receipt image.
type processReceiptRequest struct {
	ReceiptImage *multipart.FileHeader `form:"receipt_image" binding:"required"`
}

func (server *Server) processReceipt(ctx *gin.Context) {
	var req processReceiptRequest
	userIDStr := ctx.Param("userid")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid userID format")
		return
	}

	ID := int32(userID)

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid file: %s", err.Error()))
		return
	}

	// Save the uploaded file temporarily
	imagePath := "./" + req.ReceiptImage.Filename
	if err := ctx.SaveUploadedFile(req.ReceiptImage, imagePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to save the uploaded file: %s", err.Error()))
		return
	}

	// Ensure the file is deleted after processing
	defer os.Remove(imagePath)

	// Process the image using OCR
	extractedText, err := utils.ProcessReceiptImage(imagePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to process the receipt image: %s", err.Error()))
		return
	}

	args := db.CreateRawReceiptParams{
		UserID:  &ID,
		OcrText: extractedText,
	}

	rawReceipt, err := server.store.CreateRawReceipt(ctx, args)
	fmt.Println(rawReceipt)

	ctx.JSON(http.StatusOK, gin.H{
		"extracted_text": extractedText,
	})
}
