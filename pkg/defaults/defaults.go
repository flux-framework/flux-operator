package defaults

import "time"

// https://github.com/kubernetes-sigs/kueue/blob/main/pkg/constants/constants.go

const (
	QueueName         = "flux"
	JobControllerName = QueueName + "-job-controller"
	AdmissionName     = QueueName + "-admission"

	// UpdatesBatchPeriod is the batch period to hold workload updates
	// before syncing a Queue and ClusterQueue objects.
	UpdatesBatchPeriod = time.Second

	// DefaultPriority is used to set priority of workloads
	// that do not specify any priority class and there is no priority class
	// marked as default.
	DefaultPriority = 0
)
