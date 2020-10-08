package models

type Policy struct {
	User   string `json:"user" form:"user" query:"user"`
	Path   string `json:"path" form:"path" query:"path"`
	Method string `json:"method" form:"method" query:"method"`
}

type GroupPolicy struct {
	Member  string `json:"member" form:"member" query:"member"`
	Group string `json:"group" form:"group" query:"Group"`
}
