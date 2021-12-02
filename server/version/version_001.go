package version

import "db"

func getVersion_001() (versionInfo *VersionInfo) {

	versionInfo = &VersionInfo{}

	versionInfo.Version = "0.0.1"
	sqls := []db.SqlParam{}

	sqls = append(sqls, db.SqlParam{
		Sql: `
		CREATE TABLE TM_ID (
			type int(2) NOT NULL COMMENT 'ID类型',
			id bigint(20) NOT NULL COMMENT 'ID',
			PRIMARY KEY (type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ID';
		`,
		Params: []interface{}{},
	})

	versionInfo.Sqls = sqls

	return
}
