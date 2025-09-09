package cron

import (
	"log/slog"

	"gl.iteco.com/technology/go_general/errproc"
)

// SetErrProc устанавливает errProc в cron
func (c *Cron) SetErrProc(errProc *errproc.ErrProc) error {
	var err error
	c.errProc, err = errProc.NewBaseErrProc("AppCron")
	return err
}

// SetLogger устанавливает logger в cron
func (c *Cron) SetLogger(logger *slog.Logger) {
	c.logger = logger
}
