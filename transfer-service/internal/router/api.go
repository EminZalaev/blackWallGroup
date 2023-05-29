package router

import (
	"github.com/gin-gonic/gin"

	"transfer-service/internal/controller"
)

type OrdTransferController interface {
	PostTransferV1(ctx *gin.Context)
}

type TransferRouter struct {
	controller controller.TransferController
}

func NewTransferRouter(controller controller.TransferController) *TransferRouter {
	return &TransferRouter{
		controller: controller,
	}
}

// InitApiV1Group api/v1 group
func (t TransferRouter) InitApiV1Group(group *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	group.Use(handlers...)

	cash := group.Group("/cash")
	cash.POST("/transfer", t.controller.SendTransfer)
}
