import server from "@/server/index.js";
import tool from "@/tool/index.js";
import source from "@/source/index.js";

const newWorker = function (workerOption) {
    workerOption = workerOption || {};
    const worker = {
        workerId: workerOption.workerId,
        place: workerOption.place,
        placeId: workerOption.placeId,
        key: workerOption.key,
        onSocketOpen: workerOption.onSocketOpen,
        onSocketClose: workerOption.onSocketClose,
        onSocketError: workerOption.onSocketError,
        init() {
            this.build()
        },
        refresh() {
            this.build()
        },
        async build() {
            if (this.building) {
                return
            }
            this.building = true
            if (tool.isNotEmpty(this.key)) {
                if (worker.socket != null) {
                    worker.socket.close();
                }
                await this.close();
            }
            this.key = await this.newKey()
            this.newSocket();
        },
        writeData(data) {
            worker.socket.send(data);
        },
        newSocket() {
            if (worker.socket != null) {
                worker.socket.close();
            }
            let url = source.api;
            url = url.substring(url.indexOf(":"));
            url = "ws" + url + "api/terminal/websocket";
            url += "?key=" + encodeURIComponent(worker.key);
            url += "&jwt=" + encodeURIComponent(tool.getJWT());
            url += "&place=" + encodeURIComponent(worker.place);
            url += "&placeId=" + encodeURIComponent(worker.placeId);
            url += "&workerId=" + encodeURIComponent(worker.workerId);
            let socket = new WebSocket(url);
            worker.socket = socket;
            worker.socket.binaryType = "arraybuffer";
            worker.socket.onopen = () => {
                worker.onSocketOpen();
                this.building = false
            };
            socket.onmessage = (event) => {
                let data = event.data;
                worker.onSocketData(data);
            };
            socket.onclose = () => {
                worker.socket = null;
                worker.onSocketClose();
            };
            socket.onerror = () => {
                worker.onSocketError();
            };
        },
        getParam() {
            let param = {
                workerId: worker.workerId,
                place: worker.place,
                placeId: worker.placeId,
                key: worker.key,
            };
            return param;
        },
        async newKey() {
            let param = worker.getParam();
            let res = await server.terminal.key(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            return res.data;
        },
        async close() {
            let param = worker.getParam();
            let res = await server.terminal.close(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            return res.data;
        },
    };


    return worker
};
export default {
    newWorker
}