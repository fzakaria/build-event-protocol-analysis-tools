package converter

import (
	"log"

	bes "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/build_event_stream"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventConverter interface {
	Match(*bes.BuildEvent) bool
	Convert(*bes.BuildEvent) (*ParquetEventRow, error)
}

var all = []EventConverter{
	&ProgressConverter{},
	&actionExecutedConverter{},
}

func Convert(ev *bes.BuildEvent) (*ParquetEventRow, error) {
	for _, conv := range all {
		if conv.Match(ev) {
			return conv.Convert(ev)
		}
	}
	return nil, errors.New("no converter found for event type: " + toType(ev))
}

func toType(event *bes.BuildEvent) string {
	switch p := event.Id.Id.(type) {
	case *bes.BuildEventId_Unknown:
		return "UnknownBuildEvent"
	case *bes.BuildEventId_Progress:
		return "Progress"
	case *bes.BuildEventId_Started:
		return "Started"
	case *bes.BuildEventId_UnstructuredCommandLine:
		return "UnstructuredCommandLine"
	case *bes.BuildEventId_StructuredCommandLine:
		return "StructuredCommandLine"
	case *bes.BuildEventId_WorkspaceStatus:
		return "WorkspaceStatus"
	case *bes.BuildEventId_OptionsParsed:
		return "OptionsParsed"
	case *bes.BuildEventId_Fetch:
		return "Fetch"
	case *bes.BuildEventId_Pattern:
		return "Pattern"
	case *bes.BuildEventId_Workspace:
		return "Workspace"
	case *bes.BuildEventId_BuildMetadata:
		return "BuildMetadata"
	case *bes.BuildEventId_TargetConfigured:
		return "TargetConfigured"
	case *bes.BuildEventId_NamedSet:
		return "NameSet"
	case *bes.BuildEventId_Configuration:
		return "Configuration"
	case *bes.BuildEventId_TargetCompleted:
		return "TargetCompleted"
	case *bes.BuildEventId_ActionCompleted:
		return "ActionCompleted"
	case *bes.BuildEventId_UnconfiguredLabel:
		return "UnconfiguredLabel"
	case *bes.BuildEventId_ConfiguredLabel:
		return "ConfiguredLabel"
	case *bes.BuildEventId_TestResult:
		return "TestResult"
	case *bes.BuildEventId_TestProgress:
		return "TestProgress"
	case *bes.BuildEventId_TestSummary:
		return "TestSummary"
	case *bes.BuildEventId_TargetSummary:
		return "TargetSummary"
	case *bes.BuildEventId_BuildFinished:
		return "BuildFinished"
	case *bes.BuildEventId_BuildToolLogs:
		return "BuildToolLogs"
	case *bes.BuildEventId_BuildMetrics:
		return "BuildMetrics"
	case *bes.BuildEventId_ConvenienceSymlinksIdentified:
		return "ConvenienceSymlinksIdentified"
	case *bes.BuildEventId_ExecRequest:
		return "ExecRequest"
	default:
		log.Fatalf("Unknown BuildEvent type: %T", p)
		return "Unknown"
	}
}

func toTimestamp(ts *timestamppb.Timestamp) int64 {
	if ts == nil {
		return 0
	}
	// First convert seconds -> millis
	// add on the nanos
	return ts.Seconds*1_000 + int64(ts.Nanos)/1_000_000
}
