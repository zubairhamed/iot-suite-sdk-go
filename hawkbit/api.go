package hawkbit

const (
	STATUS_EXEC_CANCELED   = ExecStatus("canceled")
	STATUS_EXEC_PROCEEDING = ExecStatus("proceeding")
	STATUS_EXEC_CLOSED     = ExecStatus("closed")
)

const (
	STATUS_RESULT_NONE    = ResultStatus("none")
	STATUS_RESULT_SUCCESS = ResultStatus("success")
	STATUS_RESULT_FAILED  = ResultStatus("failed")
)

type HawkbitClient interface {
	UseProxy(proxy string) error
	UpdateActionStatus(string, ExecStatus, ResultStatus) error
	Get(string) ([]byte, error)
	GetActions() (*Action, error)
	GetDownloadableChunks() []*Chunk
	ShouldUpdateNow() (bool, string)
	DownloadArtifact(*Artifact) error
	Start()

	// Events
	OnUpdate(fn CallbackFn)
}

type ExecStatus string
type ResultStatus string
type CallbackFn func(HawkbitClient, string)
