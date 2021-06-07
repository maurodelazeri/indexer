package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/quiknode-labs/indexer/database"
	"github.com/sirupsen/logrus"
)

type rpcCall struct {
	Version string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      *json.RawMessage  `json:"id"`
}

type rpcResponse struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  interface{}      `json:"result"`
}

func formatResponse(result interface{}, call *rpcCall) *rpcResponse {
	return &rpcResponse{
		Version: "2.0",
		ID:      call.ID,
		Result:  result,
	}
}

func (q *Quiknode) getLatestBlock(ctx context.Context) (int64, error) {
	var result int64
	err := database.GetDB().DB.QueryRowContext(ctx, "SELECT max(block_number) FROM logs;").Scan(&result)
	return result, err
}

func (q *Quiknode) get_logs(ctx context.Context, rw http.ResponseWriter, r *http.Request, call *rpcCall) {
	latestBlock, err := q.getLatestBlock(ctx)
	if err != nil {
		logrus.Error("problem geting the last block")
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

	crit := filters.FilterCriteria{}
	if len(call.Params) < 1 {
		logrus.Error("missing value for required argument 0")
		cbody := json.RawMessage(`{"code":-32000,"message":"missing value for required argument 0"}`)
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
	if err := json.Unmarshal(call.Params[0], &crit); err != nil {
		logrus.Error("missing value for required argument 0")
		cbody := json.RawMessage(`{"code":-32000,"message":"format is not valid"}`)
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

	whereClause := []string{}
	indexClause := ""
	if crit.BlockHash != nil {
		whereClause = append(whereClause, fmt.Sprintf("block_hash = %s", crit.BlockHash.String()))
	}
	var fromBlock, toBlock int64
	if crit.FromBlock == nil || crit.FromBlock.Int64() < 0 {
		fromBlock = latestBlock
	} else {
		fromBlock = crit.FromBlock.Int64()
	}
	whereClause = append(whereClause, fmt.Sprintf("block_number >= %d", fromBlock))
	if crit.ToBlock == nil || crit.ToBlock.Int64() < 0 {
		toBlock = latestBlock
	} else {
		toBlock = crit.ToBlock.Int64()
	}
	whereClause = append(whereClause, fmt.Sprintf("block_number <= %d", toBlock))

	addressClause := []string{}
	for _, address := range crit.Addresses {
		addressClause = append(addressClause, fmt.Sprintf("address = %s", address.String()))
	}
	if len(addressClause) > 0 {
		whereClause = append(whereClause, fmt.Sprintf("(%v)", strings.Join(addressClause, " OR ")))
	}
	topicsClause := []string{}
	for i, topics := range crit.Topics {
		topicClause := []string{}
		for _, topic := range topics {
			topicClause = append(topicClause, fmt.Sprintf("topic%v = %s", i, topic.String()))
		}
		if len(topicClause) > 0 {
			topicsClause = append(topicsClause, fmt.Sprintf("(%v)", strings.Join(topicClause, " OR ")))
		} else {
			topicsClause = append(topicsClause, fmt.Sprintf("topic%v IS NOT NULL", i))
		}
	}
	if len(topicsClause) > 0 {
		whereClause = append(whereClause, fmt.Sprintf("(%v)", strings.Join(topicsClause, " AND ")))
	}

	query := fmt.Sprintf("SELECT address, topic0, topic1, topic2, topic3, data, block_number, transaction_hash, transaction_index, block_hash, log_index FROM logs %v WHERE %v;", indexClause, strings.Join(whereClause, " AND "))
	rows, err := database.GetDB().DB.QueryContext(ctx, query)
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

	responseBytes, err := json.Marshal(formatResponse(logs, call))
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
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(responseBytes)
}
