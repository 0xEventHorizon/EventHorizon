package config

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"strings"
)

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
	RpcWs   string
	RpcHttp string
	Webhook struct {
		Url         string
		Headers     map[string]string
		MaxAttempts uint64
	}
	Database struct {
		ConnectionString string
		Schema           string
	}
	Network   string
	Addresses []common.Address
	FullTx    bool
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
