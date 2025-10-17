package models

import (
	"fmt"
	"net"
	"strings"
)

type CreateZoneRequest struct {
	Name       string   `json:"name" binding:"required"`
	Kind       string   `json:"kind" binding:"required"`
	Masters    []string `json:"masters,omitempty"`
	AlsoNotify []string `json:"also_notify,omitempty"`
	SoaEditApi string   `json:"soa_edit_api,omitempty"`
}

func (req *CreateZoneRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !strings.HasSuffix(req.Name, ".") {
		return fmt.Errorf("name must end with a dot")
	}

	validKinds := map[string]bool{"Native": true, "Primary": true, "Secondary": true}
	if !validKinds[req.Kind] {
		return fmt.Errorf("kind must be one of: Native, Primary, Secondary")
	}

	if req.Kind == "Secondary" {
		if len(req.Masters) == 0 {
			return fmt.Errorf("masters must contain at least one IP when kind is Secondary")
		}
		for _, ip := range req.Masters {
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("masters contains invalid IP address: %s", ip)
			}
		}
	}

	if req.Kind == "Primary" && len(req.AlsoNotify) > 0 {
		for _, ip := range req.AlsoNotify {
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("also_notify contains invalid IP address: %s", ip)
			}
		}
	}

	if req.SoaEditApi != "" {
		validSoaEditApis := map[string]bool{"DEFAULT": true, "INCREASE": true, "EPOCH": true, "OFF": true}
		if !validSoaEditApis[req.SoaEditApi] {
			return fmt.Errorf("soa_edit_api must be one of: DEFAULT, INCREASE, EPOCH, OFF")
		}
	}

	return nil
}
