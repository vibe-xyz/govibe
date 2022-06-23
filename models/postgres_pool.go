package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresPool struct {
	db  *sql.DB
	dbx *sqlx.DB
}

type PostgresCfg struct {
	Host      string `toml:"host,omitzero"`
	Port      int64  `toml:"port,omitzero"`
	User      string `toml:"user,omitzero"`
	Pass      string `toml:"pass,omitzero"`
	EnableSsl bool   `toml:"ssl,omitzero"`
}

var (
	GPostgresPool *PostgresPool
)

type PostgresConnConfig struct {
	DBName         string `json:"dbname"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	Port           int64  `json:"port"`
	ConnectTimeout int64  `json:"connect_timeout"`
	EnableSsl      bool   `json:"ssl"`
}

func InitPostgresPool(Host string, User string, Pass string, Port int64, Database string) (pool *PostgresPool, err error) {
	postgres_cfg := PostgresConnConfig{
		User:      User,
		Password:  Pass,
		Host:      Host,
		Port:      Port,
		DBName:    Database,
		EnableSsl: false,
	}
	d := postgres_cfg.FormatDSN()
	db, err := sql.Open("postgres", d)
	if err != nil {
		return nil, err
	}
	db.Ping()

	dbx := sqlx.NewDb(db, "postgres")
	dbx.MapperFunc(LowerCaseWithUnderscores)

	pool = &PostgresPool{db, dbx}
	if GPostgresPool == nil {
		GPostgresPool = pool
	}

	return pool, nil
}

func (c *PostgresConnConfig) FormatDSN() string {
	dsn := ""
	if c.DBName != "" {
		dsn = fmt.Sprintf("%s dbname=%s ", dsn, c.DBName)
	}
	if c.User != "" {
		dsn = fmt.Sprintf("%s user=%s ", dsn, c.User)
	}
	if c.Host != "" {
		dsn = fmt.Sprintf("%s host=%s ", dsn, c.Host)
	} else {
		dsn = fmt.Sprintf("%s host=%s ", dsn, "localhost")
	}
	if c.Password != "" {
		dsn = fmt.Sprintf("%s password=%s ", dsn, c.Password)
	}
	if c.Port != 0 {
		dsn = fmt.Sprintf("%s port=%d ", dsn, c.Port)
	}
	if c.EnableSsl {
		dsn = fmt.Sprintf("%s sslmode=%s", dsn, "enable")
	} else {
		dsn = fmt.Sprintf("%s sslmode=%s", dsn, "disable")
	}
	return dsn
}

func (p *PostgresPool) GetConn() *sqlx.DB {
	return p.dbx
}

func (p *PostgresPool) UnGetConn(db interface{}) {
}
