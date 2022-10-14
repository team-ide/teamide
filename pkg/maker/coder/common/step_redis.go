package common

import "teamide/pkg/maker/model"

type StepRedisGenerator interface {
	GenStepDel(step *model.StepRedisModel) (err error)
	GenStepExists(step *model.StepRedisModel) (err error)
	GenStepExpire(step *model.StepRedisModel) (err error)
	GenStepPersist(step *model.StepRedisModel) (err error)
	GenStepKeys(step *model.StepRedisModel) (err error)
	GenStepPTTL(step *model.StepRedisModel) (err error)
	GenStepRename(step *model.StepRedisModel) (err error)
	GenStepRenameNX(step *model.StepRedisModel) (err error)
	GenStepMove(step *model.StepRedisModel) (err error)
	GenStepType(step *model.StepRedisModel) (err error)

	GenStepSet(step *model.StepRedisModel) (err error)
	GenStepGet(step *model.StepRedisModel) (err error)
	GenStepGetSet(step *model.StepRedisModel) (err error)
	GenStepMSet(step *model.StepRedisModel) (err error)
	GenStepMGet(step *model.StepRedisModel) (err error)
	GenStepSetNX(step *model.StepRedisModel) (err error)
	GenStepStrLen(step *model.StepRedisModel) (err error)
	GenStepIncr(step *model.StepRedisModel) (err error)
	GenStepDecr(step *model.StepRedisModel) (err error)
	GenStepIncrBy(step *model.StepRedisModel) (err error)
	GenStepAppend(step *model.StepRedisModel) (err error)
}
