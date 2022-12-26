import server from "@/server/index.js";
import tool from "@/tool/index.js";
import source from "@/source/index.js";

const newWorker = function (workerOption) {
    workerOption = workerOption || {};
    const worker = {
        workerId: workerOption.workerId,
        place: workerOption.place,
        placeId: workerOption.placeId,
        selectPaths: [],
        fileList: null,
        dir: "",
        dirNames: [],
        isInputDir: false,
        filter: "",
        fileWorkerKey: tool.generatekey(30),
        onChangeOpenDir: workerOption.onChangeOpenDir,
        async init() {

        },
        refresh() {
            this.loadFiles(this.dir);
        },
        openDir(dir) {
            this.loadFiles(dir);
        },
        toDirName(dirNameIndex) {
            tool.stopEvent();
            let dir = "";
            if (this.dirNames.length > 0) {
                this.dirNames.forEach((name, i) => {
                    if (i <= dirNameIndex) {
                        dir += "/" + name;
                    }
                });
                if (this.dirNames[0] != "") {
                    dir = dir.substring(1);
                }
            }
            this.openDir(dir);
        },
        getFile(path) {
            let file = null;
            if (this.fileList) {
                this.fileList.forEach((one) => {
                    if (one.path == path) {
                        file = one;
                    }
                });
            }
            return file;
        },
        getFileIndex(path) {
            let index = -1;
            if (this.fileList) {
                this.fileList.forEach((one, i) => {
                    if (one == path || one.path == path) {
                        index = i;
                    }
                });
            }
            return index;
        },
        getRenameFile() {
            let file = null;
            if (this.fileList) {
                this.fileList.forEach((one) => {
                    if (one.rename) {
                        file = one;
                    }
                });
            }
            return file;
        },
        getSelectFiles() {
            let files = [];
            if (this.fileList) {
                this.fileList.forEach((one) => {
                    if (one.select) {
                        files.push(one);
                    }
                });
            }
            return files;
        },
        setSelect(file) {
            if (this.selectPaths.indexOf(file.path) < 0) {
                this.selectPaths.push(file.path);
            }
            file.select = true;
        },
        setUnselect(file) {
            if (this.selectPaths.indexOf(file.path) >= 0) {
                this.selectPaths.splice(this.selectPaths.indexOf(file.path), 1);
            }
            file.select = false;
        },
        toSelectAll() {
            this.fileList.forEach((one, i) => {
                if (one.name == "..") {
                    this.setUnselect(one);
                } else {
                    this.setSelect(one);
                }
            });
        },
        selectFile(file) {
            if (this.fileList) {
                let fileIndex = this.getFileIndex(file);
                let startIndex = fileIndex;
                let endIndex = fileIndex;
                if (window.event.shiftKey) {
                    let selectFiles = this.getSelectFiles();
                    selectFiles.forEach((one) => {
                        let selectIndex = this.getFileIndex(one);
                        if (selectIndex < 0) {
                            return;
                        }
                        if (selectIndex < startIndex) {
                            startIndex = selectIndex;
                        }
                        if (selectIndex > endIndex) {
                            endIndex = selectIndex;
                        }
                    });
                } else if (window.event.ctrlKey) {
                    if (fileIndex >= 0) {
                        this.setSelect(this.fileList[fileIndex]);
                    }
                    return;
                }
                this.fileList.forEach((one, i) => {
                    if (one.name == "..") {
                        this.setUnselect(one);
                    } else if (i >= startIndex && i <= endIndex) {
                        this.setSelect(one);
                    } else {
                        this.setUnselect(one);
                    }
                });
            }
        },
        getFileByTarget(target) {
            let file = null;
            if (target) {
                let box = tool.jQuery(target).closest("[path]");
                if (box.length > 0) {
                    let path = box.attr("path");
                    file = this.getFile(path);
                }
            }
            return file;
        },
        async loadFiles(dir) {
            if (this.dir != dir) {
                if (this.onChangeOpenDir) {
                    this.onChangeOpenDir(dir)
                }
            }
            this.fileList = null;
            this.selectPaths = [];
            let data = await this.files(dir);
            this.fileList = data.files || [];
            this.dir = data.dir || "";
            this.dirNames = ("" + this.dir).split("/");

            this.filterFile();
        },
        filterFile() {
            this.selectPaths = this.selectPaths || [];
            this.fileList.forEach((one) => {
                if (
                    tool.isEmpty(this.filter) ||
                    one.name == ".." ||
                    one.name.toLowerCase().indexOf(this.filter.toLowerCase()) >= 0
                ) {
                    one.show = true;
                    if (this.selectPaths.indexOf(one.path) >= 0) {
                        one.select = true;
                    }
                } else {
                    one.show = false;
                    one.select = false;
                }
            });
        },
        getParam() {
            let param = {
                workerId: worker.workerId,
                place: worker.place,
                placeId: worker.placeId,
                fileWorkerKey: worker.fileWorkerKey,
            };
            return param;
        },
        formatFiles(files) {
            files = files || [];
            files.forEach(file => {
                this.formatFile(file);
            })
        },
        formatFile(file) {
            if (file == null) {
                return
            }
            file.rename = false;
            file.size = file.size || 0;
            file.isNew = false;
            file.show = true;
            if (!file.isDir) {
                tool.formatSize(file, "size", "unitSize", "unit");
            }
            if (file.modTime) {
                file.dateTime = tool.formatDate(
                    new Date(file.modTime),
                    "yyyy-MM-dd hh:mm:ss"
                );
            }
            file.select = false;
            if (this.selectPaths.indexOf(file.path) >= 0) {
                file.select = true;
            }
        },
        async file(path) {
            let param = worker.getParam();
            param.path = path;
            let res = await server.fileManager.file(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            let file = res.data;
            if (file) {
                this.formatFile(file);
            }
            return file;
        },
        async files(dir) {
            let param = worker.getParam();
            param.dir = dir;
            let res = await server.fileManager.files(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            let data = res.data || {}
            this.formatFiles(data.files);
            return data;
        },
        async create(path, isDir) {
            let param = worker.getParam();
            param.path = path;
            param.isDir = isDir;
            let res = await server.fileManager.create(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            let file = res.data;
            if (file) {
                this.formatFile(file);
            }
            return file;
        },
        async rename(oldPath, newPath) {
            let param = worker.getParam();
            param.oldPath = oldPath;
            param.newPath = newPath;
            let res = await server.fileManager.rename(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            let file = res.data;
            if (file) {
                this.formatFile(file);
            }
            return file;
        },
        async remove(path) {
            let param = worker.getParam();
            param.path = path;
            let res = await server.fileManager.remove(param);
            if (res.code != 0) {
                tool.error(res.msg);
                return false;
            } else {
                let fileIndex = this.getFileIndex(path);
                if (fileIndex >= 0) {
                    this.fileList.splice(fileIndex, 1);
                }
                return true;
            }
        },
        async move(oldPath, newPath) {
            let param = worker.getParam();
            param.oldPath = oldPath;
            param.newPath = newPath;
            let res = await server.fileManager.move(param);
            if (res.code != 0) {
                tool.error(res.msg);
            } else {
                let fileIndex = this.getFileIndex(oldPath);
                if (fileIndex >= 0) {
                    this.fileList.splice(fileIndex, 1);
                }
            }
            return res.data;
        },
        pathIsSubfile(path) {
            if (tool.isEmpty(path) || path.indexOf('/') < 0) {
                return false;
            }
            let dir = this.dir;
            if (dir.endsWith('/')) {
                dir = dir.substring(0, dir.length - 1);
            }
            if (dir == path.substring(0, path.lastIndexOf('/'))) {
                return true;
            }
            return false;
        },
        async uploadFile(dir, file, fullPath) {
            let param = {
                workerId: this.workerId,
                fileWorkerKey: this.fileWorkerKey,
                place: this.place,
                placeId: this.placeId,
                dir: dir,
                fullPath: fullPath,
            };
            let form = new FormData();
            for (let key in param) {
                form.append(key, param[key]);
            }
            form.append("file", file);
            let res = await server.fileManager.upload(form);
            if (res.code != 0) {
                tool.error(res.msg);
                return false;
            } else {
                // let files = res.data || [];
                // files.forEach(one => {
                //     if (!this.pathIsSubfile(one.path)) {
                //         return
                //     }
                //     this.formatFile(one);
                //     let index = this.getFileIndex(one.path);
                //     if (index >= 0) {
                //         this.fileList.splice(index, 1, one)
                //     } else {
                //         this.fileList.push(one)
                //     }
                // })
                return true;
            }
        },
        onUploadFileInfo(fileInfo) {
            if (fileInfo == null) {
                return
            }
            if (!this.pathIsSubfile(fileInfo.path)) {
                return
            }
            this.formatFile(fileInfo);
            let index = this.getFileIndex(fileInfo.path);
            if (index >= 0) {
                this.fileList.splice(index, 1, fileInfo)
            } else {
                this.fileList.push(fileInfo)
            }
        },
        async download(path) {
            let url = source.api + "fileManager/download?";
            url += "workerId=" + (this.workerId || "");
            url = tool.appendUrlBaseParam(url);
            url += "&fileWorkerKey=" + (this.fileWorkerKey || "");
            url += "&place=" + (this.place || "");
            url += "&placeId=" + (this.placeId || "");
            url += "&path=" + encodeURIComponent(path);
            window.location.href = url;
        },
        async copy(path, fromPlace, fromPlaceId, fromPath) {
            let param = worker.getParam();
            param.path = path;
            param.fromPlace = fromPlace;
            param.fromPlaceId = fromPlaceId;
            param.fromPath = fromPath;
            let res = await server.fileManager.copy(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            return res.data;
        },
        async close() {
            let param = worker.getParam();
            let res = await server.fileManager.close(param);
            if (res.code != 0) {
                tool.error(res.msg);
            }
            return res.data;
        },
        async read(path, force) {
            let param = worker.getParam();
            param.path = path
            param.force = force
            let res = await server.fileManager.read(param);
            return res;
        },
        async write(path, text) {
            let param = worker.getParam();
            param.path = path
            param.text = text
            let res = await server.fileManager.write(param);
            return res;
        },
    };


    return worker
};
export default {
    newWorker
}