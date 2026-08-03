package logs

type Log struct {
	Zone      string `bson:"zone"      json:"zone"`
	Username  string `bson:"username"  json:"username"`
	Action    string `bson:"action"    json:"action"`
	Details   string `bson:"details"   json:"details"`
	CreatedAt string `bson:"createdAt" json:"createdAt"`
}

type LogsResponse struct {
	Logs  []Log `json:"logs"`
	Total int64 `json:"total"`
}
