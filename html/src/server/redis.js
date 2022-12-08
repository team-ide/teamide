import http from '@/server/http';

let redis = {
    info(param) {
        return http.post('api/redis/info', param);
    },
    get(param) {
        return http.post('api/redis/get', param);
    },
    keys(param) {
        return http.post('api/redis/keys', param);
    },
    set(param) {
        return http.post('api/redis/set', param);
    },
    sadd(param) {
        return http.post('api/redis/sadd', param);
    },
    srem(param) {
        return http.post('api/redis/srem', param);
    },
    lpush(param) {
        return http.post('api/redis/lpush', param);
    },
    rpush(param) {
        return http.post('api/redis/rpush', param);
    },
    lset(param) {
        return http.post('api/redis/lset', param);
    },
    lrem(param) {
        return http.post('api/redis/lrem', param);
    },
    hset(param) {
        return http.post('api/redis/hset', param);
    },
    hdel(param) {
        return http.post('api/redis/hdel', param);
    },
    delete(param) {
        return http.post('api/redis/delete', param);
    },
    deletePattern(param) {
        return http.post('api/redis/deletePattern', param);
    },
    expire(param) {
        return http.post('api/redis/expire', param);
    },
    ttl(param) {
        return http.post('api/redis/ttl', param);
    },
    persist(param) {
        return http.post('api/redis/persist', param);
    },
    close(param) {
        return http.post('api/redis/close', param);
    },
};


export default redis;