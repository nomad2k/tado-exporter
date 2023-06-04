package main

import (
	"context"
	"errors"
	slackbot2 "github.com/clambin/go-common/slackbot"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	promserver "github.com/clambin/go-common/taskmanager/prometheus"
	"github.com/clambin/tado"
	"github.com/clambin/tado-exporter/collector"
	"github.com/clambin/tado-exporter/controller"
	"github.com/clambin/tado-exporter/controller/slackbot"
	"github.com/clambin/tado-exporter/controller/zonemanager/rules"
	"github.com/clambin/tado-exporter/health"
	"github.com/clambin/tado-exporter/poller"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var (
	configFilename string
	cmd            = cobra.Command{
		Use:     "tado-monitor",
		Short:   "exporter / controller for Tadoº thermostats",
		Run:     Main,
		Version: version,
	}
)

// overridden during build
var version = "change-me"

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func Main(cmd *cobra.Command, _ []string) {
	var opts slog.HandlerOptions
	if viper.GetBool("debug") {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &opts)))

	slog.Info("tado-monitor starting", "version", cmd.Version)

	// Do we have zone rules?
	r, err := GetZoneRules()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to read zone rules", "err", err)
		os.Exit(1)
	}

	mgr := taskmanager.New(makeTasks(r)...)

	// context to terminate the created go routines
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err = mgr.Run(ctx); err != nil {
		slog.Error("failed to start tado-monitor", "err", err)
	}

	slog.Info("tado-monitor stopped")
}

func makeTasks(rules []rules.ZoneConfig) []taskmanager.Task {
	var tasks []taskmanager.Task

	api := tado.New(
		viper.GetString("tado.username"),
		viper.GetString("tado.password"),
		viper.GetString("tado.clientSecret"),
	)

	// Poller
	p := poller.New(api, viper.GetDuration("poller.interval"))
	tasks = append(tasks, p)

	// Collector
	coll := &collector.Collector{Poller: p}
	prometheus.MustRegister(coll)
	tasks = append(tasks, coll)

	// Prometheus Server
	tasks = append(tasks, promserver.New(promserver.WithAddr(viper.GetString("exporter.addr"))))

	// Health Endpoint
	h := health.New(p)
	tasks = append(tasks, h)
	r := http.NewServeMux()
	r.Handle("/health", http.HandlerFunc(h.Handle))
	tasks = append(tasks, httpserver.New(viper.GetString("health.addr"), r))

	// Slackbot
	var bot slackbot.SlackBot
	if viper.GetBool("controller.tadoBot.enabled") {
		bot = slackbot2.New(viper.GetString("controller.tadoBot.token"), slackbot2.WithName("tado "+version))
		tasks = append(tasks, bot)
	}

	// Controller
	if len(rules) > 0 {
		tasks = append(tasks, controller.New(api, rules, bot, p))
	} else {
		slog.Warn("no rules found. controller will not run")
	}

	return tasks
}

func GetZoneRules() ([]rules.ZoneConfig, error) {
	f, err := os.Open(filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "rules.yaml"))
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var config struct {
		Zones []rules.ZoneConfig `yaml:"zones"`
	}

	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	for _, zone := range config.Zones {
		var kinds []string
		for _, rule := range zone.Rules {
			kinds = append(kinds, rule.Kind.String())
		}
		slog.Info("zone rules found", "zone", zone.Zone, "rules", strings.Join(kinds, ","))
	}
	return config.Zones, nil
}

func init() {
	cobra.OnInitialize(initConfig)
	cmd.Flags().StringVar(&configFilename, "config", "", "Configuration file")
	cmd.Flags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.AddConfigPath("/etc/tado-monitor/")
		viper.AddConfigPath("$HOME/.tado-monitor")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetDefault("debug", false)
	viper.SetDefault("tado.username", "")
	viper.SetDefault("tado.password", "")
	viper.SetDefault("tado.clientSecret", "")
	viper.SetDefault("exporter.addr", ":9090")
	viper.SetDefault("poller.interval", 30*time.Second)
	viper.SetDefault("health.addr", ":8080")
	viper.SetDefault("controller.tadobot.enabled", true)
	viper.SetDefault("controller.tadobot.token", "")

	viper.SetEnvPrefix("TADO_MONITOR")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
		os.Exit(1)
	}
}
