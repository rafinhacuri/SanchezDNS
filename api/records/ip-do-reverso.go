package records

import (
	"net"
	"slices"
	"strings"
)

func IPDoReverso(nomeReverso string) net.IP {
	nome := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(nomeReverso)), ".")

	if resto, ok := strings.CutSuffix(nome, ".in-addr.arpa"); ok {
		octetos := strings.Split(resto, ".")
		if len(octetos) != 4 {
			return nil
		}

		slices.Reverse(octetos)

		return net.ParseIP(strings.Join(octetos, "."))
	}

	resto, ok := strings.CutSuffix(nome, ".ip6.arpa")
	if !ok {
		return nil
	}

	nibbles := strings.Split(resto, ".")
	if len(nibbles) != 32 {
		return nil
	}

	slices.Reverse(nibbles)

	var texto strings.Builder

	for i, nibble := range nibbles {
		if i > 0 && i%4 == 0 {
			texto.WriteString(":")
		}

		texto.WriteString(nibble)
	}

	return net.ParseIP(texto.String())
}
