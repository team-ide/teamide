module db

go 1.15

replace (
	base => ../base
	config => ../config
	worker => ../worker
)

require (
	base v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	github.com/Chain-Zhang/pinyin v0.1.3 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.4
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	worker v0.0.0-00010101000000-000000000000
)
