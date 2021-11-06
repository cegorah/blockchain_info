package repo

import (
	"context"
	"github.com/cegorah/blockchain_info/models"
	"github.com/cegorah/blockchain_info/restapi/operations/blockchain_info"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func NewPsClient(ctx context.Context, conString string) (pc PostgresClient, err error) {
	pc.pool, err = pgxpool.Connect(ctx, conString)
	return
}

func (pc *PostgresClient) GetBlockData(
	ctx context.Context, params blockchain_info.GetBlockParams) ([]models.BlockInfo, error) {
	var res []models.BlockInfo
	var ftt []*models.TransactionInfo
	rows, err := pc.pool.Query(
		ctx,
		"SELECT bi.net_code, bi.id, bi.timestamp, prevb.hash, nextb.hash, bi.size, "+
			"COALESCE(ti.id,0) AS tx_id, COALESCE(ti.fee, 0) AS ti_fee, COALESCE(ti.sent_value, 0) AS ti_sv, ti.timestamp "+
			"FROM block_info bi "+
			"FULL JOIN block2tx b2t ON b2t.block_id = bi.id "+
			"FULL JOIN tx_info ti ON b2t.tx_id = ti.id "+
			"FULL JOIN block_info prevb ON bi.prev_id = prevb.id "+
			"FULL JOIN block_info nextb ON bi.prev_id = nextb.id "+
			"WHERE bi.hash = $1 and bi.net_code = $2 LIMIT 10;", params.Hash, params.NetCode)
	if err != nil {
		return res, err
	}
	var bi models.BlockInfo
	for rows.Next() {
		var tt models.TransactionInfo
		e := rows.Scan(&bi.NetCode, &bi.ID, &bi.Timestamp, &bi.PrevHash, &bi.NextHash, &bi.Size,
			&tt.ID, &tt.Fee, &tt.SentValue, &tt.Timestamp)
		if e != nil {
			return res, e
		}
		if tt.ID != 0 {
			ftt = append(ftt, &tt)
		}
	}
	if bi.ID == 0{
		return nil, nil
	}
	bi.FirstTenTransactions = ftt
	res = append(res, bi)
	return res, nil
}

func (pc *PostgresClient) GetTxData(
	ctx context.Context, params blockchain_info.GetTxInfoParams) ([]models.TransactionInfo, error) {
	var res []models.TransactionInfo
	rows, err := pc.pool.Query(ctx,
		"SELECT tx.* FROM tx_info tx INNER JOIN block2tx b2t "+
		"ON b2t.tx_id = tx.id INNER JOIN block_info bi ON bi.id = b2t.block_id "+
		"WHERE bi.hash = $1 AND bi.net_code = $2;", params.Hash, params.NetCode)
	if err != nil && err != pgx.ErrNoRows {
		return res, err
	}
	for rows.Next() {
		var ti models.TransactionInfo
		if e := rows.Scan(&ti.ID, &ti.Fee, &ti.SentValue, &ti.Timestamp); e != nil && e != pgx.ErrNoRows {
			return res, e
		}
		if ti.ID != 0 {
			res = append(res, ti)
		}
	}
	return res, nil
}

func (pc *PostgresClient) Close() error {
	pc.pool.Close()
	return nil
}
