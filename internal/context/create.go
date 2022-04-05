package context

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"strings"
	"teamide/internal/config"
	"teamide/pkg/db"
	"teamide/pkg/util"
)

type ServerConf struct {
	Server       string
	PublicKey    string
	PrivateKey   string
	IsStandAlone bool
	IsHtmlDev    bool
	RootDir      string
	UserHomeDir  string
}

func NewServerContext(serverConf ServerConf) (context *ServerContext, err error) {
	context = &ServerContext{
		IsStandAlone: serverConf.IsStandAlone,
		IsHtmlDev:    serverConf.IsHtmlDev,
		RootDir:      serverConf.RootDir,
		UserHomeDir:  serverConf.UserHomeDir,
	}
	context.HttpAesKey = "Q56hFAauWk18Gy2i"
	var serverConfig *config.ServerConfig
	serverConfig, err = config.CreateServerConfig(serverConf.Server)
	if err != nil {
		return
	}
	context.ServerConf = serverConf
	context.ServerConfig = serverConfig
	err = context.init(serverConfig)
	if err != nil {
		return
	}
	if serverConf.PublicKey != "" || serverConf.PrivateKey != "" {
		context.Decryption, err = NewDecryption(serverConf.PublicKey, serverConf.PrivateKey, context.Logger)
		if err != nil {
			return
		}
	} else {
		context.Decryption, err = NewDefaultDecryption(context.Logger)
		if err != nil {
			return
		}
	}
	return
}

//init 格式化配置，填充默认值
func (this_ *ServerContext) init(serverConfig *config.ServerConfig) (err error) {
	if !this_.IsStandAlone {
		if serverConfig.Server.Host == "" || serverConfig.Server.Port == 0 {
			err = errors.New("请检查Server配置是否正确")
			return
		}
	}
	if serverConfig.Server.Data == "" {
		serverConfig.Server.Data = this_.RootDir + "data"
	} else {
		serverConfig.Server.Data = this_.RootDir + strings.TrimPrefix(serverConfig.Server.Data, "./")
	}

	if !strings.HasSuffix(serverConfig.Server.Data, "/") {
		serverConfig.Server.Data += "/"
	}
	exist, err := util.PathExists(serverConfig.Server.Data)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.Data, 0777)
		if err != nil {
			return
		}
	}
	if serverConfig.Log.Filename == "" {
		serverConfig.Log.Filename = this_.RootDir + "log/server.log"
	} else {
		serverConfig.Log.Filename = this_.RootDir + strings.TrimPrefix(serverConfig.Log.Filename, "./")
	}

	var databaseConfig *db.DatabaseConfig
	if serverConfig.Mysql == nil || serverConfig.Mysql.Host == "" || serverConfig.Mysql.Port == 0 {
		databaseConfig = &db.DatabaseConfig{
			Type:     "sqlite",
			Database: this_.RootDir + "data/database",
		}
	} else {
		databaseConfig = &db.DatabaseConfig{
			Type:     "mysql",
			Host:     serverConfig.Mysql.Host,
			Port:     serverConfig.Mysql.Port,
			Database: serverConfig.Mysql.Database,
			Username: serverConfig.Mysql.Username,
			Password: serverConfig.Mysql.Password,
		}
	}
	this_.Logger = newZapLogger(serverConfig)

	this_.ServerContext = serverConfig.Server.Context
	if this_.ServerContext == "" || !strings.HasSuffix(this_.ServerContext, "/") {
		this_.ServerContext = this_.ServerContext + "/"
	}
	this_.ServerHost = serverConfig.Server.Host
	this_.ServerPort = serverConfig.Server.Port

	if this_.ServerHost == "" {
		this_.ServerHost = "127.0.0.1"
	}

	if this_.ServerPort == 0 {
		if this_.IsHtmlDev {
			this_.ServerPort = 21080
		} else {
			var listener net.Listener
			listener, err = net.Listen("tcp", ":0")
			if err != nil {
				this_.Logger.Error("随机端口获取失败", zap.Error(err))
				return
			}
			this_.ServerPort = listener.Addr().(*net.TCPAddr).Port
		}
	}

	if this_.ServerHost == "0.0.0.0" || this_.ServerHost == ":" || this_.ServerHost == "::" {
		this_.ServerUrl = fmt.Sprint("http://127.0.0.1:", this_.ServerPort)
	} else {
		this_.ServerUrl = fmt.Sprint("http://", this_.ServerHost, ":", this_.ServerPort)
	}

	this_.DatabaseConfig = databaseConfig
	this_.DatabaseWorker, err = db.NewDatabaseWorker(*databaseConfig)
	if err != nil {
		this_.Logger.Error("数据库连接异常", zap.Error(err))
		return
	}

	return
}
