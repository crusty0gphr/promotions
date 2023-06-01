package promotions

import (
	"github.com/google/uuid"
)

type (
	Result struct {
		ID             string  `json:"id"`
		Price          float64 `json:"price"`
		ExpirationDate string  `json:"expiration_date"`
	}

	Row struct {
		ID             int
		Key            uuid.UUID
		Price          float64
		ExpirationDate string
	}

	Promotions []Row
)
