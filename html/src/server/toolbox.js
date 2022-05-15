import http from '@/server/http';

let toolbox = {
    data(param) {
        return http.post('api/toolbox', param);
    },
    page(param) {
        return http.post('api/toolbox/page', param);
    },
    list(param) {
        return http.post('api/toolbox/list', param);
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
    delete(param) {
        return http.post('api/toolbox/delete', param);
    },
    context(param) {
        return http.post('api/toolbox/context', param);
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
};


export default toolbox;