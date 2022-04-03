<template>
  <div class="toolbox-ftp-files">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="50px">
          <div class="file-dir-input-box">
            <input
              class="file-dir-input"
              v-model="form.dir"
              @keyup.enter="inputOnEnter"
            />
          </div>
        </tm-layout>
        <tm-layout height="auto">
          <div class="drag-file" v-show="dragFile.show" :style="dragFile.style">
            <div class="drag-file-name">{{ dragFile.name }}</div>
          </div>
          <template v-if="files == null">
            <div class="text-center pdtb-10 ft-15">加载中...</div>
          </template>
          <div
            class="files-box scrollbar"
            @contextmenu.prevent="fileContextmenu"
            ref="filesBox"
            @click="filesBoxClick"
          >
            <template v-for="(one, index) in list">
              <div
                :key="'file-' + index"
                class="file-box"
                :class="{ 'file-select': one.select }"
                draggable="true"
                @dblclick="openFile(one)"
                :path="one.path"
              >
                <div class="file-icon">
                  <i v-if="one.isDir" class="mdi mdi-folder color-orange-3"></i>
                  <i v-else class="mdi mdi-file color-white"></i>
                </div>
                <div class="file-name">
                  <template v-if="one.rename">
                    <input
                      class="file-rename-input"
                      v-model="one.newname"
                      @blur="onRenameBlur(one, $event)"
                      @keyup="onRenameKeyup(one, $event)"
                    />
                  </template>
                  <template v-else>
                    {{ one.name }}
                  </template>
                </div>
                <div class="file-size" v-if="!one.isDir">
                  <span class="file-size-unitSize">
                    {{ one.unitSize }}
                  </span>
                  <span class="file-size-unit">
                    {{ one.unit }}
                  </span>
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
  props: ["source", "dir", "files", "place", "wrap"],
  data() {
    return {
      ready: false,
      list: [],
      dragFile: {
        show: false,
        style: {},
        name: "上传文件",
      },
      form: {
        dir: null,
      },
    };
  },
  computed: {},
  watch: {
    dir() {
      this.form.dir = this.dir;
    },
    files() {
      this.formatFiles();
    },
  },
  methods: {
    async init(e) {
      this.form.dir = this.dir;
      this.formatFiles();
      this.ready = true;
      this.$nextTick(() => {
        this.initEvent();
      });
    },
    inputOnEnter(e) {
      e = e || window.event;
      var charCode = e.charCode ? e.charCode : e.which ? e.which : e.keyCode;
      if (charCode == 13 || charCode == 3) {
        this.openDir(this.form.dir);
      }
    },
    getFile(path) {
      let file = null;
      if (this.list) {
        this.list.forEach((one) => {
          if (one.path == path) {
            file = one;
          }
        });
      }
      return file;
    },
    getFileIndex(path) {
      let index = -1;
      if (this.list) {
        this.list.forEach((one, i) => {
          if (one == path || one.path == path) {
            index = i;
          }
        });
      }
      return index;
    },
    getRenameFile() {
      let file = null;
      if (this.list) {
        this.list.forEach((one) => {
          if (one.rename) {
            file = one;
          }
        });
      }
      return file;
    },
    getSelectFiles() {
      let files = [];
      if (this.list) {
        this.list.forEach((one) => {
          if (one.select) {
            files.push(one);
          }
        });
      }
      return files;
    },
    selectFile(file) {
      if (this.list) {
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
            this.list[fileIndex].select = true;
          }
          return;
        }
        this.list.forEach((one, i) => {
          if (one.name == "..") {
            one.select = false;
          } else if (i >= startIndex && i <= endIndex) {
            one.select = true;
          } else {
            one.select = false;
          }
        });
      }
    },
    filesBoxClick(e) {
      e = e || window.event;
      let file = this.getFileByTarget(e.target);
      if (file && file.select) {
        file.select = false;
        return;
      }
      this.selectFile(file);
    },
    getFileByTarget(target) {
      let file = null;
      if (target) {
        let box = this.tool.jQuery(target).closest(".file-box");
        if (box.length > 0) {
          let path = box.attr("path");
          file = this.getFile(path);
        }
      }
      return file;
    },
    fileContextmenu(e) {
      e = e || window.event;
      let file = this.getFileByTarget(e.target);
      if (file && !file.select) {
        this.selectFile(file);
      }
      let files = this.getSelectFiles();
      let menus = [];

      menus.push({
        text: "刷新",
        onClick: () => {
          this.toRefresh();
        },
      });
      if (files.length == 1) {
        menus.push({
          text: "重命名",
          onClick: () => {
            this.toRename(files[0]);
          },
        });
      }
      if (files.length > 0) {
        menus.push({
          text: "删除",
          onClick: () => {
            this.toRemove(files);
          },
        });
      }

      this.tool.showContextmenu(menus);
      // e.preventDefault();
    },
    ondragstart(e) {
      e = e || window.event;
      if (this.getRenameFile() != null) {
        e.preventDefault();
        return;
      }
      let file = this.getFileByTarget(e.target);
      // e.preventDefault();
      if (file != null) {
        if (!file.select) {
          this.selectFile(file);
        }
        let files = this.getSelectFiles();

        e.dataTransfer.setData("files", JSON.stringify(files));
      } else {
        e.dataTransfer.setData("files", "");
      }
      e.dataTransfer.setData("place", this.place);
      e.dataTransfer.setData("dir", this.dir);
    },
    ondragend(e) {
      // e.preventDefault();
      console.log("ondragend", e);
    },
    ondrag(e) {
      // e.preventDefault();
      console.log("ondrag", e);
    },
    ondragleave(e) {
      this.dragFile.show = false;
      e.preventDefault();
    },
    ondragover(e) {
      e.preventDefault();
    },
    ondrop(e) {
      this.dragFile.show = false;
      e.preventDefault();
      let putFile = this.getFileByTarget(e.target);
      let putDir = null;
      if (putFile != null && putFile.isDir) {
        putDir = putFile.path;
      } else {
        putDir = this.dir;
      }
      // console.log("ondrop", e);
      let files = e.dataTransfer.getData("files");
      if (this.tool.isNotEmpty(files)) {
        let place = e.dataTransfer.getData("place");
        let dir = e.dataTransfer.getData("dir");
        if (place == this.place && dir == putDir) {
          return;
        }
        let files = e.dataTransfer.getData("files");
        files = JSON.parse(files);
        if (place != this.place) {
          files.forEach((one) => {
            this.$emit(
              "copy",
              { place: place, dir: dir, path: one.path },
              {
                place: this.place,
                dir: putDir,
                path: putDir + "/" + one.name,
              }
            );
          });
        } else {
          files.forEach((one) => {
            this.$emit(
              "rename",
              this.place,
              dir,
              one.path,
              putDir + "/" + one.name
            );
          });
        }
      } else if (
        e.dataTransfer &&
        e.dataTransfer.items &&
        e.dataTransfer.items.length > 0
      ) {
        Array.prototype.forEach.call(e.dataTransfer.items, async (one) => {
          if (one.webkitGetAsEntry) {
            let webkitGetAsEntry = one.webkitGetAsEntry();
            this.uploadEntryFile(putDir, webkitGetAsEntry);
            return;
          }
          let file = one.getAsFile();
          if (file != null) {
            this.$emit("upload", this.place, putDir, file);
          } else {
            console.log(one);
          }
        });
      }
    },
    uploadEntryFile(putDir, entry) {
      if (entry.isFile) {
        entry.file(
          (file) => {
            this.uploadFile(putDir, file, entry.fullPath);
          },
          (e) => {
            console.log(e);
          }
        );
      } else {
        let reader = entry.createReader();
        reader.readEntries(
          (entries) => {
            entries.forEach((entry) => this.uploadEntryFile(putDir, entry));
          },
          (e) => {
            console.log(e);
          }
        );
      }
    },
    uploadFile(putDir, file, fullPath) {
      let dir = putDir;
      this.$emit("upload", this.place, dir, file, fullPath);
    },
    toRefresh() {
      this.$emit("refresh", this.place, this.dir);
    },
    openFile(file) {
      this.$emit("open", this.place, this.dir, file);
    },
    openDir(dir) {
      this.$emit("openDir", this.place, dir);
    },
    toRemove(files) {
      this.$emit("remove", this.place, this.dir, files);
    },
    toRename(file) {
      file.newname = file.name;
      file.rename = true;
      this.toFocusFile(file);
      // this.$emit("rename", this.place, this.dir, file);
    },
    toFocusFile(file) {
      this.$nextTick(() => {
        if (this.$el.getElementsByClassName("file-rename-input")[0]) {
          setTimeout(() => {
            this.$el.getElementsByClassName("file-rename-input")[0].focus();
          }, 100);
        }
      });
    },
    toBlurFile(file) {
      this.$nextTick(() => {
        if (this.$el.getElementsByClassName("file-rename-input")[0]) {
          this.$el.getElementsByClassName("file-rename-input")[0].blur();
        }
      });
    },
    onRenameBlur(file, event) {
      this.doRename(file);
    },
    onRenameKeyup(file, event) {
      event = event || window.event;
      if (event.keyCode == 13 || event.keyCode == 27) {
        this.toBlurFile(file);
      }
    },
    doRename(file) {
      if (this.tool.isEmpty(file.newname) || file.name == file.newname) {
        file.rename = false;
        return;
      }
      if (file.newname.indexOf("/") >= 0 || file.newname.indexOf("\\") >= 0) {
        this.tool.error("文件名输入有误，请重新输入！");
        this.toFocusFile(file);
        return;
      }
      file.rename = false;
      this.$emit(
        "rename",
        this.place,
        this.dir,
        file.path,
        this.dir + "/" + file.newname
      );
    },
    formatFiles() {
      let files = this.files || [];
      let list = [];
      files.forEach((one) => {
        one = Object.assign({}, one);
        one.dir = this.dir;
        one.path = one.dir + "/" + one.name;
        one.select = false;
        one.rename = false;
        one.size = one.size || 0;
        if (!one.isDir) {
          this.wrap.formatSize(one, "size", "unitSize", "unit");
        }
        list.push(one);
      });
      this.list = list;
    },
    initEvent() {
      this.$nextTick(() => {
        this.$refs.filesBox.ondragover = this.ondragover;
        this.$refs.filesBox.ondragleave = this.ondragleave;
        this.$refs.filesBox.ondrop = this.ondrop;
        this.$refs.filesBox.ondragstart = this.ondragstart;
      });
      // this.$el.ondragend = this.ondragend;
      // this.$el.ondrag = this.ondrag;
    },
  },
  created() {},
  mounted() {
    document.ondragover = function (e) {
      e.preventDefault(); //只有在ondragover中阻止默认行为才能触发 ondrop 而不是 ondragleave
    };
    document.ondrop = function (e) {
      e.preventDefault(); //阻止 document.ondrop的默认行为  *** 在新窗口中打开拖进的图片
    };
    this.init();
  },
  beforeDestroy() {},
};
</script>

<style>
.toolbox-ftp-files {
  width: 100%;
  height: 100%;
  user-select: none;
  position: relative;
}
.files-box {
  width: 100%;
  height: 100%;
  user-select: none;
}
.file-box {
  display: flex;
  line-height: 20px;
  font-size: 15px;
  cursor: context-menu;
  user-select: none;
  padding: 0px 5px;
}
.file-box:hover {
  background: #4a4a4a;
}
.file-box.file-select {
  background: #4f4f4f;
}
.file-box .file-icon {
  padding: 0px 5px;
}
.file-box .file-name {
  padding: 0px 5px;
  flex: 1;
}
.file-box .file-size {
  padding: 0px 5px;
  font-size: 12px;
}
.drag-file {
  position: absolute;
  left: 0px;
  right: 0px;
}
.drag-file .drag-file-name {
  padding: 0px 5px;
  font-size: 15px;
}
.file-dir-input-box {
  padding: 5px 10px;
}
.file-dir-input {
  color: #ffffff;
  width: 100%;
  border: 0px dashed #ddd;
  border-bottom: 1px solid #ddd;
  background-color: transparent;
  height: 30px;
  line-height: 30px;
  padding-left: 0px;
  padding-right: 0px;
  box-sizing: border-box;
  outline: none;
  font-size: 15px;
}
.file-rename-input {
  color: #ffffff;
  width: 100%;
  border: 1px solid #ddd;
  background-color: transparent;
  line-height: 16px;
  box-sizing: border-box;
  outline: none;
}
</style>
