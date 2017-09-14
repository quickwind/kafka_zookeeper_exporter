package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	zkLog "log"
)

const (
	metricsRoute = "/metrics"
)

var (
	version = "unknown"

	listenAddress = flag.String("web.listen-address", ":9381", "Address to listen on for web interface and telemetry.")
	zookeeper = flag.String("zk.addr", "zookeeper:2181", "Zookeeper Address.")
	zkChroot = flag.String("zk.chroot", "/kafka/cluster", "Zookeeper path for kafka cluster.")
	serverTimeout = flag.Duration("web.timeout", 60*time.Second, "Timeout for responding to HTTP requests.")
	zkTimeout     = flag.Duration("zk.timeout", 5*time.Second, "Timeout for ZooKeeper requests")
	showVersion   = flag.Bool("version", false, "Show version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// kazoo uses ZooKeeper client that logs everything by default, so we end up
	// with duplicated logs we don't control, disable vanilla logger messages
	// and rely on logs generated by our code
	zkLog.SetOutput(ioutil.Discard)

	log.Infoln("Listening on", *listenAddress)
	log.Infoln("Zookeeper address", *zookeeper)
	log.Infoln("Zookeeper chroot", *zkChroot)
	prometheus.DefaultRegisterer.MustRegister(newCollector(*zookeeper, *zkChroot, []string{}))

	log.Fatal(http.ListenAndServe(*listenAddress, promhttp.Handler()))
}
