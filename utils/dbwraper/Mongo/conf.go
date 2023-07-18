package Mongo

type Cfg struct {
	Url         string   `json:"url" form:"url"`
	Addr        []string `json:"addr" form:"addr"`
	Direct      bool     `json:"direct" form:"direct"`
	Timeout     int      `json:"timeout" form:"timeout"`
	Database    string   `json:"database" form:"database"`
	Source      string   `json:"source" form:"source"`
	Username    string   `json:"username" form:"username"`
	Password    string   `json:"password" form:"password"`
	MaxPoolSize int      `json:"maxPoolSize" form:"maxPoolSize"`
	MinPoolSize int      `json:"minPoolSize" form:"minPoolSize"`
}

var (
	Conf Cfg
)
