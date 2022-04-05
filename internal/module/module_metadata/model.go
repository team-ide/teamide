package module_metadata

import "time"

const (
	// ModuleMetadata 元数据模块
	ModuleMetadata = "metadata"
	// TableMetadata 元数据信息表
	TableMetadata = "TM_METADATA"
)

// MetadataModel 元数据模型，和元数据表对应
type MetadataModel struct {
	MetadataId int64     `json:"metadataId,omitempty"`
	SourceType int       `json:"sourceType,omitempty"`
	SourceId   int64     `json:"sourceId,omitempty"`
	StructCode int       `json:"structCode,omitempty"`
	FieldCode  int       `json:"fieldCode,omitempty"`
	Value      string    `json:"value,omitempty"`
	ParentId   int64     `json:"parentId,omitempty"`
	Deleted    int8      `json:"deleted,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}
