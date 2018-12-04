package record

type LogRecord struct {
	Stream     string `json:"stream"`
	Time       string `json:"time"`
	Log        string `json:"log"`
	Kubernetes struct {
		Namespace     string            `json:"namespace_name"`
		PodName       string            `json:"pod_name"`
		ContainerName string            `json:"container_name"`
		Host          string            `json:"host"`
		Labels        map[string]string `json:"labels"`
	} `json:"kubernetes"`
}

type LogRecords []LogRecord
