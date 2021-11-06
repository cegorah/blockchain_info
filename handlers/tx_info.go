package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cegorah/blockchain_info/cache"
	"github.com/cegorah/blockchain_info/models"
	"github.com/cegorah/blockchain_info/repo"
	"github.com/cegorah/blockchain_info/restapi/operations/blockchain_info"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/runtime/middleware"
	"time"
)

type TxInfoImpl struct {
	DefaultTimeoutSecond int
	DBClient             repo.Repository
	Cache                cache.Cacher
	Logger               logger.Logger
}

func (txi *TxInfoImpl) Handle(params blockchain_info.GetTxInfoParams, i interface{}) middleware.Responder {
	lg := txi.Logger
	if lg == nil {
		lg = logger.StandardLogger{}
	}
	res := blockchain_info.NewGetTxInfoOK()
	cacheKey := fmt.Sprintf("ti_%s_%s", params.NetCode, params.Hash)
	dt, err := txi.Cache.GetCache(cacheKey)
	if err != nil {
		lg.Printf("%s", err)
	} else if dt != nil {
		var txs models.Transactions
		if err := json.Unmarshal(dt, &txs); err != nil {
			lg.Printf("%s", err)
			return blockchain_info.NewGetTxInfoInternalServerError()
		}
		res.Payload = txs
		lg.Debugf("cache hit")
		return res
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(txi.DefaultTimeoutSecond)*time.Second)
	defer cancel()
	data, err := txi.DBClient.GetTxData(ctx, params)
	if err != nil {
		lg.Printf("%s", err)
		return blockchain_info.NewGetTxInfoInternalServerError()
	}
	if data == nil {
		return blockchain_info.NewGetTxInfoNotFound()
	}
	var resPayload = make([]*models.TransactionInfo, len(data))
	for i, d := range data {
		tmp := d
		resPayload[i] = &tmp
	}
	cachedData, e := json.Marshal(&resPayload)
	if e != nil {
		lg.Printf("%s", err)
		return blockchain_info.NewGetTxInfoInternalServerError()
	}
	e = txi.Cache.SetCache(cacheKey, cachedData)
	if e != nil {
		lg.Printf("%s", err)
	}
	res.Payload = resPayload
	return res
}
