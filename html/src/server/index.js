import http from '@/server/http';
let server = {
    data(param) {
        return http.post('server/data', param, {
            headers: {
            }
        });
    },
    session(param) {
        return http.post('server/session', param, {
            headers: {
            }
        });
    },
    login(param) {
        return http.post('server/login', param, {
            headers: {
            }
        });
    },
    logout(param) {
        return http.post('server/logout', param, {
            headers: {
            }
        });
    },
};

export default server;