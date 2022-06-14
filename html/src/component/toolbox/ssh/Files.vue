<template>
  <div class="toolbox-ftp-files" tabindex="-1">
    <template v-if="ready">
      <tm-layout height="100%">
        <tm-layout height="70px">
          <div class="pdtb-5 pdlr-10">
            <template v-if="isInputDir">
              <input
                ref="file-dir-input"
                class="file-dir-input"
                v-model="form.dir"
                title="回车确认"
                @keyup.enter="inputOnEnter"
                @blur="setIsInputDir(false)"
              />
            </template>
            <template v-else>
              <div
                class="toolbox-dir-names-breadcrumb"
                @click="setIsInputDir(true)"
              >
                <el-breadcrumb separator="/">
                  <template v-for="(one, index) in dirNames">
                    <el-breadcrumb-item :key="index">
                      <span @click="toDirName(index)">{{ one }}</span>
                    </el-breadcrumb-item>
                  </template>
                </el-breadcrumb>
              </div>
            </template>
          </div>
          <div class="pdlr-10">
            <el-form size="mini" inline @submit.native.prevent>
              <el-form-item label="Filter" class="mgb-0">
                <el-input v-model="form.filter" style="width: 150px">
                </el-input>
              </el-form-item>
            </el-form>
          </div>
        </tm-layout>
        <tm-layout height="auto">
          <div
            class="files-box scrollbar"
            @contextmenu.prevent="fileContextmenu"
            ref="filesBox"
            @click="filesBoxClick"
          >
            <template v-if="files == null">
              <div class="text-center pdtb-10 ft-15">加载中...</div>
            </template>
            <template v-for="(one, index) in list">
              <div
                :key="'file-' + index"
                v-if="one.show"
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
                <div class="file-date">
                  {{ one.dateTime }}
                </div>
                <div class="file-mode">
                  {{ one.fileMode }}
                </div>
                <div class="file-size">
                  <template v-if="!one.isDir">
                    <span class="file-size-unitSize">
                      {{ one.unitSize }}
                    </span>
                    <span class="file-size-unit">
                      {{ one.unit }}
                    </span>
                  </template>
                </div>
              </div>
            </template>
          </div>
        </tm-layout>
      </tm-layout>
    </template>
    <input
      ref="chooseFolder"
      type="file"
      name="file"
      id="file"
      style="
        width: 0px;
        height: 0px;
        position: fixed;
        left: -100px;
        top: -100px;
        z-index: -1000;
      "
      multiple
    />
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
      form: {
        dir: null,
        pattern: null,
        filter: null,
      },
      selectPaths: [],
      dirNames: [],
      isInputDir: false,
    };
  },
  computed: {},
  watch: {
    dir() {
      this.form.dir = this.dir;
      this.dirNames = ("" + this.dir).split("/");
    },
    files() {
      this.formatFiles();
    },
    "form.filter"() {
      this.filterFile();
    },
  },
  methods: {
    filterFile() {
      this.selectPaths = this.selectPaths || [];
      this.list.forEach((one) => {
        if (
          this.tool.isEmpty(this.form.filter) ||
          one.name == ".." ||
          one.name.toLowerCase().indexOf(this.form.filter.toLowerCase()) >= 0
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
    setScrollTop(scrollTop) {
      if (scrollTop >= 0) {
        this.$nextTick(() => {
          this.tool.jQuery(this.$refs.filesBox).scrollTop(scrollTop);
        });
      }
    },
    getScrollTop() {
      var scrollTop = parseInt(
        this.tool.jQuery(this.$refs.filesBox).scrollTop()
      );
      return scrollTop;
    },
    async init(e) {
      this.form.dir = this.dir;
      this.dirNames = ("" + this.dir).split("/");
      this.formatFiles();
      this.ready = true;
      this.$nextTick(() => {
        this.initEvent();
      });
    },
    onFocus() {
      this.$el.focus();
    },
    refresh() {
      this.toRefresh();
    },
    setIsInputDir(isInputDir) {
      this.isInputDir = isInputDir;
      if (isInputDir) {
        this.$nextTick(() => {
          this.$refs["file-dir-input"].focus();
        });
      }
    },
    toDirName(dirNameIndex) {
      this.tool.stopEvent();
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
    toSearch() {
      this.openDir(this.form.dir, this.form.pattern);
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
    toChooseFolder() {
      this.$refs.chooseFolder.click();
    },
    toSaveFile(file) {
      this.toChooseFolder();
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
            this.setSelect(this.list[fileIndex]);
          }
          return;
        }
        this.list.forEach((one, i) => {
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
    filesBoxClick(e) {
      e = e || window.event;
      let file = this.getFileByTarget(e.target);
      if (file && file.select) {
        this.setUnselect(file);
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
        menus.push({
          text: "复制文件名",
          onClick: async () => {
            let res = await this.tool.clipboardWrite(files[0].name);
            if (res.success) {
              this.tool.success("复制成功");
            } else {
              this.tool.warn("复制失败，请允许访问剪贴板！");
            }
          },
        });
        menus.push({
          text: "复制文件路径",
          onClick: async () => {
            let res = await this.tool.clipboardWrite(files[0].path);
            if (res.success) {
              this.tool.success("复制成功");
            } else {
              this.tool.warn("复制失败，请允许访问剪贴板！");
            }
          },
        });
        if (!files[0].isDir) {
          // menus.push({
          //   text: "另存为",
          //   onClick: () => {
          //     this.toSaveFile(files[0]);
          //   },
          // });

          menus.push({
            text: "在线编辑",
            onClick: () => {
              this.openFile(files[0]);
            },
          });
          menus.push({
            text: "下载",
            onClick: () => {
              this.toDownload(files[0]);
            },
          });
        }
      }
      menus.push({
        text: "新建文件夹",
        onClick: () => {
          this.toInsertFile(true, files[0]);
        },
      });
      menus.push({
        text: "新建文件",
        onClick: () => {
          this.toInsertFile(false, files[0]);
        },
      });
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
      e.preventDefault();
    },
    ondragover(e) {
      e.preventDefault();
    },
    ondrop(e) {
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
          let names = [];
          files.forEach((one) => {
            names.push(one.name);
          });
          this.tool
            .confirm(
              "移动[" +
                names.join(",") +
                "]到[" +
                putDir +
                "/" +
                "]后无法恢复，确定移动？"
            )
            .then(async () => {
              files.forEach((one) => {
                this.$emit(
                  "rename",
                  this.place,
                  dir,
                  one.path,
                  putDir + "/" + one.name
                );
              });
            })
            .catch((e) => {});
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
    openDir(dir, pattern) {
      this.$emit("openDir", this.place, dir, pattern);
    },
    toRemove(files) {
      this.$emit("remove", this.place, this.dir, files);
    },
    toDownload(file) {
      this.$emit("download", this.place, this.dir, file);
    },
    toInsertFile(isDir, afterFile) {
      let index = this.list.indexOf(afterFile);
      if (index < 0) {
        index = this.list.length - 1;
      }
      let newFile = {};
      newFile.name = "新建文件";
      if (isDir) {
        newFile.name = "新建文件夹";
      }
      newFile.isDir = isDir;
      newFile.dir = this.dir;
      newFile.path = newFile.dir + "/" + newFile.name;
      newFile.rename = false;
      newFile.size = newFile.size || 0;
      newFile.isNew = true;
      newFile.show = true;
      this.list.splice(index + 1, 0, newFile);
      this.$nextTick(() => {
        this.toRename(newFile);
      });
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
      if (this.tool.isEmpty(file.newname) && file.isNew) {
        let index = this.list.indexOf(file);
        this.list.splice(index, 1);
        return;
      }
      if (this.tool.isEmpty(file.newname) || file.name == file.newname) {
        if (!file.isNew) {
          file.rename = false;
          return;
        }
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
        this.dir + "/" + file.newname,
        file.isDir,
        file.isNew
      );
    },
    formatFiles() {
      if (this.list) {
        this.list.splice(0, this.list.length);
      }
      if (this.files) {
        this.files.forEach((one) => {
          one = Object.assign({}, one);
          one.dir = this.dir;
          one.path = one.dir + "/" + one.name;
          one.rename = false;
          one.size = one.size || 0;
          one.isNew = false;
          one.show = true;
          if (!one.isDir) {
            this.tool.formatSize(one, "size", "unitSize", "unit");
          }
          if (one.modTime) {
            one.dateTime = this.tool.formatDate(
              new Date(one.modTime),
              "yyyy-MM-dd hh:mm:ss"
            );
          }
          one.select = false;
          if (this.selectPaths.indexOf(one.path) >= 0) {
            one.select = true;
          }
          this.list.push(one);
        });
        this.filterFile();
        this.selectPaths = [];
      }
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
    onKeyDown() {
      let files = this.getSelectFiles();
      if (files == null || files.length == 0) {
        return;
      }
      if (this.tool.keyIsF2()) {
        this.tool.stopEvent();
        this.toRename(files[0]);
      } else if (this.tool.keyIsDelete()) {
        this.tool.stopEvent();
        this.toRemove(files);
      }
    },
    bindEvent() {
      if (this.bindEvented) {
        return;
      }
      this.bindEvented = true;
      this.$el.addEventListener("keydown", (e) => {
        this.onKeyDown(e);
      });
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
    this.bindEvent();
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
  outline: none;
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
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.file-box .file-date {
  padding: 0px 0px;
  font-size: 12px;
  width: 130px;
}
.file-box .file-mode {
  padding: 0px 0px;
  font-size: 12px;
  width: 80px;
}
.file-box .file-size {
  padding: 0px 0px;
  font-size: 12px;
  width: 90px;
  text-align: right;
}
.file-box .file-size .file-size-unit {
  width: 20px;
  display: inline-block;
  text-align: left;
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
.file-dir-input {
  color: #ffffff;
  width: 100%;
  border: 0px dashed #ddd;
  border-bottom: 1px solid #ddd;
  background-color: transparent;
  height: 25px;
  line-height: 25px;
  padding-left: 0px;
  padding-right: 0px;
  box-sizing: border-box;
  outline: none;
  font-size: 15px;
}
.file-rename-input {
  color: #3c3c3c !important;
  width: 100%;
  border: 1px solid #ddd;
  background-color: #ddd !important;
  line-height: 16px;
  box-sizing: border-box;
  outline: none;
}
.toolbox-dir-names-breadcrumb {
  border: 0px dashed #ddd;
  border-bottom: 1px solid #ddd;
  height: 25px;
  line-height: 25px;
  font-size: 15px;
  cursor: text;
}
.toolbox-dir-names-breadcrumb .el-breadcrumb {
  height: 25px;
  line-height: 25px;
}
.toolbox-dir-names-breadcrumb .el-breadcrumb .el-breadcrumb__separator {
  margin: 0px;
  color: #ffffff;
}
.toolbox-dir-names-breadcrumb .el-breadcrumb .el-breadcrumb__inner {
  margin: 0px;
  color: #ffffff;
}
.toolbox-dir-names-breadcrumb .el-breadcrumb .el-breadcrumb__inner > span {
  cursor: pointer;
  display: inline-block;
}

.toolbox-dir-names-breadcrumb .el-breadcrumb .el-breadcrumb__inner:hover {
  color: #ddd;
}
</style>
