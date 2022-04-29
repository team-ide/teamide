import http from '@/server/http';

import application from "./application.js";
import toolbox from "./toolbox.js";
let server = {
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
    application,
    toolbox,
};

export default server;