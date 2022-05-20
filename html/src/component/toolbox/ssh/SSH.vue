<template>
  <div class="toolbox-ssh-editor">
    <div
      class="terminal-box"
      ref="terminal"
      style="
        padding-top: 0px;
        padding-right: 0px;
        padding-bottom: 0px;
        padding-left: 0px;
      "
    />
    <div class="terminal-box-back" ref="terminal_back" />
    <SSHUpload :source="source" :wrap="wrap" :token="token"></SSHUpload>
  </div>
</template>


<script>
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";

import Zmodem from "zmodem.js";
console.log(Zmodem);
import SSHUpload from "./SSHUpload.vue";
export default {
  components: { SSHUpload },
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
      rows: 40,
      cols: 100,
      style: {
        width: null,
        height: null,
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.initTerminal();
    },
    onFocus() {
      this.term && this.term.focus();
    },
    refresh() {},
    onEvent(event) {
      if (event == "shell ready") {
        this.toStart();
      } else if (event == "shell created") {
        this.$nextTick(() => {
          this.initAttachAddon();
        });
      } else if (event == "shell to upload file") {
      }
    },
    onData(data) {
      if (typeof data === "object") {
        this.zsentry.consume(data);
      } else {
        this.term.write(typeof data === "string" ? data : new Uint8Array(data));
      }
    },
    bytesHuman(bytes, precision) {
      if (!/^([-+])?|(\.\d+)(\d+(\.\d+)?|(\d+\.)|Infinity)$/.test(bytes)) {
        return "-";
      }
      if (bytes === 0) return "0";
      if (typeof precision === "undefined") precision = 1;
      const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB"];
      const num = Math.floor(Math.log(bytes) / Math.log(1024));
      const value = (bytes / Math.pow(1024, Math.floor(num))).toFixed(
        precision
      );
      return `${value} ${units[num]}`;
    },
    updateProgress(xfer, action = "upload") {
      let detail = xfer.get_details();
      let name = detail.name;
      let total = detail.size;
      let offset = xfer.get_offset();
      let percent;
      if (total === 0 || total === offset) {
        percent = 100;
      } else {
        percent = Math.round((offset / total) * 100);
      }
      // term.write("\r" + name + ": " + total + " " + offset + " " + percent + "% ");
      this.term.write(
        "\r" +
          action +
          " " +
          name +
          ": " +
          this.bytesHuman(offset) +
          " " +
          this.bytesHuman(total) +
          " " +
          percent +
          "% "
      );
    },
    uploadFile(zsession) {
      // https://github.com/leffss/gowebssh/blob/c594d8cfc48b5bb64489431fea57d39e6d159359/example/html/zmodem/zmodem.devel.js#L922
      this.wrap.showSSHUpload(
        (files) => {
          //Zmodem.Browser.send_files(zsession, files, {
          this.send_block_files(zsession, files, {
            on_offer_response: (obj, xfer) => {
              if (xfer) {
                this.term.write("\r\n");
              } else {
                this.term.write(obj.name + " was upload skipped\r\n");
                // socket.send(JSON.stringify({ type: "ignore", data: utf8_to_b64("\r\n" + obj.name + " was upload skipped\r\n") }));
              }
            },
            on_progress: (obj, xfer) => {
              this.updateProgress(xfer);
            },
            on_file_complete: (obj) => {
              this.term.write("\r\n");
              // socket.send(
              //   JSON.stringify({
              //     type: "ignore",
              //     data: utf8_to_b64(
              //       obj.name + "(" + obj.size + ") was upload success"
              //     ),
              //   })
              // );
            },
          })
            .then(zsession.close.bind(zsession), console.error.bind(console))
            .then(() => {
              // res();
            });
        },
        () => {
          this.tool.warn("取消上传");
          // zsession 每 5s 发送一个 ZACK 包，5s 后会出现提示最后一个包是 ”ZACK“ 无法正常关闭
          // 这里直接设置 _last_header_name 为 ZRINIT，就可以强制关闭了
          zsession._last_header_name = "ZRINIT";
          zsession.close();
          // this.wrap.writeEvent("shell cancel upload file");
        }
      );
    },
    send_block_files(session, files, options) {
      if (!options) options = {};

      //Populate the batch in reverse order to simplify sending
      //the remaining files/bytes components.
      var batch = [];
      var total_size = 0;
      for (var f = files.length - 1; f >= 0; f--) {
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

        return session.send_offer(cur_b).then(function after_send_offer(xfer) {
          if (options.on_offer_response) {
            options.on_offer_response(cur_b.obj, xfer);
          }

          if (xfer === undefined) {
            return promise_callback(); //skipped
          }

          return new Promise(function (res) {
            var block = 1024 * 1024;
            var fileSize = cur_b.size;
            var fileLoaded = 0;
            var reader = new FileReader();
            reader.onerror = function reader_onerror(e) {
              console.error("file read error", e);
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
                  piece = new Uint8Array(e.target.result);
                  if (session.aborted()) {
                    throw new Zmodem.Error("aborted");
                  }
                  xfer.send(piece);
                  if (options.on_progress) {
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
    onError(error) {
      this.tool.error(error);
    },
    toStart() {
      let data = {};
      data.cols = this.cols;
      data.rows = this.rows;
      data.width = 0;
      data.height = 0;
      this.wrap.writeEvent("shell start" + JSON.stringify(data));
    },
    changeSize() {
      if (this.term == null) {
        return;
      }
      let data = {};
      data.cols = this.cols;
      data.rows = this.rows;
      data.width = 0;
      data.height = 0;
      this.wrap.writeEvent("change size" + JSON.stringify(data));
    },
    initAttachAddon() {
      this.term.onData((data) => {
        this.wrap.writeData(data);
      });
      this.term.onBinary((data) => {
        this.wrap.writeData(data);
      });

      this.zsentry = new Zmodem.Sentry({
        to_terminal: (octets) => {}, //i.e. send to the terminal
        on_detect: (detection) => {
          let zsession = detection.confirm();
          if (zsession.type === "receive") {
            this.downloadFile(zsession);
          } else {
            this.uploadFile(zsession);
          }
        },
        on_retract: () => {},
        sender: (octets) => {
          this.socket.send(new Uint8Array(octets));
        },
      });
      console.log(this.zsentry);
      // this.attachAddon = new AttachAddon(this.socket);
      // this.term.loadAddon(this.attachAddon);
    },
    initTerminal() {
      if (this.term != null) {
        this.term.dispose();
      }
      this.term = new Terminal({
        useStyle: true,
        cursorBlink: true, //光标闪烁
        cursorStyle: "bar", // 光标样式 'block' | 'underline' | 'bar'
        rendererType: "canvas", //渲染类型
        width: 500,
        height: 400,
        // rows: this.rows, //行数
        // cols: this.cols, // 不指定行数，自动回车后光标从下一行开始
        // convertEol: true, //启用时，光标将设置为下一行的开头
        // // scrollback: 50, //终端中的回滚量
        // disableStdin: false, //是否应禁用输入
        // // cursorStyle: "underline", //光标样式
        // theme: {
        //   foreground: "#ECECEC", //字体
        //   background: "#000000", //背景色
        //   cursor: "help", //设置光标
        //   lineHeight: 20,
        // },
      });
      this.term.open(this.$refs.terminal, true);

      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.fitAddon.fit();

      this.term.focus();
      this.cols = this.term.cols;
      this.rows = this.term.rows;
      this.initSize();

      // this.term.on("resize", () => {
      //   this.cols = this.term.cols;
      //   this.rows = this.term.rows;
      //   this.changeSize();
      // });
      this.changeSizeTimer();
    },
    changeSizeTimer() {
      if (this.term == null) {
        return;
      }
      if (
        this.tool.jQuery(this.$refs.terminal_back).width() !=
          this.style.width ||
        this.tool.jQuery(this.$refs.terminal_back).height() != this.style.height
      ) {
        this.style.width = this.tool.jQuery(this.$refs.terminal_back).width();
        this.style.height = this.tool.jQuery(this.$refs.terminal_back).height();
        this.tool.jQuery(this.term.element).css({
          width: parseInt(this.style.width),
          height: parseInt(this.style.height),
        });
        // console.log(this.term.element);

        this.fitAddon.fit();

        if (this.term.cols != this.cols || this.term.rows != this.rows) {
          this.cols = this.term.cols;
          this.rows = this.term.rows;
          this.changeSize();
        }
      }

      // window.setTimeout(() => {
      // 窗口尺寸变化时，终端尺寸自适应

      window.setTimeout(() => {
        this.changeSizeTimer();
      }, 200);
      // }, 100);
    },
    initSize() {},
    onresize() {
      // window.setTimeout(() => {
      // this.initSize();
      // }, 500);
    },
    dispose() {
      if (this.term != null) {
        this.term.dispose();
      }
      this.term = null;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.dispose();
  },
};
</script>

<style>
.toolbox-ssh-editor {
  width: 100%;
  height: 100%;
  position: relative;
}
.toolbox-ssh-editor .terminal-box {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  position: absolute;
  top: 0px;
  left: 0px;
}
.toolbox-ssh-editor .terminal-box-back {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  position: absolute;
  top: 0px;
  left: 0px;
}
</style>
