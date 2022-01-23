
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
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "表", name: "table" },
            { text: "父", name: "parent" },
            {
                text: "字段", name: "fields", isList: true,
                fields: [
                    { text: "名称", name: "name" },
                    { text: "注释", name: "comment" },
                    { text: "数据类型", name: "dataType" },
                    { text: "是List", name: "isList" },
                    { text: "字段", name: "column" },
                    { text: "字段类型", name: "columnType" },
                    { text: "字段长度", name: "columnLength" },
                    { text: "小数长度", name: "columnDecimal" },
                    { text: "是主键", name: "primaryKey" },
                    { text: "不能为空", name: "notNull" },
                    { text: "默认", name: "default" },
                ]
            },
            {
                text: "索引", name: "indexs", isList: true,
                fields: [
                    { text: "名称", name: "name" },
                    { text: "注释", name: "comment" },
                    { text: "类型", name: "type" },
                    { text: "字段", name: "columns" },
                ]
            },
        ]
    },
    {
        name: "services", text: "服务接口",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
        ],
    },
    {
        name: "constants", text: "常量",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "数据类型", name: "dataType" },
            { text: "值", name: "value" },
            { text: "环境变量", name: "environmentVariable", comment: "优先取环境变量中的值", },
        ],
    },
    {
        name: "errors", text: "错误码",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "错误码", name: "code" },
            { text: "错误信息", name: "msg" },
        ],
    },
    {
        name: "tests", text: "测试",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
        ],
    },
    {
        name: "dictionaries", text: "数据字典",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            {
                text: "选项", name: "options", isList: true,
                fields: [
                    { text: "文案", name: "text" },
                    { text: "值", name: "value" },
                    { text: "数据类型", name: "dataType" },
                    { text: "注释", name: "comment" },
                ]
            },
        ],
    },
    {
        name: "serverWebs", text: "Web服务",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "Host", name: "host" },
            { text: "Port", name: "port" },
            { text: "ContextPath", name: "contextPath" },
            {
                text: "Token", name: "token",
                fields: [
                    { text: "验证路径", name: "include" },
                    { text: "忽略路径", name: "exclude" },
                    { text: "创建Token服务", name: "createService" },
                    { text: "验证Token服务", name: "validateService" },
                    { text: "变量名称", name: "variableName" },
                    { text: "变量数据类型", name: "variableDataType" },
                ]
            },
        ],
    },
    {
        name: "datasourceDatabases", text: "Database数据源",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "类型", name: "type" },
            { text: "Host", name: "host" },
            { text: "Port", name: "port" },
            { text: "Database", name: "database" },
            { text: "Username", name: "username" },
            { text: "Password", name: "password" },
        ],
    },
    {
        name: "datasourceRedises", text: "Redis数据源",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "Redis地址", name: "address" },
            { text: "密码", name: "auth" },
            { text: "前缀", name: "prefix", comment: "如果配置，所有key将自动拼接该前缀", },
        ],
    },
    {
        name: "datasourceKafkas", text: "Kafka数据源",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "Kafka地址", name: "address" },
            { text: "前缀", name: "prefix", comment: "如果配置，所有topic将自动拼接该前缀", },
        ],
    },
    {
        name: "datasourceZookeepers", text: "Zookeeper数据源",
        fields: [
            { text: "名称", name: "name" },
            { text: "注释", name: "comment" },
            { text: "Zookeeper地址", name: "address" },
            { text: "命名空间", name: "namespace", comment: "如果配置，则所有路径将放在该命名空间下", },
        ],
    },
];

export default application;