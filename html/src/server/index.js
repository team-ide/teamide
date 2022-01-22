import http from '@/server/http';

import application from "./application.js";
let server = {
    data(param) {
        param = param || {};
        param.origin = location.origin;
        param.pathname = location.pathname;
        return http.post('api/', param,);
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
    application,
};

export default server;