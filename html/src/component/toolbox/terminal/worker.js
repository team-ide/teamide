import server from "@/server/index.js";
import tool from "@/tool/index.js";
import source from "@/source/index.js";
import { async } from "@antv/x6/lib/registry/marker/async";

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
        building: false,
        rows: 40,
        cols: 100,
        socket: null,
        uploadSocket: null,
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
                    worker.closeSocket();
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
        closeSocket() {
            if (worker.socket != null) {
                worker.socket.close();
            }
            if (worker.uploadSocket != null) {
                worker.uploadSocket.close();
            }
        },
        newSocket() {
            this.closeSocket();
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

            socket.binaryType = "arraybuffer"
            if (worker.isWindows) {
                // socket.binaryType = "string"
            } else {
                // socket.binaryType = "blob"

            }
            socket.onopen = () => {
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
            this.newUploadSocket();
        },
        wait(timeout) {
            return new Promise((resolve, reject) => {
                setTimeout(() => { resolve() }, timeout)
            });
        },
        async waitUploading() {
            if (worker.uploadSocketSending) {
                await this.wait(100);
            }
        },
        uploadSocketSend(data) {
            if (worker.uploadSocket == null) {
                throw new Error("uploadSocket is null.")
            }
            return new Promise(async (resolve, reject) => {
                await this.waitUploading();
                worker.uploadSocketSending = true;
                worker.uploadSocketSendCallback = () => {
                    delete worker.uploadSocketSendCallback
                    delete worker.uploadSocketSending;
                    resolve();
                }
                worker.uploadSocket.send(new Uint8Array(data));
            });
        },
        newUploadSocket() {
            let url = source.api;
            url = url.substring(url.indexOf(":"));
            url = "ws" + url + "api/terminal/uploadWebsocket";
            url += "?key=" + encodeURIComponent(worker.key);
            url += "&jwt=" + encodeURIComponent(tool.getJWT());
            let uploadSocket = new WebSocket(url);
            worker.uploadSocket = uploadSocket;
            uploadSocket.binaryType = "arraybuffer"
            uploadSocket.onopen = () => {
            };
            uploadSocket.onmessage = (event) => {
                if (worker.uploadSocketSendCallback) {
                    worker.uploadSocketSendCallback();
                }
            };
            uploadSocket.onclose = () => {
                worker.uploadSocket = null;
            };
            uploadSocket.onerror = () => {
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
        async changeSize() {
            let param = worker.getParam();
            param.cols = Number(worker.cols)
            param.rows = Number(worker.rows)
            let res = await server.terminal.changeSize(param);
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