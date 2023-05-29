package controller

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"transfer-service/internal/models"
	"transfer-service/internal/router/response"
)

type KafkaService interface {
	SendTransfer(transfer *models.KafkaTransfer) error
}

type TransferController struct {
	Service KafkaService

	transferRequests TransferRequests
}

func NewTransferController(service KafkaService) TransferController {
	return TransferController{
		Service: service,
		transferRequests: TransferRequests{
			requestIds: make(map[models.TransferUserIdRequestId]struct{}),
			mtx:        &sync.Mutex{},
		},
	}
}

type TransferRequests struct {
	requestIds map[models.TransferUserIdRequestId]struct{}

	mtx *sync.Mutex
}

func (t TransferRequests) Set(id, userId string) {
	t.mtx.Lock()
	t.requestIds[models.TransferUserIdRequestId{
		UserId: userId,
		Id:     id,
	}] = struct{}{}
	t.mtx.Unlock()
}

func (t TransferRequests) Check(id, userId string) bool {
	t.mtx.Lock()
	if _, ok := t.requestIds[models.TransferUserIdRequestId{
		UserId: userId,
		Id:     id,
	}]; ok {
		return true
	}
	t.mtx.Unlock()
	return false
}

func (c TransferController) SendTransfer(ctx *gin.Context) {
	m := &models.TransferRequest{}
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Println("error bind json: ", err)
		response.ReturnError(ctx, http.StatusBadRequest, err)
		return
	}

	if c.transferRequests.Check(m.Id, m.UserId) {
		log.Printf("error duplicate transfer id: %s and user id: %s\n, %v", m.Id, m.UserId, err)
		response.ReturnError(ctx, http.StatusConflict, err)
		return
	}

	err = c.Service.SendTransfer(&models.KafkaTransfer{
		Id:     m.Id,
		Amount: m.Amount,
	})
	if err != nil {
		log.Printf("error send transfer to kafka id: %s and user id: %s\n, %v", m.Id, m.UserId, err)
		response.ReturnError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
