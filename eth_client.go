package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ethClient struct {
	httpClient *http.Client

	url string

	jsonRPC string
}

type ClientOption func(c *ethClient)

func NewEthClient(opts ...ClientOption) *ethClient {
	// Default client configuration
	ec := &ethClient{
		httpClient: http.DefaultClient,
		url:        "https://cloudflare-eth.com",
		jsonRPC:    "2.0",
	}

	for _, optFn := range opts {
		optFn(ec)
	}

	return ec
}

const (
	EthMethodBlockNumber   = "eth_blockNumber"
	EthMethodBlockByNumber = "eth_getBlockByNumber"
)

type ethRequest struct {
	ID      int64         `json:"id"`
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type ethResponse struct {
	ID      int64           `json:"id"`
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

type Block struct {
	Hash         string        `json:"hash"`
	Number       HexNumber     `json:"number"`
	ParentHash   string        `json:"parentHash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BlockNumber HexNumber `json:"blockNumber"`
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"`
}

func (ec *ethClient) GetCurrentBlock(id int64) (int64, error) {
	resp, err := ec.sendEthRequest(id, EthMethodBlockNumber)
	if err != nil {
		return 0, err
	}

	var blockNumber HexNumber
	err = json.Unmarshal(resp.Result, &blockNumber)
	if err != nil {
		return 0, fmt.Errorf("unmarshal current block response value: %v", err)
	}

	return int64(blockNumber), nil
}

func (ec *ethClient) GetBlockByNumber(id int64, number int64) (Block, error) {
	resp, err := ec.sendEthRequest(id, EthMethodBlockByNumber, HexNumber(number), true)
	if err != nil {
		return Block{}, err
	}

	var result Block

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return Block{}, fmt.Errorf("unmarshal get block by number resposne: %v", err)
	}

	return result, nil
}

func (ec *ethClient) sendEthRequest(id int64, method string, params ...interface{}) (ethResponse, error) {
	reqBodyBytes, err := json.Marshal(ethRequest{
		ID:      id,
		JsonRPC: ec.jsonRPC,
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return ethResponse{}, fmt.Errorf("marshal eth request body for method %s: %v", method, err)
	}

	req, err := http.NewRequest(http.MethodPost, ec.url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return ethResponse{}, fmt.Errorf("create eth request for method %s: %v", method, err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return ethResponse{}, fmt.Errorf("do eht request for method %s: %v", method, err)
	}

	defer func() { _ = resp.Body.Close() }()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ethResponse{}, fmt.Errorf("read the response body for method %s: %d", method, err)
	}

	var result ethResponse
	err = json.Unmarshal(respBodyBytes, &result)
	if err != nil {
		return ethResponse{}, fmt.Errorf("unmarshal eth response body for method %s: %v", method, err)
	}

	return result, nil
}

type HexNumber int64

func (i HexNumber) MarshalJSON() ([]byte, error) {
	hexNumber := fmt.Sprintf("\"0x%x\"", i)

	return []byte(hexNumber), nil
}

func (i *HexNumber) UnmarshalJSON(data []byte) error {
	var resultStr string
	err := json.Unmarshal(data, &resultStr)
	if err != nil {
		return err
	}

	hex := strings.Replace(resultStr, "0x", "", 1)

	result, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return err
	}

	*i = HexNumber(result)

	return nil
}
