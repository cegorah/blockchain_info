package main

import (
	"encoding/json"
	"fmt"
	"github.com/cegorah/blockchain_info/internal"
	"github.com/cegorah/blockchain_info/models"
	"github.com/cegorah/blockchain_info/restapi"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	Debug         = true
	RedisPassword = "admin"
	PsqlDsn       = "postgres://admin:admin_password@bc_db:5432/block_chain"
)

type apiModeler interface {
	Validate(formats strfmt.Registry) error
}

var Tests = []struct {
	Url          string
	ExpectedCode int
	ExpectedBody apiModeler
}{
	{"/v1/tx/DG/a", 422, nil},
	{"/v1/tx/BTC/aa", 200, &models.Transactions{
		{Fee: 1.4, ID: 1, SentValue: 2},
		{Fee: 2.4, ID: 2, SentValue: 3},
		{Fee: 3.5, ID: 3, SentValue: 4},
		{Fee: 2.4, ID: 4, SentValue: 3},
		{Fee: 3.5, ID: 5, SentValue: 4},
		{Fee: 2.4, ID: 6, SentValue: 3},
	}},
	{"/v1/block/DG/a", 422, nil},
	{"/v1/block/BTC/aa", 200, &models.BlockInfo{
		FirstTenTransactions: models.Transactions{
			{Fee: 1.4, ID: 1, SentValue: 2},
			{Fee: 2.4, ID: 2, SentValue: 3},
			{Fee: 3.5, ID: 3, SentValue: 4},
			{Fee: 2.4, ID: 4, SentValue: 3},
			{Fee: 3.5, ID: 5, SentValue: 4},
			{Fee: 2.4, ID: 6, SentValue: 3},
		},
		ID:        1,
		NetCode:   "BTC",
		NextHash:  "aa",
		PrevHash:  "aa",
		Size:      5,
		Timestamp: strfmt.Date{},
	}},
}

func TestMain(m *testing.M) {
	os.Setenv("PSQL_DSN", PsqlDsn)
	os.Setenv("DEBUG", fmt.Sprintf("%v", Debug))
	os.Setenv("REDIS_PASSWORD", RedisPassword)
	ret := m.Run()
	os.Exit(ret)
}

func TestMe(t *testing.T) {
	dr, e := filepath.Abs("../../config/config.json")
	if e != nil {
		panic(e)
	}
	config(dr)
	handler, err := restapi.GetApiHandler()
	if err != nil {
		t.Fatal("get api handler", err)
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()
	cl := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/v1/block/BTC/aa", nil)
	res, _ := cl.Do(req)
	internal.TestEqual(t, 401, res.StatusCode)
}

func TestAll(t *testing.T) {
	dr, e := filepath.Abs("../../config/config.json")
	internal.TestOk(t, e)
	config(dr)
	handler, err := restapi.GetApiHandler()
	if err != nil {
		t.Fatal("get api handler", err)
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()
	cl := &http.Client{}
	for _, tst := range Tests {
		req, _ := http.NewRequest("GET", ts.URL+tst.Url, nil)
		req.Header.Set("Authorization", "Bearer 1")
		res, _ := cl.Do(req)
		internal.TestEqual(t, tst.ExpectedCode, res.StatusCode)
		if tst.ExpectedBody != nil {
			dt, e := ioutil.ReadAll(res.Body)
			internal.TestOk(t, e)
			switch v := tst.ExpectedBody.(type) {
			case *models.BlockInfo:
				var bi models.BlockInfo
				internal.TestOk(t, json.Unmarshal(dt, &bi))
				assert.ElementsMatch(t, v.FirstTenTransactions, bi.FirstTenTransactions)
				internal.TestEqual(t, v, &bi)
			case *models.Transactions:
				var txs models.Transactions
				internal.TestOk(t, json.Unmarshal(dt, &txs))
				assert.ElementsMatch(t, txs, *v)
			}
		}
	}
}