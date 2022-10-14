package common

import "teamide/pkg/maker/model"

type StepZkGenerator struct {
	GenGet       func(step *model.StepZkModel) (err error)
	GenCreate    func(step *model.StepZkModel) (err error)
	GenSet       func(step *model.StepZkModel) (err error)
	GenStat      func(step *model.StepZkModel) (err error)
	GenChildren  func(step *model.StepZkModel) (err error)
	GenDelete    func(step *model.StepZkModel) (err error)
	GenExists    func(step *model.StepZkModel) (err error)
	GenGetW      func(step *model.StepZkModel) (err error)
	GenChildrenW func(step *model.StepZkModel) (err error)
}
