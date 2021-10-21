


# Build go-ethereum on the fly and delete all build tools afterwards
# RUN \
#   apk add --update bash curl jq go git make gcc musl-dev              \
#         ca-certificates linux-headers                           && \
#   git clone --depth 1 --branch $branch  https://github.com/ethereum/go-ethereum && \
#   (cd go-ethereum && make geth)                               && \
#   (cd go-ethereum                                             && \
#   echo "{}"                                                      \
#   | jq ".+ {\"repo\":\"$(git config --get remote.origin.url)\"}" \
#   | jq ".+ {\"branch\":\"$(git rev-parse --abbrev-ref HEAD)\"}"  \
#   | jq ".+ {\"commit\":\"$(git rev-parse HEAD)\"}"               \
#   > /version.json)                                            && \
#   cp go-ethereum/build/bin/geth /geth                         && \
#   apk del go git make gcc musl-dev linux-headers              && \
#   rm -rf /go-ethereum && rm -rf /var/cache/apk/*
RUN cd /go-ethereum && GO111MODULE=on go run build/ci.go install ./cmd/geth