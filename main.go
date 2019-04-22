package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ethMult = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func main() {
	var (
		argGeth      = flag.String("geth", "http://localhost:8545", "GETH endpoint")
		argAddress   = addressList{}
		argPort      = flag.Uint("port", 2112, "Port to serve metrics")
		argNamespace = flag.String("ns", "", "Prometheus metrics namespace")
		argSubsys    = flag.String("ss", "ethlevel", "Prometheus metrics subsystem")
		argPeriod    = flag.Uint("period", 30, "Check period in seconds")
	)
	flag.Var(&argAddress, "addr", "Address to observe")
	flag.Parse()

	if len(argAddress.list) == 0 {
		panic("add an address with 'addr' flag")
	}

	geth, err := ethclient.Dial(*argGeth)
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	defer geth.Close()

	mtxETH := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "ethlevel",
		Help:      "ETH amount",
		Namespace: *argNamespace, Subsystem: *argSubsys,
	}, []string{"name", "address"})

	// ---

	wg := sync.WaitGroup{}
	stopped := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%v", *argPort),
			Handler: http.DefaultServeMux,
		}
		fmt.Println("Listen on port", *argPort)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic("failed to listen: " + err.Error())
			}
		}()
		<-stopped
		server.Shutdown(nil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		period := time.Second * time.Duration(*argPeriod)
		if period < time.Second {
			period = time.Second
		}
		for {
			for _, v := range argAddress.list {
				b, err := geth.BalanceAt(context.Background(), common.HexToAddress(v.address), nil)
				if err != nil {
					fmt.Println("Failed to get ETH balance for", v.address)
				} else {
					f, _ := new(big.Float).Quo(
						big.NewFloat(0).SetInt(b),
						big.NewFloat(0).SetInt(ethMult),
					).Float64()
					mtxETH.
						WithLabelValues(v.name, v.address).
						Set(f)
				}
			}
			select {
			case _, ok := <-stopped:
				if !ok {
					return
				}
			case <-time.After(period):
			}
		}
	}()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		close(stopped)
	}()

	wg.Wait()
	fmt.Println("Stopped")
}

// ---

// ---

type address struct {
	name    string
	address string
}

type addressList struct {
	list []address
}

// Set impl.
func (aa *addressList) Set(v string) error {
	p := strings.Split(v, ":")
	switch len(p) {
	case 1:
		return aa.add("", p[0])
	case 2:
		return aa.add(p[0], p[1])
	default:
		return errors.New("valid format: '[name:]hex'")
	}
}

// String impl.
func (aa *addressList) String() string {
	return fmt.Sprintf("%#v", aa.list)
}

func (aa *addressList) add(name, hex string) error {
	if !common.IsHexAddress(hex) {
		return errors.New("invalid address")
	}
	if name == "" {
		name = hex
	}
	aa.list = append(
		aa.list,
		address{
			name:    name,
			address: hex,
		},
	)
	return nil
}
