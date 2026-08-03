package statistics

import "strconv"

func Inteiro(valores map[string]string, chave string) int {
	valor, ok := valores[chave]
	if !ok {
		return 0
	}

	inteiro, err := strconv.Atoi(valor)
	if err == nil {
		return inteiro
	}

	numero, err := strconv.ParseFloat(valor, 64)
	if err == nil {
		return int(numero)
	}

	return 0
}
