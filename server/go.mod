module teamide

go 1.15

replace (
	base => ./base
	cache => ./cache
	config => ./config
	db => ./db
	redis => ./redis
	service => ./service
	version => ./version
	web => ./web
	worker => ./worker
	zookeeper => ./zookeeper
)

require (
	base v0.0.0-00010101000000-000000000000
	cache v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	db v0.0.0-00010101000000-000000000000
	github.com/Chain-Zhang/pinyin v0.1.3 // indirect
	redis v0.0.0-00010101000000-000000000000
	service v0.0.0-00010101000000-000000000000
	version v0.0.0-00010101000000-000000000000
	web v0.0.0-00010101000000-000000000000
	worker v0.0.0-00010101000000-000000000000
	zookeeper v0.0.0-00010101000000-000000000000
)
