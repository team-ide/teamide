package module_sync

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"testing"
)

func TestRead(t *testing.T) {
	str := `
######################
说明: Team IDE 配置文件
所属: 朱亮
加密: 是
时间: 2024-09-13 13:33
密串: 214588ECADE14CAAA04761F788BABD70
签名: 30ED5E981E6F31D2B8397353958A2DBF
######################

######################
个人设置: |
  cNBG8jdhrVKh6Jxuoc7zS9wSXrzAfQLjZkT4Kj9Vc60NvpypajaxMwnceWyXZD5z

######################

######################
工具分组: |
  rL0DsCKDJUcOUraxhrl2LPteu9JAeUuz4c9/nQScjzE=
  99R8tKTji5Lh2kZbPlPt8gGoWbiviWFo0AUEGQXcvRw=

######################

######################
工具: |
  yeSJ0hMzv9dlk+mwixR3i7wPyfsLJEeqNHKfXrOkWVEmtesKXl1z771XtdfjEGEz+nvcOhxA37GxMqhfVuawKw==
  pobWX7wsgGn5RqwCwIIhGkGnLC5MQOORBQG0xep9UvluBIZFjR8vw5ZqFyToviYttnrOOBOXVFuNHzy/FnQrpw==

######################
`

	info, err := Read("123456", str)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(util.GetStringValue(info))

}

func TestWrite(t *testing.T) {
	info := &SyncInfo{
		Explain:  "Team IDE 配置文件",
		CreateBy: "朱亮",
		CreateAt: "2024-09-13 13:33",
		UserSetting: map[string]any{
			"option": map[string]any{
				"name1": "xx",
				"name2": "xxa",
			},
		},
		ToolboxGroupList: []map[string]any{
			{"name": "分组1", "groupId": 11},
			{"name": "分组2", "groupId": 33},
		},
		ToolboxList: []map[string]any{
			{"name": "工具1", "option": `{"address":"xxa"}`, "groupId": 11},
			{"name": "工具2", "option": `{"address":"xxa"}`, "groupId": 22},
		},
	}
	content, err := Gen("123456", info)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(content)

}
