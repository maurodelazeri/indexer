package main

import (
	"encoding/json"
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
