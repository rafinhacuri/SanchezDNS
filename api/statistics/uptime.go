package statistics

import (
	"fmt"
	"time"
)

func Uptime(segundos int) string {
	if segundos <= 0 {
		return "0s"
	}

	duracao := time.Duration(segundos) * time.Second
	dias := int(duracao.Hours()) / 24
	horas := int(duracao.Hours()) % 24
	minutos := int(duracao.Minutes()) % 60

	if dias > 0 {
		return fmt.Sprintf("%dd %dh %dm", dias, horas, minutos)
	}

	if horas > 0 {
		return fmt.Sprintf("%dh %dm", horas, minutos)
	}

	return fmt.Sprintf("%dm", minutos)
}
