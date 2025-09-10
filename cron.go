package cron

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
)

// RegisterJobs принимает slice cron.Job работ
//
//	type Job struct{
//		Config JobConfig
//		Fn func (context.Context) error
//	}
//
//	type JobConfig struct{
//		Spec []string - Время когда вызывать фунцкии -  "*/1 * * * * *" - sec, min, hour, day, month, day (week)
//		RunOnStart bool - запуск при инициализации
//		Allowed - Включение/отключение крона
//	}
func (c *Cron) RegisterJobs(jobs ...Job) {
	c.jobs = append(c.jobs, jobs...)
}

// Start запускает cron
func (c *Cron) Start(ctx context.Context) error {
	c.logger.Info("starting on start jobs")
	for _, job := range c.jobs {
		if job.Config.Allowed && job.Config.RunOnStart {
			c.RunJob(ctx, job.Fn)
		}
	}

	c.logger.Info("on start jobs completed")

	for _, job := range c.jobs {
		if job.Config.Allowed {
			for _, spec := range job.Config.Spec {
				_, err := c.cron.AddFunc(spec, func() {
					c.runJob(ctx, job.Fn)
				})
				if err != nil {
					return err
				}
			}
		}
	}

	c.logger.Info("AppCron has been started!")
	c.cron.Start()
	defer c.cron.Stop()

	<-ctx.Done()

	return nil
}

func (c *Cron) runJob(ctx context.Context, fn func(context.Context) error) {
	if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Error("cron job failed", "error", err)
		c.handlerError(err)
	}
}

func (c *Cron) RunJob(ctx context.Context, fn func(context.Context) error) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("cron job panicked", "recover", r, "stack", string(debug.Stack()))
			c.handlerError(fmt.Errorf("panic: %v", r))
		}
	}()

	if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Error("cron job failed", "error", err)
		c.handlerError(err)
	}
}

func (c *Cron) handlerError(err error) {
	c.errProc.NewException(err).Capture()
}
