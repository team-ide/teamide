module worker

go 1.15

replace (
	base => ../base
	config => ../config
	server => ../server
)

require (
	base v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	server v0.0.0-00010101000000-000000000000
)
