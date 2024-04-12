package javascript

import (
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

type appCoder struct {
	*modelers.Application
	codeType common.CodeType
	*modelers.LanguageJavascriptModel
}
