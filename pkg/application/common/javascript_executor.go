package common

type IJavascriptExecutor interface {
	ExecuteExpressionScript(app IApplication, invokeNamespace *InvokeNamespace, expressionScript string) (res interface{}, err error)
	ExecuteStatementScript(app IApplication, invokeNamespace *InvokeNamespace, statementScript string) (res interface{}, err error)
	ExecuteFunctionScript(app IApplication, invokeNamespace *InvokeNamespace, functionScript string) (res interface{}, err error)
}
