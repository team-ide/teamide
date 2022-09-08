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
        onSocketData: workerOption.onSocketData,
        rows: 40,
        cols: 100,
        socket: null,
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
            let keyData = await this.newKey()
            if (keyData == null || tool.isEmpty(keyData.key)) {
                this.building = false
                return
            }
            this.key = keyData.key;
            this.isWindows = keyData.isWindows;
            this.newSocket();
        },
        sendDataToWS(data) {
            if (this.isWindows) {
                // data = data.replace(/(\r\n|\n|\r|↵)/g, `\r\n`);
            }
            if (typeof data === "string") {
                worker.socket.send(data);
            } else {
                worker.socket.send(new Uint8Array(data));
            }
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
            url += "&cols=" + worker.cols;
            url += "&rows=" + worker.rows;
            let socket = new WebSocket(url);
            worker.socket = socket;
            if (!this.worker.isWindows) {
                worker.socket.binaryType = "arraybuffer";
            }
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