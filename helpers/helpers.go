package helpers

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(200, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(400, APIResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	})
}
