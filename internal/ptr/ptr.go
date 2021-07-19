package ptr

import "strings"

// StrRef ...
func StrRef(s string) *string {
	return &s
}

// StrRefDefaultNil ...
func StrRefDefaultNil(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}

	return &s
}

// StrSafeDeref ...
func StrSafeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
