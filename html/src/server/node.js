import http from '@/server/http';

let node = {
    data(param) {
        return http.post('api/node', param);
    },
    context(param) {
        return http.post('api/node/context', param);
    },
    page(param) {
        return http.post('api/node/page', param);
    },
    list(param) {
        return http.post('api/node/list', param);
    },
    insert(param) {
        return http.post('api/node/insert', param);
    },
    update(param) {
        return http.post('api/node/update', param);
    },
    updateOption(param) {
        return http.post('api/node/updateOption', param);
    },
    enable(param) {
        return http.post('api/node/enable', param);
    },
    disable(param) {
        return http.post('api/node/disable', param);
    },
    delete(param) {
        return http.post('api/node/delete', param);
    },
    netProxy: {
        insert(param) {
            return http.post('api/node/netProxy/insert', param);
        },
        update(param) {
            return http.post('api/node/netProxy/update', param);
        },
        updateOption(param) {
            return http.post('api/node/netProxy/updateOption', param);
        },
        monitorData(param) {
            return http.post('api/node/netProxy/monitorData', param);
        },
        enable(param) {
            return http.post('api/node/netProxy/enable', param);
        },
        disable(param) {
            return http.post('api/node/netProxy/disable', param);
        },
        delete(param) {
            return http.post('api/node/netProxy/delete', param);
        },
    },
    system: {
        info(param) {
            return http.post('api/node/system/info', param);
        },
        queryMonitorData(param) {
            return http.post('api/node/system/queryMonitorData', param);
        },
        cleanMonitorData(param) {
            return http.post('api/node/system/cleanMonitorData', param);
        },
    },
};


export default node;