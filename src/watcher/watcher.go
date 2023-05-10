package watcher

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"github.com/k0rean-rand0m/event-horizon/src/utils"
	"log"
)

//var httpClient *ethclient.Client

func Run() (err error) {
	// Dealing Node
	//httpClient, err := ethclient.Dial(config.Config.HttpRpc)
	wsClient, err := ethclient.Dial(config.Config.WsRpc)
	if err != nil {
		return err
	}

	// Preparing filter query
	topics := make([]common.Hash, len(config.Config.Events))
	for hash, i := range config.Config.Topics {
		topics[i] = hash
	}
	query := ethereum.FilterQuery{
		Addresses: config.Config.Addresses,
		Topics:    [][]common.Hash{topics},
	}
	log.Println(query.Topics)

	// Creating logs channel
	logs := make(chan types.Log)

	// Subscribing to events
	sub, err := wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	// Running watcher goroutine
	go watch(sub, logs)

	// Incrementing WaitGroup counter
	utils.Wg.Add(1)

	return nil
}

func watch(sub ethereum.Subscription, logs <-chan types.Log) {
	// Listening to events
	for {
		select {
		case err := <-sub.Err():
			// Websocket reconnection
			for err != nil {
				log.Println("Websocket reconnecting caused by:", err)
				err = Run()
			}
			// If success release WaitGroup
			utils.Wg.Done()
			return
		case logEntry := <-logs:
			go process(logEntry)
		}
	}
}
