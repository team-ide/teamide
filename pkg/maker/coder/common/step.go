package common

import "teamide/pkg/maker/model"

type StepGenerator struct {
	GenStep      func(step *model.StepModel) (err error)
	GenStepVar   func(step *model.StepVarModel) (err error)
	GenStepError func(step *model.StepErrorModel) (err error)

	GenStepCacheGet          func(step *model.StepCacheModel) (err error)
	GenStepCacheGets         func(step *model.StepCacheModel) (err error)
	GenStepCacheSet          func(step *model.StepCacheModel) (err error)
	GenStepCacheSets         func(step *model.StepCacheModel) (err error)
	GenStepCacheRemove       func(step *model.StepCacheModel) (err error)
	GenStepCacheRemoves      func(step *model.StepCacheModel) (err error)
	GenStepCacheSetIfAbsent  func(step *model.StepCacheModel) (err error)
	GenStepCacheSetIfAbsents func(step *model.StepCacheModel) (err error)
	GenStepCacheClean        func(step *model.StepCacheModel) (err error)

	GenStepCommand func(step *model.StepCommandModel) (err error)

	GenStepDao func(step *model.StepDaoModel) (err error)

	GenStepFileGet      func(step *model.StepDaoModel) (err error)
	GenStepFileRead     func(step *model.StepDaoModel) (err error)
	GenStepFileWrite    func(step *model.StepDaoModel) (err error)
	GenStepFileAppend   func(step *model.StepDaoModel) (err error)
	GenStepFileRemove   func(step *model.StepDaoModel) (err error)
	GenStepFileReadLine func(step *model.StepDaoModel) (err error)

	GenStepDbSelect     func(step *model.StepDbModel) (err error)
	GenStepDbSelectOne  func(step *model.StepDbModel) (err error)
	GenStepDbSelectPage func(step *model.StepDbModel) (err error)
	GenStepDbInsert     func(step *model.StepDbModel) (err error)
	GenStepDbUpdate     func(step *model.StepDbModel) (err error)
	GenStepDbDelete     func(step *model.StepDbModel) (err error)
	GenStepDbCustom     func(step *model.StepDbModel) (err error)

	GenStepHttpGet     func(step *model.StepHttpModel) (err error)
	GenStepHttpPose    func(step *model.StepHttpModel) (err error)
	GenStepHttpPut     func(step *model.StepHttpModel) (err error)
	GenStepHttpDelete  func(step *model.StepHttpModel) (err error)
	GenStepHttpHead    func(step *model.StepHttpModel) (err error)
	GenStepHttpTrace   func(step *model.StepHttpModel) (err error)
	GenStepHttpOptions func(step *model.StepHttpModel) (err error)

	GenSteLockLock   func(step *model.StepLockModel) (err error)
	GenSteLockUnlock func(step *model.StepLockModel) (err error)

	GenSteMqPush func(step *model.StepMqModel) (err error)
	GenSteMqPull func(step *model.StepMqModel) (err error)

	GenSteScript  func(step *model.StepScriptModel) (err error)
	GenSteService func(step *model.StepServiceModel) (err error)

	GenSteEsDeleteIndex func(step *model.StepEsModel) (err error)
	GenSteEsCreateIndex func(step *model.StepEsModel) (err error)
	GenSteEsIndexNames  func(step *model.StepEsModel) (err error)
	GenSteEsGetMapping  func(step *model.StepEsModel) (err error)
	GenSteEsPutMapping  func(step *model.StepEsModel) (err error)
	GenSteEsSearch      func(step *model.StepEsModel) (err error)
	GenSteEsInsert      func(step *model.StepEsModel) (err error)
	GenSteEsUpdate      func(step *model.StepEsModel) (err error)
	GenSteEsDelete      func(step *model.StepEsModel) (err error)
	GenSteEsReindex     func(step *model.StepEsModel) (err error)
	GenSteEsScroll      func(step *model.StepEsModel) (err error)
}
