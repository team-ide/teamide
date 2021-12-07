module web

go 1.15

replace (
	base => ../base
	config => ../config
	server => ../server
	service => ../service
)

require (
	base v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
)
