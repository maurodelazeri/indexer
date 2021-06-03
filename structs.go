package main

import (
	"encoding/json"
	"math/big"
)

// Log struct
type Log struct {
	LogIndex         int64  `json:"log_index,omitempty"`
	TransactionHash  string `json:"jsonrpc,omitempty"`
	TransactionIndex int64  `json:"transaction_index,omitempty"`
	Address          string `json:"address,omitempty"`
	Data             string `json:"data,omitempty"`
	Topic0           string `json:"topic0,omitempty"`
	Topic1           string `json:"topic1,omitempty"`
	Topic2           string `json:"topic2,omitempty"`
	Topic3           string `json:"topic3,omitempty"`
	BlockTimestamp   string `json:"block_timestamp,omitempty"`
	BlockNumber      int64  `json:"block_number,omitempty"`
	BlockHash        string `json:"block_hash,omitempty"`
}

// FilterQuery contains options for contract log filtering.
type FilterQuery struct {
	BlockHash string   // used by eth_getLogs, return logs only from block with this hash
	FromBlock *big.Int // beginning of the queried range, nil means genesis block
	ToBlock   *big.Int // end of the range, nil means latest block
	Addresses []string // restricts matches to events created by specific contracts
	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position AND B in second position
	// {{A}, {B}}         matches topic A in first position AND B in second position
	// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
	Topics [][]string
}

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage
}

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}
