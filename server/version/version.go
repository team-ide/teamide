package version

import (
	"base"
	"db"
)

type VersionInfo struct {
	Version string        `json:"version"`
	Sqls    []db.SqlParam `json:"sqls"`
}

type VersionDetail struct {
	Detail string `column:"detail"`
}

func appendVersion() {

}

func InstallVersion(versionInfo *VersionInfo) {
	if versionInfo == nil {
		return
	}
	var err error
	var res []interface{}
	res, err = db.DBService.Query(db.SqlParam{
		Sql:    "SELECT detail FROM  " + db.TABLE_INSTALL + " WHERE version=? ",
		Params: []interface{}{versionInfo.Version},
	}, func() interface{} {
		return &VersionDetail{}
	})
	if err != nil {
		panic(err)
	}
	var installedVersion *VersionInfo
	if len(res) > 0 {
		VersionDetail := res[0].(*VersionDetail)
		if VersionDetail.Detail != "" {
			err = base.ToBean([]byte(VersionDetail.Detail), &installedVersion)
			if err != nil {
				panic(err)
			}
		}
	}
	UpgradeVersion(installedVersion, versionInfo)
}

func UpgradeVersion(installedVersion *VersionInfo, versionInfo *VersionInfo) {
	var needInstallSqls []db.SqlParam
	var needSaveVersion = VersionInfo{
		Version: versionInfo.Version,
		Sqls:    []db.SqlParam{},
	}
	if installedVersion != nil {
		for index, sql := range versionInfo.Sqls {
			if index >= len(installedVersion.Sqls) {
				needInstallSqls = append(needInstallSqls, sql)
			} else {
				needSaveVersion.Sqls = append(needSaveVersion.Sqls, sql)
			}
		}
	} else {
		needInstallSqls = versionInfo.Sqls
	}
	if len(needInstallSqls) == 0 {
		return
	}

	var sqlErr error
	for _, sqlParam := range needInstallSqls {
		if sqlParam.Sql != "" {
			// println("upgrade version sql:", sqlParam.Sql)
			_, sqlErr = db.DBService.Exec(sqlParam)
			if sqlErr != nil {
				println(sqlErr)
				break
			}
		}
		needSaveVersion.Sqls = append(needSaveVersion.Sqls, sqlParam)
	}
	var detail = base.ToJSON(needSaveVersion)

	var dbErr error
	var dbSqlParam db.SqlParam
	if installedVersion != nil {
		dbSqlParam = db.SqlParam{
			Sql:    "UPDATE " + db.TABLE_INSTALL + " SET detail=?, updateTime=? WHERE version=? ",
			Params: []interface{}{detail, base.Now(), versionInfo.Version},
		}
	} else {
		dbSqlParam = db.SqlParam{
			Sql:    "INSERT INTO  " + db.TABLE_INSTALL + " (version, detail, createTime) VALUES (?, ?, ?) ",
			Params: []interface{}{versionInfo.Version, detail, base.Now()},
		}
	}
	_, dbErr = db.DBService.Exec(dbSqlParam)
	if dbErr != nil {
		panic(dbErr)
	}
	if sqlErr != nil {
		panic("version sql execute error")
	}
}
