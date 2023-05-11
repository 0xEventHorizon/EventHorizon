# Event Horizon
One-click solution to index EVM networks events

### Installation
1. Install Golang
2. Install dependencies
```
go mod tidy
```

### Configure
All configuration is made in `horizon.config.json` file.

In a current version only 1 watcher per service is supported.
That means that one microservice is used to index one network and one set of addresses.

Here's the refference for the config file:

```json
{
  // Websocket RPC endpoint
  "RpcWs": "****",
  // RPC HTTP endpoint
  "RpcHttp": "****",
  // Network label
  "network": "polygon",
  
  // HTTP RPC endpoint
  "webhook": {
    "url": "http://localhost:8081/eventhorizon",
    // If specified, events will be sent to this webhook
    "headers": {
      // Optional headers for the webhook
      "X-Auth-EventHorizon": "1234567890"
    },
    // Max attempts to send event to the webhook
    "maxAttempts": 3
  },

  // PostgreSQL configuration (WIP)
  "database": {  
    // PostgreSQL connection string
    "connectionString": "",
    // PostgreSQL schema to use for indexer (will be created automatically)
    "schema": "event_horizon"
  },

  // If true, the full transaction details will be indexed
  "fullTx": false,

  // Addresses to watch (Whitelist of emitters)
  "addresses": [
    "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
  ],
  
  // Events to listen (example for ERC20 Transfer event)
  "events": [
    {
      // Event label
      "label": "Transfer",
      // Event arguments (can be obtained from the contract ABI)
      "arguments": [
        {
          // Argument label
          "label": "from",
          // Argument solidity type
          "type": "address",
          // Is event argument indexed
          "indexed": true
        },
        {
          "label": "to",
          "type": "address",
          "indexed": true
        },
        {
          "label": "value",
          "type": "uint256",
          "indexed": false
        }
      ]
    }
  ]
}
```

### Local run
`make run`

### Deploying with Docker
Work in progress