package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
)

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

// https://github.com/ethereum/go-ethereum/blob/d8ff53dfb8a516f47db37dbc7fd7ad18a1e8a125/ethclient/ethclient.go#L287
func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

type LogsFilterQuery struct {
	BlockHash *string  // used by eth_getLogs, return logs only from block with this hash
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

func (q *QuiknodeIndexer) get_logs(rw http.ResponseWriter, r *http.Request, request_payload Request) {
	filter := LogsFilterQuery{}
	result, _ := json.Marshal(request_payload.Params)
	err := json.Unmarshal(result, &filter)
	if err != nil {
		logrus.Error("problem to unmarshal filter.", err.Error())
		return
	}

	rows, err := GetDB().DB.Query(`select log_index,transaction_hash,transaction_index,address,data,
	topic0,topic1,topic2,topic3,block_timestamp,block_number,block_hash from logs limit 1`)
	if err != nil {
		logrus.Error("problem with select.", err.Error())
		cbody := json.RawMessage(`{"code":-32000,"message":"please try again"}`)
		var cresp = Response{
			Jsonrpc: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		io.WriteString(rw, string(cresb))
		return
	}

	var logs []Log
	for rows.Next() {
		var LogIndex int64
		var TransactionHash sql.NullString
		var TransactionIndex int64
		var Address sql.NullString
		var Data sql.NullString
		var Topic0 sql.NullString
		var Topic1 sql.NullString
		var Topic2 sql.NullString
		var Topic3 sql.NullString
		var BlockTimestamp sql.NullString
		var BlockNumber int64
		var BlockHash sql.NullString
		err := rows.Scan(&LogIndex, &TransactionHash, &TransactionIndex, &Address,
			&Data, &Topic0, &Topic1, &Topic2, &Topic3, &BlockTimestamp,
			&BlockNumber, &BlockHash)
		if err != nil {
			logrus.Error("problem with select Scan.", err.Error())
			cbody := json.RawMessage(`{"code":-32000,"message":"please try again"}`)
			var cresp = Response{
				Jsonrpc: "2.0",
				Error:   cbody,
			}
			cresb, _ := json.Marshal(cresp)
			rw.Header().Set("Content-Type", "application/json")
			rw.Header().Set("X-Content-Type-Options", "nosniff")
			io.WriteString(rw, string(cresb))
			return
		}
		var log Log
		log.LogIndex = LogIndex
		log.TransactionHash = TransactionHash.String
		log.TransactionIndex = TransactionIndex
		log.Address = Address.String
		log.Data = Data.String
		log.Topic0 = Topic0.String
		log.Topic1 = Topic1.String
		log.Topic2 = Topic2.String
		log.Topic3 = Topic3.String
		log.BlockTimestamp = BlockTimestamp.String
		log.BlockNumber = BlockNumber
		log.BlockHash = BlockHash.String
		logs = append(logs, log)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(rw).Encode(logs)
}
