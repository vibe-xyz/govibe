package models

import (
	"fmt"
	"strconv"
)

// Table : select table
func (dba *Sqler) PgTable(table string) *Sqler {
	table = fmt.Sprintf("\"%s\"", table)
	dba.table = table
	return dba
}

func (dba *Sqler) PgSelect(args ...interface{}) string {
	if len(args) > 0 {
		dba.fields = args[0].(string)
	}
	return utils_RetStr(dba.PgBuildQuery())
}

func (dba *Sqler) PgBuildQuery() (string, error) {
	// Agg
	unionArr := []string{
		dba.count,
		dba.sum,
		dba.avg,
		dba.max,
		dba.min,
	}
	var union string
	for _, item := range unionArr {
		if item != "" {
			union = item
			break
		}
	}
	// distinct
	distinct := If(dba.distinct, "distinct ", "")
	// fields
	fields := If(dba.fields == "", "*", dba.fields).(string)
	// table
	table := dba.table
	// join
	parseJoin, err := dba.parseJoin()
	if err != nil {
		return "", err
	}
	join := parseJoin
	// where
	// beforeParseWhereData = dba.where
	parseWhere, err := dba.parseWhere(dba.where)
	if err != nil {
		return "", err
	}
	where := If(parseWhere == "", "", " WHERE "+parseWhere).(string)
	// group
	group := If(dba.group == "", "", " GROUP BY "+dba.group).(string)
	// having
	having := If(dba.having == "", "", " HAVING "+dba.having).(string)
	// order
	order := If(dba.order == "", "", " ORDER BY "+dba.order).(string)
	// limit
	limit := If(dba.limit == 0, "", " LIMIT "+strconv.Itoa(dba.limit))
	// offset
	offset := If(dba.offset == 0, "", " OFFSET "+strconv.Itoa(dba.offset))

	//sqlstr := "select " + fields + " from " + table + " " + where + " " + order + " " + limit + " " + offset
	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		distinct, If(union != "", union, fields), table, join, where, group, having, order, limit, offset)

	return sqlstr, nil
}
