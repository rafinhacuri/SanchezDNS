package models

import (
	"fmt"
	"strings"
)

type soa struct {
	StartOfAuthority string `json:"startOfAuthority"`
	Email            string `json:"email"`
	Refresh          int    `json:"refresh"`
	Retry            int    `json:"retry"`
	Expire           int    `json:"expire"`
	NegativeCacheTtl int    `json:"negativeCacheTtl"`
}

func (s *soa) Validate() error {
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
	Soa    soa    `json:"soa" binding:"required"`
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
