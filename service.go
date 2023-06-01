package promotions

import (
	"context"
	"log"
	"math"
	"strings"
	"time"
)

const filename = "promotions.csv"

type Service struct {
	logger  *log.Logger
	storage DB
}

func NewService(l *log.Logger) Service {
	return Service{
		logger:  l,
		storage: NewDBConn(),
	}
}

func (s Service) getOne(ctx context.Context, id int) (Result, error) {
	r, err := s.storage.getById(ctx, id)
	if err != nil {
		return Result{}, err
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700 MST", r.ExpirationDate)
	if err != nil {
		return Result{}, err
	}

	fnRoundUp := func(val float64, precision int) float64 {
		return math.Ceil(val*(math.Pow10(precision))) / math.Pow10(precision)
	}
	return Result{
		ID:             strings.ToUpper(r.Key.String()),
		Price:          fnRoundUp(r.Price, 2),
		ExpirationDate: date.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s Service) updateData(ctx context.Context) error {
	if err := s.storage.clearDB(ctx); err != nil {
		return err
	}

	file, err := s.openFileFromBucket(filename)
	if err != nil {
		s.logger.Fatal(err)
		return err
	}
	return s.storage.batchInsert(ctx, file)
}
