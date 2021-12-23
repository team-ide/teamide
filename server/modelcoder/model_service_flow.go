package modelcoder

type ServiceFlowModel struct {
	Name       string                 `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string                 `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel       `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel         `json:"result,omitempty"`     // 结果配置
	Steps      []ServiceFlowStepModel `json:"steps,omitempty"`
}

func (this_ *ServiceFlowModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowModel) GetType() *ServiceModelType {
	return SERVICE_FLOW
}

func (this_ *ServiceFlowModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *ServiceFlowModel) GetResult() *VariableModel {
	return this_.Result
}

func (this_ *ServiceFlowModel) GetStartStep() ServiceFlowStepModel {
	if len(this_.Steps) == 0 {
		return nil
	}
	for _, one := range this_.Steps {
		if one.GetType() == SERVICE_FLOW_STEP_START {
			return one
		}
	}
	return nil
}

func (this_ *ServiceFlowModel) GetStep(name string) ServiceFlowStepModel {
	if len(this_.Steps) == 0 {
		return nil
	}
	for _, one := range this_.Steps {
		if one.GetName() == name {
			return one
		}
	}
	return nil
}
