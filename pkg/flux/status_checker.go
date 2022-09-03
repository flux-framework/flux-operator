package flux

// StatusChecker checks status of clusterQueue.
type StatusChecker interface {
	// ClusterQueueActive returns whether the clusterQueue is active.
	ClusterQueueActive(name string) bool
}
