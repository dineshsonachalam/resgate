package httpapi

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Resource holds a resource information to be sent to the client
type Resource struct {
	APIPath    string
	RID        string
	Model      interface{}
	Collection []interface{}
	Error      error
}

// MarshalJSON implements the json.Marshaler interface
func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		HREF       string        `json:"href"`
		Model      interface{}   `json:"model,omitempty"`
		Collection []interface{} `json:"collection,omitempty"`
		Error      error         `json:"error,omitempty"`
	}{
		HREF:       RIDToPath(r.RID, r.APIPath),
		Model:      r.Model,
		Collection: r.Collection,
		Error:      r.Error,
	})
}

// PathToRID parses a raw URL path and returns the resource ID.
// The prefix is the beginning of the path which is not part of the
// resource ID, and it should both start and end with /. Eg. "/api/"
func PathToRID(path, query, prefix string) string {
	if len(path) == len(prefix) || !strings.HasPrefix(path, prefix) {
		return ""
	}

	path = path[len(prefix):]

	// Dot separator not allowed in path
	if strings.ContainsRune(path, '.') {
		return ""
	}

	if path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part, err := url.PathUnescape(parts[i])
		if err != nil {
			return ""
		}
		parts[i] = part
	}

	rid := strings.Join(parts, ".")
	if query != "" {
		rid += "?" + query
	}

	return rid
}

// PathToRIDAction parses a raw URL path and returns the resource ID and action.
// The prefix is the beginning of the path which is not part of the
// resource ID, and it should both start and end with /. Eg. "/api/"
func PathToRIDAction(path, query, prefix string) (string, string) {
	if len(path) == len(prefix) || !strings.HasPrefix(path, prefix) {
		return "", ""
	}

	path = path[len(prefix):]

	// Dot separator not allowed in path
	if strings.ContainsRune(path, '.') {
		return "", ""
	}

	if path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", ""
	}

	for i := len(parts) - 1; i >= 0; i-- {
		part, err := url.PathUnescape(parts[i])
		if err != nil {
			return "", ""
		}
		parts[i] = part
	}

	rid := strings.Join(parts[:len(parts)-1], ".")
	if query != "" {
		rid += "?" + query
	}

	return rid, parts[len(parts)-1]
}

// RIDToPath converts a resource ID to a URL path string.
// The prefix is the part of the path that should be prepended
// to the resource ID path, and it should both start and end with /. Eg. "/api/".
func RIDToPath(rid, prefix string) string {
	return prefix + strings.Replace(url.PathEscape(rid), ".", "/", -1)
}
