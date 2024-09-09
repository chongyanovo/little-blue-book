package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(server *gin.Engine)
}
