import server from "../server";
import tool from "../tool";
import form from "../form";

let toolbox = {
    async initContext() {
        let res = await server.toolbox.data();
        if (res.code != 0) {
            tool.error(res.msg);
        } else {
            let data = res.data || {};
            toolbox.quickCommandTypes = data.quickCommandTypes;
            toolbox.databaseTypes = data.databaseTypes||[];
            toolbox.types = data.types;
            data.types.forEach((one) => {
                form.toolbox[one.name] = one.configForm;
                if (one.otherForm) {
                    for (let formName in one.otherForm) {
                        form.toolbox[one.name][formName] = one.otherForm[formName];
                    }
                }
            });
            toolbox.sqlConditionalOperations =
                data.sqlConditionalOperations;
        }

        let param = {};
        res = await server.toolbox.context(param);
        if (res.code != 0) {
            tool.error(res.msg);
        } else {
            let context = res.data.context || {};
            let groups = res.data.groups || [];
            toolbox.groups = groups;
            toolbox.context = context;
        }
    },

    getToolboxData(toolboxData) {
        let res = null;
        if (toolbox.context) {
            for (let type in toolbox.context) {
                if (toolbox.context[type] == null) {
                    continue;
                }
                toolbox.context[type].forEach((one) => {
                    if (
                        one == toolboxData ||
                        one.toolboxId == toolboxData ||
                        one.toolboxId == toolboxData.toolboxId
                    ) {
                        res = one;
                    }
                });
            }
        }
        return res;
    },
    getToolboxType(type) {
        let res = null;
        if (toolbox.types) {
            toolbox.types.forEach((one) => {
                if (one == type || one.name == type || one.name == type.name) {
                    res = one;
                }
            });
        }
        return res;
    },
    getQuickCommandType(name) {
        if (toolbox.quickCommandTypes == null) {
            return null;
        }
        let res = null;
        toolbox.quickCommandTypes.forEach((one) => {
            if (one.name == name) {
                res = one;
            }
        });
        return res;
    },
    getOptionJSON(option) {
        let json = {};
        if (tool.isNotEmpty(option)) {
            json = JSON.parse(option);
        }
        return json;
    },
    async open(toolboxId, extend, createTime) {
        let extendStr = null;
        if (extend != null) {
            extendStr = JSON.stringify(extend);
        }
        let param = {
            toolboxId: toolboxId,
            extend: extendStr,
        };
        if (createTime) {
            param.createTime = createTime;
        }
        let res = await server.toolbox.open(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
        let openData = res.data.open;
        return openData
    },
    createToolboxDataTab(openData) {
        let toolboxType = openData.toolboxType;
        let toolboxData = openData.toolboxData;
        let extend = openData.extend;
        let title = toolboxType.text + ":" + toolboxData.name;
        let name = toolboxData.name;

        extend = extend || {};
        this.formatExtend(toolboxType, toolboxData, extend);
        if (extend.isFTP) {
            title = "FTP:" + toolboxData.name;
            name = toolboxData.name;
        } else if (toolboxType.name == "ssh") {
            title = "SSH:" + toolboxData.name;
            name = toolboxData.name;
        }
        let tab = {
            toolboxData,
            title,
            name,
            toolboxType,
            extend,
            openId: openData.openId,
            openData: openData,
            comment: "",
        };
        tab.active = false;
        return tab;
    },
    async loadOpens() {
        let param = {};
        let res = await server.toolbox.queryOpens(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
        let opens = res.data.opens || [];
        return opens;
    },
    async closeOpen(openId) {
        let param = {
            openId: openId,
        };
        let res = await server.toolbox.close(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
    },
    async activeOpen(openId) {
        let param = {
            openId: openId,
        };
        let res = await server.toolbox.open(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
    },

    async openTab(openData, extend) {
        let extendStr = null;
        if (extend != null) {
            extendStr = JSON.stringify(extend);
        }
        let param = {
            openId: openData.openId,
            toolboxId: openData.toolboxId,
            extend: extendStr,
        };
        let res = await server.toolbox.openTab(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
        let tabData = res.data.tab;
        return tabData
    },
    createOpenTabTab(tabData) {
        let extend = tabData.extend;
        let title = extend.title;
        let name = extend.name;

        let tab = {
            tabData,
            title,
            name,
            extend,
            tabId: tabData.tabId,
            openId: tabData.openId,
            comment: "",
        };
        tab.active = false;
        return tab;
    },
    async loadOpenTabs(openId) {
        let param = {
            openId: openId,
        };
        let res = await server.toolbox.queryOpenTabs(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
        let tabs = res.data.tabs || [];
        return tabs;
    },
    async closeOpenTab(tabId) {
        let param = {
            tabId: tabId,
        };
        let res = await server.toolbox.closeTab(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
    },
    async activeOpenTab(tabId) {
        let param = {
            tabId: tabId,
        };
        let res = await server.toolbox.openTab(param);
        if (res.code != 0) {
            tool.error(res.msg);
        }
    },
};

toolbox.toolboxTypeOpens = [];

toolbox.tabs = [];
toolbox.activeTab = null;
toolbox.context = null;
toolbox.groups = null;
toolbox.quickCommands = null;
toolbox.quickCommandSSHCommands = null;

toolbox.formatExtend = (toolboxType, data, extend) => {
    if (toolboxType.name == "ssh") {
        toolbox.formatSSHExtend(toolboxType, data, extend);
    }
};
toolbox.formatSSHExtend = (toolboxType, data, extend) => {
    extend.local = extend.local || {};
    extend.remote = extend.remote || {};
    extend.local.dir = extend.local.dir || "";
    extend.remote.dir = extend.remote.dir || "";
};

toolbox.databaseTypes = null;
toolbox.sqlConditionalOperations = null;
toolbox.types = null;
toolbox.mysqlColumnTypeInfos = null;

export default toolbox;