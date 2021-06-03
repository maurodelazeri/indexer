# Indexer proxy

## eth_getLogs

Object - The filter options:

fromBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.

toBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
address: DATA|Array, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.

topics: Array of DATA, - (optional) Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with "or" options.

blockhash: DATA, 32 Bytes - (optional) With the addition of EIP-234 (Geth >= v1.8.13 or Parity >= v2.1.0), blockHash is a new filter option which restricts the logs returned to the single block with the 32-byte hash blockHash. Using blockHash is equivalent to fromBlock = 

toBlock = the block number with hash blockHash. If blockHash is present in the filter criteria, then neither fromBlock nor toBlock are allowed.

// Request

```
 curl -XPOST -H "content-type: application/json" --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}],"id":74}' http://localhost:5000
```



```
 curl -XPOST -H "content-type: application/json" --data '{"jsonrpc":"2.0","method":"eth_getLogs","params": [],"id":1}' http://localhost:5000
 ```