# Event Horizon
One-click solution to index EVM networks events

### Installation
1. Install Golang
2. Install dependencies
```
go mod tidy
```
3. Install Docker (for production deployment)

### Configure
All configuration is made in `horizon.config.json` file.

In a current version only 1 watcher per service is supported.

Here's the refference for the config file:
```
{
  "RpcWs": "****",              // Websocket RPC endpoint
  "RpcHttp": "****",            // HTTP RPC endpoint
  "webhook": "****",            // Webhook endpoint (will be triggered by each captured event)
  "dbConnectionString": "****", // PostgreSQL conn string (Work in proggress)
  "network": "polygon",         // Network label
  "addresses": [                // Addresses to watch (might be empty to watch events from all addresses)
    "0xc2132D05D31c914a87C6611C10748AEb04B58e8F" // Polygon USDT as an example
  ],
  // Additional TX parsing options
  "parse": {
    "network": true,        // Include network label in the event
    "emitter": true,        // Include emitter contract address in the event
    "to": false,            // Include TX "to" address in the event
    "from": false,          // Include TX "from" address in the event
    "nonce":  false,        // Include TX nonce in the event
    "amount": false,        // Include TX amount in the event
    "gasLimit":  false,     // Include TX gas limit in the event
    "gasPrice": false,      // Include TX gas price in the event
    "gasFeeCap":  false,    // Include TX gas fee cap in the event
    "gasTipCap": false,     // Include TX gas tip cap in the event
    "data":  false,         // Include TX data in the event
    "accessList": false,    // Include TX access list in the event
    "isFake":  false        // Include TX fake flag in the event
  },
  // Events to listen (example for ERC20 Transfer event)
  "events": [
    {
      "label": "Transfer", // Event label
      // Event arguments (can be obtained from the contract ABI)
      "arguments": [
        {"label": "from", "type": "address"},
        {"label": "to", "type": "address"},
        {"label": "value", "type": "uint256"}
      ]
    }
  ]
}
```

### Local run
`make run`

### Deploying with Docker
Work in progress