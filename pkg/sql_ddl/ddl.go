package sql_ddl

import (
	"errors"
	"teamide/pkg/db"
)

func ToDatabaseDDL(database string, databaseType string) (sqls []string, err error) {

	switch db.GetDatabaseType(databaseType) {
	case db.DatabaseTypeMySql:
		sqls, err = ToDatabaseDDLForMySql(database)
		break
	case db.DatabaseTypeSqlite:
		sqls, err = ToDatabaseDDLForSqlite(database)
		break
	case db.DatabaseTypeOracle:
		sqls, err = ToDatabaseDDLForOracle(database)
		break
	case db.DatabaseTypeShenTong:
		sqls, err = ToDatabaseDDLForShenTong(database)
		break
	case db.DatabaseTypeDM:
		sqls, err = ToDatabaseDDLForDaMeng(database)
		break
	case db.DatabaseTypeKingBase:
		sqls, err = ToDatabaseDDLForKingBase(database)
		break
	case nil:
		err = errors.New("数据库类型[" + databaseType + "]暂不支持")
		break

	}

	return
}

func ToTableDDL(databaseType string, table *TableDetailInfo) (sqls []string, err error) {

	switch db.GetDatabaseType(databaseType) {
	case db.DatabaseTypeMySql:
		sqls, err = ToTableDDLForMySql(table)
		break
	case db.DatabaseTypeSqlite:
		sqls, err = ToTableDDLForSqlite(table)
		break
	case db.DatabaseTypeOracle:
		sqls, err = ToTableDDLForOracle(table)
		break
	case db.DatabaseTypeShenTong:
		sqls, err = ToTableDDLForShenTong(table)
		break
	case db.DatabaseTypeDM:
		sqls, err = ToTableDDLForDaMeng(table)
		break
	case db.DatabaseTypeKingBase:
		sqls, err = ToTableDDLForKingBase(table)
		break
	case nil:
		err = errors.New("数据库类型[" + databaseType + "]暂不支持")
		break
	}
	return
}
