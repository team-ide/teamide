<template>
  <div class="toolbox-ssh-editor">
    <tm-layout height="100%">
      <tm-layout height="auto">
        <tm-layout width="50%">
          <ToolboxFTPFiles
            ref="localToolboxFTPFiles"
            :source="source"
            place="local"
            :dir="extend.local.dir"
            :files="localFiles"
            :wrap="wrap"
            @open="openFile"
            @openDir="openDir"
            @upload="doUploadFile"
            @download="doDownloadFile"
            @remove="doRemoveFile"
            @rename="doRenameFile"
            @refresh="toRefresh"
            @copy="toCopy"
          ></ToolboxFTPFiles>
        </tm-layout>
        <tm-layout-bar right></tm-layout-bar>
        <tm-layout width="auto">
          <ToolboxFTPFiles
            ref="remoteToolboxFTPFiles"
            :source="source"
            place="remote"
            :dir="extend.remote.dir"
            :files="remoteFiles"
            :wrap="wrap"
            @open="openFile"
            @openDir="openDir"
            @upload="doUploadFile"
            @download="doDownloadFile"
            @remove="doRemoveFile"
            @rename="doRenameFile"
            @refresh="toRefresh"
            @copy="toCopy"
          ></ToolboxFTPFiles>
        </tm-layout>
      </tm-layout>
      <tm-layout-bar top></tm-layout-bar>
      <tm-layout height="100px">
        <div class="works-box scrollbar" @contextmenu.prevent="workContextmenu">
          <template v-for="(one, index) in works">
            <div :key="'work-' + index" class="work-box">
              <div class="work-text">
                <span class="pdr-5">{{ one.startTime }}</span>
                <template v-if="one.place == 'local'">
                  <span class="pdr-5 color-grey-2">本地</span>
                </template>
                <template v-else-if="one.place == 'remote'">
                  <span class="pdr-5 color-grey-2">远程</span>
                </template>
                <template v-if="one.work == 'rename'">
                  <span class="pdr-5">重命名</span>
                  <span class="pdr-5">修改：{{ one.oldPath }}</span>
                  <span class="pdr-5">为：{{ one.newPath }}</span>
                </template>
                <template v-else-if="one.work == 'copy'">
                  <span class="pdr-5">复制</span>
                  <template v-if="one.fromFile.place == 'local'">
                    <span class="pdr-5 color-grey-2">本地</span>
                  </template>
                  <template v-else-if="one.fromFile.place == 'remote'">
                    <span class="pdr-5 color-grey-2">远程</span>
                  </template>
                  <span class="pdr-5">{{ one.fromFile.path }}</span>
                  <span class="pdr-5">到</span>
                  <template v-if="one.toFile.place == 'local'">
                    <span class="pdr-5 color-grey-2">本地</span>
                  </template>
                  <template v-else-if="one.toFile.place == 'remote'">
                    <span class="pdr-5 color-grey-2">远程</span>
                  </template>
                  <span class="pdr-5">{{ one.toFile.path }}</span>
                </template>
                <template v-else-if="one.work == 'upload'">
                  <span class="pdr-5">文件上传</span>
                  <span class="pdr-5">上传文件：{{ one.fileName }}</span>
                  <span class="pdr-5">目录：{{ one.dir }}</span>
                </template>
                <template v-else-if="one.work == 'remove'">
                  <span class="pdr-5">删除文件</span>
                  <span class="pdr-5">{{ one.path }}</span>
                </template>
              </div>
              <div class="work-status">
                <template v-if="one.progress != null">
                  <template v-if="one.progress.count != null">
                    <span class="color-grey-2 mgr-5">
                      文件：
                      {{ one.progress.count }}
                      /
                      {{ one.progress.successCount }}
                    </span>
                  </template>
                  <template v-if="one.progress.unitSize != null">
                    <span class="color-grey-2 mgr-5">
                      大小：
                      {{ one.progress.unitSize }}
                      {{ one.progress.unit }}
                      <template v-if="one.progress.unitSuccessSize != null">
                        /
                        {{ one.progress.unitSuccessSize }}
                        {{ one.progress.unitSuccess }}
                      </template>
                    </span>
                  </template>
                  <template v-if="one.progress.unitSleepSize != null">
                    <span class="color-grey-2 mgr-5">
                      速度：
                      {{ one.progress.unitSleepSize }}
                      {{ one.progress.unitSleep }}
                      / 秒
                    </span>
                  </template>
                  <template v-if="one.progress.percentage != null">
                    <span class="color-grey-2 mgr-5">
                      进度：
                      {{ one.progress.percentage }}
                    </span>
                  </template>
                </template>
                <template v-if="one.msg">
                  <span class="color-red">{{ one.msg }}</span>
                </template>
                <template v-else-if="one.isEnd">
                  <span class="color-green">完成</span>
                </template>
                <template v-else>
                  <span class="color-orange">执行中</span>
                </template>
              </div>
            </div>
          </template>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>

<script>
export default {
  components: {},
  props: [
    "source",
    "toolboxType",
    "toolbox",
    "option",
    "extend",
    "wrap",
    "token",
    "socket",
  ],
  data() {
    return {
      localFiles: null,
      remoteFiles: null,
      works: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.wrap.formatSize = this.formatSize;
    },
    onFocus() {
      this.$refs.remoteToolboxFTPFiles.onFocus();
    },
    refresh() {
      this.$refs.localToolboxFTPFiles.refresh();
      this.$refs.remoteToolboxFTPFiles.refresh();
    },
    formatSize(data, name, sizeName, sizeUnitName) {
      data[name] = data[name] || 0;
      let sizeMap = [
        { size: 1024 * 1024 * 1024 * 1024, unit: "TB" },
        { size: 1024 * 1024 * 1024, unit: "GB" },
        { size: 1024 * 1024, unit: "MB" },
        { size: 1024, unit: "KB" },
      ];

      sizeMap.forEach((one) => {
        if (!data[sizeUnitName] && data[name] >= one.size) {
          data[sizeName] = Number(data[name] / one.size).toFixed(2);
          data[sizeUnitName] = one.unit;
        }
      });
      if (!data[sizeUnitName]) {
        data[sizeName] = data[name];
        data[sizeUnitName] = "B";
      }
    },
    workContextmenu(e) {
      e = e || window.event;
      let menus = [];

      menus.push({
        text: "清理已完成",
        onClick: () => {
          this.cleanWork();
        },
      });

      this.tool.showContextmenu(menus);
    },
    cleanWork() {
      let ws = [];
      this.works.forEach((one) => {
        if (!one.isEnd) {
          ws.push(one);
        }
      });
      this.works = ws;
    },
    getWork(workId) {
      let res = null;
      this.works.forEach((one) => {
        if (one.workId == workId) {
          res = one;
        }
      });
      return res;
    },
    addWork(data) {
      data.workId = "work-" + this.tool.getNumber();
      let work = {};
      Object.assign(work, data);
      work.status = "working";
      work.isEnd = false;
      work.progress = null;
      work.startTime = this.tool.formatDate(new Date(), "yyyy-MM-dd hh:mm:ss");
      this.works.push(work);
    },
    onWorkSuccess(data) {
      let work = this.getWork(data.workId);
      if (work == null) {
        return;
      }
      work.msg = data.msg;
      work.status = "worked";
      work.isEnd = true;
    },
    onWorkProgress(data) {
      let work = this.getWork(data.workId);
      if (work == null) {
        return;
      }
      work.progress = data.progress;
      if (work.progress.size) {
        this.formatSize(work.progress, "size", "unitSize", "unit");

        if (work.progress.successSize) {
          this.formatSize(
            work.progress,
            "successSize",
            "unitSuccessSize",
            "unitSuccess"
          );

          if (work.progress.size > 0 && work.progress.successSize > 0) {
            let percentage = Number(
              (work.progress.successSize / work.progress.size) * 100
            ).toFixed(0);
            work.progress.percentage = percentage + "%";
          } else {
            if (work.progress.size == work.progress.successSize) {
              work.progress.percentage = "100%";
            } else {
              work.progress.percentage = "0%";
            }
          }

          if (work.progress.startTime && work.progress.timestamp) {
            let usetime = work.progress.timestamp - work.progress.startTime;
            if (work.progress.endTime) {
              usetime = work.progress.endTime - work.progress.startTime;
            }
            let sleepSize = Number(
              (work.progress.successSize / usetime) * 1000
            ).toFixed(2);
            work.progress.sleepSize = sleepSize;

            this.formatSize(
              work.progress,
              "sleepSize",
              "unitSleepSize",
              "unitSleep"
            );
          }
        }
      }
    },
    openFile(place, dir, file) {
      if (file.isDir) {
        this.openDir(place, file.path);
      }
    },
    openDir(place, dir) {
      this.loadFiles(place, dir);
    },
    toRefresh(place, dir) {
      this.openDir(place, dir);
    },
    loadFiles(place, dir) {
      this.doLoadFiles(place, dir);
    },
    doLoadFiles(place, dir) {
      let scrollTop = 0;
      if (place == "local") {
        this.localFiles = null;
        scrollTop = this.$refs.localToolboxFTPFiles.getScrollTop();
      } else if (place == "remote") {
        this.remoteFiles = null;
        scrollTop = this.$refs.remoteToolboxFTPFiles.getScrollTop();
      }
      let request = {
        place: place,
        dir: dir,
        work: "files",
        scrollTop: Number(scrollTop),
      };

      this.send(request);
    },
    doRenameFile(place, dir, oldPath, newPath) {
      let request = {
        place: place,
        dir: dir,
        newPath: newPath,
        oldPath: oldPath,
        work: "rename",
      };
      this.addWork(request);
      this.send(request);
    },
    toCopy(fromFile, toFile) {
      let request = {
        place: toFile.place,
        dir: toFile.dir,
        fromFile: fromFile,
        toFile: toFile,
        work: "copy",
      };
      this.addWork(request);
      this.send(request);
    },
    doRemoveFile(place, dir, files) {
      let names = [];
      files.forEach((one) => {
        names.push(one.name);
      });
      this.tool
        .confirm("删除[" + names.join(",") + "]后无法恢复，确定删除？")
        .then(async () => {
          files.forEach((one) => {
            let request = {
              place: place,
              dir: dir,
              path: one.path,
              work: "remove",
            };

            this.addWork(request);
            this.send(request);
          });
        })
        .catch((e) => {});
    },
    send(request) {
      let message = JSON.stringify(request);
      this.wrap.writeData(message);
    },
    onConfirm(response) {
      let confirm = response.confirm;
      let okTitle = "确认";
      let cancelTitle = "取消";
      if (response.isFileExist) {
        confirm = "文件[" + response.path + "]已存在，请选择操作！";
        okTitle = "覆盖";
        cancelTitle = "跳过";
      }
      this.tool
        .confirm(confirm, okTitle, cancelTitle)
        .then(() => {
          let request = {
            confirmId: response.confirmId,
            work: "confirmResult",
            isOk: true,
          };
          this.send(request);
        })
        .catch((e) => {
          let request = {
            confirmId: response.confirmId,
            work: "confirmResult",
            isCancel: true,
          };
          this.send(request);
        });
    },
    onResponse(response) {
      if (this.tool.isNotEmpty(response.msg)) {
        this.tool.error(response.msg);
      }
      if (response.isProgress) {
        this.onWorkProgress(response);
        return;
      }
      if (response.isConfirm) {
        this.onConfirm(response);
        return;
      }
      if (response.work == "files") {
        if (response.place == "local") {
          if (this.extend.local.dir != response.dir) {
            this.wrap.updateExtend(["local", "dir"], response.dir);
          }
          this.localFiles = response.files || [];
          this.$refs.localToolboxFTPFiles.setScrollTop(response.scrollTop);
        } else if (response.place == "remote") {
          if (this.extend.remote.dir != response.dir) {
            this.wrap.updateExtend(["remote", "dir"], response.dir);
          }
          if (response.dir.split("/").length > 3) {
            let ss = response.dir.split("/");
            let d = "..";
            for (var i = ss.length - 2; i < ss.length; i++) {
              d += "/" + ss[i];
            }
            this.wrap.updateComment(d);
          } else {
            this.wrap.updateComment(response.dir);
          }
          this.remoteFiles = response.files || [];
          this.$refs.remoteToolboxFTPFiles.setScrollTop(response.scrollTop);
        }
      }
      this.onWorkSuccess(response);
      if (
        response.work == "upload" ||
        response.work == "remove" ||
        response.work == "copy" ||
        response.work == "rename"
      ) {
        if (response.place == "local") {
          if (response.dir == this.extend.local.dir) {
            this.loadFiles("local", this.extend.local.dir);
          }
        } else if (response.place == "remote") {
          if (response.dir == this.extend.remote.dir) {
            this.loadFiles("remote", this.extend.remote.dir);
          }
        }
      }
    },
    async doUploadFile(place, dir, file, fullPath) {
      let request = {
        work: "upload",
        place: place,
        dir: dir,
        token: this.token,
        fileName: fullPath || file.name,
        fullPath: fullPath,
      };
      this.addWork(request);
      let form = new FormData();
      for (let key in request) {
        form.append(key, request[key]);
      }
      form.append("file", file);
      let res = await this.server.toolbox.ssh.ftp.upload(form);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      return true;
    },
    async doDownloadFile(place, dir, file) {
      let url =
        this.source.api +
        "api/toolbox/ssh/ftp/download?token=" +
        encodeURIComponent(this.token) +
        "&jwt=" +
        encodeURIComponent(this.tool.getJWT()) +
        "&place=" +
        encodeURIComponent(place) +
        "&dir=" +
        encodeURIComponent(dir) +
        "&path=" +
        encodeURIComponent(file.path);
      window.location.href = url;
    },
    onEvent(event) {
      if (event == "ftp ready") {
        this.toStart();
      } else if (event == "ftp created") {
        this.doLoadFiles("local", this.extend.local.dir);
        this.doLoadFiles("remote", this.extend.remote.dir);
      }
    },
    onError(error) {
      this.tool.error(error);
    },
    onData(data) {
      let response = JSON.parse(data);
      this.onResponse(response);
    },
    toStart() {
      this.wrap.writeEvent("ftp start");
    },
    dispose() {},
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeUpdate() {},
  beforeDestroy() {
    this.dispose();
  },
};
</script>

<style>
.toolbox-ftp-editor {
  width: 100%;
  height: 100%;
}
.works-box {
  width: 100%;
  height: 100%;
}
.work-box {
  display: flex;
  line-height: 20px;
  font-size: 12px;
  padding: 0px 5px;
}
.work-box .work-icon {
  padding: 0px 5px;
}
.work-box .work-text {
  padding: 0px 5px;
  flex: 1;
}
.work-box .work-size {
  padding: 0px 5px;
  font-size: 12px;
}
</style>
