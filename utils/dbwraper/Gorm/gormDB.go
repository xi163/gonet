package Gorm

import (
	"fmt"
	"time"

	"github.com/cwloo/gonet/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	DB *gorm.DB
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (s *DB) Init(cfg Cfg) {
	//user:passwd@tcp(localhost:9910)/db?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Addr[0],
		cfg.Database)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			Writer{},
			logger.Config{
				SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.LogLevel(cfg.LogLevel),                       // Log level
				IgnoreRecordNotFoundError: true,                                                // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                                                // Disable color
			},
		),
	})
	// gormDB, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN: dsn, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
	// 	// DefaultStringSize:        256,  // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
	// 	DisableDatetimePrecision: true, // disable datetime precision support, which not supported before MySQL 5.6
	// 	// DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
	// 	DontSupportRenameIndex:    true,  // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
	// 	DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
	// 	SkipInitializeWithVersion: false, // smart configure based on used version
	// }), &gorm.Config{
	// 	Logger: logger.New(
	// 		Writer{},
	// 		logger.Config{
	// 			SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond, // Slow SQL threshold
	// 			LogLevel:                  logger.LogLevel(cfg.LogLevel),                       // Log level
	// 			IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
	// 			Colorful:                  true,                                                                // Disable color
	// 		},
	// 	),
	// })
	if err != nil {
		logs.Fatalf(err.Error() + ":" + dsn)
	}
	// gormDB.Set("gorm:table_options", "CHARSET=utf8")
	// gormDB.Set("gorm:table_options", "collation=utf8_unicode_ci")
	db, err := gormDB.DB()
	if err != nil {
		logs.Fatalf(err.Error())
	}
	db.SetMaxOpenConns(cfg.MaxConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	db.Ping()
	s.DB = gormDB
	logs.Debugf("ok")
}
