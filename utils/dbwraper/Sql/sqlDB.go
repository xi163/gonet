package Sql

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/cwloo/gonet/logs"
)

type DB struct {
	DB *sql.DB
}

func (s *DB) Init(conf Cfg) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logs.Fatalf(err.Error())
	}
	//user:passwd@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local
	cfg := mysql.Config{
		Net: "tcp",
		// Addr:                 "127.0.0.1:3306",
		Collation:            "utf8mb4",
		AllowNativePasswords: true,
		ParseTime:            true, Loc: loc}
	cfg.Addr = conf.Addr[0]
	cfg.User = conf.Username
	cfg.Passwd = conf.Password
	cfg.DBName = conf.Database
	cfg.Timeout = 5 * time.Second
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logs.Fatalf(err.Error())
	}
	db.SetMaxOpenConns(conf.MaxConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)
	db.Ping()
	s.DB = db
	logs.Debugf("ok")
}

// 插入数据库通用函数 不可以使用bool 可以使用 int int32 int64 string  where => "uid=1 and cid=2"
func (s *DB) UpdateOne(table string, data map[string]any, where string) (int64, error) {
	Id := int64(0)
	count := len(data)
	if count <= 0 {
		panic("error")
	}
	arrKeys := make([]string, 0, count)
	arrValues := make([]any, 0, count)
	for k, v := range data {
		k = k + "=?"
		arrKeys = append(arrKeys, k)
		tmp, ok := v.(string)
		if ok {
			//tmp = strings.Replace(tmp, "?", "？", -1)
			v = tmp
		}
		arrValues = append(arrValues, v)
	}
	sqlStr := fmt.Sprintf("update %s set %s", table, strings.Join(arrKeys, ","))
	if where == "" {
		logs.Fatalf("error")
	}
	// start := time.Now()
	sqlStr = sqlStr + " WHERE " + where
	stmt, err := s.DB.Prepare(sqlStr)
	if err != nil {
		return Id, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arrValues...)
	if err != nil {
		return Id, err
	}
	Id, err = res.RowsAffected()
	if err != nil {
		return Id, err
	}
	return Id, nil
}

// 插入数据库通用函数 不可以使用bool 可以使用 int int32 int64 string
func (s *DB) InsertOne(table string, data map[string]any) (int64, error) {
	Id := int64(0)
	count := len(data)
	if count <= 0 {
		logs.Fatalf("error")
	}
	arrKeys := make([]string, 0, count)
	arrVValues := make([]string, 0, count)
	arrValues := make([]any, 0, count)
	for k, v := range data {
		arrKeys = append(arrKeys, k)
		arrVValues = append(arrVValues, "?")
		tmp, ok := v.(string)
		if ok {
			// tmp = strings.Replace(tmp, "?", "'?", -1)
			v = tmp
		}
		arrValues = append(arrValues, v)
	}
	sqlStr := fmt.Sprintf("Insert into %s (%s) values(%s)", table, strings.Join(arrKeys, ","), strings.Join(arrVValues, ","))
	// start := time.Now()
	stmt, err := s.DB.Prepare(sqlStr)
	if err != nil {
		return Id, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arrValues...)
	if err != nil {
		return Id, err
	}
	Id, err = res.LastInsertId()
	if err != nil {
		return Id, err
	}
	return Id, nil
}

// 批量插入数据库通用函数 不可以使用bool 可以使用 int int32 int64 string
func (s *DB) InsertMore(table string, datas ...map[string]any) (int64, int64, error) {
	Id := int64(0)
	row := int64(0)
	sqlStr := ""
	arrValues := []any{}
	for i, data := range datas[0:] {
		count := len(data)
		if count <= 0 {
			panic("error")
		}

		arrKeys := make([]string, 0, count)
		arrVValues := make([]string, 0, count)
		// arrValues = make([]interface{}, 0, count)
		var keys []string
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := data[k]
			arrKeys = append(arrKeys, k)
			arrVValues = append(arrVValues, "?")
			tmp, ok := v.(string)
			if ok {
				// tmp = strings.Replace(tmp, "?", "'?", -1)
				v = tmp
			}
			arrValues = append(arrValues, v)
		}
		// for k, v := range data {
		// 	arrKeys = append(arrKeys, k)
		// 	arrVValues = append(arrVValues, "?")
		// 	// arrVValues = append(arrVValues, v)
		// 	tmp, ok := v.(string)
		// 	if ok {
		// 		// tmp = strings.Replace(tmp, "?", "'?", -1)
		// 		v = tmp
		// 	}
		// 	arrValues = append(arrValues, v)
		// }
		// logs.Debugf("keys=%+v", arrKeys)
		if i == 0 {
			sqlStr = fmt.Sprintf("Insert into %s (%s) values(%s)", table, strings.Join(arrKeys, ","), strings.Join(arrVValues, ","))
		} else {
			sqlStr += fmt.Sprintf(",(%s)", strings.Join(arrVValues, ","))
		}
	}
	// logs.Debugf("values=%+v", arrValues)
	// logs.Debugf("sqlStr=%s", sqlStr)
	// start := time.Now()
	stmt, err := s.DB.Prepare(sqlStr)
	if err != nil {
		return Id, 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arrValues...)
	if err != nil {
		logs.Errorf("%s %+v %+v", sqlStr, res, err)
		return Id, 0, err
	}
	Id, err = res.LastInsertId()
	row, err = res.RowsAffected()
	// logs.Debugf("Id=%+v row=%+v", Id, row)
	if err != nil {
		return Id, row, err
	}
	return Id, 0, nil
}

// 对参数做防注入过滤
func (s *DB) SafeExec(format string, args ...interface{}) (sql.Result, error) {
	ret, err, _ := s.SafeExecSql(format, args...)
	return ret, err
}

func (s *DB) SafeExecSql(format string, args ...interface{}) (sql.Result, error, string) {
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case string:
			args[i] = strings.Replace(args[i].(string), "'", "''", -1)
			args[i] = strings.Replace(args[i].(string), "\\", "\\\\", -1)
		default:
		}
	}
	sql := fmt.Sprintf(format, args...)
	// logs.Debugf(sql)
	start := time.Now()
	ret, err := s.DB.Exec(sql)
	tmpStr := sql
	ind := strings.Index(tmpStr, "where")
	if ind >= 0 {
		tmpStr = tmpStr[:ind]
	}
	if time.Now().Sub(start) > 5*time.Second { //5秒以上 告警
	}
	return ret, err, sql
}

// 对参数做防注入过滤
func (s *DB) SafeQuery(format string, args ...interface{}) (*sql.Rows, error) {
	ret, err, _ := s.SafeQuerySql(format, args...)
	return ret, err
}

func (s *DB) SafeQuerySql(format string, args ...interface{}) (*sql.Rows, error, string) {
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case string:
			args[i] = strings.Replace(args[i].(string), "'", "''", -1)
			args[i] = strings.Replace(args[i].(string), "\\", "\\\\", -1)
		default:
		}
	}
	sql := fmt.Sprintf(format, args...)
	// logs.Debugf(sql)
	start := time.Now()
	ret, err := s.DB.Query(sql)
	//get tableName
	tmpStr := strings.ToLower(sql)
	ind := strings.Index(tmpStr, "from ") + 5
	if ind >= 0 {
		tmpStr = tmpStr[ind:]
	}
	tmpStr = strings.TrimLeft(tmpStr, " ")
	ind = strings.Index(tmpStr, " ")
	if ind > 0 {
		tmpStr = tmpStr[:ind]
	}
	if time.Now().Sub(start) > 5*time.Second { //5秒以上 告警
	}
	return ret, err, sql
}

// 直接查询 对于 有%% 才可以用 有闲在用
func (s *DB) SafeQueryByString(sql string) (*sql.Rows, error) {
	start := time.Now()
	ret, err := s.DB.Query(sql)
	//get tableName
	tmpStr := strings.ToLower(sql)
	ind := strings.Index(tmpStr, "from ") + 5
	if ind >= 0 {
		tmpStr = tmpStr[ind:]
	}
	tmpStr = strings.TrimLeft(tmpStr, " ")
	ind = strings.Index(tmpStr, " ")
	if ind > 0 {
		tmpStr = tmpStr[:ind]
	}
	if time.Now().Sub(start) > 5*time.Second { //5秒以上 告警
	}
	return ret, err
}

func (s *DB) DelTable(table string, format string, args ...any) (int64, error) {
	row := int64(0)
	DelSql := fmt.Sprintf("DELETE FROM %s WHERE ", table)
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			switch args[i].(type) {
			case string:
				{
					args[i] = strings.Replace(args[i].(string), "'", "''", -1)
					args[i] = strings.Replace(args[i].(string), "\\", "\\\\", -1)
				}
			default:
			}
		}

		DelSql = fmt.Sprintf(DelSql+" "+format, args...)
	} else {
		DelSql = DelSql + " " + format
	}
	start := time.Now()
	res, err := s.DB.Exec(DelSql)
	if err != nil {
		return row, err
	}
	row, err = res.RowsAffected()
	if err != nil {
		return row, err
	}
	if time.Now().Sub(start) > 5*time.Second { //5秒以上 告警
	}
	return row, nil
}
