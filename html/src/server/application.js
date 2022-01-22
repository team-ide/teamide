import http from '@/server/http';

let application = {
    page(param) {
        return http.post('api/application/page', param);
    },
    list(param) {
        return http.post('api/application/list', param);
    },
    insert(param) {
        return http.post('api/application/insert', param);
    },
    update(param) {
        return http.post('api/application/update', param);
    },
    rename(param) {
        return http.post('api/application/rename', param);
    },
    delete(param) {
        return http.post('api/application/delete', param);
    },
    context: {
        load(param) {
            return http.post('api/application/context/load', param);
        },
        save(param) {
            return http.post('api/application/context/save', param);
        },
    },
};


export default application;