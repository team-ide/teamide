package javascript

import (
	"teamide/pkg/maker"
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

type appCoder struct {
	*maker.Application
	codeType common.CodeType
	*modelers.LanguageJavascriptModel
}
