import http from '@/server/http';

let user = {
    get(param) {
        return http.post('api/user/get', param||{});
    },
    update(param) {
        return http.post('api/user/update', param);
    },
    updatePassword(param) {
        return http.post('api/user/updatePassword', param);
    },
};


export default user;