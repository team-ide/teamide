import http from '@/server/http';
let server = {
    data(param) {
        param = param || {};
        param.origin = location.origin;
        param.pathname = location.pathname;
        return http.post('api/', param, {
            headers: {
            }
        });
    },
    session(param) {
        return http.post('api/session', param, {
            headers: {
            }
        });
    },
    login(param) {
        return http.post('api/login', param, {
            headers: {
            }
        });
    },
    logout(param) {
        return http.post('api/logout', param, {
            headers: {
            }
        });
    },
};

export default server;