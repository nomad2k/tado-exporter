package stack

import (
	"context"
	"fmt"
	"github.com/clambin/httpserver"
	"github.com/clambin/tado"
	"github.com/clambin/tado-exporter/collector"
	"github.com/clambin/tado-exporter/configuration"
	"github.com/clambin/tado-exporter/controller"
	"github.com/clambin/tado-exporter/health"
	"github.com/clambin/tado-exporter/pkg/slackbot"
	"github.com/clambin/tado-exporter/poller"
	"github.com/clambin/tado-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"time"
)

// Stack groups all components, so they can be easily started/stopped
type Stack struct {
	Poller     poller.Poller
	Health     *health.Health
	Collector  *collector.Collector
	TadoBot    slackbot.SlackBot
	Controller *controller.Controller
	HTTPServer *httpserver.Server
	cfg        *configuration.Configuration
	wg         sync.WaitGroup
}

func New(cfg *configuration.Configuration) (stack *Stack, err error) {
	username := os.Getenv("TADO_USERNAME")
	password := os.Getenv("TADO_PASSWORD")
	clientSecret := os.Getenv("TADO_CLIENT_SECRET")

	if username == "" || password == "" {
		return nil, fmt.Errorf("TADO_USERNAME/TADO_PASSWORD environment variables not set")
	}

	API := tado.New(username, password, clientSecret)
	p := poller.New(API)
	h := &health.Health{Poller: p}
	stack = &Stack{
		Poller:    p,
		cfg:       cfg,
		Health:    h,
		Collector: collector.New(p),
		HTTPServer: &httpserver.Server{
			Application: httpserver.Application{
				Name: "tado-exporter",
				Port: cfg.Port,
				Handlers: []httpserver.Handler{
					{Path: "/health", Handler: http.HandlerFunc(h.Handle)},
				},
			},
			Prometheus: httpserver.Prometheus{Port: cfg.Exporter.Port},
		},
	}

	if stack.cfg.Controller.Enabled {
		stack.TadoBot = slackbot.New("tado "+version.BuildVersion, stack.cfg.Controller.TadoBot.Token, nil)
		stack.Controller = controller.New(API, &stack.cfg.Controller, stack.TadoBot, stack.Poller)
	}

	return
}

func (s *Stack) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		s.Poller.Run(ctx, s.cfg.Interval)
		s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		s.Health.Run(ctx)
		s.wg.Done()
	}()

	if s.Collector != nil {
		s.wg.Add(1)
		go func() {
			s.Collector.Run(ctx)
			s.wg.Done()
		}()
		prometheus.MustRegister(s.Collector)
	}

	if s.TadoBot != nil {
		s.wg.Add(1)
		go func() {
			if err := s.TadoBot.Run(ctx); err != nil {
				log.WithError(err).Fatal("tadoBot failed to start")
			}
			s.wg.Done()
		}()
	}

	if s.Controller != nil {
		s.wg.Add(1)
		go func() {
			s.Controller.Run(ctx, 30*time.Second)
			s.wg.Done()
		}()
	}

	s.wg.Add(1)
	go func() {
		log.Info("HTTP server started")
		if errs := s.HTTPServer.Run(); len(errs) > 0 {
			log.WithError(errs[0]).Fatal("failed to start HTTP server")
		}
		log.Info("HTTP server stopped")
		s.wg.Done()
	}()
}

func (s *Stack) Stop() {
	if errs := s.HTTPServer.Shutdown(30 * time.Second); len(errs) > 0 {
		log.WithError(errs[0]).Warning("failed to stop HTTP Server")
	}
	s.wg.Wait()
}
