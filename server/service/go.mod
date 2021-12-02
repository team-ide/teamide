module service

go 1.15

replace (
	base => ../base
	config => ../config
	db => ../db
	redis => ../redis
	worker => ../worker
	zookeeper => ../zookeeper
)

require (
	base v0.0.0-00010101000000-000000000000
	db v0.0.0-00010101000000-000000000000
	redis v0.0.0-00010101000000-000000000000
)
