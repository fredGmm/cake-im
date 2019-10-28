package config

import (
	"github.com/Terry-Mao/goim/pkg/time"
	"github.com/Unknwon/goconfig"
	"github.com/alecthomas/kingpin"
	"imChong/toolkit/path"
	"strconv"
)

var (
	dev     bool
	version = "0.0.1"
)

type Config struct {
	Host string
	Ip   string
	Port string

	TCP       *TCP
	Websocket *Websocket
	//Protocol  *Protocol
	Bucket *Bucket
	//RPCClient *RPCClient
	//RPCServer *RPCServer
	Whitelist *Whitelist
	Redis     *Redis
}

func Init() (c *Config, error error) {
	configPath, err := path.GetConfigPath("")
	if err != nil {
		return nil, err
	}
	parseFlag()
	c, e := parseConf(configPath)
	return c, e
}

func parseFlag() {
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version)
	//kingpin.Flag("configPath", "配置文件地址").Default("").StringVar(&configPath)
	kingpin.Flag("dev", "环境选择").Default("true").BoolVar(&dev)
	kingpin.Parse() // first parse conf
	return
}

func parseConf(confPath string) (c *Config, err error) {
	cfg, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	host, err2 := cfg.GetValue("local", "host")
	if err2 != nil {
		return nil, err2
	}
	ip, err3 := cfg.GetValue("local", "ip")
	if err3 != nil {
		return nil, err3
	}
	port, err4 := cfg.GetValue("local", "port")
	if err4 != nil {
		return nil, err4
	}

	redisHost, err5 := cfg.GetValue("redis", "host")
	if err5 != nil {
		return nil, err5
	}

	redisPortStr, err6 := cfg.GetValue("redis", "port")
	if err6 != nil {
		return nil, err6
	}
	redisPort, err9 := strconv.Atoi(redisPortStr)
	if err9 != nil {
		return nil, err9
	}

	redisPwd, err7 := cfg.GetValue("redis", "password")
	if err7 != nil {
		return nil, err7
	}
	redisDatabaseStr, err8 := cfg.GetValue("redis", "database")
	if err8 != nil {
		return nil, err8
	}
	redisDatabase, err10 := strconv.Atoi(redisDatabaseStr)
	if err10 != nil {
		return nil, err10
	}

	conf := &Config{
		Host: host,
		Ip:   ip,
		Port: port,
		TCP: &TCP{
			Bind:         []string{":8888"},
			SendBuf:      4096,
			RcvBuf:       4096,
			KeepAlive:    false,
			Reader:       32,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       32,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		},
		Websocket: &Websocket{
			Bind: []string{":3102"},
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			Room:          1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
		Redis: &Redis{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPwd,
			Database: redisDatabase,
		},
	}
	return conf, nil
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    time.Duration
	Timeout time.Duration
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string
	Addr              string
	Timeout           time.Duration
	IdleTimeout       time.Duration
	MaxLifeTime       time.Duration
	ForceCloseWait    time.Duration
	KeepAliveInterval time.Duration
	KeepAliveTimeout  time.Duration
}

type TCP struct {
	Bind         []string
	SendBuf      int
	RcvBuf       int
	KeepAlive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Websocket is websocket config.
type Websocket struct {
	Bind        []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int
	TimerSize        int
	SvrProto         int
	CliProto         int
	HandshakeTimeout time.Duration
}

// Bucket is bucket config.
type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

// Whitelist is white list config.
type Whitelist struct {
	Whitelist []int64
	WhiteLog  string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

func (r *Redis) GetAddr() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}
