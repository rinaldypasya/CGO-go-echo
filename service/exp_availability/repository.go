package exp_availability

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpId(ctx context.Context,expId string)([]*models.ExpAvailability ,error)
}