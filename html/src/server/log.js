import http from '@/server/http';

let database = {
    queryPage(param) {
        return http.post('api/log/queryPage', param);
    },
    clean(param) {
        return http.post('api/log/clean', param);
    },
};


export default database;