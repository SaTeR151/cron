package cron

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"
	"gl.iteco.com/technology/go_general/errproc"
)

const (
	_defaultJobsCount = 8
)

type (
	Cron struct {
		cron    *cron.Cron
		jobs    []Job
		logger  *slog.Logger
		errProc *errproc.BaseErrProc
	}

	Job struct {
		Config JobConfig
		Fn     func(context.Context) error
	}

	JobConfig struct {
		Spec       []string // Spec Время когда вызывать фунцкии
		RunOnStart bool     // Allowed Включение/отключение крона
		Allowed    bool     // RunOnStart запуск при инициализации
	}
)

// NewAppCron возвращает cron с установленным errProc и базовыми для работы настройками
func NewAppCron(errProc *errproc.ErrProc) (*Cron, error) {
	baseErrProc, err := errProc.NewBaseErrProc("AppCron")
	if err != nil {
		return nil, err
	}

	logger := slog.With("service", "AppCron")

	cLogger := NewLogger(logger)

	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cLogger),
			cron.SkipIfStillRunning(cLogger),
		),
	)

	return &Cron{
		cron:    c,
		jobs:    make([]Job, 0, _defaultJobsCount),
		logger:  logger,
		errProc: baseErrProc,
	}, nil
}

// NewAppCronWithOptions обертка для cron.New() + slog
//
// Подрбнее на https://github.com/robfig/cron
func NewAppCronWithOptions(opts ...cron.Option) *Cron {
	logger := slog.With("service", "AppCron")

	c := cron.New()

	for _, opt := range opts {
		opt(c)
	}

	return &Cron{
		cron:   c,
		jobs:   make([]Job, 0, _defaultJobsCount),
		logger: logger,
	}
}
