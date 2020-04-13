package harbors

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetAllWithJoinCPC(ctx context.Context, page *int, size *int, search string) ([]*models.HarborsWCPCDto, error)
}
