package models

type TransferRequest struct {
	Id     string  `json:"id"`
	Amount float64 `json:"amount"`
	UserId string  `json:"userId"`
}

type TransferUserIdRequestId struct {
	Id     string
	UserId string
}
