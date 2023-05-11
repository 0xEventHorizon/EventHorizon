package watcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"github.com/k0rean-rand0m/event-horizon/src/db"
	"github.com/k0rean-rand0m/event-horizon/src/utils"
	"log"
	"strconv"
	"strings"
)

func processDb(logEntry types.Log) {

	if !db.Instance.InUse() {
		return
	}

	// Search for event index by topic
	event, err := utils.EventByTopic(logEntry.Topics[0])
	if err != nil {
		return
	}

	// Process event
	labels := []string{"hash", "network"}
	placeholders := []string{"$1", "$2"}
	values := []interface{}{logEntry.TxHash.Hex(), config.Config.Network}

	for i, argument := range event.Arguments {
		if !argument.Indexed {
			continue
		}
		labels = append(labels, "\""+argument.Label+"\"")
		placeholders = append(placeholders, "$"+strconv.Itoa(i+3))
		// TODO: handle all types
		if argument.Type == "address" {
			values = append(values, common.HexToAddress(logEntry.Topics[i+1].Hex()).Hex())
			continue
		}
		values = append(values, logEntry.Topics[i+1])
	}
	query := "insert into event_horizon." + event.Table + " (" + strings.Join(labels, ", ") + ") values (" + strings.Join(placeholders, ", ") + ")"
	log.Println(query)
	_, err = db.Instance.Execute(query, values...)
	if err != nil {
		//TODO: handle error
		// some events might not have same network-hash (several events in one TX)
		log.Println(err)
	}
}
