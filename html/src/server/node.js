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
    delete(param) {
        return http.post('api/node/delete', param);
    },
};


export default node;