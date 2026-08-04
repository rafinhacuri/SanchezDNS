package statistics

import "time"

func IniciadoEm(segundos int) time.Time {
	if segundos <= 0 {
		return time.Now().UTC()
	}

	return time.Now().UTC().Add(-time.Duration(segundos) * time.Second)
}
