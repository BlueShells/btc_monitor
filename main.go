package main

import (
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	blockHeightGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "btc_block_height",
			Help: "Current block height of the Bitcoin node",
		},
		[]string{"network", "url"},
	)
)

func init() {
	prometheus.MustRegister(blockHeightGauge)
}

func checkBlockHeight(client *rpcclient.Client, wg *sync.WaitGroup, network, url string) {
	defer wg.Done()

	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Printf("Error fetching block count: %v", err)
		return
	}

	blockHeightGauge.WithLabelValues(network, url).Set(float64(blockCount))
}

func main() {
	// 从环境变量读取用户名和密码
	rpcUser := os.Getenv("BTC_RPC_USER")
	rpcPass := os.Getenv("BTC_RPC_PASS")
	rpcUrl := os.Getenv("BTC_RPC_URL") //test.btc.com:443
	if rpcUser == "" || rpcPass == "" || rpcUrl == "" {
		log.Fatalf("RPC_USER and RPC_PASS environment variables must be set")
	}
	// 配置比特币RPC客户端
	connCfg := &rpcclient.ConnConfig{
		Host:     rpcUrl,
		User:     rpcUser,
		Endpoint: "ws",
		Pass:     rpcPass,
		//HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		//DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating new BTC RPC client: %v", err)
	}
	defer client.Shutdown()

	var wg sync.WaitGroup
	network := "mainnet" // 你可以根据你的网络类型修改这个值
	url := connCfg.Host  // 使用连接配置中的 Host 作为 url 标签值

	go func() {
		for {
			wg.Add(1)
			checkBlockHeight(client, &wg, network, url)
			wg.Wait()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9091", nil))
}
