package dbwraper

import (
	"github.com/cwloo/gonet/utils/dbwraper/Gorm"
	"github.com/cwloo/gonet/utils/dbwraper/Mongo"
	"github.com/cwloo/gonet/utils/dbwraper/Redis"
	"github.com/cwloo/gonet/utils/dbwraper/Sql"
	"github.com/cwloo/gonet/utils/json"
)

var Wrap = &Wraper{}

type Wraper struct {
	Redis Redis.DB
	Mongo Mongo.DB
	Sql   Sql.DB
	Gorm  Gorm.DB
}

func InitRedis(RedisConf any) {
	json.Parse(json.Bytes(RedisConf), &Redis.Conf)
	Wrap.Redis.Init(Redis.Conf)
}

func InitMongo(MongoConf any) {
	json.Parse(json.Bytes(MongoConf), &Mongo.Conf)
	Wrap.Mongo.Init(Mongo.Conf)
}

func InitMySql(SqlConf any) {
	json.Parse(json.Bytes(SqlConf), &Sql.Conf)
	Wrap.Sql.Init(Sql.Conf)
}

func InitMyGorm(GormConf any) {
	json.Parse(json.Bytes(GormConf), &Gorm.Conf)
	Wrap.Gorm.Init(Gorm.Conf)
}

func Init(RedisConf, MongoConf, SqlConf, GormConf any) {
	json.Parse(json.Bytes(RedisConf), &Redis.Conf)
	json.Parse(json.Bytes(MongoConf), &Mongo.Conf)
	json.Parse(json.Bytes(SqlConf), &Sql.Conf)
	json.Parse(json.Bytes(GormConf), &Gorm.Conf)
	Wrap.Redis.Init(Redis.Conf)
	Wrap.Mongo.Init(Mongo.Conf)
	Wrap.Sql.Init(Sql.Conf)
	Wrap.Gorm.Init(Gorm.Conf)
}
