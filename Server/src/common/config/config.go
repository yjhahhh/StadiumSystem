package config

type Config struct {
	Addr         string // ip:port
	Mod          string // 运行模型 开发模式为dev
	Secret       string // 加密密钥
	LogPath      string // 日志目录
	LogLevel     string // panic, fatal, error, warn, info, debug, trace
	AllowOrigins []string // 允许访问的域名配置
	MySQLConf    MySQLConfig
	RedisConf    RedisConfig
}

type MySQLConfig struct {
	UserName  string // 用户名
	Password  string // 密码
	Protocol  string // 协议 tcp/udp
	Address   string
	DBName    string // 数据库名
	Charset   string // 字符集
}

type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

