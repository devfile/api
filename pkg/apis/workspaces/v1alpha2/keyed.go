package v1alpha2

type Keyed interface {
	Key() string
}

type TopLevelLists map[string][]string

type TopLevelListContainer interface {
	GetToplevelLists() TopLevelLists
}
