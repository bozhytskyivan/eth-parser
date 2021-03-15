# Ethereum notifier service

## Goal

Implement Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Problem

Users not able to receive push notifications for incoming/outgoing transactions. By Implementing `Parser` interface we would be able to hook this up to notifications service to notify about any incoming/outgoing transactions.

## Limitations

- Use Go Language
- Avoid usage of external libraries
- Use Ethereum JSONRPC to interact with Ethereum Blockchain
- Use memory storage for storing any data (should be easily extendable to support any storage in the future)

Expose public interface for external usage either via code or command line that will include supported list of operations defined in the `Parser` interface

```golang
type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}
```

### Endpoint

URL: [https://cloudflare-eth.com](https://cloudflare-eth.com/)

Request example

```shell
// Request
curl -X POST 'https://cloudflare-eth.com' --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}'

// Result
{
  "id":83,
  "jsonrpc": "2.0",
  "result": "0x4b7" // 1207
}
```

### References

- [Ethereum JSON RPC Interface](https://eth.wiki/json-rpc/API)