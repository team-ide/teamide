module base

go 1.15

replace (
	config => ../config
	server => ../server
)

require (
	config v0.0.0-00010101000000-000000000000
	github.com/Chain-Zhang/pinyin v0.1.3
	github.com/json-iterator/go v1.1.11
	github.com/satori/go.uuid v1.2.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	server v0.0.0-00010101000000-000000000000
)
