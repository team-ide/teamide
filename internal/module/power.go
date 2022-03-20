package module

import (
	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
	"strings"
	base2 "teamide/internal/server/base"
)

func (this_ *Api) checkPower(api *base2.ApiWorker, JWT *base2.JWTBean, c *gin.Context) bool {
	if base2.IsStandAlone {
		if api.Power.AllowNative {
			return true
		}
	}
	if api.Power.ShouldLogin && (JWT == nil || JWT.UserId == 0) {
		base2.ResponseJSON(nil, base2.ShouldLoginError, c)
		return false
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

	base2.ResponseJSON(nil, base2.NoPowerError, c)
	return find
}

func (this_ *Api) getPowersByJWT(JWT *base2.JWTBean) (powers []*base2.PowerAction) {
	var userId int64 = 0
	if JWT != nil {
		userId = JWT.UserId
	}
	ps := base2.GetPowers()
	psStrs := this_.getPowersByUserId(userId)
	if len(psStrs) == 0 {
		return
	}
	for _, power := range ps {
		if arrays.ContainsString(psStrs, power.Action) >= 0 {
			powers = append(powers, power)
		}
	}
	return
}

func (this_ *Api) getPowersByUserId(userId int64) (powers []string) {
	ps := base2.GetPowers()
	for _, power := range ps {
		if !power.ShouldLogin {
			powers = append(powers, power.Action)
		}
	}
	if userId != 0 {
		var userPowers []string
		for _, power := range ps {
			if !power.ShouldLogin {
				continue
			}
			if arrays.ContainsString(userPowers, power.Action) >= 0 {
				powers = append(powers, power.Action)
			} else {
				if strings.Index(power.Action, "user_") == 0 {
					powers = append(powers, power.Action)
				}
				if strings.Index(power.Action, "manage_") == 0 {
					powers = append(powers, power.Action)
				}
				if strings.Index(power.Action, "workspace_") == 0 {
					powers = append(powers, power.Action)
				}
				if strings.Index(power.Action, "toolbox_") == 0 {
					powers = append(powers, power.Action)
				}
			}
		}
	}
	return
}
