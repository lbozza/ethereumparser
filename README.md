
# ETH Parser API

ETH Parser API is a tool designed to parse and analyze Ethereum blockchain data, specifically targeting transactions related to specified Ethereum addresses. This API provides a convenient interface for users to subscribe to addresses of interest, retrieve transaction data, and obtain information about the latest processed blocks on the Ethereum network.

## API Endpoints

- **/subscribe?address={address}** - Subscribe to an Ethereum address to receive notifications about transactions related to that address. *(HTTP Method: POST)*
- **/transactions?address={address}** - Retrieve a list of transactions associated with the specified Ethereum address. *(HTTP Method: GET)*
- **/block** - Get information about the latest processed block on the Ethereum network. *(HTTP Method: GET)*

## Usage

 Application runs on localhost:8080

### Subscribing to an Address

To subscribe to an Ethereum address, send a POST request to the `/subscribe` endpoint with the address specified in the query parameter.

Example:

POST /subscribe?address=0x1234567890123456789012345678901234567890


### Retrieving Transactions

To retrieve transactions associated with a specific Ethereum address, send a GET request to the `/transactions` endpoint with the address specified in the query parameter.

Example:

GET /transactions?address=0x1234567890123456789012345678901234567890


### Getting the Latest Processed Block

To obtain information about the latest processed block on the Ethereum network, send a GET request to the `/block` endpoint.

Example:

GET /block

## Business Requirements


![Screenshot 2024-03-27 at 11 05 05](https://github.com/lbozza/ethereumparser/assets/21343976/d323307a-0105-4332-b9c1-ab72088025ae)
