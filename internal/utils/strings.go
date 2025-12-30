package utils

import (
	"strings"

	"github.com/google/uuid"
)

func FirstCharToLowerCase(str string) string {
	firstChar := str[:1]
	return strings.ToLower(firstChar) + str[1:]
}

func IdsToUUUIds(ids []string) []uuid.UUID {
	out := make([]uuid.UUID, len(ids))
	for i, v := range ids {
		out[i] = uuid.MustParse(v)
	}
	return out
}
