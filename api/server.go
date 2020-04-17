package api

import (
	"bytes"
	"coin-distribution/model"
	"coin-distribution/storage"
	"coin-distribution/storage/leveldb"
	"context"
	"encoding/gob"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/jsonrpc"
	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"log"
	"math/big"
	"net/http"
)

//
type CoinDisServer struct {
	server *http.Server
	db     storage.Storage

	fullNode api.FullNode
	closer   jsonrpc.ClientCloser

	miners map[string]*model.MtMiner
}

//
func (cds *CoinDisServer) Run() {
	if err := cds.server.ListenAndServe(); err != nil {
		log.Printf("MtsServer.Run:%s", err)
	}
}

//
func NewServer(addr string, port int) (*CoinDisServer, error) {
	server := &CoinDisServer{}
	router := gin.Default()

	addrPort := fmt.Sprintf(addr+":%d", port)
	httpServer := &http.Server{
		Addr:    addrPort,
		Handler: router,
	}

	ldb, err := leveldb.NewLevelDbStorage("./data")
	if err == nil {
		return nil, err
	}

	db := storage.NewStorageWithCache(ldb, 256)

	server.server = httpServer
	server.db = db
	server.miners = make(map[string]*model.MtMiner)

	return server, nil
}

//
func (cds *CoinDisServer) AutoDistributeCoins() {
	tipSet, err := cds.fullNode.ChainHead(context.TODO())
	if err != nil {
		log.Println("ERROR", err.Error())
		return
	}

	if tipSet == nil {
		log.Println("ERROR", "tip set nil")
		return
	}

	h := tipSet.Height()
	if h%4000 == 0 {
		for minerAddr := range cds.miners {
			var record model.MtMiner
			mb, err := cds.db.Get(minerAddr)
			if err != nil {
				continue
			}
			var buff bytes.Buffer
			decoder := gob.NewDecoder(bytes.NewBuffer([]byte(mb)))
			encoder := gob.NewEncoder(&buff)

			_ = decoder.Decode(&record)
			amount := big.NewInt(10000000000)
			cds.sendCoin(minerAddr, "t3rj7bmiqcwe25l2eghaw4mwiilx4k2vjlkqzac4syasfs4jp4qrv4wvu5q5qmesmujkzm7eodl4zhzjhj5cqa", amount)
			record.Address = minerAddr
			record.Balance = new(big.Int).Sub(record.Balance, amount)
			_ = encoder.Encode(&record)
			_ = cds.db.Set(minerAddr, buff.String())
		}
	}
}

//
func (cds *CoinDisServer) calcProfit() {
	var total int64
	mp := make(map[string]int64)

	for _, miner := range cds.miners {
		for _, worker := range miner.Workers {
			total += worker.Power
			mp[worker.Address] = worker.Power
		}
	}

	// TODO: get the total profit
	var totalProfit int64 = 10000000

	for addr, power := range mp {
		profit := power * 1000000000 / total * totalProfit

		var buff bytes.Buffer

		var worker model.MtWorker

		data, _ := cds.db.Get(addr)
		dec := gob.NewDecoder(bytes.NewBuffer([]byte(data)))
		_ = dec.Decode(&worker)
		worker.Balance = new(big.Int).Add(worker.Balance, big.NewInt(profit))

		enc := gob.NewEncoder(&buff)
		_ = enc.Encode(&worker)

		cds.db.Set(addr, buff.String())
	}

}

//
func (cds *CoinDisServer) sendCoin(sender, receiver string, amount *big.Int) {
	for minerAddr := range cds.miners {
		from, _ := address.NewFromString(minerAddr)
		to, _ := address.NewFromString("t3rj7bmiqcwe25l2eghaw4mwiilx4k2vjlkqzac4syasfs4jp4qrv4wvu5q5qmesmujkzm7eodl4zhzjhj5cqa")
		msg := &types.Message{
			From:  from,
			To:    to,
			Value: types.BigInt{big.NewInt(1000000000000000)},
		}
		cds.fullNode.MpoolPushMessage(context.TODO(), msg)
	}
}

func (cds *CoinDisServer) initRoute(r gin.IRouter) {
	r.POST("/api/v0/distributeCoin", cds.DistributeCoins)
}
