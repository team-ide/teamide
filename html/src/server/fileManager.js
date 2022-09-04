import http from '@/server/http';

let fileManager = {
    file(param) {
        return http.post('api/file_manager/file', param);
    },
    files(param) {
        return http.post('api/file_manager/files', param);
    },
    create(param) {
        return http.post('api/file_manager/create', param);
    },
    read(param) {
        return http.post('api/file_manager/read', param);
    },
    write(param) {
        return http.post('api/file_manager/write', param);
    },
    rename(param) {
        return http.post('api/file_manager/rename', param);
    },
    move(param) {
        return http.post('api/file_manager/move', param);
    },
    remove(param) {
        return http.post('api/file_manager/remove', param);
    },
    copy(param) {
        return http.post('api/file_manager/copy', param);
    },
    callAction(param) {
        return http.post('api/file_manager/callAction', param);
    },
    upload(param) {
        return http.post('api/file_manager/upload', param, { headers: { 'Content-Type': 'multipart/form-data' } });
    },
    download(param) {
        return http.post('api/file_manager/download', param, { responseType: "blob" });
    },
};


export default fileManager;