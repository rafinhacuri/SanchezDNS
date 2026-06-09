package records

import (
	"fmt"
	"strings"
)

func normalizeRecordValue(
	tipo,
	vl string,
	priority,
	weight,
	port *int,
	target string,
	svcPriority *int,
	targetName,
	svcParams string,
) string {
	if tipo == "TXT" && vl != "" && !strings.HasPrefix(vl, "\"") {
		vl = fmt.Sprintf("\"%s\"", vl)
	}

	if (tipo == "CNAME" || tipo == "NS" || tipo == "ALIAS" || tipo == "MX" || tipo == "PTR") && vl != "" {
		if !strings.HasSuffix(vl, ".") {
			vl += "."
		}
	}

	if tipo == "MX" && priority != nil {
		vl = fmt.Sprintf("%d %s", *priority, vl)
	}

	if tipo == "CAA" && vl != "" && !strings.Contains(vl, "issue") && !strings.Contains(vl, "iodef") {
		vl = fmt.Sprintf("0 issue \"%s\"", vl)
	}

	if tipo == "SRV" {
		p := 0
		w := 0
		prt := 0
		tgt := target

		if priority != nil {
			p = *priority
		}

		if weight != nil {
			w = *weight
		}

		if port != nil {
			prt = *port
		}

		if tgt != "" && !strings.HasSuffix(tgt, ".") {
			tgt += "."
		}

		vl = fmt.Sprintf("%d %d %d %s", p, w, prt, tgt)
	}

	if tipo == "HTTPS" {
		sp := 0
		tn := targetName
		params := svcParams

		if svcPriority != nil {
			sp = *svcPriority
		}

		if tn == "" {
			tn = "."
		}

		if !strings.HasSuffix(tn, ".") {
			tn += "."
		}

		params = strings.ReplaceAll(params, `"`, "")
		if params == "" {
			params = "alpn=h2"
		}

		vl = fmt.Sprintf("%d %s %s", sp, tn, params)
	}

	return vl
}
