package context

import (
	"github.com/team-ide/cron"
	"github.com/team-ide/go-tool/db"
	"go.uber.org/zap"
	"teamide/internal/config"
)

type ServerContext struct {
	Version        string
	ServerContext  string
	ServerHost     string
	ServerPort     int
	ServerUrl      string
	ServerConfig   *config.ServerConfig
	DatabaseWorker db.IService
	DatabaseConfig *db.Config  `json:"-" yaml:"-"`
	Logger         *zap.Logger `json:"-" yaml:"-"`
	LoggerP1       *zap.Logger `json:"-" yaml:"-"`
	LoggerP2       *zap.Logger `json:"-" yaml:"-"`
	Decryption     *Decryption `json:"-" yaml:"-"`
	HttpAesKey     string      `json:"-" yaml:"-"`
	JWTAesKey      string      `json:"-" yaml:"-"`
	IsServer       bool
	IsHtmlDev      bool
	IsServerDev    bool
	RootDir        string
	UserHomeDir    string
	Setting        *Setting
	CronHandler    *cron.Cron
	IsHttps        bool
	ServerProtocol string
}

func (this_ *ServerContext) GetFilesDir() string {
	return this_.ServerConfig.Server.Data + "files/"
}

func (this_ *ServerContext) GetFilesFile(path string) string {
	return this_.GetFilesDir() + path
}
