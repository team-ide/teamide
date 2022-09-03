import http from '@/server/http';

import toolbox from "./toolbox.js";
import node from "./node.js";
import fileManager from "./fileManager.js";

import tool from '../tool/index.js';
import source from '../source/index.js';
let server = {
    toolbox,
    node,
    fileManager,
    data(param) {
        param = param || {};
        param.origin = location.origin;
        param.pathname = location.pathname;
        return http.post('api/data', param,);
    },
    session(param) {
        return http.post('api/session', param);
    },
    login(param) {
        return http.post('api/login', param,);
    },
    logout(param) {
        return http.post('api/logout', param,);
    },
    register(param) {
        return http.post('api/register', param,);
    },
    upload(param) {
        return http.post('api/upload', param, { headers: { 'Content-Type': 'multipart/form-data' } });
    },
    download(param) {
        return http.post('api/download', param, { responseType: "blob" });
    },
    updateCheck(param) {
        return http.post('api/updateCheck', param,);
    },
    openWebsocket() {
        if (serverSocket != null) { return }
        if (source.login.user == null) { return; }
        let url = source.api;
        url = url.substring(url.indexOf(":"));
        url = "ws" + url + "websocket";
        url += "?id=" + tool.md5("serverSocket:" + new Date().getTime());
        url += "&jwt=" + encodeURIComponent(tool.getJWT());
        serverSocket = new WebSocket(url);
        serverSocket.onopen = () => {
            serverSocketIsOpen = true;
            serverSocketOnOpen();
        };
        serverSocket.onmessage = (event) => {
            serverSocketOnMessage(event)
        };
        serverSocket.onclose = () => {
            serverSocketIsOpen = false;
            serverSocket = null;
            serverSocketOnClose();
            server.openWebsocket()
        };
        serverSocket.onerror = () => {
            serverSocketOnError();
        };

    },
    closeWebsocket() {
        if (serverSocket == null) { return }
        serverSocket.close()
    },
    addServerSocketOnEvent(event, call) {
        if (call == null) { return }
        serverSocketOnEventList[event] = serverSocketOnEventList[event] || []
        serverSocketOnEventList[event].push(call)
    },
    removeServerSocketOnEvent(event, call) {
        if (call == null) { return }
        let list = serverSocketOnEventList[event] || []
        let newList = [];
        list.forEach(one => {
            if (one != call) {
                newList.push(one)
            }
        })
        serverSocketOnEventList[event] = newList
    },
    addServerSocketOnOpen(call) {
        if (call == null) { return }
        if (serverSocketIsOpen) {
            call()
        }
        serverSocketOnOpenList.push(call)
    },
};

let serverSocket = null;
let serverSocketIsOpen = false;
const serverSocketOnOpenList = []
const serverSocketOnEventList = {}
const serverSocketOnCloseList = []
const getServerSocketOnEventList = (event) => {
    return serverSocketOnEventList[event] || []
}

const serverSocketOnOpen = () => {
    serverSocketOnOpenList.forEach(one => {
        if (one == null) { return }
        one()
    })
}

const serverSocketOnClose = () => {
    serverSocketOnCloseList.forEach(one => {
        if (one == null) { return }
        one()
    })
}

const serverSocketOnError = () => {

}

const serverSocketOnMessage = (event) => {
    let message = event.data;
    let json = JSON.parse(message)
    if (json == null) { return }
    if (json.isMessage) {
        if (tool.isNotEmpty(json.errorMessage)) {
            tool.error(json.errorMessage)
        }
        if (tool.isNotEmpty(json.warnMessage)) {
            tool.warn(json.warnMessage)
        }
        if (tool.isNotEmpty(json.infoMessage)) {
            tool.info(json.infoMessage)
        }
        if (tool.isNotEmpty(json.successMessage)) {
            tool.success(json.successMessage)
        }
    }
    if (json.isEvent) {
        let list = getServerSocketOnEventList(json.event,)
        list.forEach(one => {
            if (one == null) { return }
            one(json.data)
        })
    }
}

export default server;