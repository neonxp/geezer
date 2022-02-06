package geezer

type Values map[string][]string

func (v Values) Get(key string) string {
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v Values) Has(key string) bool {
	_, has := v[key]
	return has
}
