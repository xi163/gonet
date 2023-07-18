package Mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cwloo/gonet/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type DB struct {
	DB *mongo.Client
}

func (s *DB) Init(conf Cfg) {
	if conf.Url != "" {
		c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.Url))
		if err != nil {
			logs.Fatalf(err.Error() + ":" + conf.Url)
		}
		s.DB = c
	} else {
		c, err := mongo.Connect(context.TODO(), options.Client().
			SetAuth(options.Credential{
				// AuthMechanism: "SCRAM-SHA-1",
				AuthSource: conf.Source,
				Username:   conf.Username,
				Password:   conf.Password,
			}).
			SetConnectTimeout(5*time.Second).
			SetSocketTimeout(5*time.Second).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
			SetHosts(conf.Addr).
			SetMaxPoolSize(uint64(conf.MaxPoolSize)).
			// SetMinPoolSize(uint64(conf.MinPoolSize)).
			SetReadPreference(readpref.Primary()).
			// SetReplicaSet("replicaSet").
			SetRetryWrites(true).SetDirect(conf.Direct).
			SetTimeout(time.Duration(conf.Timeout)*time.Second).
			SetServerSelectionTimeout(5*time.Second))
		if err != nil {
			logs.Fatalf(err.Error())
		}
		s.DB = c
	}
	if err := s.DB.Ping(context.TODO(), readpref.Primary()); err != nil {
		logs.Fatalf(err.Error())
	}
	logs.Debugf("ok")
}

func (s *DB) CreateIndex(dbname, tblname string, unique bool, keys ...string) error {
	db := s.DB.Database(dbname).Collection(tblname)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}
	// 复合索引
	for _, key := range keys {
		switch strings.HasPrefix(key, "-") {
		case true:
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		default:
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	// 创建索引
	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(unique),
		},
		opts,
	)
	if result == "" || err != nil {
		logs.Fatalf(err.Error())
	}
	return nil
}

func (s *DB) FindOne(dbname, tblname string, cols bson.M, where bson.M) (bson.M, error) {
	coll := s.DB.Database(dbname).Collection(tblname)
	cols["_id"] = 0
	opts := options.FindOne().SetProjection(cols)
	model := bson.M{}
	if err := coll.FindOne(context.TODO(), where, opts).Decode(&model); err != nil {
		// logs.Errorf(err.Error())
		return nil, err
	}
	return model, nil
}

func (s *DB) UpdateOne(dbname, tblname string, update, where bson.M) (*mongo.UpdateResult, error) {
	coll := s.DB.Database(dbname).Collection(tblname)
	opts := options.Update()
	result, err := coll.UpdateOne(context.TODO(), where, update, opts)
	if err != nil {
		logs.Errorf(err.Error())
		return result, err
	}
	fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	fmt.Println("UpdateOne() result ModifiedCount:", result.ModifiedCount)
	fmt.Println("UpdateOne() result UpsertedCount:", result.UpsertedCount)
	fmt.Println("UpdateOne() result UpsertedID:", result.UpsertedID)
	return result, err
}

func (s *DB) InsertOne(dbname, tblname string, doc any) (*mongo.InsertOneResult, error) {
	coll := s.DB.Database(dbname).Collection(tblname)
	opts := options.InsertOne()
	result, err := coll.InsertOne(context.TODO(), doc, opts)
	if err != nil {
		logs.Errorf(err.Error())
		return result, err
	}
	return result, err
}

func (s *DB) FindOneAndUpdate(dbname, tblname string, cols bson.M, update bson.M, where bson.M) (bson.M, error) {
	coll := s.DB.Database(dbname).Collection(tblname)
	cols["_id"] = 0
	opts := options.FindOneAndUpdate().SetProjection(cols)
	model := bson.M{}
	if err := coll.FindOneAndUpdate(context.TODO(), where, update, opts).Decode(&model); err != nil {
		logs.Errorf(err.Error())
		return nil, err
	}
	return model, nil
}

// https://haicoder.net/mongodb/mongodb-aggregate-max.html
func (s *DB) FindMinMaxAggregate(dbname, tblname string, name string) (int64, error) {
	opt := "$min"
	switch strings.HasPrefix(name, "-") {
	case true:
		opt = "$max"
		name = strings.TrimLeft(name, "-")
	}
	coll := s.DB.Database(dbname).Collection(tblname)
	opts := options.Aggregate()
	where := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: ""},
			{Key: name, Value: bson.D{
				{Key: opt, Value: "$" + name}}},
		}}}
	project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: name, Value: 1},
		}},
	}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{where, project}, opts)
	if err != nil {
		logs.Errorf(err.Error())
		return 0, err
	}
	model := []bson.M{}
	err = cursor.All(context.TODO(), &model)
	cursor.Close(context.TODO())
	if err != nil {
		logs.Errorf(err.Error())
		return 0, err
	}
	return model[0][name].(int64), err
}

func (s *DB) FindMinMax(dbname, tblname string, cols, sort bson.M) (bson.M, error) {
	coll := s.DB.Database(dbname).Collection(tblname)
	cols["_id"] = 0
	opts := options.Find().SetProjection(cols).SetSort(sort).SetSkip(0).SetLimit(1)
	opts.SetMaxTime(5 * time.Second)
	cursor, err := coll.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		logs.Errorf(err.Error())
		return nil, err
	}
	model := []bson.M{}
	err = cursor.All(context.TODO(), &model)
	cursor.Close(context.TODO())
	if err != nil {
		logs.Errorf(err.Error())
		return nil, err
	}
	if len(model) > 0 {
		return model[0], err
	}
	return nil, errors.New("emtpy")
}
