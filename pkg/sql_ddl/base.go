package sql_ddl

import (
	"regexp"
	"strings"
)

type TableDetailInfo struct {
	Name    string            `json:"name"`
	Comment string            `json:"comment"`
	Columns []TableColumnInfo `json:"columns"`
	Indexs  []TableIndexInfo  `json:"indexs"`
}

type TableColumnInfo struct {
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Type       string `json:"type"`
	Length     int    `json:"length"`
	Decimal    int    `json:"decimal"`
	PrimaryKey bool   `json:"primaryKey"`
	NotNull    bool   `json:"notNull"`
	Default    string `json:"default"`
}

type TableIndexInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Columns string `json:"columns"`
	Comment string `json:"comment"`
}

func DatabaseIsMySql(databaseType string) bool {
	return strings.ToLower(databaseType) == "mysql"
}
func DatabaseIsOracle(databaseType string) bool {
	return strings.ToLower(databaseType) == "oracle"
}
func DatabaseIsShenTong(databaseType string) bool {
	return strings.ToLower(databaseType) == "shentong"
}
func DatabaseIsDaMeng(databaseType string) bool {
	return strings.ToLower(databaseType) == "dm" || strings.ToLower(databaseType) == "dameng"
}
func DatabaseIsKingbase(databaseType string) bool {
	return strings.ToLower(databaseType) == "kingbase"
}

func formatSql(sql string, data map[string]string) (foramtSql string, err error) {
	var re *regexp.Regexp
	re, err = regexp.Compile(`\[(.+?)\]`)
	if err != nil {
		return
	}
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	var sql_ string
	var formatValueSql_ string
	var find bool = true
	for _, indexs := range indexsList {
		sql_ = sql[lastIndex:indexs[0]]
		formatValueSql_, find = formatValueSql(sql_, data)
		if find {
			foramtSql += formatValueSql_
		}

		lastIndex = indexs[1]

		sql_ = sql[indexs[0]+1 : indexs[1]-1]

		if !strings.Contains(sql_, `{`) {
			if data[strings.TrimSpace(sql_)] != "" {
				foramtSql += sql_
			}
		} else {
			formatValueSql_, find = formatValueSql(sql_, data)
			if find {
				foramtSql += formatValueSql_
			}
		}
	}
	sql_ = sql[lastIndex:]
	formatValueSql_, find = formatValueSql(sql_, data)
	if find {
		foramtSql += formatValueSql_
	}
	return
}

func formatValueSql(sql string, data map[string]string) (res string, find bool) {
	var re *regexp.Regexp
	re, _ = regexp.Compile(`{(.+?)}`)
	find = true
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	for _, indexs := range indexsList {
		res += sql[lastIndex:indexs[0]]

		lastIndex = indexs[1]

		key := sql[indexs[0]+1 : indexs[1]-1]
		value := data[key]
		if value == "" {
			find = false
			return
		}
		res += value
	}
	res += sql[lastIndex:]
	return
}
