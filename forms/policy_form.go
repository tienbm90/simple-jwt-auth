package forms

type Policy struct {
	Subject string `json:"subject" forms:"Subject" query:"Subject"`
	Object  string `json:"object" forms:"Object" query:"Object"`
	Action  string `json:"act" forms:"act" query:"act"`
}

func (f Policy) HasSubject() bool {
	return f.Subject != "" && len(f.Subject) <= 255
}

func (f Policy) HasObject() bool {
	return f.Object != "" && len(f.Object) <= 255
}

func (f Policy) HasAction() bool {
	return f.Action != "" && len(f.Action) <= 255
}

type GroupPolicy struct {
	Member string `json:"member" forms:"member" query:"member"`
	Group  string `json:"group" forms:"group" query:"Group"`
}

func (f GroupPolicy) HasMemeber() bool {
	return f.Member != "" && len(f.Member) <= 255
}

func (f GroupPolicy) HasGroup() bool {
	return f.Group != "" && len(f.Group) <= 255
}
