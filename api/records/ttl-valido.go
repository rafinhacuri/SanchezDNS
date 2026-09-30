package records

const TTLMinimo = 60

func TTLValido(ttl int) bool {
	return ttl == 0 || ttl >= TTLMinimo
}
