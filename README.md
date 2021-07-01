## Setup

1. build `go-ethereum`
```
make geth
```
2. run archive node
```
./build/bin/geth --syncmode full --gcmode archive
```
Note: This step requires 10T+ ssd space.


## Modified files

The modified go-ethereum is forked from commit `9c653ff6625df1ff0f6c6612958b2a1da0021e40` of upstream go-ethereum.

The main logic of mutator is located at the directory `core/teller/`. The mutation rules were written in `core/teller/mutateRule.go`. The remaining code mostly remain the same as the upstream version. We create a new api `debug_mutateTraceTransaction` which the analyzer communicate with. Please see the deatils of the api implementation at `eth/tracers/api.go`.  
