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

type BlockInfoImpl struct {
	DefaultTimeoutSecond int
	DBClient             repo.Repository
	Cache                cache.Cacher
	Logger               logger.Logger
}

func (bInfo *BlockInfoImpl) Handle(params blockchain_info.GetBlockParams, i interface{}) middleware.Responder {
	lg := bInfo.Logger
	if lg == nil {
		lg = logger.StandardLogger{}
	}
	res := blockchain_info.NewGetBlockOK()
	cacheKey := fmt.Sprintf("bi_%s_%s", params.NetCode, params.Hash)
	dt, err := bInfo.Cache.GetCache(cacheKey)
	if err != nil {
		lg.Printf("%s", err)
	} else if dt != nil {
		blInfo := models.BlockInfo{}
		if err := json.Unmarshal(dt, &blInfo); err != nil {
			lg.Printf("%s", err)
			return blockchain_info.NewGetBlockInternalServerError()
		}
		res.Payload = &blInfo
		lg.Debugf("cache hit")
		return res
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(bInfo.DefaultTimeoutSecond)*time.Second)
	defer cancel()
	data, err := bInfo.DBClient.GetBlockData(ctx, params)
	if err != nil {
		lg.Printf("%s", err)
		return blockchain_info.NewGetBlockInternalServerError()
	}
	if data == nil {
		return blockchain_info.NewGetBlockNotFound()
	}
	bd := &data[0]
	cachedData, e := json.Marshal(&bd)
	if e != nil {
		lg.Printf("%s", err)
		return blockchain_info.NewGetBlockInternalServerError()
	}
	e = bInfo.Cache.SetCache(cacheKey, cachedData)
	if e != nil {
		lg.Printf("%s", err)
	}
	res.Payload = bd
	return res
}
