import server from "@/server/index.js";
import tool from "@/tool/index.js";
const newToolboxWorker = function (workerOption) {
    workerOption = workerOption || {};
    const worker = {
        toolboxId: workerOption.toolboxId,
        openId: workerOption.openId,
        toolboxType: workerOption.toolboxType,
        extend: workerOption.extend,

        async work(work, data) {
            let param = {
                toolboxId: worker.toolboxId,
                work: work,
                data: data,
            };
            let res = await server.toolbox.work(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            return res;
        },
        onRemoveTab(tab) {
            this.toolbox.closeOpenTab(tab.tabId);
        },
        onActiveTab(tab) {
            this.toolbox.activeOpenTab(tab.tabId);
        },
        doActiveTab(tab) {
            this.wrap.doActiveTab(tab);
        },
        updateExtend(keyValueMap) {
            this.updateOpenExtend(this.openId, keyValueMap);
        },
        updateComment(comment) {
            this.updateOpenComment(this.openId, comment);
        },
        async updateOpenTabExtend(tabId, keyValueMap) {
            let tab = this.wrap.getTab("" + tabId);
            if (tab == null) {
                return;
            }
            if (keyValueMap == null) {
                return;
            }
            if (Object.keys(keyValueMap) == 0) {
                return;
            }
            let obj = tab.extend;
            for (let key in keyValueMap) {
                let value = keyValueMap[key];
                let names = key.split(".");
                names.forEach((name, index) => {
                    if (index < names.length - 1) {
                        obj[name] = obj[name] || {};
                        obj = obj[name];
                    } else {
                        obj[name] = value;
                    }
                });
            }

            let param = {
                tabId: tabId,
                extend: JSON.stringify(tab.extend),
            };
            let res = await this.server.toolbox.updateOpenTabExtend(param);
            if (res.code != 0) {
                this.tool.error(res.msg);
            }
        },
        async updateOpenExtend(keyValueMap) {
            if (keyValueMap == null) {
                return;
            }
            if (Object.keys(keyValueMap) == 0) {
                return;
            }
            let obj = worker.extend;
            for (let key in keyValueMap) {
                let value = keyValueMap[key];
                let names = key.split(".");
                names.forEach((name, index) => {
                    if (index < names.length - 1) {
                        obj[name] = obj[name] || {};
                        obj = obj[name];
                    } else {
                        obj[name] = value;
                    }
                });
            }
            let param = {
                openId: worker.openId,
                extend: JSON.stringify(worker.extend),
            };
            let res = await server.toolbox.updateOpenExtend(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
        },
    };


    return worker
};
export default {
    newToolboxWorker
}