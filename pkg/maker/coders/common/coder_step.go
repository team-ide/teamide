package common

import (
	"teamide/pkg/maker/modelers"
)

type IStepCoder interface {
	Gen(code *Code, step *modelers.StepModel) (err error)
}

type IStepCacheCoder interface {
	GenGet(code *Code, step *modelers.StepCacheModel) (err error)
	GenGets(code *Code, step *modelers.StepCacheModel) (err error)
	GenSet(code *Code, step *modelers.StepCacheModel) (err error)
	GenSets(code *Code, step *modelers.StepCacheModel) (err error)
	GenRemove(code *Code, step *modelers.StepCacheModel) (err error)
	GenRemoves(code *Code, step *modelers.StepCacheModel) (err error)
	GenSetIfAbsent(code *Code, step *modelers.StepCacheModel) (err error)
	GenSetIfAbsents(code *Code, step *modelers.StepCacheModel) (err error)
	GenClean(code *Code, step *modelers.StepCacheModel) (err error)
}

type IStepCommandCoder interface {
	Gen(code *Code, step *modelers.StepCommandModel) (err error)
}

type IStepDaoCoder interface {
	Gen(code *Code, step *modelers.StepDaoModel) (err error)
}

type IStepDbCoder interface {
	GenSelect(code *Code, step *modelers.StepDbModel) (err error)
	GenSelectOne(code *Code, step *modelers.StepDbModel) (err error)
	GenSelectPage(code *Code, step *modelers.StepDbModel) (err error)
	GenInsert(code *Code, step *modelers.StepDbModel) (err error)
	GenUpdate(code *Code, step *modelers.StepDbModel) (err error)
	GenDelete(code *Code, step *modelers.StepDbModel) (err error)
	GenCustom(code *Code, step *modelers.StepDbModel) (err error)
}

type IStepErrorCoder interface {
	Gen(code *Code, step *modelers.StepErrorModel) (err error)
}

type IStepEsCoder interface {
	GenDeleteIndex(code *Code, step *modelers.StepEsModel) (err error)
	GenCreateIndex(code *Code, step *modelers.StepEsModel) (err error)
	GenIndexNames(code *Code, step *modelers.StepEsModel) (err error)
	GenGetMapping(code *Code, step *modelers.StepEsModel) (err error)
	GenPutMapping(code *Code, step *modelers.StepEsModel) (err error)
	GenSearch(code *Code, step *modelers.StepEsModel) (err error)
	GenInsert(code *Code, step *modelers.StepEsModel) (err error)
	GenUpdate(code *Code, step *modelers.StepEsModel) (err error)
	GenDelete(code *Code, step *modelers.StepEsModel) (err error)
	GenReindex(code *Code, step *modelers.StepEsModel) (err error)
	GenScroll(code *Code, step *modelers.StepEsModel) (err error)
}

type IStepFileCoder interface {
	GenGet(code *Code, step *modelers.StepFileModel) (err error)
	GenRead(code *Code, step *modelers.StepFileModel) (err error)
	GeneWrite(code *Code, step *modelers.StepFileModel) (err error)
	GenAppend(code *Code, step *modelers.StepFileModel) (err error)
	GenRemove(code *Code, step *modelers.StepFileModel) (err error)
	GenReadLine(code *Code, step *modelers.StepFileModel) (err error)
}

type IStepFuncCoder interface {
	Gen(code *Code, step *modelers.StepFuncModel) (err error)
}

type IStepHttpCoder interface {
	GenGet(code *Code, step *modelers.StepHttpModel) (err error)
	GenPose(code *Code, step *modelers.StepHttpModel) (err error)
	GenPut(code *Code, step *modelers.StepHttpModel) (err error)
	GenDelete(code *Code, step *modelers.StepHttpModel) (err error)
	GenHead(code *Code, step *modelers.StepHttpModel) (err error)
	GenTrace(code *Code, step *modelers.StepHttpModel) (err error)
	GenOptions(code *Code, step *modelers.StepHttpModel) (err error)
}

type IStepLockCoder interface {
	GenLock(code *Code, step *modelers.StepLockModel) (err error)
	GenUnlock(code *Code, step *modelers.StepLockModel) (err error)
}

type IStepMqCoder interface {
	GenPush(code *Code, step *modelers.StepMqModel) (err error)
	GenPull(code *Code, step *modelers.StepMqModel) (err error)
}

type IStepRedisCoder interface {
	GenStepDel(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepExists(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepExpire(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepPersist(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepKeys(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepPTTL(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepRename(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepRenameNX(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMove(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepType(code *Code, step *modelers.StepRedisModel) (err error)

	GenStepSet(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepGet(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepGetSet(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMSet(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMGet(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepSetNX(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepStrLen(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepIncr(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepDecr(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepIncrBy(code *Code, step *modelers.StepRedisModel) (err error)
	GenStepAppend(code *Code, step *modelers.StepRedisModel) (err error)
}

type IStepScriptCoder interface {
	Gen(code *Code, step *modelers.StepScriptModel) (err error)
}

type IStepServiceCoder interface {
	Gen(code *Code, step *modelers.StepServiceModel) (err error)
}

type IStepVarCoder interface {
	Gen(code *Code, step *modelers.StepVarModel) (err error)
}

type IStepZkCoder struct {
	GenGet       func(code *Code, step *modelers.StepZkModel) (err error)
	GenCreate    func(code *Code, step *modelers.StepZkModel) (err error)
	GenSet       func(code *Code, step *modelers.StepZkModel) (err error)
	GenStat      func(code *Code, step *modelers.StepZkModel) (err error)
	GenChildren  func(code *Code, step *modelers.StepZkModel) (err error)
	GenDelete    func(code *Code, step *modelers.StepZkModel) (err error)
	GenExists    func(code *Code, step *modelers.StepZkModel) (err error)
	GenGetW      func(code *Code, step *modelers.StepZkModel) (err error)
	GenChildrenW func(code *Code, step *modelers.StepZkModel) (err error)
}
