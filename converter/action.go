package converter

import (
	"errors"
	"log"

	bes "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/build_event_stream"
	fd "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/failure_details"
)

type actionExecutedConverter struct{}

func (c *actionExecutedConverter) Match(ev *bes.BuildEvent) bool {
	_, ok := ev.Payload.(*bes.BuildEvent_Action)
	return ok
}

func (c *actionExecutedConverter) Convert(ev *bes.BuildEvent) (*ParquetEventRow, error) {
	action := ev.GetAction()
	if action == nil {
		return nil, errors.New("missing ActionExecuted payload")
	}

	var failureDetails *FailureDetail = nil
	if action.FailureDetail != nil {
		failureDetails = &FailureDetail{
			Message:  action.FailureDetail.Message,
			Category: toCategory(action.FailureDetail),
		}
	}

	return &ParquetEventRow{
		Type: "Progress",
		Action: &Action{
			Type:     action.Type,
			ExitCode: action.ExitCode,
			Success:  action.Success,
			// TODO: Malloy seems to not handle the LIST logical type
			//CommandLine:   action.CommandLine,
			StartTime:     toTimestamp(action.StartTime),
			EndTime:       toTimestamp(action.EndTime),
			FailureDetail: failureDetails,
		},
	}, nil
}

func toCategory(failure_detail *fd.FailureDetail) string {
	if failure_detail.Category == nil {
		return "Unknown"
	}
	switch p := failure_detail.Category.(type) {
	case *fd.FailureDetail_Interrupted:
		return "Interrupted"
	case *fd.FailureDetail_ExternalRepository:
		return "ExternalRepository"
	case *fd.FailureDetail_BuildProgress:
		return "BuildProgress"
	case *fd.FailureDetail_RemoteOptions:
		return "RemoteOptions"
	case *fd.FailureDetail_ClientEnvironment:
		return "ClientEnvironment"
	case *fd.FailureDetail_Crash:
		return "Crash"
	case *fd.FailureDetail_SymlinkForest:
		return "SymlinkForest"
	case *fd.FailureDetail_PackageOptions:
		return "PackageOptions"
	case *fd.FailureDetail_RemoteExecution:
		return "RemoteExecution"
	case *fd.FailureDetail_Execution:
		return "Execution"
	case *fd.FailureDetail_Workspaces:
		return "Workspaces"
	case *fd.FailureDetail_CrashOptions:
		return "CrashOptions"
	case *fd.FailureDetail_Filesystem:
		return "Filesystem"
	case *fd.FailureDetail_ExecutionOptions:
		return "ExecutionOptions"
	case *fd.FailureDetail_Command:
		return "Command"
	case *fd.FailureDetail_Spawn:
		return "Spawn"
	case *fd.FailureDetail_GrpcServer:
		return "GrpcServer"
	case *fd.FailureDetail_CanonicalizeFlags:
		return "CanonicalizeFlags"
	case *fd.FailureDetail_BuildConfiguration:
		return "BuildConfiguration"
	case *fd.FailureDetail_InfoCommand:
		return "InfoCommand"
	case *fd.FailureDetail_MemoryOptions:
		return "MemoryOptions"
	case *fd.FailureDetail_Query:
		return "Query"
	case *fd.FailureDetail_LocalExecution:
		return "LocalExecution"
	case *fd.FailureDetail_ActionCache:
		return "ActionCache"
	case *fd.FailureDetail_FetchCommand:
		return "FetchCommand"
	case *fd.FailureDetail_SyncCommand:
		return "SyncCommand"
	case *fd.FailureDetail_Sandbox:
		return "Sandbox"
	case *fd.FailureDetail_IncludeScanning:
		return "IncludeScanning"
	case *fd.FailureDetail_TestCommand:
		return "TestCommand"
	case *fd.FailureDetail_ActionQuery:
		return "ActionQuery"
	case *fd.FailureDetail_TargetPatterns:
		return "TargetPatterns"
	case *fd.FailureDetail_CleanCommand:
		return "CleanCommand"
	case *fd.FailureDetail_ConfigCommand:
		return "ConfigCommand"
	case *fd.FailureDetail_ConfigurableQuery:
		return "ConfigurableQuery"
	case *fd.FailureDetail_DumpCommand:
		return "DumpCommand"
	case *fd.FailureDetail_HelpCommand:
		return "HelpCommand"
	case *fd.FailureDetail_MobileInstall:
		return "MobileInstall"
	case *fd.FailureDetail_ProfileCommand:
		return "ProfileCommand"
	case *fd.FailureDetail_RunCommand:
		return "RunCommand"
	case *fd.FailureDetail_VersionCommand:
		return "VersionCommand"
	case *fd.FailureDetail_PrintActionCommand:
		return "PrintActionCommand"
	case *fd.FailureDetail_WorkspaceStatus:
		return "WorkspaceStatus"
	case *fd.FailureDetail_JavaCompile:
		return "JavaCompile"
	case *fd.FailureDetail_ActionRewinding:
		return "ActionRewinding"
	case *fd.FailureDetail_CppCompile:
		return "CppCompile"
	case *fd.FailureDetail_StarlarkAction:
		return "StarlarkAction"
	case *fd.FailureDetail_NinjaAction:
		return "NinjaAction"
	case *fd.FailureDetail_DynamicExecution:
		return "DynamicExecution"
	case *fd.FailureDetail_FailAction:
		return "FailAction"
	case *fd.FailureDetail_SymlinkAction:
		return "SymlinkAction"
	case *fd.FailureDetail_CppLink:
		return "CppLink"
	case *fd.FailureDetail_LtoAction:
		return "LtoAction"
	case *fd.FailureDetail_TestAction:
		return "TestAction"
	case *fd.FailureDetail_Worker:
		return "Worker"
	case *fd.FailureDetail_Analysis:
		return "Analysis"
	case *fd.FailureDetail_PackageLoading:
		return "PackageLoading"
	case *fd.FailureDetail_Toolchain:
		return "Toolchain"
	case *fd.FailureDetail_StarlarkLoading:
		return "StarlarkLoading"
	case *fd.FailureDetail_ExternalDeps:
		return "ExternalDeps"
	case *fd.FailureDetail_DiffAwareness:
		return "DiffAwareness"
	case *fd.FailureDetail_ModCommand:
		return "ModCommand"
	case *fd.FailureDetail_BuildReport:
		return "BuildReport"
	case *fd.FailureDetail_Skyfocus:
		return "Skyfocus"
	case *fd.FailureDetail_RemoteAnalysisCaching:
		return "RemoteAnalysisCaching"
	default:
		log.Fatalf("Unknown Category type: %T", p)
		return "Unknown"
	}
}
