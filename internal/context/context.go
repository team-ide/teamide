package context

import (
	"go.uber.org/zap"
	"teamide/internal/config"
	"teamide/pkg/db"
)

type ServerContext struct {
	ServerContext  string
	ServerHost     string
	ServerPort     int
	ServerUrl      string
	ServerConfig   *config.ServerConfig
	DatabaseWorker db.DatabaseWorker
	DatabaseConfig *db.DatabaseConfig `json:"-" yaml:"-"`
	Logger         *zap.Logger        `json:"-" yaml:"-"`
	Decryption     *Decryption        `json:"-" yaml:"-"`
	HttpAesKey     string             `json:"-" yaml:"-"`
	IsServer       bool
	IsHtmlDev      bool
	IsServerDev    bool
	RootDir        string
	UserHomeDir    string
}

func (this_ *ServerContext) GetFilesDir() string {
	return this_.ServerConfig.Server.Data + "files/"
}

func (this_ *ServerContext) GetFilesFile(path string) string {
	return this_.GetFilesDir() + path
}
