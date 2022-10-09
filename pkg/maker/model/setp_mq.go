package model

type StepMqModel struct {
	*StepModel

	Mq    string `json:"mq,omitempty"`    // MQ操作
	Topic string `json:"topic,omitempty"` // MQ topic

}
