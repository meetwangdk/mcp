package v1

const (
	// ContainerDiskPressure means the pod container is under disk pressure.
	ContainerDiskPressure PodConditionType = "ContainerDiskPressure"
	// ContainersLivenessProbePassed indicates whether all containers pass Livenessprobe.
	ContainersLivenessProbePassed PodConditionType = "ContainersLivenessProbePassed"
)
