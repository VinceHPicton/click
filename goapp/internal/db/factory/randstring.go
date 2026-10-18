package factory

import (
	"database/sql"
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randNullString(n int) sql.NullString {
	return sql.NullString{
		String: randString(n),
		Valid:  true,
	}
}
