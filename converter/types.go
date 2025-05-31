package converter

type ParquetEventRow struct {
	Type     string    `parquet:"type"`
	Progress *Progress `parquet:"progress,optional"`
	Action   *Action   `parquet:"action,optional"`
}

type Progress struct {
	Stdout string `parquet:"stdout"`
	Stderr string `parquet:"stderr"`
}

type FailureDetail struct {
	Message  string `parquet:"message"`
	Category string `parquet:"category"`
}

type Action struct {
	Type     string `parquet:"type"`
	ExitCode int32  `parquet:"exit_code"`
	Success  bool   `parquet:"success"`
	// TODO: Malloy seems to not handle the LIST logical type
	//CommandLine   []string       `parquet:"command_line,list"`
	StartTime     int64          `parquet:"start_time,timestamp,optional"`
	EndTime       int64          `parquet:"end_time,timestamp,optional"`
	FailureDetail *FailureDetail `parquet:"failure_detail,optional"`
}
