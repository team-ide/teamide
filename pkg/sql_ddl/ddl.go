package sql_ddl

func ToDatabaseDDL(database string, databaseType string) (sqls []string, err error) {

	if DatabaseIsMySql(databaseType) {
		sqls, err = ToDatabaseDDLForMySql(database)
	} else if DatabaseIsOracle(databaseType) {
		sqls, err = ToDatabaseDDLForOracle(database)
	} else if DatabaseIsShenTong(databaseType) {
		sqls, err = ToDatabaseDDLForShenTong(database)
	} else if DatabaseIsDaMeng(databaseType) {
		sqls, err = ToDatabaseDDLForDaMeng(database)
	} else if DatabaseIsKingbase(databaseType) {
		sqls, err = ToDatabaseDDLForKingBase(database)
	}

	return
}

func ToTableDDL(databaseType string, table TableDetailInfo) (sqls []string, err error) {
	if DatabaseIsMySql(databaseType) {
		sqls, err = ToTableDDLForMySql(table)
	} else if DatabaseIsOracle(databaseType) {
		sqls, err = ToTableDDLForOracle(table)
	} else if DatabaseIsShenTong(databaseType) {
		sqls, err = ToTableDDLForShenTong(table)
	} else if DatabaseIsDaMeng(databaseType) {
		sqls, err = ToTableDDLForDaMeng(table)
	} else if DatabaseIsKingbase(databaseType) {
		sqls, err = ToTableDDLForKingBase(table)
	}
	return
}
