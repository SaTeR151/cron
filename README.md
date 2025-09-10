# cron
**cron** - пакет для удобной работы с appCron. 

## Быстрый старт
```
go get github.com/SaTeR151/cron
```

## Пример работы
```go
import (
	"context"
	"github.com/SaTeR151/cron"
	"github.com/sirupsen/logrus"
	"gl.iteco.com/technology/go_general/errproc"
	"log/slog"
)

type Config struct {
	Spec       []string
	RunOnStart bool
	Allowed    bool
}

func GetConfig() Config {
	return Config{
		Spec:       []string{"*/1 * * * * *"},
		RunOnStart: true,
		Allowed:    true,
	}
}

type Fn struct{}

func main() {
	errProc, err := errproc.NewErrProc(logrus.New())
	if err != nil {
		slog.Error(err.Error())
	}
	c, err := cron.NewAppCron(errProc)
	if err != nil {
		slog.Error(err.Error())
	}
	config := GetConfig()
	fn := Fn{}
	c.RegisterJobs(
		cron.Job{
			Config: cron.JobConfig(config),
			Fn:     fn.Foo,
		},
	)
	go c.Start(context.Background())
}

func (fn Fn) Foo(ctx context.Context) error {
	return nil
}
```
