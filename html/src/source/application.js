import tool from "../tool";

let application = {};


application.apps = null;
application.app = null;
application.context = null;

application.groupOpens = [];

application.tabs = [];
application.activeTab = null;

application.groups = [
    {
        name: "structs", text: "结构体",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "表", name: "table" },
            { text: "父", name: "parent", type: "select", isStructOption: true, },
            {
                text: "字段", name: "fields", isList: true,
                fields: [
                    { text: "名称", name: "name" },
                    { text: "注释", name: "comment" },
                    { text: "数据类型", name: "dataType", type: "select", isDataTypeOption: true, },
                    { text: "是List", name: "isList", type: "switch" },
                    { text: "字段", name: "column", ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "字段类型", name: "columnType", type: "select", isColumnTypeOption: true, ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "字段长度", name: "columnLength", isNumber: true, ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "小数长度", name: "columnDecimal", isNumber: true, ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "是主键", name: "primaryKey", type: "switch", ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "不能为空", name: "notNull", type: "switch", ifScript(data) { return tool.isNotEmpty(data.table) }, },
                    { text: "默认", name: "default", ifScript(data) { return tool.isNotEmpty(data.table) }, },
                ]
            },
            {
                text: "索引", name: "indexs", isList: true,
                fields: [
                    { text: "名称", name: "name" },
                    { text: "注释", name: "comment" },
                    { text: "类型", name: "type", type: "select", isIndexTypeOption: true, },
                    { text: "字段", name: "columns" },
                ]
            },
        ]
    },
    {
        name: "actions", text: "服务接口", isAction: true,
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
        ],
    },
    {
        name: "constants", text: "常量",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "数据类型", name: "dataType", type: "select", isDataTypeOption: true, },
            { text: "值", name: "value" },
            { text: "环境变量", name: "environmentVariable", comment: "优先取环境变量中的值", },
        ],
    },
    {
        name: "errors", text: "错误码",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "错误码", name: "code" },
            { text: "错误信息", name: "msg" },
        ],
    },
    {
        name: "tests", text: "测试",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
        ],
    },
    {
        name: "dictionaries", text: "数据字典",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            {
                text: "选项", name: "options", isList: true,
                fields: [
                    { text: "文案", name: "text" },
                    { text: "值", name: "value" },
                    { text: "数据类型", name: "dataType", type: "select", isDataTypeOption: true, },
                    { text: "注释", name: "comment" },
                ]
            },
        ],
    },
    {
        name: "serverWebs", text: "Web服务",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "Host", name: "host" },
            { text: "Port", name: "port" },
            { text: "ContextPath", name: "contextPath" },
            {
                text: "Token", name: "token",
                fields: [
                    { text: "验证路径", name: "include" },
                    { text: "忽略路径", name: "exclude" },
                    { text: "创建Token服务", name: "createAction", type: "select", isActionOption: true, },
                    { text: "验证Token服务", name: "validateAction", type: "select", isActionOption: true, },
                    { text: "变量名称", name: "variableName" },
                    { text: "变量数据类型", name: "variableDataType", type: "select", isDataTypeOption: true, },
                ]
            },
        ],
    },
    {
        name: "datasourceDatabases", text: "Database数据源",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "类型", name: "type", type: "select", isDatabaseTypeOption: true, },
            { text: "Host", name: "host" },
            { text: "Port", name: "port", isNumber: true, },
            { text: "Database", name: "database" },
            { text: "Username", name: "username" },
            { text: "Password", name: "password" },
        ],
    },
    {
        name: "datasourceRedises", text: "Redis数据源",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "Redis地址", name: "address" },
            { text: "密码", name: "auth" },
            { text: "前缀", name: "prefix", comment: "如果配置，所有key将自动拼接该前缀", },
        ],
    },
    {
        name: "datasourceKafkas", text: "Kafka数据源",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "Kafka地址", name: "address" },
            { text: "前缀", name: "prefix", comment: "如果配置，所有topic将自动拼接该前缀", },
        ],
    },
    {
        name: "datasourceZookeepers", text: "Zookeeper数据源",
        fields: [
            { text: "名称", name: "name", readonly: true, },
            { text: "注释", name: "comment" },
            { text: "Zookeeper地址", name: "address" },
            { text: "命名空间", name: "namespace", comment: "如果配置，则所有路径将放在该命名空间下", },
        ],
    },
];

var databaseTypeOptions = [];
databaseTypeOptions.push({ value: "Mysql", text: "MySql" });
application.getDatabaseTypeOptions = function (context) {
    return databaseTypeOptions;

};

var columnTypeOptions = [];
columnTypeOptions.push({ value: "varchar", text: "varchar" });
columnTypeOptions.push({ value: "bigint", text: "bigint" });
columnTypeOptions.push({ value: "int", text: "int" });
columnTypeOptions.push({ value: "datetime", text: "datetime" });
columnTypeOptions.push({ value: "number", text: "number" });
application.getColumnTypeOptions = function (context) {
    return columnTypeOptions;
};

var indexTypeOptions = [];
indexTypeOptions.push({ value: "unique", text: "普通索引" });
application.getIndexTypeOptions = function (context) {
    return indexTypeOptions;
};

var dataTypeOptions = [];
dataTypeOptions.push({ value: "string", text: "字符串" });
dataTypeOptions.push({ value: "int", text: "整形" });
dataTypeOptions.push({ value: "long", text: "长整型" });
dataTypeOptions.push({ value: "boolean", text: "布尔型" });
dataTypeOptions.push({ value: "byte", text: "字节型" });
dataTypeOptions.push({ value: "date", text: "日期" });
dataTypeOptions.push({ value: "short", text: "短整型" });
dataTypeOptions.push({ value: "double", text: "双精度浮点型" });
dataTypeOptions.push({ value: "float", text: "浮点型" });
dataTypeOptions.push({ value: "map", text: "集合" });
application.getDataTypeOptions = function (context) {
    let options = [];
    dataTypeOptions.forEach(one => {
        options.push(one);
    });
    if (context && context.structs) {
        context.structs.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;

};

application.getStructOptions = function (context) {
    let options = [];
    if (context && context.structs) {
        context.structs.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;
};

application.getActionOptions = function (context) {
    let options = [];
    if (context && context.actions) {
        context.actions.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;

};
export default application;