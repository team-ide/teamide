module db

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
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.4
	worker v0.0.0-00010101000000-000000000000
)
