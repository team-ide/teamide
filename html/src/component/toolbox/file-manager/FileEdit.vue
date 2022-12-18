<template>
  <el-dialog
    ref="modal"
    :title="`文件：${path}`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    :fullscreen="true"
    width="96%"
    custom-class="toolbox-file-manager-edit-file-dialog"
    :destroy-on-close="true"
  >
    <div class="toolbox-file-manager-edit-file" v-loading="loading">
      <template v-if="isImage">
        <el-image
          style="width: 100%; height: 100%"
          :src="openUrl"
          fit="scale-down"
        ></el-image>
      </template>
      <template v-else-if="isVideo">
        <video :src="openUrl" controls width="100%">
          您的浏览器不支持 video 标签。 Internet Explorer 9+, Firefox, Opera,
          Chrome 以及 Safari 支持 video 标签。
        </video>
      </template>

      <template v-else>
        <div class="teamide-editor-box">
          <template v-if="textReady">
            <Editor
              ref="Editor"
              :source="source"
              :value="text"
              :language="language"
            ></Editor>
          </template>
        </div>
        <div class="pdt-10 pdl-10">
          <div
            class="tm-btn bg-green ft-13"
            @click="toSave"
            :class="{ 'tm-disabled': saveIng }"
          >
            保存
          </div>
          <div
            class="tm-btn bg-grey ft-13"
            @click="hide"
            :class="{ 'tm-disabled': saveIng }"
          >
            关闭
          </div>
        </div>
      </template>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "fileWorker"],
  data() {
    return {
      showDialog: false,
      path: null,
      saveIng: false,
      loading: false,
      isImage: false,
      isVideo: false,
      language: "txt",
      textReady: false,
      openUrl: null,
      text: null,
    };
  },
  computed: {},
  watch: {
    showDialog() {
      if (!this.showDialog) {
        this.textReady = false;
      }
    },
  },
  methods: {
    show(file) {
      if (file == null || file.name == null) {
        this.tool.error("文件不存在");
        return;
      }
      let name = file.name;
      this.type = null;
      if (name.lastIndexOf(".") > 0 && name.lastIndexOf(".") < name.length) {
        this.type = name.substring(name.lastIndexOf(".") + 1);
      }
      this.language = this.type;
      this.isImage = this.tool.isImageByType(this.type);
      this.isVideo = this.tool.isVideoByType(this.type);
      this.file = file;
      this.path = file.path;

      if (this.isImage) {
        this.openUrl = this.getFileOpenUrl();
        this.showDialog = true;
      } else if (this.isVideo) {
        this.openUrl = this.getFileOpenUrl();
        this.showDialog = true;
      } else {
        this.toLoad();
      }
    },
    getFileOpenUrl() {
      let url = this.source.api + "fileManager/open?";
      url += "workerId=" + (this.fileWorker.workerId || "");
      url = this.tool.appendUrlBaseParam(url);
      url += "&fileWorkerKey=" + (this.fileWorker.fileWorkerKey || "");
      url += "&place=" + (this.fileWorker.place || "");
      url += "&placeId=" + (this.fileWorker.placeId || "");
      url += "&path=" + encodeURIComponent(this.path);
      return url;
    },
    hide() {
      this.openUrl = null;
      this.showDialog = false;
    },
    async toLoad(force) {
      try {
        this.loading = true;
        let res = await this.fileWorker.read(this.path, force);
        if (res.code != 0) {
          if (res.code == 5001) {
            this.tool
              .confirm("文件过大，是否打开？")
              .then(() => {
                this.toLoad(true);
              })
              .catch((e) => {
                this.showDialog = false;
              });
          } else {
            this.tool.error(res.msg || res.error);
          }
        } else {
          this.showDialog = true;
          let data = res.data || {};
          this.text = data.text;
          this.textReady = true;
        }
      } catch (e) {}
      this.loading = false;
    },
    async toSave() {
      if (this.$refs.Editor == null) {
        return;
      }
      this.saveIng = true;
      let text = this.$refs.Editor.getValue();
      let res = await this.fileWorker.write(this.path, text);
      if (res.code == 0) {
        this.tool.success("保存成功!");
      } else {
        this.tool.error(res.msg);
      }
      this.saveIng = false;
    },
    init() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-file-manager-edit-file-dialog {
  background-color: #3c3c3c;
  color: white;
  user-select: text;
}
.toolbox-file-manager-edit-file-dialog .el-dialog__header {
  padding: 10px 20px 10px;
}
.toolbox-file-manager-edit-file-dialog .el-dialog__header .el-dialog__title {
  color: white;
}
.toolbox-file-manager-edit-file-dialog
  .el-dialog__header
  .el-dialog__headerbtn {
  top: 15px;
  right: 15px;
}
.toolbox-file-manager-edit-file-dialog .el-dialog__body {
  height: calc(100% - 44px);
  padding: 0px;
  overflow: hidden;
}
.toolbox-file-manager-edit-file {
  height: 100%;
  width: 100%;
}
.toolbox-file-manager-edit-file .teamide-editor-box {
  height: calc(100% - 50px);
  width: 100%;
}
</style>
