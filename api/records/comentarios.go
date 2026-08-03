package records

func Comentarios(rr *RRSet) map[string]string {
	porConteudo := make(map[string]string)
	if rr == nil {
		return porConteudo
	}

	for i, rec := range rr.Records {
		if i < len(rr.Comments) {
			porConteudo[rec.Content] = rr.Comments[i].Content
		}
	}

	return porConteudo
}
