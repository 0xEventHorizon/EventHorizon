package utils

// Sharable helper functions and entities

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"sync"
)

// Wg - Global WaitGroup
var Wg sync.WaitGroup

// EventByTopic - Returns config.Event by given topic
func EventByTopic(topic common.Hash) (config.Event, error) {
	eventIndex, exists := config.Config.Topics[topic]
	if exists {
		return config.Config.Events[eventIndex], nil
	}

	return config.Event{}, errors.New("event not found for provided topic")
}
