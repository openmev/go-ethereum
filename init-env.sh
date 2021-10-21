#!/bin/sh
# Configure RPC.
FLAGS="$FLAGS --http --http.addr=0.0.0.0 --http.port=8545 --http.api=admin,debug,eth,miner,net,personal,txpool,web3"
FLAGS="$FLAGS --ws --ws.addr=0.0.0.0 --ws.origins \"*\" --ws.api=admin,debug,eth,miner,net,personal,txpool,web3"
if [ "$OPENMEV_GRAPHQL_ENABLED" != "" ]; then
	FLAGS="$FLAGS --graphql"
fi
# used for the graphql to allow submission of unprotected tx
if [ "$OPENMEV_ALLOW_UNPROTECTED_TX" != "" ]; then
 	FLAGS="$FLAGS --rpc.allow-unprotected-txs"
fi

# Run the go-ethereum implementation with the requested flags.
FLAGS="$FLAGS --nat=none"
echo "Running go-ethereum with flags $FLAGS"
$geth "$FLAGS"