package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"

	"github.com/SchoIsles/kruise-state-metrics/pkg/metricshandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"k8s.io/kube-state-metrics/pkg/options"
	"k8s.io/kube-state-metrics/pkg/util/proc"
	"k8s.io/kube-state-metrics/pkg/version"
	"k8s.io/kube-state-metrics/pkg/whiteblacklist"

	"context"

	"github.com/SchoIsles/kruise-state-metrics/internal/store"
	clientset "github.com/SchoIsles/kruise-state-metrics/pkg/client/clientset/versioned"
)

const (
	metricsPath = "/metrics"
	healthzPath = "/healthz"
)

// promLogger implements promhttp.Logger
type promLogger struct{}

func (pl promLogger) Println(v ...interface{}) {
	klog.Error(v...)
}

func main() {
	opts := options.NewOptions()
	opts.AddFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := opts.Parse()
	if err != nil {
		klog.Fatalf("Error: %s", err)
	}

	if opts.Version {
		fmt.Printf("%#v\n", version.GetVersion())
		os.Exit(0)
	}

	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}

	storeBuilder := store.NewBuilder()
	ksmMetricsRegistry := prometheus.NewRegistry()
	storeBuilder.WithMetrics(ksmMetricsRegistry)

	var collectors = []string{"clonesets"}
	if err := storeBuilder.WithEnabledResources(collectors); err != nil {
		klog.Fatalf("Failed to set up collectors: %v", err)
	}

	if len(opts.Namespaces) == 0 {
		klog.Info("Using all namespace")
		storeBuilder.WithNamespaces(options.DefaultNamespaces)
	} else {
		if opts.Namespaces.IsAllNamespaces() {
			klog.Info("Using all namespace")
		} else {
			klog.Infof("Using %s namespaces", opts.Namespaces)
		}
		storeBuilder.WithNamespaces(opts.Namespaces)
	}

	whiteBlackList, err := whiteblacklist.New(opts.MetricWhitelist, opts.MetricBlacklist)
	if err != nil {
		klog.Fatal(err)
	}

	err = whiteBlackList.Parse()
	if err != nil {
		klog.Fatalf("error initializing the whiteblack list : %v", err)
	}

	storeBuilder.WithWhiteBlackList(whiteBlackList)

	proc.StartReaper()

	config, err := clientcmd.BuildConfigFromFlags(opts.Apiserver, opts.Kubeconfig)
	if err != nil {
		panic(err)
	}
	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	storeBuilder.WithKubeClient(kubeClient)
	ksmMetricsRegistry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)
	go telemetryServer(ksmMetricsRegistry, opts.TelemetryHost, opts.TelemetryPort)

	serveMetrics(ctx, kubeClient, storeBuilder, opts, opts.Host, opts.Port, opts.EnableGZIPEncoding)

}

func telemetryServer(registry prometheus.Gatherer, host string, port int) {
	// Address to listen on for web interface and telemetry
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))

	klog.Infof("Starting kube-state-metrics self metrics server: %s", listenAddress)

	mux := http.NewServeMux()

	// Add metricsPath
	mux.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: promLogger{}}))
	// Add index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Kube-State-Metrics Metrics Server</title></head>
             <body>
             <h1>Kube-State-Metrics Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
			 </ul>
             </body>
             </html>`))
	})
	klog.Fatal(http.ListenAndServe(listenAddress, mux))
}

func serveMetrics(ctx context.Context, kubeClient clientset.Interface, storeBuilder *store.Builder, opts *options.Options, host string, port int, enableGZIPEncoding bool) {
	// Address to listen on for web interface and telemetry
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))

	klog.Infof("Starting metrics server: %s", listenAddress)

	mux := http.NewServeMux()

	// TODO: This doesn't belong into serveMetrics
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	m := metricshandler.New(
		opts,
		kubeClient,
		storeBuilder,
		enableGZIPEncoding,
	)
	go m.Run(ctx)
	mux.Handle(metricsPath, m)

	// Add healthzPath
	mux.HandleFunc(healthzPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
	// Add index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Kube Metrics Server</title></head>
             <body>
             <h1>Kube Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
             <li><a href='` + healthzPath + `'>healthz</a></li>
			 </ul>
             </body>
             </html>`))
	})
	klog.Fatal(http.ListenAndServe(listenAddress, mux))
}
