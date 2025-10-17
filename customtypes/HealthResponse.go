package CustomTypes

type Health struct {
	Jobs    Jobs    `json:"jobs"`
	Workers Workers `json:"workers"`
}

type Jobs struct {
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	InProgress int `json:"inProgress"`
	InQueue    int `json:"inQueue"`
	Retried    int `json:"retried"`
}

type Workers struct {
	Idle         int `json:"idle"`
	Initializing int `json:"initializing"`
	Ready        int `json:"ready"`
	Running      int `json:"running"`
	Throttled    int `json:"throttled"`
	Unhealthy    int `json:"unhealthy"`
}

type RequestResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type JobStatus struct {
	DelayTime      int32  `json:"delayTime"`
	ID             string `json:"id"`
	ConversationID string `json:"conversationID"`
	Status         string `json:"status"`
}
