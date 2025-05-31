package converter

import (
	"errors"

	bes "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/build_event_stream"
)

type ProgressConverter struct{}

func (c *ProgressConverter) Match(ev *bes.BuildEvent) bool {
	_, ok := ev.Payload.(*bes.BuildEvent_Progress)
	return ok
}

func (c *ProgressConverter) Convert(ev *bes.BuildEvent) (*ParquetEventRow, error) {
	progress := ev.GetProgress()
	if progress == nil {
		return nil, errors.New("missing progress payload")
	}
	return &ParquetEventRow{
		Type: "Progress",
		Progress: &Progress{
			Stdout: progress.Stdout,
			Stderr: progress.Stderr,
		},
	}, nil
}
