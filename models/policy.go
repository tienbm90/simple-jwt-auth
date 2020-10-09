package models

type Policy struct {
	User   string `json:"user" forms:"user" query:"user"`
	Path   string `json:"path" forms:"path" query:"path"`
	Method string `json:"method" forms:"method" query:"method"`
}

type GroupPolicy struct {
	Member  string `json:"member" forms:"member" query:"member"`
	Group string `json:"group" forms:"group" query:"Group"`
}
