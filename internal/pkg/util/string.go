package util

func trimString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func TrimLog(s string) string {
	return trimString(s, 100) + "..."
}
