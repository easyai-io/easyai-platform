package schema

import (
	"fmt"
	"sort"
	"time"
)

// Job defines the job
type Job struct {

	// for db only
	ID     uint32 `json:"id" yaml:"id"`
	UkUUID string `json:"uk_uuid" yaml:"uk_uuid"`

	// ============= common config define for a single-pod job or any-framework's worker/default config ==================
	Name            string            `json:"name" yaml:"name" binding:"required" example:"a demo 作业"`
	Namespace       string            `json:"namespace" yaml:"namespace" example:"default"`
	Cluster         string            `json:"cluster" yaml:"cluster" example:"dev"`
	Description     string            `json:"description" yaml:"description" example:"desc作业描述"`
	Owner           string            `json:"owner" yaml:"owner"`
	NodeSelectors   map[string]string `json:"node_selectors" yaml:"node_selectors" example:"key:value"`
	Tolerations     []string          `json:"tolerations" yaml:"tolerations" example:"key=value"`
	Image           string            `json:"image" yaml:"image" binding:"required" example:"registry.cn-shanghai.aliyuncs.com/easyai-io/easyai-demo:tf2.4.3-gpu-jupyter-lab"`
	ImagePullPolicy string            `json:"image_pull_policy" yaml:"image_pull_policy" example:"IfNotPresent"`
	Envs            map[string]string `json:"envs" yaml:"envs" example:"key1:value1"`
	// resources
	CPU       float32 `json:"cpu" yaml:"cpu" binding:"required" example:"2.5"`
	Memory    float32 `json:"memory" yaml:"memory" binding:"required" example:"8.0"`
	GPU       float32 `json:"gpu" yaml:"gpu" example:"0"`
	GPUMemory float32 `json:"gpu_memory" yaml:"gpu_memory" example:"0"`
	Hardware  string  `json:"hardware" yaml:"hardware"` // cpu or gpu
	// code logic
	Workspace      string            `json:"workspace" yaml:"workspace" example:"/workspace"`
	EntrypointType string            `json:"entrypoint_type" yaml:"entrypoint_type" binding:"required" example:"bash -c"`
	Entrypoint     string            `json:"entrypoint" yaml:"entrypoint" example:"echo hello world"`
	RetryCount     int32             `json:"retry_count" yaml:"retry_count" example:"0"`
	WorkerCount    int32             `json:"worker_count" yaml:"worker_count" binding:"required" example:"1"`
	MaxRetry       int32             `json:"max_retry" yaml:"max_retry" example:"0"`
	VolumeMounts   map[string]string `json:"volume_mounts" yaml:"volume_mounts"`
	IsNonRoot      bool              `json:"is_non_root" yaml:"is_non_root" example:"false"`
	// labels && annotations
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
	MagicFlags  map[string]string `json:"magic_flags" yaml:"magic_flags"`
	TrainArgs   []string          `json:"train_args" yaml:"train_args" example:"lr=0.1,async"`
	// ============= common config end ==================
	// for job status
	State         JobState   `json:"state" yaml:"state"`
	Status        JobStatus  `json:"status" yaml:"status"`
	Message       string     `json:"message" yaml:"message"`
	Reason        string     `json:"reason" yaml:"reason"`
	Result        string     `json:"result" yaml:"result"`
	IsDeleted     DeleteFlag `json:"is_deleted" yaml:"is_deleted"`
	Duration      int32      `json:"duration" yaml:"duration"`
	StartTime     time.Time  `json:"start_time" yaml:"start_time"`
	SchedulerTime time.Time  `json:"scheduler_time" yaml:"scheduler_time"`
	CreatedAt     time.Time  `json:"created_at" yaml:"created_at"`
	ModifiedAt    time.Time  `json:"modified_at" yaml:"modified_at"`
	// for ml framework config, e.g. distribute tfjob
	Framework       JobFramework    `json:"framework" yaml:"framework" binding:"required" example:"tensorflow"`
	FrameworkConfig FrameworkConfig `json:"framework_config" yaml:"framework_config"`

	// the task list of a job
	Tasks []*Task `json:"tasks" yaml:"tasks"`

	// for frontend web use only
	FrontendField `json:",inline"`
}

// Task task
type Task struct {
	ID          uint32    `json:"id"`
	JobID       uint32    `json:"job_id"`
	Role        string    `json:"role"`
	Index       int8      `json:"index"`
	Command     string    `json:"command"`
	Args        string    `json:"args"`
	OverSale    float32   `json:"over_sale"`
	CPU         float32   `json:"cpu"`
	Memory      float32   `json:"memory"`
	GPU         float32   `json:"gpu"`
	Status      JobStatus `json:"status"`
	Reason      string    `json:"reason"`
	Message     string    `json:"message"`
	ExitCode    string    `json:"exit_code"`
	PodName     string    `json:"pod_name"`
	PodIP       string    `json:"pod_ip"`
	NodeName    string    `json:"node_name"`
	NodeIP      string    `json:"node_ip"`
	ContainerID string    `json:"container_id"`
	Devices     string    `json:"devices"`
	StartTime   time.Time `json:"start_time"`
	Port        int32     `json:"port"`
	NeedCheck   bool      `json:"need_check"`
	RetryCount  int8      `json:"retry_count"`
	MaxRetry    int8      `json:"max_retry"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	// for frontend web use only
	Services  map[string]string `json:"services,omitempty"`
	Monitor   map[string]string `json:"monitor,omitempty"`
	ClientOps map[string]string `json:"client_operation,omitempty"`
}

// RangeTasks ...
func (j *Job) RangeTasks(f func(t *Task)) {
	for _, item := range j.Tasks {
		f(item)
	}
}

// WrapFrontendFiled for fe filed to present, depreciated
func (j *Job) WrapFrontendFiled() *Job {
	f := FrontendField{
		FeTasks:    map[string][]*Task{},
		FeResource: map[string]string{},
	}
	for _, t := range j.Tasks {
		f.FeTasks[t.Role] = append(f.FeTasks[t.Role], t)
	}
	for role, tg := range f.FeTasks {
		f.FeResource[role] = fmt.Sprintf("%d * (CPU:%f Memory:%fGB GPU: %f)", len(tg), tg[0].CPU, tg[0].Memory, tg[0].GPU)
	}
	j.FrontendField = f
	return j
}

type FrameworkConfig interface {
	Framework() JobFramework
}

func (n *JobFramework) Framework() JobFramework { return *n }

// TFJobConfig tf job config, framework=="tf"
type TFJobConfig struct {
	JobFramework      `json:"-"`
	RoleNodeSelectors map[string]map[string]string `json:"role_node_selectors,omitempty"`
	CleanPolicy       string                       `json:"clean_policy,omitempty"`

	Port int32 `json:"port,omitempty"`
	// worker config
	WorkerImage         string  `json:"worker_image,omitempty"`
	WorkerPort          int32   `json:"worker_port,omitempty"`
	WorkerCPU           float32 `json:"worker_cpu,omitempty"`
	WorkerMemory        float32 `json:"worker_memory,omitempty"`
	WorkerGPU           float32 `json:"worker_gpu,omitempty"`
	WorkerGPUMemory     float32 `json:"worker_gpu_memory,omitempty"`
	WorkerCount         int32   `json:"worker_count,omitempty"`
	WorkerRestartPolicy string  `json:"worker_restart_policy,omitempty"`
	// ps config
	PsImage         string  `json:"ps_image,omitempty"`
	PsPort          int32   `json:"ps_port,omitempty"`
	PsCPU           float32 `json:"ps_cpu,omitempty"`
	PsMemory        float32 `json:"ps_memory,omitempty"`
	PsGPU           float32 `json:"ps_gpu,omitempty"`
	PsCount         int32   `json:"ps_count,omitempty"`
	PsRestartPolicy string  `json:"ps_restart_policy,omitempty"`
	// chief config
	UseChief           bool    `json:"use_chief,omitempty"`
	ChiefImage         string  `json:"chief_image,omitempty"`
	ChiefPort          int32   `json:"chief_port,omitempty"`
	ChiefCPU           float32 `json:"chief_cpu,omitempty"`
	ChiefMemory        float32 `json:"chief_memory,omitempty"`
	ChiefGPU           float32 `json:"chief_gpu,omitempty"`
	ChiefGPUMemory     float32 `json:"chief_gpu_memory,omitempty"`
	ChiefCount         int32   `json:"chief_count,omitempty"`
	ChiefRestartPolicy string  `json:"chief_restart_policy,omitempty"`
	// evaluator config
	UseEvaluator           bool    `json:"use_evaluator,omitempty"`
	EvaluatorImage         string  `json:"evaluator_image,omitempty"`
	EvaluatorPort          int32   `json:"evaluator_port,omitempty"`
	EvaluatorCPU           float32 `json:"evaluator_cpu,omitempty"`
	EvaluatorMemory        float32 `json:"evaluator_memory,omitempty"`
	EvaluatorGPU           float32 `json:"evaluator_gpu,omitempty"`
	EvaluatorGPUMemory     float32 `json:"evaluator_gpu_memory,omitempty"`
	EvaluatorRestartPolicy string  `json:"evaluator_restart_policy,omitempty"`
	// master config  for now, I don't know which distribute case should use master
	MasterImage         string  `json:"master_image,omitempty"`
	MasterPort          int32   `json:"master_port,omitempty"`
	MasterCPU           float32 `json:"master_cpu,omitempty"`
	MasterMemory        float32 `json:"master_memory,omitempty"`
	MasterGPU           float32 `json:"master_gpu,omitempty"`
	MasterGPUMemory     float32 `json:"master_gpu_memory,omitempty"`
	MasterCount         int32   `json:"master_count,omitempty"`
	MasterRestartPolicy string  `json:"master_restart_policy,omitempty"`
}

// PytorchJobConfig pytorch job config, framework=="pytorch"
type PytorchJobConfig struct {
	JobFramework `json:"-"`
	// todo implement me!
}

// NotebookJobConfig Notebook job config, framework=="pytorch"
type NotebookJobConfig struct {
	// todo implement me!
	JobFramework    `json:"-"`
	JupyterServerID uint32 `json:"jupyter_server_id"`
}

// TFJobConfig tf config
func (j *Job) TFJobConfig() (TFJobConfig, error) {
	var config TFJobConfig
	return config, nil
}

// NotebookJobConfig Notebook Job Config
func (j *Job) NotebookJobConfig() (NotebookJobConfig, error) {
	var config NotebookJobConfig
	return config, nil
}

// PytorchJobConfig Pytorch Job Config
func (j *Job) PytorchJobConfig() (PytorchJobConfig, error) {
	var config PytorchJobConfig
	return config, nil
}

// MagicFlag 魔法参数
type MagicFlag struct {
	Key    string `json:"key"`
	DValue string `json:"default_value"`
	Alias  string `json:"alias"`
	Desc   string `json:"description"`
}

// FrameworkItem framework
type FrameworkItem struct {
	Name JobFramework `json:"name"`
	Desc string       `json:"description"`
}

// FrontendField for web view
type FrontendField struct {
	FeTasks    map[string][]*Task `json:"fe_tasks"`
	FeResource map[string]string  `json:"fe_resource"`
	JobMonitor string             `json:"job_monitor"`
}

// SortTasks sort tasks
func (ff FrontendField) SortTasks() {
	for _, tasks := range ff.FeTasks {
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Index < tasks[j].Index
		})
	}
}

// InputJobConfig 用于创建job
type InputJobConfig struct {
	// scheme.CommonConfig
	Name            string            `json:"name" yaml:"name" binding:"required"`
	Cluster         string            `json:"cluster" yaml:"cluster"`
	Namespace       string            `json:"namespace" yaml:"namespace"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	NodeSelectors   map[string]string `json:"node_selectors" yaml:"node_selectors" example:"key:value"`
	Tolerations     []string          `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Image           string            `json:"image" yaml:"image"`
	ImagePullPolicy string            `json:"image_pull_policy,omitempty" yaml:"image_pull_policy,omitempty"`
	Envs            map[string]string `json:"envs,omitempty" yaml:"envs,omitempty"`
	// resources
	CPU       float32 `json:"cpu" yaml:"cpu"`
	Memory    float32 `json:"memory" yaml:"memory"`
	GPU       float32 `json:"gpu,omitempty" yaml:"gpu,omitempty"`
	GPUMemory float32 `json:"gpu_memory,omitempty" yaml:"gpu_memory,omitempty"`
	// code logic
	Workspace      string            `json:"workspace,omitempty" yaml:"workspace,omitempty"`
	EntrypointType string            `json:"entrypoint_type" yaml:"entrypoint_type"`
	Entrypoint     string            `json:"entrypoint" yaml:"entrypoint"`
	TrainArgs      []string          `json:"train_args,omitempty" yaml:"train_args,omitempty"`
	MagicFlags     map[string]string `json:"magic_flags,omitempty" yaml:"magic_flags,omitempty"`
	WorkerCount    int32             `json:"worker_count" yaml:"worker_count"`
	MaxRetry       int32             `json:"max_retry,omitempty" yaml:"max_retry,omitempty"`
	VolumeMounts   map[string]string `json:"volume_mounts,omitempty" yaml:"volume_mounts,omitempty"`
	IsNonRoot      bool              `json:"is_non_root,omitempty" yaml:"is_non_root,omitempty"`
	// labels && annotations
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	// for job
	Framework       JobFramework `json:"framework" yaml:"framework"`
	FrameworkConfig interface{}  `json:"framework_config,omitempty" yaml:"framework_config,omitempty"`
}
