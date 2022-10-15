package common

import (
	"teamide/pkg/maker/modelers"
)

type IStepCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepModel) (err error)
}

type IStepCacheCoder interface {
	GenGet(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenGets(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenSet(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenSets(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenRemove(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenRemoves(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenSetIfAbsent(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenSetIfAbsents(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
	GenClean(appCoder IAppCoder, code *Code, step *modelers.StepCacheModel) (err error)
}

type IStepCommandCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepCommandModel) (err error)
}

type IStepDaoCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepDaoModel) (err error)
}

type IStepDbCoder interface {
	GenSelect(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenSelectOne(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenSelectPage(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenInsert(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenUpdate(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenDelete(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
	GenCustom(appCoder IAppCoder, code *Code, step *modelers.StepDbModel) (err error)
}

type IStepErrorCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepErrorModel) (err error)
}

type IStepEsCoder interface {
	GenDeleteIndex(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenCreateIndex(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenIndexNames(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenGetMapping(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenPutMapping(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenSearch(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenInsert(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenUpdate(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenDelete(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenReindex(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
	GenScroll(appCoder IAppCoder, code *Code, step *modelers.StepEsModel) (err error)
}

type IStepFileCoder interface {
	GenGet(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
	GenRead(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
	GeneWrite(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
	GenAppend(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
	GenRemove(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
	GenReadLine(appCoder IAppCoder, code *Code, step *modelers.StepFileModel) (err error)
}

type IStepFuncCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepFuncModel) (err error)
}

type IStepHttpCoder interface {
	GenGet(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenPose(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenPut(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenDelete(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenHead(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenTrace(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
	GenOptions(appCoder IAppCoder, code *Code, step *modelers.StepHttpModel) (err error)
}

type IStepLockCoder interface {
	GenLock(appCoder IAppCoder, code *Code, step *modelers.StepLockModel) (err error)
	GenUnlock(appCoder IAppCoder, code *Code, step *modelers.StepLockModel) (err error)
}

type IStepMqCoder interface {
	GenPush(appCoder IAppCoder, code *Code, step *modelers.StepMqModel) (err error)
	GenPull(appCoder IAppCoder, code *Code, step *modelers.StepMqModel) (err error)
}

type IStepRedisCoder interface {
	GenStepDel(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepExists(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepExpire(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepPersist(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepKeys(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepPTTL(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepRename(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepRenameNX(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMove(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepType(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)

	GenStepSet(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepGet(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepGetSet(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMSet(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepMGet(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepSetNX(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepStrLen(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepIncr(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepDecr(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepIncrBy(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
	GenStepAppend(appCoder IAppCoder, code *Code, step *modelers.StepRedisModel) (err error)
}

type IStepScriptCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepScriptModel) (err error)
}

type IStepServiceCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepServiceModel) (err error)
}

type IStepVarCoder interface {
	Gen(appCoder IAppCoder, code *Code, step *modelers.StepVarModel) (err error)
}

type IStepZkCoder struct {
	GenGet       func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenCreate    func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenSet       func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenStat      func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenChildren  func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenDelete    func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenExists    func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenGetW      func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
	GenChildrenW func(appCoder IAppCoder, code *Code, step *modelers.StepZkModel) (err error)
}
