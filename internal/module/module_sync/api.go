package module_sync

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"net/http"
	"net/url"
	"os"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/module/module_user"
	"teamide/pkg/base"
)

type api struct {
	toolboxService     *module_toolbox.ToolboxService
	userSettingService *module_user.UserSettingService
	userService        *module_user.UserService
}

func NewApi(toolboxService *module_toolbox.ToolboxService, userService *module_user.UserService, userSettingService *module_user.UserSettingService) *api {
	return &api{
		toolboxService:     toolboxService,
		userSettingService: userSettingService,
		userService:        userService,
	}
}

var (
	Power      = base.AppendPower(&base.PowerAction{Action: "sync", Text: "同步", ShouldLogin: true, StandAlone: true})
	exportFile = base.AppendPower(&base.PowerAction{Action: "exportFile", Text: "导出文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	checkFile  = base.AppendPower(&base.PowerAction{Action: "checkFile", Text: "检测文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	importFile = base.AppendPower(&base.PowerAction{Action: "importFile", Text: "导入文件", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: exportFile, Do: this_.exportFile})
	apis = append(apis, &base.ApiWorker{Power: checkFile, Do: this_.checkFile})
	apis = append(apis, &base.ApiWorker{Power: importFile, Do: this_.importFile})

	return
}

type BaseRequest struct {
	Key         string `json:"key"`
	Path        string `json:"path"`
	UserSetting bool   `json:"userSetting"`
	Toolbox     bool   `json:"toolbox"`
	ExistsDo    int    `json:"existsDo"`
}

func (this_ *api) checkFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	param := &BaseRequest{}
	if !base.RequestJSON(param, c) {
		return
	}

	filePath := this_.toolboxService.GetFilesFile(param.Path)
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	res, err = Read(param.Key, string(bs))
	if err != nil {
		return
	}

	return
}
func (this_ *api) importFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	param := &BaseRequest{}

	if !base.RequestJSON(param, c) {
		return
	}
	filePath := this_.toolboxService.GetFilesFile(param.Path)
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	info, err := Read(param.Key, string(bs))
	if err != nil {
		return
	}

	if param.UserSetting && info.UserSettingSize > 0 {
		var setting map[string]string
		setting, err = this_.userSettingService.Query(requestBean.JWT.UserId)
		if err != nil {
			return
		}
		option, ok := info.UserSetting["option"].(map[string]string)
		if ok {
			for k, v := range option {
				_, find := setting[k]
				if find && param.ExistsDo == 1 {
					continue
				}
				setting[k] = v
			}
			_, err = this_.userSettingService.Save(requestBean.JWT.UserId, setting)
			if err != nil {
				return
			}
		}
	}

	if param.Toolbox && info.ToolboxSize > 0 {
		var groupIdCache = map[int64]int64{}

		for _, d := range info.ToolboxGroupList {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			g := &module_toolbox.ToolboxGroupModel{}
			err = json.Unmarshal(bs, g)
			if err != nil {
				return
			}
			if g.GroupId == 0 {
				continue
			}
			if g.Name == "" {
				continue
			}
			gId := g.GroupId
			g.GroupId = 0
			find, _ := this_.toolboxService.GetUserGroupByName(requestBean.JWT.UserId, g.Name)
			if find != nil {
				groupIdCache[gId] = find.GroupId
				if param.ExistsDo == 1 {
					continue
				}
			} else {
				find = g
				find.UserId = requestBean.JWT.UserId
				_, err = this_.toolboxService.InsertGroup(find)
				if err != nil {
					return
				}
				groupIdCache[gId] = find.GroupId
			}
		}

		// 修改parentId
		for _, d := range info.ToolboxGroupList {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			g := &module_toolbox.ToolboxGroupModel{}
			err = json.Unmarshal(bs, g)
			if err != nil {
				return
			}
			if g.ParentId == 0 || groupIdCache[g.ParentId] == 0 {
				continue
			}
			_, err = this_.toolboxService.UpdateGroup(&module_toolbox.ToolboxGroupModel{
				GroupId:  g.GroupId,
				ParentId: groupIdCache[g.ParentId],
			})
		}

		var toolboxIdCache = map[int64]int64{}

		for _, d := range info.ToolboxList {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			t := &module_toolbox.ToolboxModel{}
			err = json.Unmarshal(bs, t)
			if err != nil {
				return
			}
			if t.ToolboxType != "" && t.Option != "" {
				toolboxType := module_toolbox.GetToolboxType(t.ToolboxType)
				if toolboxType == nil {
					continue
				}
			}
			tId := t.ToolboxId
			t.ToolboxId = 0
			if t.ToolboxType == "" {
				continue
			}
			if t.Name == "" {
				continue
			}
			find, _ := this_.toolboxService.GetUserToolboxByName(t.ToolboxType, t.Name, requestBean.JWT.UserId)
			var gId int64
			if t.GroupId != 0 {
				gId = groupIdCache[t.GroupId]
			}
			if find != nil {
				toolboxIdCache[tId] = find.ToolboxId
				if param.ExistsDo == 1 {
					continue
				}
				find.Option = t.Option
				if gId > 0 {
					find.GroupId = gId
				}
				_, err = this_.toolboxService.Update(find)
				if err != nil {
					return
				}
			} else {
				find = t
				if gId > 0 {
					find.GroupId = gId
				}
				find.UserId = requestBean.JWT.UserId
				_, err = this_.toolboxService.Insert(find)
				if err != nil {
					return
				}
				toolboxIdCache[tId] = find.ToolboxId
			}
		}

		for _, d := range info.ToolboxExtendList {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			t := &module_toolbox.ToolboxExtendModel{}
			err = json.Unmarshal(bs, t)
			if err != nil {
				return
			}
			t.ExtendId = 0
			if t.Value == "" || t.Name == "" || t.ExtendType == "" {
				continue
			}
			if t.ToolboxId != 0 {
				t.ToolboxId = toolboxIdCache[t.ToolboxId]
				if t.ToolboxId == 0 {
					continue
				}
			}
			find, _ := this_.toolboxService.QueryExtends(&module_toolbox.ToolboxExtendModel{
				UserId:     requestBean.JWT.UserId,
				ExtendType: t.ExtendType,
				Name:       t.Name,
				Value:      t.Value,
				ToolboxId:  t.ToolboxId,
			})
			if len(find) > 0 {
			} else {
				t.UserId = requestBean.JWT.UserId
				err = this_.toolboxService.SaveExtend(t)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (this_ *api) exportFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	this_.toolboxService.Logger.Info("下载执行记录 start")
	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	param := map[string]string{}

	err = c.Bind(&param)
	if err != nil {
		return
	}
	request := &BaseRequest{}
	request.Key = param["key"]
	request.UserSetting = param["userSetting"] == "true"
	request.Toolbox = param["toolbox"] == "true"

	info, content, err := this_.genContent(requestBean.JWT.UserId, request)
	if err != nil {
		return
	}

	fileName := "" + info.Explain + "-" + info.CreateBy + "-" + info.CreateAt + ".yaml"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileName)))

	// 此处不设置 文件大小，如果设置文件大小，将无法终止下载
	//c.Header("Content-Length", fmt.Sprint(fileInfo.Size))
	c.Header("download-file-name", fileName)

	_, err = c.Writer.WriteString(content)
	c.Status(http.StatusOK)
	return
}
func (this_ *api) genContent(userId int64, r *BaseRequest) (info *SyncInfo, content string, err error) {
	user, err := this_.userService.Get(userId)
	if err != nil {
		return
	}
	info = &SyncInfo{
		Explain:  "Team IDE 配置文件",
		CreateBy: user.Name,
		CreateAt: util.GetNowFormat(),
	}
	if r.UserSetting {
		info.UserSetting = make(map[string]any)
		info.UserSetting["option"], err = this_.userSettingService.Query(userId)
		if err != nil {
			return
		}
	}
	var bs []byte

	if r.Toolbox {
		var groups []*module_toolbox.ToolboxGroupModel
		if groups, err = this_.toolboxService.QueryGroup(&module_toolbox.ToolboxGroupModel{
			UserId: userId,
		}); err != nil {
			return
		}
		for _, d := range groups {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			one := map[string]any{}
			err = json.Unmarshal(bs, &one)
			if err != nil {
				return
			}
			info.ToolboxGroupList = append(info.ToolboxGroupList, one)
		}

		var list []*module_toolbox.ToolboxModel
		if list, err = this_.toolboxService.Query(&module_toolbox.ToolboxModel{
			UserId: userId,
		}); err != nil {
			return
		}
		for _, d := range list {
			_ = this_.toolboxService.FormatOption(d, true)
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			one := map[string]any{}
			err = json.Unmarshal(bs, &one)
			if err != nil {
				return
			}
			info.ToolboxList = append(info.ToolboxList, one)
		}

		var extendList []*module_toolbox.ToolboxExtendModel
		if extendList, err = this_.toolboxService.QueryExtends(&module_toolbox.ToolboxExtendModel{
			UserId: userId,
		}); err != nil {
			return
		}
		for _, d := range extendList {
			bs, err = json.Marshal(d)
			if err != nil {
				return
			}
			one := map[string]any{}
			err = json.Unmarshal(bs, &one)
			if err != nil {
				return
			}
			info.ToolboxExtendList = append(info.ToolboxExtendList, one)
		}
	}

	content, err = Gen(r.Key, info)
	if err != nil {
		return
	}
	return
}
