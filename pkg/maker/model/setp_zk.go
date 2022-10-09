package model

type StepZkModel struct {
	*StepModel

	Zk string `json:"zk,omitempty"` // ZK操作
}
