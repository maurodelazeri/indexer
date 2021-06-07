package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/quiknode-labs/indexer/database"
	"github.com/sirupsen/logrus"
)

type QuiknodeExporter struct {
	RPC string
}

type TaskRanges struct {
	start uint64
	end   uint64
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	Blocknumber      string   `json:"blockNumber"`
	Transactionhash  string   `json:"transactionHash"`
	Transactionindex string   `json:"transactionIndex"`
	Blockhash        string   `json:"blockHash"`
	Logindex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type RPCRerturn struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  []struct {
		Address          string   `json:"address"`
		Topics           []string `json:"topics"`
		Data             string   `json:"data"`
		Blocknumber      string   `json:"blockNumber"`
		Transactionhash  string   `json:"transactionHash"`
		Transactionindex string   `json:"transactionIndex"`
		Blockhash        string   `json:"blockHash"`
		Logindex         string   `json:"logIndex"`
		Removed          bool     `json:"removed"`
	} `json:"result"`
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *QuiknodeExporter) httpPost(url string, data []byte) (string, error) {
	timeout := time.Second * 300
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logrus.Error("Problem to create request", err.Error())
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Problem making request: ", err.Error())
		logrus.Infoln(string(data))
		time.Sleep(10 * time.Second)
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Problem to read response", err.Error())
		return "", err
	}
	return string(respBody), nil
}

// EncodeUint64 encodes i as a hex string with 0x prefix.
func (e *QuiknodeExporter) EncodeUint64(i uint64) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, i, 16))
}

func (e *QuiknodeExporter) hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	output, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	return strconv.FormatInt(output, 10)
}

func (e *QuiknodeExporter) exportLogsWorker(ranges TaskRanges) {
	ss := e.EncodeUint64(ranges.start)
	ee := e.EncodeUint64(ranges.end)
	data := []byte("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getLogs\",\"params\":[{\"fromBlock\": \"" + ss + "\", \"toBlock\":\"" + ee + "\"}],\"id\":1}")

	respBody, err := e.httpPost(e.RPC, data)
	if err != nil {
		logrus.Error("Problem with http post: ", err.Error())
	}
	rpc_return := RPCRerturn{}
	err = json.Unmarshal([]byte(respBody), &rpc_return)
	if err != nil {
		logrus.Error("Problem with the parse rpc return: ", err.Error(), "\n", "data: ", string(respBody))
	}
	if len(rpc_return.Error.Message) > 0 {
		logrus.Error("Problem with the rpc return: ", rpc_return.Error.Message)
	}
	if len(rpc_return.Result) == 0 {
		return
	}
	logs := []Log{}
	result, _ := json.Marshal(rpc_return.Result)
	err = json.Unmarshal(result, &logs)
	if err != nil {
		logrus.Error("Problem with the parse results return: ", err.Error())
	}
	query := "INSERT INTO logs (log_index,transaction_hash,transaction_index,address,data,topic0,topic1,topic2,topic3,block_number,block_hash) VALUES "
	for _, log := range logs {
		topics := []string{"", "", "", ""}
		for i, topic := range log.Topics {
			topics[i] = topic
		}
		query += "(" + e.hexaNumberToInteger(log.Logindex) + ",'" + log.Transactionhash + "'," + e.hexaNumberToInteger(log.Transactionindex) + ",'" + log.Address + "','" + log.Data + "','" + topics[0] + "','" + topics[1] + "','" + topics[2] + "','" + topics[3] + "'," + e.hexaNumberToInteger(log.Blocknumber) + ",'" + log.Blockhash + "'),"
	}

	query = strings.TrimRight(query, ",")
	query += ";"

	_, err = database.GetDB().DB.Exec(query)
	if err != nil {
		logrus.Error("Problem to insert query ", err.Error())
	}
}

func (e *QuiknodeExporter) Exporterlogs() {
	var beginBlock uint64
	fmt.Print("Enter the initial block number ")
	fmt.Scanln(&beginBlock)

	var endBlock uint64
	fmt.Print("Enter the final block number ")
	fmt.Scanln(&endBlock)

	var blocksPerRequest uint64
	fmt.Print("Enter Desired the number of blocks per request ")
	fmt.Scanln(&blocksPerRequest)

	var i uint64
	ranges := []TaskRanges{}
	wp := workerpool.New(60)
	for i = beginBlock; i < endBlock; {
		task_range := TaskRanges{i, i + blocksPerRequest}
		ranges = append(ranges, task_range)
		i = i + blocksPerRequest + 1
	}
	for _, r := range ranges {
		r := r
		wp.Submit(func() {
			e.exportLogsWorker(r)
		})
	}
	wp.StopWait()
}
