package modelcoder

type ServiceFlow struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *ServiceFlow) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlow) GetType() *ServiceModelType {
	return SERVICE_FLOW
}
