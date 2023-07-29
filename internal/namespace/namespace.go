package namespace

type Namespace struct {
	Name string `json:"name"`
}

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
}
