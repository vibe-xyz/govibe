package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func SqlGet(sqlstr string, dest interface{}) (err error) {
	conn := GMysqlPool.GetConn()
	defer GMysqlPool.UnGetConn(conn)
	err = conn.Get(dest, sqlstr)
	if nil != err && ErrNoRows != err {
		err = NewError("%v sql[%s]", err, sqlstr)
	}
	return
}

func NewSqlTxConn() (conn *sqlx.DB, tx *sql.Tx, err error) {
	conn = GMysqlPool.GetConn()
	tx, err = conn.Begin()
	return
}

func CloseSqlTxConn(conn *sqlx.DB, tx *sql.Tx, inoutErr *error) {
	GMysqlPool.UnGetConn(conn)

	SqlTxProc(tx, inoutErr)
}

func SqlTxProc(sqltx *sql.Tx, inoutErr *error) {
	if *inoutErr != nil {
		err := sqltx.Rollback()
		if err != nil {
			*inoutErr = NewError("sqltx.Rollback err %v *perr %v", err, *inoutErr)
		}
		return
	}

	if err := sqltx.Commit(); err != nil {
		*inoutErr = NewError("sqltx.Commit err %v", err)
	}
}

func SqlTxExec(sqltx *sql.Tx, sqlstr string) (err error) {
	_, err = sqltx.Exec(sqlstr)
	if nil != err {
		err = NewError("%v sql[%s]", err, sqlstr)
		return
	}
	return
}
