package records

import (
	"fmt"
	"strings"
)

func NormalizarValor(reg Registro) string {
	vl := reg.VL

	switch reg.Type {
	case "TXT":
		if vl != "" && !strings.HasPrefix(vl, "\"") {
			vl = fmt.Sprintf("\"%s\"", vl)
		}
	case "CNAME", "NS", "ALIAS", "PTR":
		if vl != "" && !strings.HasSuffix(vl, ".") {
			vl += "."
		}
	case "MX":
		if vl != "" && !strings.HasSuffix(vl, ".") {
			vl += "."
		}

		if reg.Priority != nil {
			vl = fmt.Sprintf("%d %s", *reg.Priority, vl)
		}
	case "CAA":
		if vl != "" && !strings.Contains(vl, "issue") && !strings.Contains(vl, "iodef") {
			vl = fmt.Sprintf("0 issue \"%s\"", vl)
		}
	case "SRV":
		vl = valorSRV(reg)
	case "HTTPS":
		vl = valorHTTPS(reg)
	}

	return vl
}

func valorSRV(reg Registro) string {
	target := reg.Target
	if target != "" && !strings.HasSuffix(target, ".") {
		target += "."
	}

	return fmt.Sprintf("%d %d %d %s", inteiro(reg.Priority), inteiro(reg.Weight), inteiro(reg.Port), target)
}

func valorHTTPS(reg Registro) string {
	targetName := reg.TargetName
	if targetName == "" {
		targetName = "."
	}

	if !strings.HasSuffix(targetName, ".") {
		targetName += "."
	}

	params := strings.ReplaceAll(reg.SvcParams, `"`, "")
	if params == "" {
		params = "alpn=h2"
	}

	return fmt.Sprintf("%d %s %s", inteiro(reg.SvcPriority), targetName, params)
}

func inteiro(valor *int) int {
	if valor == nil {
		return 0
	}

	return *valor
}
