package module_ssh_tunnel

import "time"

const (
	// ModuleTunnel 工具箱模块
	ModuleTunnel = "tunnel"
	// TableTunnel 隧道表
	TableTunnel        = "TM_TUNNEL"
	TableTunnelComment = "隧道"
)

// TunnelModel 隧道模型，和隧道表对应
type TunnelModel struct {
	TunnelId     int64     `json:"tunnelId,omitempty"`
	TunnelType   string    `json:"tunnelType,omitempty"`
	Name         string    `json:"name,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	Option       string    `json:"option,omitempty"`
	UserId       int64     `json:"userId,omitempty"`
	DeleteUserId int64     `json:"deleteUserId,omitempty"`
	Deleted      int8      `json:"deleted,omitempty"`
	CreateTime   time.Time `json:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`
	DeleteTime   time.Time `json:"deleteTime,omitempty"`
}

func (entity *TunnelModel) GetTableName() string {
	return TableTunnel
}

func (entity *TunnelModel) GetPKColumnName() string {
	return ""
}
