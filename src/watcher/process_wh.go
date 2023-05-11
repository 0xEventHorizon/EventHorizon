package watcher

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"github.com/k0rean-rand0m/event-horizon/src/utils"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var httpClient = &http.Client{}
var rpcClient *ethclient.Client
var once sync.Once

type requestBodyTx struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Nonce     uint64 `json:"nonce"`
	Amount    string `json:"amount"`
	GasLimit  string `json:"gasLimit"`
	GasPrice  string `json:"gasPrice"`
	GasFeeCap string `json:"gasFeeCap"`
	GasTipCap string `json:"gasTipCap"`
	Data      string `json:"data"`
	//TODO AccessList
	IsFake bool `json:"isFake"`
}

type requestBodyEventArgument struct {
	Label string      `json:"label"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type requestBodyEvent struct {
	Label     string                     `json:"label"`
	Arguments []requestBodyEventArgument `json:"arguments"`
	Data      string                     `json:"data"`
}

type requestBody struct {
	Network     string           `json:"network"`
	Emitter     string           `json:"emitter"`
	Transaction *requestBodyTx   `json:"transaction"`
	Event       requestBodyEvent `json:"event"`
}

func processWh(logEntry types.Log) {
	if config.Config.Webhook.Url == "" {
		return
	}

	var err error
	var bodyTx *requestBodyTx

	// Search for event index by topic
	event, err := utils.EventByTopic(logEntry.Topics[0])
	if err != nil {
		return
	}

	if config.Config.FullTx {
		// Move to global, process_db might need it
		once.Do(func() {
			rpcClient, err = ethclient.Dial(config.Config.RpcHttp)
			if err != nil {
				log.Fatal(err)
			}
		})

		// Get full transaction by hash
		tx, _, err := rpcClient.TransactionByHash(context.Background(), logEntry.TxHash)
		// TODO: handle error
		if err != nil {
			return
		}

		// Get transaction message
		message, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
		// TODO: handle error
		if err != nil {
			return
		}

		bodyTx = &requestBodyTx{
			To:        message.To().Hex(),
			From:      message.From().Hex(),
			Nonce:     message.Nonce(),
			Amount:    message.Value().String(),
			GasLimit:  strconv.FormatUint(message.Gas(), 10),
			GasPrice:  message.GasPrice().String(),
			GasFeeCap: message.GasFeeCap().String(),
			GasTipCap: message.GasTipCap().String(),
			Data:      "0x" + hex.EncodeToString(tx.Data()),
			//AccessList: json.Marshal(tx.AccessList()),
			IsFake: message.IsFake(),
		}
	}

	// Prepare body
	body := requestBody{
		Network:     config.Config.Network,
		Emitter:     logEntry.Address.Hex(),
		Transaction: bodyTx,
		Event: requestBodyEvent{
			Label:     event.Label,
			Arguments: []requestBodyEventArgument{},
			Data:      "0x" + hex.EncodeToString(logEntry.Data),
		},
	}

	// Fill event arguments
	var i int
	for _, argument := range event.Arguments {
		if !argument.Indexed {
			continue
		}
		i += 1
		var value interface{} = logEntry.Topics[i].Hex()
		if argument.Type == "address" {
			value = common.HexToAddress(logEntry.Topics[i].Hex()).Hex()
		}
		body.Event.Arguments = append(body.Event.Arguments,
			requestBodyEventArgument{
				Label: argument.Label,
				Type:  argument.Type,
				Value: value,
			})
	}

	data, _ := json.Marshal(body)
	go request(data, 1)
}

func request(body []byte, attempt uint64) {

	if attempt > config.Config.Webhook.MaxAttempts {
		return
	}

	log.Println("Sending data. Attempt #" + strconv.FormatUint(attempt, 10))

	// Prepare request
	bodyBuffer := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", config.Config.Webhook.Url, bodyBuffer)
	// TODO handle error
	if err != nil {
		return
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range config.Config.Webhook.Headers {
		req.Header.Set(key, value)
	}

	// Sending the request
	resp, err := httpClient.Do(req)
	if err != nil {
		go request(body, attempt+1)
		return
	}
	_ = resp.Body.Close() // TODO handle error
	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		go request(body, attempt+1)
	}
}
