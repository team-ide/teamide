package modelcoder

type ServiceModel interface {
	GetName() string            // 名称，同一个应用中唯一
	GetType() *ServiceModelType // 类型，sql、http、redis等
}

type ServiceModelType struct {
	Value   string
	Text    string
	Execute func(application *Application, service ServiceModel, variable *invokeVariable) (res interface{}, err error)
}

var (
	serviceModelTypes []*ServiceModelType

	SERVICE_FLOW = newServiceModelType("FLOW", "流程", invokeServiceFlow)
)

func newServiceModelType(value, text string, execute func(application *Application, service ServiceModel, variable *invokeVariable) (res interface{}, err error)) *ServiceModelType {
	res := &ServiceModelType{
		Value:   value,
		Text:    text,
		Execute: execute,
	}
	serviceModelTypes = append(serviceModelTypes, res)
	return res
}

func GetServiceModelType(value string) *ServiceModelType {
	for _, one := range serviceModelTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}
