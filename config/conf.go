package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"github.com/alecthomas/kingpin"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	version = "0.0.1"
	defaultConfigPath string
)

func Test()  {
	log.Print("config test")
}

type Config struct {
	Websocket *Websocket
	TCP *TCP
	Redis *Redis
	//Grpc *Grpc
	Xrpc *Xrpc
	Exit chan struct{}
}

type Websocket struct {
	Bind []string
	TLSOpen     bool
	TLSBind     []string
	CertFile    string
	PrivateFile string
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

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}
func (r *Redis) GetAddr() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}


type Grpc struct {
	ServerId string
	ServerAddr string
}

type Xrpc struct {
	Addr string
}

func parseFlag()  {
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version)
	//kingpin.Flag("dev", "环境选择").Default("true").BoolVar(&dev)
	kingpin.Flag("configPath", "配置文件路径").Short('c').Default("E:\\code\\golang\\module\\cake-im\\config\\config.ini").StringVar(&defaultConfigPath)
	//kingpin.Flag("configPath", "配置文件路径").Short('c').Default("/mnt/hgfs/code/golang/module/cake-im/config/config.ini").StringVar(&defaultConfigPath)
	kingpin.Parse()
	return
}

func Init()(config *Config, error error)  {
	parseFlag()
	path, err := getConfigPath(defaultConfigPath)
	if err != nil {
		println(err)
		return nil,err
	}
	config,err = parseConf(path)
	return config,err
}

func parseConf(confPath string) (c *Config, err error) {
	cfg, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		return nil, err
	}
	ip, err := cfg.GetValue("local", "ip")
	if err != nil {
		return nil, err
	}
	//ip, err := cfg.GetValue("local", "ip")
	//if err != nil {
	//	return nil, err
	//}
	//port, err := cfg.GetValue("local", "port")
	//if err != nil {
	//	return nil, err
	//}


	tcpPort,err := cfg.GetValue("tcp", "port")
	if err != nil {
		return nil, err
	}
	tcpBind := ip + ":" + tcpPort

	websocketIp, err := cfg.GetValue("websocket", "port")
	if err != nil {
		return nil, err
	}
	websocketBind := ip + ":" + websocketIp

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

	//grpcAddr, err9 := cfg.GetValue("grpc", "serverAddr")
	//if err9 != nil {
	//	return nil, err9
	//}

	//grpcServerId, err7 := cfg.GetValue("grpc", "serverId")
	//if err7 != nil {
	//	return nil, err7
	//}

	XrpcAddr, err := cfg.GetValue("xrpc", "addr")
	if err != nil {
		return nil, err
	}


	conf := &Config{
		//Host: host,
		//Ip:   ip,
		//Port: port,
		TCP: &TCP{
			Bind:         []string{tcpBind},
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
			Bind: []string{websocketBind},
		},
		//Bucket: &Bucket{
		//	Size:          32,
		//	Channel:       1024,
		//	Room:          1024,
		//	RoutineAmount: 32,
		//	RoutineSize:   1024,
		//},
		Redis: &Redis{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPwd,
			Database: redisDatabase,
		},
		//Grpc:&Grpc{
		//	ServerId:grpcServerId,
		//	ServerAddr:grpcAddr,
		//},
		Xrpc:&Xrpc{
			Addr: XrpcAddr,
		},
	}
	return conf, nil
}


// getConfigPath get config from config path,if error, default would be give
//config path is  root_path/config/Config.ini
//targetPath Specified path
func getConfigPath(targetPath string) (configPath string, error error) {
	_, err := os.Stat(targetPath)
	if err == nil {
		return  targetPath, nil
	}
	if os.IsNotExist(err) {
		abPath, _ := filepath.Abs(os.Args[0])
		dir := filepath.Dir(abPath)
		log.Printf(dir)
		separator := string(filepath.Separator)
		configPath := dir + separator + "config" + separator + "Config.ini"
		_, err := os.Stat(configPath)
		if err == nil {
			return configPath, nil
		}
	}
	return "", errors.New("find config path unknown error")
}