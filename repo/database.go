package repo

import (
	"context"
	"github.com/cegorah/blockchain_info/models"
	"github.com/cegorah/blockchain_info/restapi/operations/blockchain_info"
)

type Repository interface {
	GetBlockData(ctx context.Context, params blockchain_info.GetBlockParams) ([]models.BlockInfo, error)
	GetTxData(ctx context.Context, params blockchain_info.GetTxInfoParams) ([]models.TransactionInfo, error)
	Close() error
}
