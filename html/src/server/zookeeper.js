import http from '@/server/http';

let zookeeper = {
    info(param) {
        return http.post('api/zookeeper/info', param);
    },
    get(param) {
        return http.post('api/zookeeper/get', param);
    },
    save(param) {
        return http.post('api/zookeeper/save', param);
    },
    getChildren(param) {
        return http.post('api/zookeeper/getChildren', param);
    },
    delete(param) {
        return http.post('api/zookeeper/delete', param);
    },
    close(param) {
        return http.post('api/zookeeper/close', param);
    },
};


export default zookeeper;