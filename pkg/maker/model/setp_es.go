package model

type StepEsModel struct {
	*StepModel

	Es string `json:"es,omitempty"` // ES操作

}
