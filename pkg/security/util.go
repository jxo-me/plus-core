package security

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func BuildSigningString(method, path string, query url.Values, body []byte, timestamp int64) string {
	sortedQuery := buildSortedQueryString(query)
	return strings.Join([]string{
		strings.ToUpper(method),
		path,
		sortedQuery,
		string(body),
		FormatTimestamp(timestamp),
	}, "\n")
}

func FormatTimestamp(ts int64) string {
	return strconv.FormatInt(ts, 10)
}

func ParseTimestamp(ts string) (int64, error) {
	return strconv.ParseInt(ts, 10, 64)
}

func IsTimestampValid(ts int64) bool {
	now := time.Now().Unix()
	return ts >= now-300 && ts <= now+300
}

func buildSortedQueryString(params url.Values) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, k := range keys {
		for _, v := range params[k] {
			if i > 0 || b.Len() > 0 {
				b.WriteByte('&')
			}
			b.WriteString(url.QueryEscape(k))
			b.WriteByte('=')
			b.WriteString(url.QueryEscape(v))
		}
	}
	return b.String()
}
