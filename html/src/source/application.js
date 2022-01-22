
let application = {};


application.apps = null;
application.app = null;
application.context = null;

application.groups = [
    { name: "serverWebs", text: "Web服务" },
    { name: "constants", text: "常量" },
    { name: "errors", text: "错误码" },
    { name: "structs", text: "结构体" },
    { name: "services", text: "服务接口" },
    { name: "tests", text: "测试" },
    { name: "dictionaries", text: "数据字典" },
    { name: "datasourceDatabases", text: "Database数据源" },
    { name: "datasourceRedises", text: "Redis数据源" },
    { name: "datasourceKafkas", text: "Kafka数据源" },
    { name: "datasourceZookeepers", text: "Zookeeper数据源" },
];

export default application;