package debezium

import "encoding/json"

// Source is the source of the event
type Source struct {
	Version     string `json:"version"`
	Connector   string `json:"connector"`
	Name        string `json:"name"`
	TimestampMS int64  `json:"ts_ms"`
	Snapshot    bool   `json:"snapshot"`
	Database    string `json:"db"`
	ReplicaSet  string `json:"rs"`
	Collection  string `json:"collection"`
	Ordinal     int    `json:"ord"`
	Hash        int64  `json:"h"`
}

// Payload is the payload of the debezium event
type Payload struct {
	After       string `json:"after"`
	Before      string `json:"before"`
	Op          string `json:"op"`
	Patch       string `json:"patch"`
	Filter      string `json:"filter"`
	Source      Source `json:"source"`
	TimestampMs int64  `json:"ts_ms"`
}

// MongoDBEvent represents a complete MongoDB Debezium connector Event
type MongoDBEvent struct {
	Payload Payload         `json:"payload"`
	Schema  json.RawMessage `json:"schema"`
}
