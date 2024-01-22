package worker

import (
	"context"
	"time"

	"github.com/riverqueue/river"
)

type ExampleJobArgs struct {
	Email string `json:"email"`
}

func (ExampleJobArgs) Kind() string {
	return "example"
}

func (ExampleJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 4 * time.Hour,
		},
	}
}

type ExampleWorker struct {
	river.WorkerDefaults[ExampleJobArgs]
}

func (w *ExampleWorker) Work(
	ctx context.Context,
	job *river.Job[ExampleJobArgs],
) error {
	return nil
}

func (w *ExampleWorker) NextRetry(
	job *river.Job[ExampleJobArgs],
) time.Time {
	return time.Now().Add(30 * time.Second)
}

func NewExampleWorker() *ExampleWorker {
	return &ExampleWorker{}
}
