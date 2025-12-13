package docker

type ContainerCheckResult struct {
	Exist       bool
	Message     string
	ContainerID string
}

const (
	ContainerExists   = "Container exists"
	ContainerNotFound = "Container not found"
)
