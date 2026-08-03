package records

import "time"

func Comentario(comentarios []Comment, index, total int) (string, string, int64) {
	var escolhido Comment

	switch {
	case len(comentarios) == total && index < len(comentarios):
		escolhido = comentarios[index]
	case len(comentarios) == 1:
		escolhido = comentarios[0]
	default:
		return "", "", time.Now().Unix()
	}

	if escolhido.ModifiedAt == 0 {
		return escolhido.Content, escolhido.Account, time.Now().Unix()
	}

	return escolhido.Content, escolhido.Account, escolhido.ModifiedAt
}
