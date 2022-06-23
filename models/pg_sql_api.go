package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func PgGet(sqlstr string, dest interface{}) (err error) {
	conn := GPostgresPool.GetConn()
	defer GPostgresPool.UnGetConn(conn)
	err = conn.Get(dest, sqlstr)
	if nil != err && ErrNoRows != err {
		err = NewError("%v sql[%s]", err, sqlstr)
	}
	return
}

func NewPgSqlTxConn() (conn *sqlx.DB, tx *sql.Tx, err error) {
	conn = GPostgresPool.GetConn()
	tx, err = conn.Begin()
	return
}

func ClosePgSqlTxConn(conn *sqlx.DB, tx *sql.Tx, inoutErr *error) {
	GPostgresPool.UnGetConn(conn)

	SqlTxProc(tx, inoutErr)
}

func PgSelect(sqlstr string, dest interface{}) (err error) {
	conn := GPostgresPool.GetConn()
	defer GPostgresPool.UnGetConn(conn)

	err = conn.Select(dest, sqlstr)
	if nil != err && ErrNoRows != err {
		err = NewError("%v sql[%s]", err, sqlstr)
	}
	return
}
