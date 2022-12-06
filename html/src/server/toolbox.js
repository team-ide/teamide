import http from '@/server/http';

let toolbox = {
    list(param) {
        return http.post('api/toolbox/list', param);
    },
    count(param) {
        return http.post('api/toolbox/count', param);
    },
    insert(param) {
        return http.post('api/toolbox/insert', param);
    },
    update(param) {
        return http.post('api/toolbox/update', param);
    },
    rename(param) {
        return http.post('api/toolbox/rename', param);
    },
    moveGroup(param) {
        return http.post('api/toolbox/moveGroup', param);
    },
    delete(param) {
        return http.post('api/toolbox/delete', param);
    },
    group: {
        list(param) {
            return http.post('api/toolbox/group/list', param);
        },
        insert(param) {
            return http.post('api/toolbox/group/insert', param);
        },
        update(param) {
            return http.post('api/toolbox/group/update', param);
        },
        delete(param) {
            return http.post('api/toolbox/group/delete', param);
        },
    },
    quickCommand: {
        query(param) {
            return http.post('api/toolbox/quickCommand/query', param);
        },
        insert(param) {
            return http.post('api/toolbox/quickCommand/insert', param);
        },
        update(param) {
            return http.post('api/toolbox/quickCommand/update', param);
        },
        delete(param) {
            return http.post('api/toolbox/quickCommand/delete', param);
        },
    },
    work(param) {
        return http.post('api/toolbox/work', param);
    },
    open(param) {
        return http.post('api/toolbox/open', param);
    },
    close(param) {
        return http.post('api/toolbox/close', param);
    },
    getOpen(param) {
        return http.post('api/toolbox/getOpen', param);
    },
    queryOpens(param) {
        return http.post('api/toolbox/queryOpens', param);
    },
    updateOpenExtend(param) {
        return http.post('api/toolbox/updateOpenExtend', param);
    },
    openTab(param) {
        return http.post('api/toolbox/openTab', param);
    },
    closeTab(param) {
        return http.post('api/toolbox/closeTab', param);
    },
    queryOpenTabs(param) {
        return http.post('api/toolbox/queryOpenTabs', param);
    },
    updateOpenTabExtend(param) {
        return http.post('api/toolbox/updateOpenTabExtend', param);
    },
    ssh: {
        ftp: {
            upload(param) {
                return http.post('api/toolbox/ssh/ftp/upload', param, { headers: { 'Content-Type': 'multipart/form-data' } });
            },
            download(param) {
                return http.post('api/toolbox/ssh/ftp/download', param, { responseType: "blob" });
            },
        },
    },
    database: {
        upload(param) {
            return http.post('api/toolbox/ssh/ftp/upload', param, { headers: { 'Content-Type': 'multipart/form-data' } });
        },
        export: {
            download(param) {
                return http.post('api/toolbox/database/export/download', param, { responseType: "blob" });
            },
        },
    }
};


export default toolbox;