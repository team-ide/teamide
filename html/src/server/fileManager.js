import http from '@/server/http';

let fileManager = {
    close(param) {
        return http.post('api/fileManager/close', param);
    },
    file(param) {
        return http.post('api/fileManager/file', param);
    },
    files(param) {
        return http.post('api/fileManager/files', param);
    },
    create(param) {
        return http.post('api/fileManager/create', param);
    },
    read(param) {
        return http.post('api/fileManager/read', param);
    },
    write(param) {
        return http.post('api/fileManager/write', param);
    },
    rename(param) {
        return http.post('api/fileManager/rename', param);
    },
    move(param) {
        return http.post('api/fileManager/move', param);
    },
    remove(param) {
        return http.post('api/fileManager/remove', param);
    },
    copy(param) {
        return http.post('api/fileManager/copy', param);
    },
    callAction(param) {
        return http.post('api/fileManager/callAction', param);
    },
    callStop(param) {
        return http.post('api/fileManager/callStop', param);
    },
    upload(param) {
        return http.post('api/fileManager/upload', param, { headers: { 'Content-Type': 'multipart/form-data' } });
    },
    download(param) {
        return http.post('api/fileManager/download', param, { responseType: "blob" });
    },
};


export default fileManager;