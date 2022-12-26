import http from '@/server/http';

import toolbox from "./toolbox.js";
import node from "./node.js";
import fileManager from "./fileManager.js";
import terminal from "./terminal.js";
import user from "./user.js";

import database from "./database.js";
import elasticsearch from "./elasticsearch.js";
import kafka from "./kafka.js";
import zookeeper from "./zookeeper.js";
import redis from "./redis.js";

import log from "./log.js";
import power from "./power.js";

import tool from '../tool/index.js';
import source from '../source/index.js';
let server = {
    toolbox,
    node,
    fileManager,
    terminal,
    user,
    database,
    elasticsearch,
    kafka,
    zookeeper,
    redis,
    log,
    power,
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
    listen(param) {
        return http.post('api/listen', param || {},);
    },
    addListenOnEvent(event, call) {
        if (call == null) { return }
        listenOnEventList[event] = listenOnEventList[event] || []
        listenOnEventList[event].push(call)
    },
    removeListenOnEvent(event, call) {
        if (call == null) { return }
        let list = listenOnEventList[event] || []
        let newList = [];
        list.forEach(one => {
            if (one != call) {
                newList.push(one)
            }
        })
        listenOnEventList[event] = newList
    },
};
var listenStartInt = false
const listenStart = async (errorCount) => {
    if (listenStartInt) {
        return
    }
    listenStartInt = true
    var isError = false
    var isDataError = false
    try {
        let res = await server.listen({})
        if (res.code != 0) {
            isDataError = true
        } else {
            let data = res.data || {}
            let events = data.events || []
            try {
                listenOnEvents(events)
            } catch (error) {

            }
        }
    } catch (error) {
        isError = true
    }
    if (isError || isDataError) {
        errorCount = errorCount || 0
        errorCount++
        window.setTimeout(() => {
            listenStartInt = false
            listenStart(errorCount)
        }, errorCount * 1000 * 5)
    } else {
        listenStartInt = false
        listenStart()
    }
}
server.listenStart = listenStart;

const listenOnEventList = {}
const getListenOnEventList = (event) => {
    return listenOnEventList[event] || []
}

const listenOnEvents = (events) => {
    events = events || []
    events.forEach(event => {
        event = event || {}
        let list = getListenOnEventList(event.event,)
        list.forEach(one => {
            if (one == null) { return }
            one(event.data)
        })
    })
}

export default server;