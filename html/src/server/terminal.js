import http from '@/server/http';

let fileManager = {
    close(param) {
        return http.post('api/terminal/close', param);
    },
    key(param) {
        return http.post('api/terminal/key', param);
    },
    changeSize(param) {
        return http.post('api/terminal/changeSize', param);
    },
};


export default fileManager;