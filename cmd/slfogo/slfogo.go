package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	edl "github.com/forklifts-for-great-justice/goforkliftit/etcdlib"
	"github.com/forklifts-for-great-justice/slfogo/pkg/slfogolib"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

const (
	envPORT      = "PORT"
	envRABBIT    = "RABBIT"
	envRUSER     = "RUSER"
	envRPASS     = "RPASS"
	envRPORT     = "RPORT"
	envREXCHANGE = "REXCHANGE"
	envETCD      = "ETCD"
	envMPORT     = "MPORT"
	envFQDN      = "FQDN"
	metricsURL   = "/metrics"
)

func handleEtcd(ctx context.Context, servers []string, myIp string, svrInfo edl.ServiceInfo, g *errgroup.Group) {
	ea, err := edl.NewEtcdAgent(
		ctx,
		servers,
	)
	if err != nil {
		log.Fatal(err)
	}

	info, err := json.Marshal(svrInfo)
	if err != nil {
		log.Fatal(err)
	}
	ea.Put(
		ctx,
		fmt.Sprintf("/service/slfogo/%s", myIp),
		string(info),
	)
	g.Go(func() error {
		return ea.KeepAlive(context.Background())
	})
}

func handleUptime(upMet prometheus.Counter) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		upMet.Inc()
	}
}

func run(ctx context.Context, getEnv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	slog.InfoContext(ctx, "slgogo online")
	svrPort, err := strconv.Atoi(getEnv(envPORT))
	if err != nil {
		return err
	}

	mtxPort, err := strconv.Atoi(getEnv(envMPORT))
	if err != nil {
		return err
	}

	upGauge := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "ForkliftsForGreatJustice",
			Subsystem: "slfogo",
			Name:      "uptime",
			Help:      "uptime for service",
		})

	etcServer := getEnv(envETCD)
	fqdn := getEnv(envFQDN)
	g := new(errgroup.Group)
	svrInfo := edl.ServiceInfo{
		ServicePort: svrPort,
		MetricsPort: mtxPort,
	}
	handleEtcd(ctx, []string{etcServer}, fqdn, svrInfo, g)

	connectStr := slfogolib.BuildConnectString(getEnv(envPORT))
	slog.InfoContext(ctx, "Listening on:", "connectStr", connectStr)

	server, lpChan := slfogolib.BuildServer()
	server.ListenTCP(connectStr)
	server.Boot()

	// op := slfogolib.NewWriterOutputProcessor(os.Stdout)
	rabbitServer := getEnv(envRABBIT)
	rabbitPort, err := strconv.Atoi(getEnv(envRPORT))
	if err != nil {
		return err
	}
	rabbitExchange := getEnv(envREXCHANGE)
	rabbitUser := getEnv(envRUSER)
	rabbitPass := getEnv(envRPASS)
	op := slfogolib.NewRabbitMQOutputProcessor(
		rabbitUser,
		rabbitPass,
		rabbitServer,
		rabbitPort,
		rabbitExchange,
	)
	if err := op.Connect(); err != nil {
		return err
	}

	mth := slfogolib.NewMetricHolder()
	mh := slfogolib.NewMessageHandler(op, lpChan, mth)

	go mth.HandleMetrics(ctx)
	go mh.HandleMessages(ctx)

	reg := prometheus.NewRegistry()
	reg.MustRegister(upGauge)
	reg.MustRegister(mth.GetGauge())

	go handleUptime(upGauge)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry: reg}))
	go http.ListenAndServe(":8889", nil)

	<-ctx.Done()
	cancel()
	server.Kill()
	server.Wait()

	close(lpChan)
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func main() {
	godotenv.Load()
	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
}
