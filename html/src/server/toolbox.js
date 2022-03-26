import http from '@/server/http';

let toolbox = {
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
};


export default toolbox;