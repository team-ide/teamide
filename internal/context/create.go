package context

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strings"
	"teamide/internal/config"
	"teamide/pkg/db"
	"teamide/pkg/node"
	"teamide/pkg/util"
)

type ServerConf struct {
	Version     string
	Server      string
	PublicKey   string
	PrivateKey  string
	IsServer    bool
	IsHtmlDev   bool
	IsServerDev bool
	RootDir     string
	UserHomeDir string
}

func NewServerContext(serverConf ServerConf) (context *ServerContext, err error) {
	context = &ServerContext{
		IsServer:    serverConf.IsServer,
		IsHtmlDev:   serverConf.IsHtmlDev,
		IsServerDev: serverConf.IsServerDev,
		RootDir:     serverConf.RootDir,
		UserHomeDir: serverConf.UserHomeDir,
		Version:     serverConf.Version,
	}
	context.HttpAesKey = "Q56hFAauWk18Gy2i"
	var serverConfig *config.ServerConfig
	serverConfig, err = config.CreateServerConfig(serverConf.Server)
	if err != nil {
		return
	}
	//context.ServerConf = serverConf
	context.ServerConfig = serverConfig
	err = context.Init(serverConfig)
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

// Init 格式化配置，填充默认值
func (this_ *ServerContext) Init(serverConfig *config.ServerConfig) (err error) {
	if this_.IsServer {
		if serverConfig.Server.Port == 0 {
			err = errors.New("请检查Server配置是否正确")
			return
		}
	}

	if serverConfig.Server.Host == "" {
		serverConfig.Server.Host = "0.0.0.0"
	}
	if this_.IsHtmlDev {
		serverConfig.Server.Host = "127.0.0.1"
		serverConfig.Server.Port = 21080
	}
	if serverConfig.Server.Port == 0 {
		var listener net.Listener
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			this_.Logger.Error("随机端口获取失败", zap.Error(err))
			return
		}
		serverConfig.Server.Port = listener.Addr().(*net.TCPAddr).Port
		err = listener.Close()
		if err != nil {
			return
		}
	}

	if this_.IsServer {
		if serverConfig.Server.Data == "" {
			serverConfig.Server.Data = this_.RootDir + "data"
		} else {
			serverConfig.Server.Data = this_.RootDir + strings.TrimPrefix(serverConfig.Server.Data, "./")
		}
		if serverConfig.Log.Filename == "" {
			serverConfig.Log.Filename = this_.RootDir + "log/server.log"
		} else {
			serverConfig.Log.Filename = this_.RootDir + strings.TrimPrefix(serverConfig.Log.Filename, "./")
		}
	} else {
		if this_.UserHomeDir == "" {
			err = errors.New("用户目录读取失败")
		}
		TeamIDEDir := this_.UserHomeDir + "/TeamIDE/"

		if serverConfig.Server.Data == "" {
			serverConfig.Server.Data = TeamIDEDir + "data"
		} else {
			serverConfig.Server.Data = TeamIDEDir + strings.TrimPrefix(serverConfig.Server.Data, "./")
		}
		if serverConfig.Log.Filename == "" {
			serverConfig.Log.Filename = TeamIDEDir + "log/server.log"
		} else {
			serverConfig.Log.Filename = TeamIDEDir + strings.TrimPrefix(serverConfig.Log.Filename, "./")
		}
	}
	serverConfig.Server.Data = util.FormatPath(serverConfig.Server.Data)

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

	if serverConfig.Server.TempDir == "" {
		serverConfig.Server.TempDir = serverConfig.Server.Data + "temp/"
	}
	if !strings.HasSuffix(serverConfig.Server.TempDir, "/") {
		serverConfig.Server.TempDir += "/"
	}
	exist, err = util.PathExists(serverConfig.Server.TempDir)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.TempDir, 0777)
		if err != nil {
			return
		}
	}
	if serverConfig.Server.BackupsDir == "" {
		serverConfig.Server.BackupsDir = serverConfig.Server.Data + "backups/"
	}
	if !strings.HasSuffix(serverConfig.Server.BackupsDir, "/") {
		serverConfig.Server.BackupsDir += "/"
	}
	exist, err = util.PathExists(serverConfig.Server.BackupsDir)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.BackupsDir, 0777)
		if err != nil {
			return
		}
	}

	if this_.IsServerDev {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Development = false
		this_.Logger, err = loggerConfig.Build()
		if err != nil {
			return
		}
	} else {
		this_.Logger = newZapLogger(serverConfig)
	}
	util.Logger = this_.Logger
	util.TempDir = serverConfig.Server.TempDir
	node.Logger = this_.Logger
	db.FileUploadDir = this_.GetFilesDir()

	this_.ServerContext = serverConfig.Server.Context
	if this_.ServerContext == "" || !strings.HasSuffix(this_.ServerContext, "/") {
		this_.ServerContext = this_.ServerContext + "/"
	}
	this_.ServerHost = serverConfig.Server.Host
	this_.ServerPort = serverConfig.Server.Port

	if this_.ServerHost == "0.0.0.0" || this_.ServerHost == ":" || this_.ServerHost == "::" {
		this_.ServerUrl = fmt.Sprint("http://localhost:", this_.ServerPort)
	} else {
		this_.ServerUrl = fmt.Sprintf("%s://%s:%d", "http", this_.ServerHost, this_.ServerPort)
	}

	var databaseConfig *db.DatabaseConfig
	if serverConfig.Mysql == nil || serverConfig.Mysql.Host == "" || serverConfig.Mysql.Port == 0 {
		databaseConfig = &db.DatabaseConfig{
			Type:         "sqlite",
			DatabasePath: serverConfig.Server.Data + "database",
		}
		err = this_.backupSqlite(serverConfig, databaseConfig)
		if err != nil {
			return
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

	this_.DatabaseConfig = databaseConfig
	this_.DatabaseWorker, err = db.NewDatabaseWorker(databaseConfig)
	if err != nil {
		this_.Logger.Error("数据库连接异常", zap.Error(err))
		return
	}

	initServerWebsocket()

	return
}

// backupSqlite 备份
func (this_ *ServerContext) backupSqlite(serverConfig *config.ServerConfig, databaseConfig *db.DatabaseConfig) (err error) {
	databasePath := databaseConfig.DatabasePath
	exist, err := util.PathExists(databasePath)
	if err != nil {
		return
	}
	if !exist {
		return
	}

	backupPath := serverConfig.Server.BackupsDir + "/版本-" + this_.Version + "-升级之前备份-数据库"

	exist, err = util.PathExists(backupPath)
	if err != nil {
		return
	}
	if exist {
		return
	}

	databaseFile, err := os.Open(databasePath)
	if err != nil {
		return
	}
	defer func() {
		_ = databaseFile.Close()
	}()

	backupFile, err := os.Create(backupPath)
	if err != nil {
		return
	}
	defer func() {
		_ = backupFile.Close()
	}()
	_, err = io.Copy(backupFile, databaseFile)

	return
}
