package handlers

import (
	"context"
	"github.com/cegorah/blockchain_info/internal"
	"github.com/cegorah/blockchain_info/models"
	"github.com/cegorah/blockchain_info/restapi/operations/blockchain_info"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"os"
	"testing"
	"time"
)

const TTL = 2

type apiModeler interface {
	Validate(formats strfmt.Registry) error
}

type BasicTest struct {
	TestId int
}

type BITest struct {
	BasicTest
	TTL      int
	Params   blockchain_info.GetBlockParams
	Expected middleware.Responder
}

type TxTest struct {
	BasicTest
	TTL      int
	Params   blockchain_info.GetTxInfoParams
	Expected middleware.Responder
}

var BiTests []BITest
var TxTests []TxTest
var BI BlockInfoImpl
var TI TxInfoImpl
var CacheMock cacheMock
var RepoMock repoMock
var TestMap map[string][]apiModeler

type repoMock struct {
}

func (rm *repoMock) Close() error {
	panic("implement me")
}

func (rm *repoMock) GetBlockData(ctx context.Context, params blockchain_info.GetBlockParams) ([]models.BlockInfo, error) {
	var res []models.BlockInfo
	for _, v := range TestMap[params.Hash] {
		res = append(res, *v.(*models.BlockInfo))
	}
	return res, nil
}
func (rm *repoMock) GetTxData(ctx context.Context, params blockchain_info.GetTxInfoParams) ([]models.TransactionInfo, error) {
	var res []models.TransactionInfo
	for _, v := range TestMap[params.Hash] {
		res = append(res, *v.(*models.TransactionInfo))
	}
	return res, nil
}

type cacheMock struct {
	ttl    int
	ts     time.Time
	cache  map[string][]byte
	cached bool
}

func (cm *cacheMock) GetCache(key string) ([]byte, error) {
	tn := time.Now()
	if tn.After(cm.ts) {
		delete(cm.cache, key)
		cm.cached = false
		return nil, nil
	}
	cm.cached = true
	return cm.cache[key], nil
}

func (cm *cacheMock) SetCache(key string, value []byte) error {
	cm.cache[key] = value
	cm.ts = time.Now().Add(time.Second * time.Duration(cm.ttl))
	return nil
}

func TestMain(m *testing.M) {
	setup()
	teardown()
	ret := m.Run()
	os.Exit(ret)
}

func setup() {
	RepoMock = repoMock{}
	CacheMock = cacheMock{ttl: TTL, cache: map[string][]byte{}}
	TestMap = map[string][]apiModeler{}
	BiTestsSetup()
	TxTestsSetup()
}

func BiTestsSetup() {
	BiTests = append(BiTests, BITest{
		TTL: TTL,
		Params: blockchain_info.GetBlockParams{
			Hash:    "BLOCK_AA",
			NetCode: "BTC",
		},
		Expected: &blockchain_info.GetBlockOK{Payload: &models.BlockInfo{
			NetCode:  "BTC",
			NextHash: "next_AA",
			PrevHash: "prev_AA",
		}},
	})
	TestMap["BLOCK_AA"] = []apiModeler{&models.BlockInfo{
		NetCode:  "BTC",
		NextHash: "next_AA",
		PrevHash: "prev_AA",
	}}
	BiTests = append(BiTests, BITest{
		TTL: TTL,
		Params: blockchain_info.GetBlockParams{
			Hash:    "BLOCK_BB",
			NetCode: "LTC",
		},
		Expected: &blockchain_info.GetBlockOK{Payload: &models.BlockInfo{
			NetCode:  "LTC",
			NextHash: "next_BB",
			PrevHash: "prev_BB",
		}},
	})
	TestMap["BLOCK_BB"] = []apiModeler{&models.BlockInfo{
		NetCode:  "LTC",
		NextHash: "next_BB",
		PrevHash: "prev_BB",
	}}
	BiTests = append(BiTests, BITest{
		TTL: TTL,
		Params: blockchain_info.GetBlockParams{
			Hash:    "BLOCK_CC",
			NetCode: "LTC",
		},
		Expected: &blockchain_info.GetBlockNotFound{},
	})
}

func TxTestsSetup() {
	TxTests = append(TxTests, TxTest{
		TTL: TTL,
		Params: blockchain_info.GetTxInfoParams{
			Hash:    "TX_AA",
			NetCode: "BTC",
		},
		Expected: &blockchain_info.GetTxInfoOK{Payload: models.Transactions{
			{ID: 1, Fee: 2, SentValue: 3},
		}},
	})
	TestMap["TX_AA"] = []apiModeler{&models.TransactionInfo{
		ID:        1,
		Fee:       2,
		SentValue: 3,
	}}
	TxTests = append(TxTests, TxTest{
		TTL: TTL,
		Params: blockchain_info.GetTxInfoParams{
			Hash:    "BLOCK_CC",
			NetCode: "LTC",
		},
		Expected: &blockchain_info.GetTxInfoNotFound{},
	})
}

func teardown() {
}

func TestBlockInfoImpl_Handle(t *testing.T) {
	BI.DBClient = &RepoMock
	BI.Cache = &CacheMock
	for _, tst := range BiTests {
		res := BI.Handle(tst.Params, nil)
		internal.TestEqual(t, tst.Expected, res)
		time.Sleep(time.Duration(tst.TTL-2) * time.Second)
		res = BI.Handle(tst.Params, nil)
		internal.TestEqual(t, tst.Expected, res)
		internal.TestAssert(t, CacheMock.cached, "")
		time.Sleep(time.Duration(tst.TTL+2) * time.Second)
		res = BI.Handle(tst.Params, nil)
		internal.TestEqual(t, tst.Expected, res)
		internal.TestAssert(t, !CacheMock.cached, "")
	}
}

func TestTxInfoImpl_Handle(t *testing.T) {
	TI.DBClient = &RepoMock
	TI.Cache = &CacheMock
	for _, tst := range TxTests {
		res := TI.Handle(tst.Params, nil)
		internal.TestEqual(t, res, tst.Expected)
		time.Sleep(time.Duration(tst.TTL-2) * time.Second)
		res = TI.Handle(tst.Params, nil)
		internal.TestEqual(t, res, tst.Expected)
		internal.TestAssert(t, CacheMock.cached, "")
		time.Sleep(time.Duration(tst.TTL+2) * time.Second)
		res = TI.Handle(tst.Params, nil)
		internal.TestEqual(t, res, tst.Expected)
		internal.TestAssert(t, !CacheMock.cached, "")
	}
}
