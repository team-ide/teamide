<template>
  <el-dialog
    ref="modal"
    :title="`文件上传`"
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
        class="terminal-upload-box"
        @drop="onDrop"
        @dragover="onDragover"
        @dragleave="onDragleave"
        draggable="true"
      >
        <i class="mdi mdi-upload ft-20"></i>
        <div class="mgt-10 ft-16">
          将文件拖到此处 ，或
          <div class="tm-link color-green" @click="toClickUpload">点击上传</div>
        </div>
      </div>
    </div>
    <input
      type="file"
      multiple
      id="input-for-upload"
      @change="uploadInputChange"
      ref="input-for-upload"
    />
  </el-dialog>
</template>

<script>
import Zmodem from "zmodem.js";
export default {
  components: {},
  props: ["source", "toolboxWorker", "worker"],
  data() {
    return {
      showDialog: false,
      isFileSelect: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    show(zsession, term, callback) {
      this.isFileSelect = false;
      this.zsession = zsession;
      this.term = term;
      this.callback = callback;
      this.showDialog = true;
    },
    hide() {
      if (!this.isFileSelect) {
        // this.tool
        //   .confirm("是否取消上传？")
        //   .then(async () => {
        this.showDialog = false;
        this.onCancel && this.onCancel();
        // })
        // .catch((e) => {});
      } else {
        this.showDialog = false;
      }
      this.callback && this.callback();
    },
    onSuccess() {
      this.tool.success("文件上传完成");
      this.callback && this.callback();
    },
    onCancel() {
      this.term.write("\r\n取消上传\r\n");
      this.tool.warn("取消上传");
      // zsession 每 5s 发送一个 ZACK 包，5s 后会出现提示最后一个包是 ”ZACK“ 无法正常关闭
      // 这里直接设置 _last_header_name 为 ZRINIT，就可以强制关闭了
      this.zsession._last_header_name = "ZRINIT";
      this.zsession.close();
      // this.toolboxWorker.writeEvent("shell cancel upload file");
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
      this.isFileSelect = true;
      this.showDialog = false;

      this.send_block_files(this.zsession, files, {
        on_offer_response: (obj, xfer) => {
          if (xfer) {
            this.term.write("\r\n");
          } else {
            this.term.write(
              "\r\n文件(" +
                obj.upload_number +
                "/" +
                obj.upload_file_count +
                ") " +
                obj.name +
                " 已跳过（请检查文件是否已存在）\r\n"
            );
            if (obj.upload_number == obj.upload_file_count) {
              this.onSuccess();
            }
          }
        },
        on_progress: (obj, xfer) => {
          this.updateProgress(xfer);
        },
        on_file_complete: (obj) => {
          this.term.write(
            "\r\n文件(" +
              obj.upload_number +
              "/" +
              obj.upload_file_count +
              ") " +
              obj.name +
              " 上传成功\r\n"
          );
          if (obj.upload_number == obj.upload_file_count) {
            this.onSuccess();
          }
        },
      })
        .then(this.zsession.close.bind(this.zsession))
        .then(() => {
          // res();
        });
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
    bytesHuman(bytes, precision) {
      if (!/^([-+])?|(\.\d+)(\d+(\.\d+)?|(\d+\.)|Infinity)$/.test(bytes)) {
        return "-";
      }
      if (bytes === 0) return "0";
      if (typeof precision === "undefined") precision = 2;
      const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB"];
      const num = Math.floor(Math.log(bytes) / Math.log(1024));
      const value = (bytes / Math.pow(1024, Math.floor(num))).toFixed(
        precision
      );
      let res = `${value} ${units[num]}`;
      let fSize = 10 - res.length;
      for (let i = 0; i < fSize; i++) {
        res = " " + res;
      }
      return res;
    },
    updateProgress(xfer) {
      let detail = xfer.get_details();
      let name = detail.name;
      let total = detail.size;
      let offset = xfer.get_offset();
      let percent;

      if (total === 0 || total <= offset) {
        percent = 100;
      } else {
        percent = ((offset / total) * 100).toFixed(2);
      }
      let percentStr = "" + percent;
      let fSize = 10 - percentStr.length;
      for (let i = 0; i < fSize; i++) {
        percentStr = " " + percentStr;
      }
      this.term.write(
        "\r" +
          "上传文件(" +
          detail.obj.upload_number +
          "/" +
          detail.obj.upload_file_count +
          ") " +
          name +
          " " +
          this.bytesHuman(total) +
          " " +
          this.bytesHuman(offset) +
          " " +
          percentStr +
          "%"
      );
    },
    send_block_files(session, files, options) {
      if (!options) options = {};

      //Populate the batch in reverse order to simplify sending
      //the remaining files/bytes components.
      var batch = [];
      var total_size = 0;
      for (var f = files.length - 1; f >= 0; f--) {
        files[f].upload_number = f + 1;
        files[f].upload_file_count = files.length;
        var fobj = files[f];
        total_size += fobj.size;
        batch[f] = {
          obj: fobj,
          name: fobj.name,
          size: fobj.size,
          mtime: new Date(fobj.lastModified),
          files_remaining: files.length - f,
          bytes_remaining: total_size,
        };
      }

      var file_idx = 0;
      function promise_callback() {
        var cur_b = batch[file_idx];

        if (!cur_b) {
          return Promise.resolve(); //batch done!
        }

        file_idx++;

        return session.send_offer(cur_b).then((xfer) => {
          if (options.on_offer_response) {
            options.on_offer_response(cur_b.obj, xfer);
          }

          if (xfer === undefined) {
            return promise_callback(); //skipped
          }

          return new Promise(function (res) {
            var block = 1 * 1024;
            var fileSize = cur_b.size;
            var fileLoaded = 0;
            var reader = new FileReader();
            reader.onerror = function reader_onerror(e) {
              throw "File read error: " + e;
            };
            function readBlob() {
              var blob;
              if (cur_b.obj.slice) {
                blob = cur_b.obj.slice(fileLoaded, fileLoaded + block + 1);
              } else if (cur_b.obj.mozSlice) {
                blob = cur_b.obj.mozSlice(fileLoaded, fileLoaded + block + 1);
              } else if (cur_b.obj.webkitSlice) {
                blob = cur_b.obj.webkitSlice(
                  fileLoaded,
                  fileLoaded + block + 1
                );
              } else {
                blob = cur_b.obj;
              }
              reader.readAsArrayBuffer(blob);
            }
            var piece;
            reader.onload = function reader_onload(e) {
              fileLoaded += e.total;
              if (fileLoaded < fileSize) {
                if (e.target.result) {
                  if (session.aborted()) {
                    throw new Zmodem.Error("aborted");
                  }
                  piece = new Uint8Array(e.target.result);
                  xfer.send(piece);
                  if (options.on_progress && piece.length) {
                    options.on_progress(cur_b.obj, xfer, piece);
                  }
                }
                readBlob();
              } else {
                //
                if (e.target.result) {
                  piece = new Uint8Array(e.target.result);
                  if (session.aborted()) {
                    throw new Zmodem.Error("aborted");
                  }
                  xfer.end(piece).then(function () {
                    if (options.on_progress && piece.length) {
                      options.on_progress(cur_b.obj, xfer, piece);
                    }
                    if (options.on_file_complete) {
                      options.on_file_complete(cur_b.obj, xfer);
                    }
                    res(promise_callback());
                  });
                }
              }
            };
            readBlob();
          });
        });
      }

      return promise_callback();
    },
  },
  created() {},
  mounted() {
    this.toolboxWorker.showUpload = this.show;
    this.toolboxWorker.hideUpload = this.show;
  },
};
</script>

<style>
.terminal-upload-box {
  width: 100%;
  margin: 0px auto;
  text-align: center;
  border: 2px dashed #ddd;
  padding: 200px;
}
#input-for-upload {
  width: 0px;
  height: 0px;
  position: fixed;
  left: -100px;
  top: -100px;
}
</style>
