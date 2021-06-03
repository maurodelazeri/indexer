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

func get_logs(rw http.ResponseWriter, r *http.Request, request_payload Request) {

	a := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Topics:    [][]common.Hash{},
	},

	// field, ok := FilterQuery(request_payload.Params)

	// var filter *Filter
	// if crit.BlockHash != nil {
	// 	// Block filter requested, construct a single-shot filter
	// 	filter = NewBlockFilter(api.backend, *crit.BlockHash, crit.Addresses, crit.Topics)
	// } else {
	// 	// Convert the RPC block numbers into internal representations
	// 	begin := rpc.LatestBlockNumber.Int64()
	// 	if crit.FromBlock != nil {
	// 		begin = crit.FromBlock.Int64()
	// 	}
	// 	end := rpc.LatestBlockNumber.Int64()
	// 	if crit.ToBlock != nil {
	// 		end = crit.ToBlock.Int64()
	// 	}
	// 	// Construct the range filter
	// 	filter = NewRangeFilter(api.backend, begin, end, crit.Addresses, crit.Topics)
	// }

	rows, err := quiknode_proxy.db.Query(`select log_index,transaction_hash,transaction_index,address,data,
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
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	json.NewEncoder(rw).Encode(logs)
}
