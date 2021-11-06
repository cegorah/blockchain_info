package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cegorah/blockchain_info/internal"
	"github.com/cegorah/blockchain_info/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
	"time"
)

type simpleTest struct {
	ttl   int32
	key   string
	value []byte
}

var trans []*models.TransactionInfo
var testData models.BlockInfo
var redisServer RedisServer
var rClient *redis.Client

var tests []simpleTest

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	tearDown()
	os.Exit(ret)
}

func setup() {
	for i := 0; i < 10; i++ {
		trans = append(trans, &models.TransactionInfo{
			Timestamp: strfmt.Date(time.Now()),
			Fee:       0,
			SentValue: 0,
			ID:        0,
		})
	}
	var testData models.BlockInfo
	v, e := json.Marshal(&testData)
	if e != nil {
		panic(e)
	}
	tests = append(tests, simpleTest{
		5, "BTC_000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf", v},
	)
	rs := map[string]interface{}{
		"addr":     "localhost:6379",
		"username": "",
		"password": "admin",
		"db":       0,
	}
	rClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v", rs["addr"]),
		Password: fmt.Sprintf("%v", rs["password"]),
		DB:       rs["db"].(int),
	})
	redisServer = NewRedisServer(rs)
}

func tearDown() {
}

func TestRedisServer_GetCache(t *testing.T) {
	for _, tst := range tests {
		r := rClient.Set(context.Background(), tst.key, tst.value, time.Duration(tst.ttl)*time.Second)
		internal.TestOk(t, r.Err())
		d, e := redisServer.GetCache(tst.key)
		internal.TestOk(t, e)
		internal.TestEqual(t, d, tst.value)
	}
}

func TestMemcacheServer_SetCache(t *testing.T) {
	for _, tst := range tests {
		e := redisServer.SetCache(tst.key, tst.value)
		internal.TestOk(t, e)
		it, e := rClient.Get(context.Background(), tst.key).Bytes()
		internal.TestOk(t, e)
		internal.TestEqual(t, tst.value, it)
		time.Sleep(time.Duration(tst.ttl+1) * time.Second)
		r, e := rClient.Get(context.Background(), tst.key).Result()
		internal.TestAssert(t, r == "" && e == redis.Nil, "")
	}
}
