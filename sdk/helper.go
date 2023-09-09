package sdk

import "strings"

func collectionTopic(slug string) string {
	var builder strings.Builder

	builder.WriteString("collection:")
	builder.WriteString(slug)

	return builder.String()
}
