package schema

// JobFramework framework
type JobFramework string

const (
	// FrameworkTensorflow tf
	FrameworkTensorflow JobFramework = "tensorflow"
	// FrameworkPytorch pytorch
	FrameworkPytorch JobFramework = "pytorch"
	// FrameworkMPI mpi
	FrameworkMPI JobFramework = "mpi"
	// FrameworkMxnet mxnet
	FrameworkMxnet JobFramework = "mxnet"
	// FrameworkXgboost xgb
	FrameworkXgboost JobFramework = "xgboost"

	// FrameworkVolcano vc
	FrameworkVolcano JobFramework = "volcano"
	// FrameworkVCTensorflow vc-tf
	FrameworkVCTensorflow JobFramework = "vc-tensorflow"

	// FrameworkNotebook nb
	FrameworkNotebook JobFramework = "notebook"
	// FrameworkTensorboard tb
	FrameworkTensorboard JobFramework = "tensorboard"
	// FrameworkHVD hvd
	FrameworkHVD JobFramework = "horovod"
	// FrameworkPaddle paddle
	FrameworkPaddle JobFramework = "paddle"
	// FrameworkAny any
	FrameworkAny JobFramework = "any"
)

// JobStatus job status
type JobStatus string // job status

const (
	// JobStatusQueued  queue
	JobStatusQueued JobStatus = "Queued"
	// JobStatusPending pending
	JobStatusPending JobStatus = "Pending"
	// JobStatusRunning  running
	JobStatusRunning JobStatus = "Running"
	// JobStatusSucceeded ...
	JobStatusSucceeded JobStatus = "Succeeded"
	// JobStatusFailed ...
	JobStatusFailed JobStatus = "Failed"
	// JobStatusUnknown ...
	JobStatusUnknown JobStatus = "Unknown"
	// JobStatusStopped ...
	JobStatusStopped JobStatus = "Stopped"
	// JobStatusRestarting ...
	JobStatusRestarting JobStatus = "Restarting"
)

// JobState job state
type JobState string // job state

const (
	// JobStateInit ...
	JobStateInit JobState = "Init"
	// JobStateStart ...
	JobStateStart JobState = "Start"
	// JobStateStop ...
	JobStateStop JobState = "Stop"
	// JobStateFinish ...
	JobStateFinish JobState = "Finish"
)

// JsStatus jupyter server status
type JsStatus string

// jupyter server status
const (
	JsStopped    JsStatus = "Stopped"    // 已停止/未运行, Stopped 只能变成 Submitted
	JsSubmitted  JsStatus = "Submitted"  // 已提交
	JsScheduling JsStatus = "Scheduling" // 等待调度
	JsStarting   JsStatus = "Starting"   // 启动中
	JsRunning    JsStatus = "Running"    // 运行中
)

// JsImageType jupyter server 镜像类型
type JsImageType string

// for jupyter server image
const (
	JsImageTraining  JsImageType = "training_image"  // 训练任务镜像
	JsImageInferring JsImageType = "inferring_image" // 推理服务镜像
	JsImageCustom    JsImageType = "custom_image"    // 自定义镜像
	JsImageSnapshot  JsImageType = "snapshot_image"  // 快照镜像
)

// NamespaceType cluster namespace type
type NamespaceType string

// namespace type
const (
	NsTraining NamespaceType = "training"
	NsServing  NamespaceType = "serving"
)
