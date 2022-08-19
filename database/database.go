package database

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	// mysql sqlite3 postgres and custom
	Dialect string `yaml:"dialect" json:"dialect"`
	// dsn
	Dsn string `yaml:"dsn" json:"dsn"`
	// MaxIdleConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxIdleConn int `yaml:"maxIdleConn" json:"maxIdleConn"`
	// MaxOpenConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConn int `yaml:"maxOpenConn" json:"maxOpenConn"`
	// MaxLifetime sets the maximum amount of time a connection may be reused.
	// If d <= 0, connections are not closed due to a connection's age.
	MaxLifetime time.Duration `yaml:"maxLifetime" json:"maxLifetime"`
	// MaxIdleTime sets the maximum amount of time a connection may be idle.
	// If d <= 0, connections are not closed due to a connection's idle time.
	MaxIdleTime time.Duration `yaml:"maxIdleTime" json:"maxIdleTime"`
	// EnableLog enabled log flag  use by user
	EnableLog bool `yaml:"enableLog" json:"enableLog"`
}

func New(c Config, config *gorm.Config, dialectorNews ...func(c Config) gorm.Dialector) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch c.Dialect {
	case "mysql":
		dialect = newMysql(&c)
	case "postgres":
		dialect = newPostgres(&c)
	case "sqlite3":
		dsn := c.Dsn
		if !strings.HasSuffix(dsn, ".db") {
			dsn += ".db"
		}
		if !isPathExist(dsn) {
			if err := os.MkdirAll(path.Dir(dsn), os.ModePerm); err != nil {
				return nil, fmt.Errorf("database mkdir (%s), %+v", dsn, err)
			}
			if _, err := os.Create(dsn); err != nil {
				return nil, fmt.Errorf("database create DB(%s), %+v", dsn, err)
			}
		}
		dialect = newSqlite3(dsn)
	case "custom":
		if len(dialectorNews) == 0 {
			panic("select option dialector should give a dialector new function")
		}
		dialectorNew := dialectorNews[0]
		dialect = dialectorNew(c)
	default:
		panic("please select database driver one of [mysql|postgres|sqlite3|custom], if use sqlite3, build tags with mysql|postgres|sqlite3!")
	}

	db, err := gorm.Open(dialect, config)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}
	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// SetDBLogger set db logger
func SetDBLogger(db *gorm.DB, l logger.Interface) {
	db.Logger = l
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// isPathExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isPathExist(paths string) bool {
	_, err := os.Stat(paths)
	return err == nil || os.IsExist(err)
}
