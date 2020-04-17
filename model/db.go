package model

import "math/big"

//
type MtMiner struct {
	Address string      `json:"address"`
	Balance *big.Int    `json:"balance"`
	Workers []*MtWorker `json:"workers"`
}

//
type MtWorker struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
	Power   int64    `json:"power"`
}
