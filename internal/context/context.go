package context

import (
	"go.uber.org/zap"
	"teamide/internal/config"
	"teamide/pkg/db"
)

type ServerContext struct {
	ServerContext string
	ServerHost    string
	ServerPort    int
	ServerUrl     string
	//ServerConf     ServerConf
	ServerConfig   *config.ServerConfig
	DatabaseWorker db.DatabaseWorker
	DatabaseConfig *db.DatabaseConfig `json:"-" yaml:"-"`
	Logger         *zap.Logger        `json:"-" yaml:"-"`
	Decryption     *Decryption        `json:"-" yaml:"-"`
	HttpAesKey     string             `json:"-" yaml:"-"`
	IsStandAlone   bool
	IsHtmlDev      bool
	IsServerDev    bool
	RootDir        string
	UserHomeDir    string
}
