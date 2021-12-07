module redis

go 1.15

replace (
	base => ../base
	config => ../config
	server => ../server
	worker => ../worker
)

require (
	base v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.3
	github.com/gomodule/redigo v1.8.5
	worker v0.0.0-00010101000000-000000000000
)
