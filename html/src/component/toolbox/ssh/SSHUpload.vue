<template>
  <el-dialog
    ref="modal"
    :title="`SSH 文件上传`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="900px"
  >
    <div>
      <div
        class="ssh-upload-box"
        @drop="onDrop"
        @dragover="onDragover"
        @dragleave="onDragleave"
        draggable="true"
      >
        <i class="mdi mdi-upload ft-20"></i>
        <div class="mgt-10 ft-16">
          将文件拖到此处
          <!-- ，或
          <div class="tm-link color-green" @click="toClickUpload">点击上传</div> -->
        </div>
      </div>
    </div>
    <input
      type="file"
      id="input-for-upload"
      @change="uploadInputChange"
      ref="input-for-upload"
    />
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "wrap", "token"],
  data() {
    return {
      showDialog: false,
      isSuccess: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    show(onUpload, onCancel) {
      this.isSuccess = false;
      this.onUpload = onUpload;
      this.onCancel = onCancel;
      this.showDialog = true;
    },
    hide() {
      if (!this.isSuccess) {
        this.tool
          .confirm("是否取消上传？")
          .then(async () => {
            this.showDialog = false;
            this.onCancel && this.onCancel();
          })
          .catch((e) => {});
      } else {
        this.showDialog = false;
      }
    },
    toClickUpload() {
      this.$refs["input-for-upload"].value = null;
      this.$refs["input-for-upload"].click();
    },
    uploadInputChange() {
      let upload = this.$refs["input-for-upload"];

      this.doUpload(upload.files);
    },
    doUpload(files) {
      this.isSuccess = true;
      this.hide();
      this.onUpload && this.onUpload(files);
    },
    onDragover(e) {
      this.tool.stopEvent(e);
    },
    onDragleave(e) {
      this.tool.stopEvent(e);
    },
    onDrop(e) {
      this.tool.stopEvent(e);
      let files = [];
      let endCall = () => {
        this.doUpload(files);
      };
      if (
        e.dataTransfer &&
        e.dataTransfer.items &&
        e.dataTransfer.items.length > 0
      ) {
        let itemsLength = e.dataTransfer.items.length;
        Array.prototype.forEach.call(
          e.dataTransfer.items,
          async (one, index) => {
            if (one.webkitGetAsEntry) {
              let webkitGetAsEntry = one.webkitGetAsEntry();
              this.uploadEntryFile(files, webkitGetAsEntry, () => {
                if (index == itemsLength - 1) {
                  endCall();
                }
              });
              return;
            }
            let file = one.getAsFile();
            if (file != null) {
              files.push(file);
              if (index == e.dataTransfer.items.length - 1) {
                endCall();
              }
            }
          }
        );
      }
    },
    uploadEntryFile(files, entry, endCall) {
      if (entry.isFile) {
        entry.file(
          (file) => {
            files.push(file);
            endCall();
          },
          (e) => {
            console.log(e);
          }
        );
      } else {
        let reader = entry.createReader();
        reader.readEntries(
          (entries) => {
            let entriesLength = entries.length;
            entries.forEach((entry, index) =>
              this.uploadEntryFile(files, entry, () => {
                if (index == entriesLength - 1) {
                  endCall();
                }
              })
            );
          },
          (e) => {
            console.log(e);
          }
        );
      }
    },
  },
  created() {},
  mounted() {
    this.wrap.showSSHUpload = this.show;
    this.wrap.hideSSHUpload = this.show;
  },
};
</script>

<style>
.ssh-upload-box {
  width: 400px;
  margin: 0px auto;
  text-align: center;
  border: 2px dashed #ddd;
  padding: 70px;
}
#input-for-upload {
  width: 0px;
  height: 0px;
  position: fixed;
  left: -100px;
  top: -100px;
}
</style>
