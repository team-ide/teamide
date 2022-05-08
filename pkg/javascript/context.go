package javascript

import (
	"teamide/pkg/util"
)

func GetContext() map[string]interface{} {
	baseContext := map[string]interface{}{}
	baseContext["_$now"] = util.Now
	baseContext["_$nowTime"] = util.GetNowTime
	baseContext["_$uuid"] = util.GenerateUUID
	baseContext["_$randomInt"] = util.RandomInt
	baseContext["_$randomString"] = util.RandomString
	baseContext["_$randomUserName"] = util.RandomUserName
	baseContext["_$toPinYin"] = util.ToPinYin

	return baseContext
}
