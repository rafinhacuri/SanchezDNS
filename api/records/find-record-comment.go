package records

import "time"

func findCommentForRecordContent(comments []Comment, index, totalRecords int) (string, string, int64) {
	if len(comments) == totalRecords && index < len(comments) {
		c := comments[index]

		mod := c.ModifiedAt
		if mod == 0 {
			mod = time.Now().Unix()
		}

		return c.Content, c.Account, mod
	}

	if len(comments) == 1 {
		c := comments[0]

		mod := c.ModifiedAt
		if mod == 0 {
			mod = time.Now().Unix()
		}

		return c.Content, c.Account, mod
	}

	return "", "", time.Now().Unix()
}
