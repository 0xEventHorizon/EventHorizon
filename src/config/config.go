package config

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"strings"
)

type parsingConfiguration struct {
	Network    bool
	Emitter    bool
	To         bool
	From       bool
	Nonce      bool
	Amount     bool
	GasLimit   bool
	GasPrice   bool
	GasFeeCap  bool
	GasTipCap  bool
	Data       bool
	AccessList bool
	IsFake     bool
}

type Event struct {
	Table     string
	Label     string
	Arguments []struct {
		Label   string
		Type    string
		Indexed bool
	}
}

type config struct {
	WsRpc     string
	HttpRpc   string
	Webhook   string
	DbUrl     string
	Network   string
	Addresses []common.Address
	Parse     parsingConfiguration
	Events    []Event
	Topics    map[common.Hash]int // a map of event signatures to event indexes
}

var Config config

func Setup() {
	configData, err := os.ReadFile("horizon.config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configData, &Config)
	if err != nil {
		log.Fatal(err)
	}

	// Building topics
	Config.Topics = make(map[common.Hash]int)
	for i, event := range Config.Events {
		arguments := make([]string, len(event.Arguments))
		for j, argument := range event.Arguments {
			arguments[j] = argument.Type
		}
		argumentsInline := strings.Join(arguments, ",")
		signature := crypto.Keccak256Hash(
			[]byte(event.Label + "(" + argumentsInline + ")"),
		)
		log.Println(event.Label + "(" + argumentsInline + ")")
		Config.Topics[signature] = i
	}
}
