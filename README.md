ETH balance checker with Prometheus metrics exposed.

## Makefile
- `make build` builds binary into /build/bin (linux/amd64, Docker is used on Windows);
- `make dockerize` builds a Docker image (see /build/linux_amd64.dockerfile);

## Usage

`--addr` flag is repeatable and has next format: `[name:]hex`

#### With Docker
```sh
docker run -d \
  --name ethlevel \
  --restart always \
  -p 52112:2112 \
  ethlevel:latest \
  /app/ethlevel --geth "https://mainnet.infura.io/v3/yoursecret" \
  --addr=my_address:0x0000000000000000000000000000000000000000
```

#### Without Docker
```sh
./ethlevel --geth "https://mainnet.infura.io/v3/yoursecret" \
  --addr=my_address:0x0000000000000000000000000000000000000000
```

## Args
`go run main.go -help`
```
-addr value
    Address to observe
-geth string
    GETH endpoint (default "http://localhost:8545")
-ns string
    Prometheus metrics namespace
-period uint
    Check period in seconds (default 30)
-port uint
    Port to serve metrics (default 2112)
-ss string
    Prometheus metrics subsystem (default "ethlevel")
```

