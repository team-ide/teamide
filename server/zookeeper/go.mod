module zookeeper

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
	github.com/go-zookeeper/zk v1.0.2
	worker v0.0.0-00010101000000-000000000000
)
