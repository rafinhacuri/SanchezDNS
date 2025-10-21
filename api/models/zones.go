package models

import (
	"fmt"
	"strings"
)

type Soa struct {
	StartOfAuthority string `json:"startOfAuthority"`
	Email            string `json:"email"`
	Refresh          int    `json:"refresh"`
	Retry            int    `json:"retry"`
	Expire           int    `json:"expire"`
	NegativeCacheTtl int    `json:"negativeCacheTtl"`
}

func (s *Soa) Validate() error {
	if s.StartOfAuthority == "" {
		return fmt.Errorf("start of authority is required")
	}
	if s.Email == "" {
		return fmt.Errorf("email is required")
	}
	if s.Refresh <= 0 {
		return fmt.Errorf("refresh must be a positive integer")
	}
	if s.Retry <= 0 {
		return fmt.Errorf("retry must be a positive integer")
	}
	if s.Expire <= 0 {
		return fmt.Errorf("expire must be a positive integer")
	}
	if s.NegativeCacheTtl <= 0 {
		return fmt.Errorf("negative cache ttl must be a positive integer")
	}
	return nil
}

type CreateZoneRequest struct {
	Domain string `json:"domain" binding:"required"`
	Soa    Soa    `json:"soa" binding:"required"`
}

func (req *CreateZoneRequest) Validate() error {
	if req.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if !strings.HasSuffix(req.Domain, ".") {
		req.Domain = req.Domain + "."
	}

	if err := req.Soa.Validate(); err != nil {
		return fmt.Errorf("soa: %w", err)
	}

	return nil
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type rrsetRecord struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	TTL      int      `json:"ttl"`
	Comments []string `json:"comments"`
	Records  []Record `json:"records"`
}

type Zone struct {
	Name         string        `json:"name"`
	RRSets       []rrsetRecord `json:"rrsets"`
	Serial       int64         `json:"serial"`
	EditedSerial int64         `json:"edited_serial"`
}

type Simplified struct {
	Zone        string  `json:"zone"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	VL          string  `json:"vl"`
	TTL         int     `json:"ttl"`
	Comment     string  `json:"comment,omitempty"`
	UpdatedAt   string  `json:"updatedAt,omitempty"`
	SVCPriority *int    `json:"svcPriority,omitempty"`
	TargetName  *string `json:"targetName,omitempty"`
	SVCParams   *string `json:"svcParams,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Target      *string `json:"target,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}
