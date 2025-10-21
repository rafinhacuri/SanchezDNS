package models

type AddRecordRequest struct {
	Zone        string `json:"zone"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	VL          string `json:"vl,omitempty"`
	TTL         int    `json:"ttl"`
	Comment     string `json:"comment,omitempty"`
	SvcPriority *int   `json:"svcPriority,omitempty"`
	TargetName  string `json:"targetName,omitempty"`
	SvcParams   string `json:"svcParams,omitempty"`
	Weight      *int   `json:"weight,omitempty"`
	Port        *int   `json:"port,omitempty"`
	Target      string `json:"target,omitempty"`
	Priority    *int   `json:"priority,omitempty"`
}
