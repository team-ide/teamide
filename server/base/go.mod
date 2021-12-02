module base

go 1.15

replace (
	config => ./config
)

require (
	github.com/json-iterator/go v1.1.11
	github.com/satori/go.uuid v1.2.0
	config v0.0.0-00010101000000-000000000000
)
