import server from "@/server/index.js";
import tool from "@/tool/index.js";
const newToolboxWorker = function (workerOption) {
    workerOption = workerOption || {};
    let itemsWorker = tool.newItemsWorker({
        onRemoveItem() {

        },
        onActiveItem() {

        },
    });
    const worker = {
        toolboxId: workerOption.toolboxId,
        openId: workerOption.openId,
        toolboxType: workerOption.toolboxType,
        extend: workerOption.extend,
        itemsWorker: itemsWorker,

        async init() {
            await this.initOpenTabs()
        },
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
        async updateExtend(keyValueMap) {
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
        async openTabByExtend(extend, fromItem) {
            let param = {
                openId: worker.openId,
                toolboxId: worker.toolboxId,
                extend: JSON.stringify(extend || {}),
            };
            let res = await server.toolbox.openTab(param);
            let openTab = null;
            if (res.code != 0) {
                tool.error(res.msg);
            } else {
                openTab = res.data.tab;
            }
            if (openTab == null) {
                return;
            }
            let item = await this.addOpenItem(openTab, fromItem);
            if (item != null) {
                itemsWorker.toActiveItem(item);
            }
        },
        async addOpenItem(data, fromItem) {
            let extend = tool.getOptionJSON(data.extend);

            let item = {
                key: data.tabId,
                title: extend.title,
                name: extend.name,
                extend: extend,
                tabId: data.tabId,
                openId: worker.openId,
                comment: "",
            };

            itemsWorker.addItem(item, fromItem);

            data.item = item;
            return item;
        },
        async initOpenTabs() {
            let openTabs = [];
            let res = await server.toolbox.queryOpenTabs({
                openId: worker.openId,
            });
            if (res.code != 0) {
                tool.error(res.msg);
            } else {
                openTabs = res.data.tabs || [];
            }

            await openTabs.forEach(async (one) => {
                await worker.addOpenItem(one);
            });

            // 激活最后
            let activeOne = null;
            openTabs.forEach(async (one) => {
                if (activeOne == null) {
                    activeOne = one;
                } else {
                    if (
                        new Date(one.openTime).getTime() >
                        new Date(activeOne.openTime).getTime()
                    ) {
                        activeOne = one;
                    }
                }
            });
            if (activeOne != null) {
                itemsWorker.toActiveItem(activeOne.item);
            }
        },
    };


    return worker
};
export default {
    newToolboxWorker
}