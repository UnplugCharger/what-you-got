package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	db "hackathon/db/sqlc"
	"hackathon/utils"
	"mime/multipart"
	"net/http"
	"regexp"
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
	//defer os.Remove(imagePath)

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

	// pass this to gpt to jsonify it well
	rawReceipt, err := server.store.CreateRawReceipt(ctx, args)
	//fmt.Println(rawReceipt)

	structured, err := utils.ExtractReceiptInfo(rawReceipt.OcrText)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(structured)

	// Convert the structured map to a JSON string
	jsonBytes, err := json.Marshal(structured)
	if err != nil {
		// Handle the error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert structured data to JSON"})
		return
	}

	// If you need the JSON as a string
	//jsonString := string(jsonBytes)

	args2 := db.CreateProcessedReceiptParams{
		ReceiptID: &rawReceipt.ReceiptID,
		UserID:    &ID,
		Data:      jsonBytes,
	}

	digitised, err := server.store.CreateProcessedReceipt(ctx, args2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to process the receipt text: %s", err.Error()))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"rawReceipt":     rawReceipt,
		"extracted_text": cleanText(rawReceipt.OcrText),
		"digitised":      digitised,
	})

	// Ensure the file is deleted after processing
	//defer os.Remove(imagePath)
}

func cleanText(input string) string {
	// Replace multiple whitespaces with a single space
	re := regexp.MustCompile(`[\s]+`)
	cleaned := re.ReplaceAllString(input, " ")

	// Structure the address for better readability
	addressPattern := `(?i)(Shipping Address:|amazoncomau TAX INVOICE|Order Number:|Order Date:|Event-Driven Architecture)`
	addressRepl := "\n\n$1 "
	re = regexp.MustCompile(addressPattern)
	cleaned = re.ReplaceAllString(cleaned, addressRepl)

	// Remove unwanted characters
	cleaned = regexp.MustCompile(`\f`).ReplaceAllString(cleaned, "")

	return cleaned
}
