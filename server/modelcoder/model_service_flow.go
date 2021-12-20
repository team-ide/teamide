package modelcoder

type ServiceFlowModel struct {
	Name  string                 `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type  string                 `json:"type,omitempty"` // 类型
	Steps []ServiceFlowStepModel `json:"steps,omitempty"`
}

func (this_ *ServiceFlowModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowModel) GetType() *ServiceModelType {
	return SERVICE_FLOW
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

type ServiceFlowStepStartModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
	Next string `json:"next,omitempty"` // 下个阶段
}

func (this_ *ServiceFlowStepStartModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepStartModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_START
}

type ServiceFlowStepDataModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
	Next string `json:"next,omitempty"` // 下个阶段
}

func (this_ *ServiceFlowStepDataModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepDataModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_DATA
}

type ServiceFlowStepDecisionModel struct {
	Name      string                                `json:"name,omitempty"`      // 名称，同一个应用中唯一
	Type      string                                `json:"type,omitempty"`      // 类型
	Condition string                                `json:"condition,omitempty"` // if条件
	IfNext    string                                `json:"ifNext,omitempty"`    // 满足 if 条件 进行下一步
	ElseIfs   []*ServiceFlowStepDecisionElseIfModel `json:"elseIfs,omitempty"`   // 不满足 if 条件 进入else if
	ElseNext  string                                `json:"elseNext,omitempty"`  // 不满足 if 和 else if 条件 进入else
}

type ServiceFlowStepDecisionElseIfModel struct {
	Condition string `json:"condition,omitempty"` // 条件
	Next      string `json:"next,omitempty"`      // 下个阶段
}
type ServiceFlowStepDecisionElseModel struct {
	Next string `json:"next,omitempty"` // 下个阶段
}

func (this_ *ServiceFlowStepDecisionModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepDecisionModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_DECISION
}

type ServiceFlowStepDaoModel struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Type    string `json:"type,omitempty"`    // 类型
	DaoName string `json:"daoName,omitempty"` // 数据层名称
	Next    string `json:"next,omitempty"`    // 下个阶段
}

func (this_ *ServiceFlowStepDaoModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepDaoModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_DAO
}

type ServiceFlowStepServiceModel struct {
	Name        string `json:"name,omitempty"`        // 名称，同一个应用中唯一
	Type        string `json:"type,omitempty"`        // 类型
	ServiceName string `json:"serviceName,omitempty"` // 服务名称
	Next        string `json:"next,omitempty"`        // 下个阶段
}

func (this_ *ServiceFlowStepServiceModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepServiceModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_SERVICE
}

type ServiceFlowStepEndModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *ServiceFlowStepEndModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepEndModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_END
}

type ServiceFlowStepErrorModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
	Code string `json:"code,omitempty"` // 异常码码
	Msg  string `json:"msg,omitempty"`  // 异常信息
}

func (this_ *ServiceFlowStepErrorModel) GetName() string {
	return this_.Name
}

func (this_ *ServiceFlowStepErrorModel) GetType() *ServiceFlowStepModelType {
	return SERVICE_FLOW_STEP_ERROR
}

type ServiceFlowStepModel interface {
	GetName() string                    // 名称，同一个应用中唯一
	GetType() *ServiceFlowStepModelType // 类型，sql、http、redis等
}

type ServiceFlowStepModelType struct {
	Value   string                                                                                                                                   `json:"value,omitempty"`
	Text    string                                                                                                                                   `json:"text,omitempty"`
	Execute func(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) `json:"-"`
}

var (
	serviceFlowStepModelTypes []*ServiceFlowStepModelType

	SERVICE_FLOW_STEP_START    = newServiceFlowStepModelType("start", "开始", invokeServiceFlowStepStart)
	SERVICE_FLOW_STEP_DATA     = newServiceFlowStepModelType("data", "Data", invokeServiceFlowStepData)
	SERVICE_FLOW_STEP_DECISION = newServiceFlowStepModelType("decision", "策略", invokeServiceFlowStepDecision)
	SERVICE_FLOW_STEP_DAO      = newServiceFlowStepModelType("dao", "Dao", invokeServiceFlowStepDao)
	SERVICE_FLOW_STEP_SERVICE  = newServiceFlowStepModelType("service", "Service", invokeServiceFlowStepService)
	SERVICE_FLOW_STEP_END      = newServiceFlowStepModelType("end", "结束", invokeServiceFlowStepEnd)
	SERVICE_FLOW_STEP_ERROR    = newServiceFlowStepModelType("error", "Error", invokeServiceFlowStepError)
)

func newServiceFlowStepModelType(value, text string, execute func(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error)) *ServiceFlowStepModelType {
	res := &ServiceFlowStepModelType{
		Value:   value,
		Text:    text,
		Execute: execute,
	}
	serviceFlowStepModelTypes = append(serviceFlowStepModelTypes, res)
	return res
}

func GetServiceFlowStepModelType(value string) *ServiceFlowStepModelType {
	for _, one := range serviceFlowStepModelTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}
