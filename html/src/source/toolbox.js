import server from "../server";
import tool from "../tool";

let toolbox = {
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

export default toolbox;