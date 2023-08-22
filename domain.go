package promotions

import (
	"github.com/google/uuid"
)

type (
	Result struct {
		ID             string  `json:"id"`
		ExpirationDate string  `json:"expiration_date"`
		Price          float64 `json:"price"`
	}

	Row struct {
		ExpirationDate string
		ID             int
		Price          float64
		Key            uuid.UUID
	}

	Promotions []Row
)
