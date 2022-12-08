import http from '@/server/http';

let kafka = {
    info(param) {
        return http.post('api/kafka/info', param);
    },
    topics(param) {
        return http.post('api/kafka/topics', param);
    },
    commit(param) {
        return http.post('api/kafka/commit', param);
    },
    pull(param) {
        return http.post('api/kafka/pull', param);
    },
    push(param) {
        return http.post('api/kafka/push', param);
    },
    reset(param) {
        return http.post('api/kafka/reset', param);
    },
    deleteTopic(param) {
        return http.post('api/kafka/deleteTopic', param);
    },
    createTopic(param) {
        return http.post('api/kafka/createTopic', param);
    },
    createPartitions(param) {
        return http.post('api/kafka/createPartitions', param);
    },
    deleteRecords(param) {
        return http.post('api/kafka/deleteRecords', param);
    },
    close(param) {
        return http.post('api/kafka/close', param);
    },
};


export default kafka;