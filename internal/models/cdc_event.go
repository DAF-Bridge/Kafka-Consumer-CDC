package models

type CDCEvent struct {
	Schema  map[string]interface{} `json:"schema"`
	Payload CDCPayload             `json:"payload"`
}

type CDCPayload struct {
	Before map[string]interface{} `json:"before"` // Data before the change (for update/delete)
	After  map[string]interface{} `json:"after"`  // Data after the change (for create/update)
	Op     string                 `json:"op"`     // Operation type: "c" (create), "u" (update), "d" (delete), "r" (read)
	Source CDCSource              `json:"source"` // Metadata about the change
	TsMs   int64                  `json:"ts_ms"`  // Timestamp of the change event
}

type CDCSource struct {
	Table string `json:"table"` // The name of the table where the change occurred
}