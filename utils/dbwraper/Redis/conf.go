package Redis

type Cfg struct {
	Addr        []string `json:"addr" form:"addr"`
	MaxIdle     int      `json:"maxIdle" form:"maxIdle"`
	MaxActive   int      `json:"maxActive" form:"maxActive"`
	IdleTimeout int      `json:"idleTimeout" form:"idleTimeout"`
	Username    string   `json:"username" form:"username"`
	Password    string   `json:"password" form:"password"`
	Cluster     bool     `json:"cluster" form:"cluster"`
}

var (
	Conf Cfg
)
