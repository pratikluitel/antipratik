package store

import "strings"

// ── Helpers ───────────────────────────────────────────────────────────────────

func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

func stringsToAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

func extractIDs(rows []baseRow) []string {
	ids := make([]string, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	return ids
}

func groupByType(rows []baseRow) map[string][]string {
	m := make(map[string][]string)
	for _, r := range rows {
		m[r.Type] = append(m[r.Type], r.ID)
	}
	return m
}

func coalesceStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
