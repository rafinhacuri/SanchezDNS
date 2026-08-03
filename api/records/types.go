package records

import "net"

type Registro struct {
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

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int64  `json:"modified_at,omitempty"` //nolint:tagliatelle
}

type RRSet struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	TTL      int       `json:"ttl"`
	Comments []Comment `json:"comments"`
	Records  []Record  `json:"records"`
}

type Zone struct {
	Name         string  `json:"name"`
	RRSets       []RRSet `json:"rrsets"`
	Serial       int64   `json:"serial"`
	EditedSerial int64   `json:"edited_serial"` //nolint:tagliatelle
}

type ZoneInfo struct {
	Name string `json:"name"`
}

type PDNSZonePatchRequest struct {
	RRSets []PDNSRRSetChange `json:"rrsets"`
}

type PDNSRRSetChange struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	TTL        *int      `json:"ttl,omitempty"`
	ChangeType string    `json:"changetype"`
	Records    []Record  `json:"records,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Simplified struct {
	Zone        string  `json:"zone"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	VL          string  `json:"vl"`
	TTL         int     `json:"ttl"`
	Comment     string  `json:"comment,omitempty"`
	SVCPriority *int    `json:"svcPriority,omitempty"`
	TargetName  *string `json:"targetName,omitempty"`
	SVCParams   *string `json:"svcParams,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Target      *string `json:"target,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

type RegistroIP struct {
	Name string
	Type string
	IP   net.IP
}

type ReversoAusente struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	IP          string `json:"ip"`
	NomeReverso string `json:"nomeReverso"`
	ZonaReversa string `json:"zonaReversa"`
}

type ReversoOrfao struct {
	NomeReverso string `json:"nomeReverso"`
	IP          string `json:"ip"`
	Alvo        string `json:"alvo"`
}

type Soa struct {
	StartOfAuthority string `json:"startOfAuthority"`
	Email            string `json:"email"`
	Refresh          int    `json:"refresh"`
	Retry            int    `json:"retry"`
	Expire           int    `json:"expire"`
	NegativeCacheTtl int    `json:"negativeCacheTtl"`
}
