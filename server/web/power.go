package web

import (
	"server/base"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
)

func checkPower(api *base.ApiWorker, JWT *base.JWTBean, c *gin.Context) bool {
	if api.Power.ShouldLogin && (JWT == nil || JWT.UserId == 0) {
		base.ResponseJSON(nil, base.ShouldLoginError, c)
		return false
	}
	ps := getPowersByJWT(JWT)

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

	base.ResponseJSON(nil, base.NoPowerError, c)
	return find
}

func getPowersByJWT(JWT *base.JWTBean) (powers []*base.PowerAction) {
	var userId int64 = 0
	if JWT != nil {
		userId = JWT.UserId
	}
	ps := base.GetPowers()
	psStrs := getPowersByUserId(userId)
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

func getPowersByUserId(userId int64) (powers []string) {
	ps := base.GetPowers()
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
