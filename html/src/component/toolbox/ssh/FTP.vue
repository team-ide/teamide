<template>
  <div class="toolbox-ssh-editor">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="auto">
          <tm-layout width="50%">
            <ToolboxFTPFiles
              :source="source"
              place="local"
              :dir="localDir"
              :files="localFiles"
              :wrap="wrap"
              @open="openFile"
              @openDir="openDir"
              @upload="doUploadFile"
              @remove="doRemoveFile"
              @rename="doRenameFile"
              @refresh="toRefresh"
              @copy="toCopy"
            ></ToolboxFTPFiles>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="auto">
            <ToolboxFTPFiles
              :source="source"
              place="remote"
              :dir="remoteDir"
              :files="remoteFiles"
              :wrap="wrap"
              @open="openFile"
              @openDir="openDir"
              @upload="doUploadFile"
              @remove="doRemoveFile"
              @rename="doRenameFile"
              @refresh="toRefresh"
              @copy="toCopy"
            ></ToolboxFTPFiles>
          </tm-layout>
        </tm-layout>
        <tm-layout-bar bottom></tm-layout-bar>
        <tm-layout height="100px">
          <div
            class="works-box scrollbar"
            @contextmenu.prevent="workContextmenu"
          >
            <template v-for="(one, index) in works">
              <div :key="'work-' + index" class="work-box">
                <div class="work-text">
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
    </template>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
      token: null,
      localDir: null,
      localFiles: null,
      remoteDir: null,
      remoteFiles: null,
      works: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.wrap.formatSize = this.formatSize;
      this.ready = true;
      await this.initToken();
      this.initSocket();
    },
    formatSize(data, name, sizeName, sizeUnitName) {
      data[name] = data[name] || 0;
      let sizeMap = [
        { size: 1024 * 1024 * 1024 * 1024, unit: "TB" },
        { size: 1024 * 1024 * 1024, unit: "GB" },
        { size: 1024 * 1024, unit: "MB" },
        { size: 1024, unit: "kb" },
      ];

      sizeMap.forEach((one) => {
        if (!data[sizeUnitName] && data[name] >= one.size) {
          data[sizeName] = Number(data[name] / one.size).toFixed(2);
          data[sizeUnitName] = one.unit;
        }
      });
      if (!data[sizeUnitName]) {
        data[sizeName] = data[name];
        data[sizeUnitName] = "b";
      }
    },
    async initToken() {
      if (this.tool.isEmpty(this.token)) {
        let param = {};
        let res = await this.wrap.work("createToken", param);
        res.data = res.data || {};
        this.token = res.data.token;
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

      // data.workId = this.tool.getNumber();
      // let work = {};
      // Object.assign(work, data);
      // this.works.push(work);
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
      // if (place == "local") {
      //   this.localDir = dir;
      //   this.localFiles = null;
      // } else if (place == "remote") {
      //   this.remoteDir = dir;
      //   this.remoteFiles = null;
      // }
      let request = {
        place: place,
        dir: dir,
        work: "files",
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
      this.socket.send(message);
    },
    onResponse(response) {
      if (this.tool.isNotEmpty(response.msg)) {
        this.tool.error(response.msg);
      }
      if (response.work == "files") {
        if (response.place == "local") {
          this.localDir = response.dir;
          this.localFiles = response.files || [];
        } else if (response.place == "remote") {
          this.remoteDir = response.dir;
          this.remoteFiles = response.files || [];
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
          if (response.dir == this.localDir) {
            this.loadFiles("local", this.localDir);
          }
        } else if (response.place == "remote") {
          if (response.dir == this.remoteDir) {
            this.loadFiles("remote", this.remoteDir);
          }
        }
      }
    },
    async doUploadFile(place, dir, file, fullPath) {
      let request = {
        work: "upload",
        type: "sftp",
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
      let res = await this.server.upload(form);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      return true;
    },
    initSocket() {
      if (this.socket != null) {
        this.socket.close();
      }

      let url = this.source.api;
      url = url.substring(url.indexOf(":"));
      url = "ws" + url + "ws/toolbox/sfpt/connection?token=" + this.token;
      this.socket = new WebSocket(url);

      this.socket.onopen = () => {
        this.loadFiles("local");
        this.loadFiles("remote");
      };
      this.socket.onmessage = (event) => {
        // 接收推送的消息
        let data = event.data.toString();
        let response = JSON.parse(data);
        this.onResponse(response);
      };
      this.socket.onclose = () => {
        console.log("close socket");
      };
      this.socket.onerror = () => {
        console.log("socket error");
      };
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    if (this.socket != null) {
      this.socket.close();
    }
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
