package Sql

type Cfg struct {
	Addr          []string `json:"addr" form:"addr"`
	Username      string   `json:"username" form:"username"`
	Password      string   `json:"password" form:"password"`
	Database      string   `json:"database" form:"database"`
	Tablename     string   `json:"tablename" form:"tablename"`
	MaxConn       int      `json:"maxConn" form:"maxConn"`
	MaxIdleConn   int      `json:"maxIdleConn" form:"maxIdleConn"`
	MaxLifeTime   int      `json:"maxLifeTime" form:"maxLifeTime"`
	LogLevel      int      `json:"logLevel" form:"logLevel"`
	SlowThreshold int      `json:"slowThreshold" form:"slowThreshold"`
}

var (
	Conf Cfg
)
