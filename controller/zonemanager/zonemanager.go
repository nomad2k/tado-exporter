package zonemanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/tado"
	"github.com/clambin/tado-exporter/configuration"
	"github.com/clambin/tado-exporter/pkg/scheduler"
	"github.com/clambin/tado-exporter/pkg/slackbot"
	"github.com/clambin/tado-exporter/poller"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Manager struct {
	Updates chan *poller.Update
	config  configuration.ZoneConfig
	task    *Task
	api     tado.API
	poster  Poster
	poller  poller.Poller
	loaded  bool
	lock    sync.RWMutex
}

func New(api tado.API, p poller.Poller, bot slackbot.SlackBot, cfg configuration.ZoneConfig) *Manager {
	return &Manager{
		Updates: make(chan *poller.Update, 1),
		config:  cfg,
		api:     api,
		poster:  Poster{SlackBot: bot},
		poller:  p,
	}
}

func (m *Manager) Run(ctx context.Context, interval time.Duration) {
	m.poller.Register(m.Updates)

	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case update := <-m.Updates:
			if err := m.processUpdate(ctx, update); err != nil {
				log.WithError(err).WithField("zone", m.config.ZoneName).Error("failed to process tado update")
			}
		case <-ticker.C:
			if err := m.processResult(); err != nil {
				log.WithError(err).WithField("zone", m.config.ZoneName).Error("failed to set next state")
			}
		}
	}
	ticker.Stop()

	m.poller.Unregister(m.Updates)
}

func (m *Manager) processUpdate(ctx context.Context, update *poller.Update) (err error) {
	if err = m.load(update); err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}

	// next state
	current, next := m.getNextState(update)

	if current != next.State {
		m.scheduleJob(ctx, next)
	} else {
		m.cancelJob()
	}

	return
}

func (m *Manager) scheduleJob(ctx context.Context, next NextState) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// if the same state is already scheduled for an earlier time, don't schedule it again.
	if m.task != nil &&
		m.task.nextState.State == next.State &&
		m.task.firesBefore(next.Delay) {
		return
	}

	m.task = newTask(m.api, next)
	m.task.job = scheduler.Schedule(ctx, m.task, next.Delay)
	if next.Delay > 0 {
		m.poster.NotifyQueued(next)
	}
}

func (m *Manager) cancelJob() {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.task != nil {
		m.task.job.Cancel()
	}
}

func (m *Manager) processResult() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.task == nil {
		return nil
	}

	completed, err := m.task.job.Result()
	if !completed {
		return nil
	}

	if err == nil {
		m.poster.NotifyAction(m.task.nextState)
	} else if errors.Is(err, scheduler.ErrCanceled) {
		m.poster.NotifyCanceled(m.task.nextState)
		err = nil
	}
	// TODO: reschedule task if it failed?

	m.task = nil

	return err
}

func (m *Manager) Scheduled() (NextState, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.task == nil {
		return NextState{}, false
	}
	return m.task.nextState, true
}

func (m *Manager) ReportTask() (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.task == nil {
		return "", false
	}

	return m.task.Report(), true
}

type Managers []*Manager

func (m Managers) GetScheduled() (states []NextState) {
	for _, mgr := range m {
		if state, scheduled := mgr.Scheduled(); scheduled {
			states = append(states, state)
		}
	}
	return
}

func (m Managers) ReportTasks() ([]string, bool) {
	var tasks []string
	for _, mgr := range m {
		if task, ok := mgr.ReportTask(); ok {
			tasks = append(tasks, task)
		}
	}
	return tasks, len(tasks) > 0
}
