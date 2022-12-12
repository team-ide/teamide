import http from '@/server/http';

let elasticsearch = {
    info(param) {
        return http.post('api/elasticsearch/info', param);
    },
    indexes(param) {
        return http.post('api/elasticsearch/indexes', param);
    },
    createIndex(param) {
        return http.post('api/elasticsearch/createIndex', param);
    },
    deleteIndex(param) {
        return http.post('api/elasticsearch/deleteIndex', param);
    },
    getMapping(param) {
        return http.post('api/elasticsearch/getMapping', param);
    },
    putMapping(param) {
        return http.post('api/elasticsearch/putMapping', param);
    },
    search(param) {
        return http.post('api/elasticsearch/search', param);
    },
    scroll(param) {
        return http.post('api/elasticsearch/scroll', param);
    },
    insertData(param) {
        return http.post('api/elasticsearch/insertData', param);
    },
    updateData(param) {
        return http.post('api/elasticsearch/updateData', param);
    },
    deleteData(param) {
        return http.post('api/elasticsearch/deleteData', param);
    },
    reindex(param) {
        return http.post('api/elasticsearch/reindex', param);
    },
    import(param) {
        return http.post('api/elasticsearch/import', param);
    },
    export(param) {
        return http.post('api/elasticsearch/export', param);
    },
    sync(param) {
        return http.post('api/elasticsearch/sync', param);
    },
    taskStatus(param) {
        return http.post('api/elasticsearch/taskStatus', param);
    },
    taskStop(param) {
        return http.post('api/elasticsearch/taskStop', param);
    },
    taskClean(param) {
        return http.post('api/elasticsearch/taskClean', param);
    },
    close(param) {
        return http.post('api/elasticsearch/close', param);
    },
};


export default elasticsearch;