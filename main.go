package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alandiegosantos/http-random-stress/internal"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
)

var (
	// GitRepository is the repository of the project
	GitRepository string
	// GitCommit is the commit of this build
	GitCommit string
)

func main() {

	url := flag.String("url", "", "URL")
	rate := flag.Int("rate", 1000, "Average requests per second")
	timeout := flag.Int("timeout", 0, "Run test for this amount of seconds (0 - run forever)")
	nWorkers := flag.Int("workers", runtime.NumCPU(), "Number of workers")
	version := flag.Bool("version", false, "Show Version")
	debug := flag.Bool("debug", false, "Set debug")
	enablePrometheus := flag.Bool("enablePrometheus", false, "Enable Prometheus HTTP server to monitor the test")
	prometheusAddr := flag.String("prometheusAddress", "0.0.0.0:9090", "Prometheus HTTP Server Address")

	flag.Parse()

	if len(*url) <= 0 && strings.Contains(*url, "http") {
		logrus.Fatal("Please provide an URL")
	}

	if *version {
		fmt.Printf("Version:\t\t%s\n", "0.1")
		fmt.Printf("Go Version:\t\t%v\n", runtime.Version())
		fmt.Printf("Repository:\t\t%v\n", GitRepository)
		fmt.Printf("Commit:\t\t\t%v\n", GitCommit)
		os.Exit(0)
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {

		select {
		case <-ctx.Done():

		case <-s:

			logrus.Info("Signal Received. Exiting")

			cancel()
		}

		wg.Done()

	}()

	if *timeout > 0 {
		ticker := time.NewTicker(time.Duration(*timeout) * time.Second)

		wg.Add(1)
		go func() {

			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return

				case <-ticker.C:
					cancel()
				}
			}

		}()
	}

	//Starting prometheus

	if *enablePrometheus {

		logrus.Infof("Starting prometheus at %s", *prometheusAddr)

		wg.Add(1)
		go func() {

			server := http.Server{
				Addr:    *prometheusAddr,
				Handler: promhttp.Handler(),
			}

			go func() {

				<-ctx.Done()

				server.Shutdown(ctx)

			}()

			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				// Error starting or closing listener:
				logrus.Errorf("HTTP server ListenAndServe: %v", err)
			}

			wg.Done()

		}()

	}

	rpsPerWorker := float64(*rate / *nWorkers)

	if *rate <= *nWorkers {
		nWorkers = rate
		rpsPerWorker = 1
	}

	for i := 0; i < *nWorkers; i++ {

		wg.Add(1)
		go func(ID int) {

			logrus.Debugf("Starting worker %d", ID)

			internal.StartWorker(ctx, rpsPerWorker, *url)

			wg.Done()

		}(i)

	}

	wg.Wait()

	fmt.Printf("Report\n")

	for i := 0; i < 6; i++ {

		fmt.Printf("%dxx responses: %v\n", i, internal.GetHttpReqCounter(i, *url))

	}

}
