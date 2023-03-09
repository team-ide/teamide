package module

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/base"
)

func (this_ *Api) checkPower(api *base.ApiWorker, JWT *base.JWTBean, c *gin.Context) bool {

	if api.Power.ShouldLogin && (JWT == nil || JWT.UserId == 0) {
		this_.Logger.Error("权限验证失败", zap.Error(base.ShouldLoginError))
		base.ResponseJSON(nil, base.ShouldLoginError, c)
		return false
	}
	if !this_.IsServer && api.Power.StandAlone {
		return true
	}
	if !api.Power.ShouldPower {
		return true
	}
	ps := this_.getPowersByJWT(JWT)

	find := false
	for _, power := range ps {
		if power == api.Power {
			find = true
			break
		}
	}
	if find {
		return find
	}
	this_.Logger.Error("权限验证失败", zap.Error(base.NoPowerError))
	base.ResponseJSON(nil, base.NoPowerError, c)
	return find
}

func (this_ *Api) getPowersByJWT(JWT *base.JWTBean) (powers []*base.PowerAction) {
	var userId int64 = 0
	if JWT != nil {
		userId = JWT.UserId
	}
	ps := base.GetPowers()
	psStr := this_.getPowersByUserId(userId)
	if len(psStr) == 0 {
		return
	}
	for _, power := range ps {
		if util.StringIndexOf(psStr, power.Action) >= 0 {
			powers = append(powers, power)
		}
	}
	return
}

func (this_ *Api) getPowersByUserId(userId int64) (powers []string) {

	ps := base.GetPowers()
	if !this_.IsServer {
		for _, power := range ps {
			if power.StandAlone {
				powers = append(powers, power.Action)
			}
		}
		return
	}

	for _, power := range ps {
		if !power.ShouldLogin {
			powers = append(powers, power.Action)
		}
	}
	if userId != 0 {
		var userPowers []string

		var isSuperRole bool
		roles, _ := this_.powerUserService.QueryPowerRolesByUserId(userId)
		for _, role := range roles {
			if role.RoleType == base.SuperRoleType {
				isSuperRole = true
			} else {

			}
		}
		for _, power := range ps {
			if util.StringIndexOf(powers, power.Action) >= 0 {
				continue
			}
			if !power.ShouldPower || isSuperRole {
				powers = append(powers, power.Action)
				continue
			}
			if util.StringIndexOf(userPowers, power.Action) >= 0 {
				powers = append(powers, power.Action)
			}
		}
	}
	return
}
