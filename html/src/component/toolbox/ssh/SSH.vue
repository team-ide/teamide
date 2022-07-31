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
    <div class="pdt-2 pdlr-10">
      <el-dropdown ref="quickCommandDropdown" size="mini" trigger="click">
        <span class="el-dropdown-link ft-12 mglr-5 color-orange tm-pointer">
          <span>快速指令</span>
        </span>
        <el-dropdown-menu
          slot="dropdown"
          class="pd-0 bd-0 toolbox-ssh-quickCommand-box"
        >
          <div class="toolbox-ssh-quickCommand-header">
            <div class="toolbox-ssh-quickCommand-search-box">
              <input
                class="toolbox-ssh-quickCommand-search"
                v-model="quickCommandSearch"
                placeholder="输入过滤"
              />
            </div>
            <div
              class="ft-12 tm-link color-green mgl-10 mgt-3"
              @click="toolboxWorker.toInsertSSHCommand()"
            >
              添加
            </div>
          </div>
          <div class="toolbox-ssh-quickCommand-list scrollbar">
            <template
              v-if="
                source.quickCommandSSHCommands != null &&
                source.quickCommandSSHCommands.length > 0
              "
            >
              <template v-for="(one, index) in source.quickCommandSSHCommands">
                <div
                  :key="index"
                  v-if="
                    tool.isEmpty(quickCommandSearch) ||
                    one.name
                      .toLowerCase()
                      .indexOf(quickCommandSearch.toLowerCase()) >= 0
                  "
                  class="toolbox-ssh-quickCommand-one"
                >
                  <div class="toolbox-ssh-quickCommand-name">
                    {{ one.name }}
                  </div>
                  <div class="toolbox-ssh-quickCommand-btn-group">
                    <div
                      class="ft-12 tm-link color-grey mgl-10"
                      @click="toExecCommand(one, false)"
                    >
                      填充，不执行
                    </div>
                    <div
                      class="ft-12 tm-link color-orange mgl-10"
                      @click="toExecCommand(one, true)"
                    >
                      填充，并执行
                    </div>
                    <div
                      class="ft-12 tm-link color-blue mgl-10"
                      @click="toolboxWorker.toUpdateSSHCommand(one)"
                    >
                      修改
                    </div>
                    <div
                      class="ft-12 tm-link color-red mgl-10"
                      @click="toolboxWorker.toDeleteSSHCommand(one)"
                    >
                      删除
                    </div>
                  </div>
                </div>
              </template>
            </template>
            <template v-else>
              <div class="pd-20 pdlr-50 ft-15 color-grey text-center">
                暂无指令
              </div>
            </template>
          </div>
        </el-dropdown-menu>
      </el-dropdown>
      <div class="ft-12 tm-link color-grey mglr-5" @click="openFtpWindow()">
        FTP
      </div>
    </div>
    <SSHUpload :source="source" :toolboxWorker="toolboxWorker"></SSHUpload>
    <SSHDownload :source="source" :toolboxWorker="toolboxWorker"></SSHDownload>
    <template v-if="isOpenFTP">
      <div
        class="toolbox-ssh-editor-ftp-box"
        :class="{ 'toolbox-ssh-editor-ftp-box-show': isShowFTP }"
        :style="{
          width: `${ftpWidth}px`,
          height: `${ftpHeight}px`,
        }"
      >
        <div
          ref="ftpBoxTopLine"
          class="toolbox-ssh-editor-ftp-box-top-line"
        ></div>
        <div
          ref="ftpBoxLeftLine"
          class="toolbox-ssh-editor-ftp-box-left-line"
        ></div>
        <div class="toolbox-ssh-editor-ftp-box-header">
          <div class=""></div>
          <div
            style="
              display: inline-block;
              position: absolute;
              top: 0px;
              right: 3px;
            "
          >
            <span
              title="关闭"
              class="tm-link color-write mgr-0"
              @click="hideFTP()"
            >
              <i class="mdi mdi-close ft-21"></i>
            </span>
          </div>
        </div>
        <div class="toolbox-ssh-editor-ftp-box-body">
          <FTP
            :source="source"
            :toolboxWorker="toolboxWorker"
            :extend="extend"
            :isFromSSH="true"
          >
          </FTP>
        </div>
      </div>
    </template>
  </div>
</template>


<script>
import FTP from "./FTP";
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
// import { AttachAddon } from "xterm-addon-attach";

import Zmodem from "zmodem.js";
import SSHUpload from "./SSHUpload.vue";
import SSHDownload from "./SSHDownload.vue";
export default {
  components: { SSHUpload, SSHDownload, FTP },
  props: ["source", "extend", "toolboxWorker", "initToken", "initSocket"],
  data() {
    return {
      quickCommandSearch: null,
      rows: 40,
      cols: 100,
      style: {
        width: null,
        height: null,
      },
      ftpWidth: 900,
      ftpHeight: 600,

      isOpenFTP: false,
      isShowFTP: true,
      sshSessionClosed: true,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      await this.initTerminal();
    },
    onFocus() {
      this.term && this.term.focus();
    },
    refresh() {},
    openFtpWindow() {
      if (this.isOpenFTP && this.isShowFTP) {
        this.hideFTP();
      } else {
        this.isOpenFTP = true;
        this.isShowFTP = true;
      }
      this.$nextTick(this.bindDrapFTPEvent);
    },
    bindDrapFTPEvent() {
      if (this.bindDrapFTPEvented) {
        return;
      }
      this.bindDrapFTPEvented = true;
      let leftLine = this.$refs.ftpBoxLeftLine;
      leftLine.addEventListener("mousedown", (e) => {
        this.lineLeftClientX = e.clientX;
        document.addEventListener("mouseup", this.documentMouseupEvent);
        document.addEventListener("mousemove", this.documentMousemoveLeftEvent);
      });
      let topLine = this.$refs.ftpBoxTopLine;
      topLine.addEventListener("mousedown", (e) => {
        this.lineTopClientY = e.clientY;
        document.addEventListener("mouseup", this.documentMouseupEvent);
        document.addEventListener("mousemove", this.documentMousemoveTopEvent);
      });
    },
    documentMouseupEvent() {
      document.removeEventListener("mouseup", this.documentMouseupEvent);
      document.removeEventListener(
        "mousemove",
        this.documentMousemoveLeftEvent
      );
      document.removeEventListener("mousemove", this.documentMousemoveTopEvent);
    },
    documentMousemoveLeftEvent(e) {
      let clientX = e.clientX;
      this.ftpWidth =
        Number(this.ftpWidth) - Number(clientX - this.lineLeftClientX);
      this.lineLeftClientX = clientX;
    },
    documentMousemoveTopEvent(e) {
      let clientY = e.clientY;
      this.ftpHeight =
        Number(this.ftpHeight) - Number(clientY - this.lineTopClientY);
      this.lineTopClientY = clientY;
    },
    hideFTP() {
      this.isShowFTP = false;
    },
    onEvent(event) {
      if (event == "shell ready") {
        this.toStart();
      } else if (event == "shell created") {
        this.sshSessionClosed = false;
        delete this.startIng;
        this.$nextTick(() => {
          this.initAttachAddon();
          // this.tool.success("SSH会话连接成功");
        });
      } else if (event == "shell create error") {
        this.sshSessionClosed = true;
        delete this.startIng;
        this.term.write("\r\nSSH创建失败，输入回车重新连接！");
        this.$nextTick(() => {
          this.initAttachAddon();
          // this.tool.error("SSH会话连接失败");
        });
      } else if (event == "ssh session closed") {
        delete this.startIng;
        this.sshSessionClosed = true;
        this.term.write("\r\nSSH会话已关闭，输入回车重新连接！");
      } else if (event == "shell to upload file") {
      }
    },
    onData(data) {
      if (typeof data === "string") {
        this.term.write(data);
        return;
      } else if (data instanceof ArrayBuffer) {
        try {
          if (data.byteLength >= this.source.sshTeamIDEBinaryStartBytesLength) {
            let bs = data.slice(
              0,
              this.source.sshTeamIDEBinaryStartBytesLength
            );
            let eq = true;
            bs = new Uint8Array(bs);
            this.source.sshTeamIDEBinaryStartBytes.forEach((a, i) => {
              if (eq && a != bs[i]) {
                eq = false;
              }
            });
            if (eq) {
              data = data.slice(
                this.source.sshTeamIDEBinaryStartBytesLength,
                data.length
              );
              this.term.write(new Uint8Array(data));
              return;
            }
          }
        } catch (error) {}
        this.zsentry.consume(data);
      }
    },
    onError(error) {
      this.tool.error(error);
    },
    toExecCommand(quickCommand, exec) {
      if (quickCommand == null) {
        this.tool.error("快速指令为空");
        return;
      }
      let option = this.tool.getOptionJSON(quickCommand.option);
      if (option == null || this.tool.isEmpty(option.command)) {
        this.tool.error("未配置命令");
        return;
      }
      this.term && this.term.focus();
      let command = option.command;

      // this.term.write(command);
      this.writeData(command);

      if (exec) {
        // this.term.write(`\n`);
        this.writeData(`\n`);
      }
      if (
        this.$refs.quickCommandDropdown &&
        this.$refs.quickCommandDropdown.hide
      ) {
        this.$refs.quickCommandDropdown.hide();
      }
    },
    reStart() {
      this.toStart();
    },
    toStart() {
      this.startIng = true;
      let data = {};
      data.cols = this.cols;
      data.rows = this.rows;
      data.width = 0;
      data.height = 0;
      this.writeEvent("shell start" + JSON.stringify(data));
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
      this.writeEvent("change size" + JSON.stringify(data));
    },
    initAttachAddon() {
      if (this.initAttachAddoned) {
        return;
      }
      this.initAttachAddoned = true;
      this.term.onData((data) => {
        this.writeData(data);
      });
      this.term.onBinary((data) => {
        this.writeData(data);
      });

      this.zsentry = new Zmodem.Sentry({
        to_terminal: (octets) => {}, //i.e. send to the terminal
        on_detect: (detection) => {
          let zsession = detection.confirm();
          if (zsession.type === "receive") {
            this.toolboxWorker.showSSHDownload(zsession, this.term, () => {
              this.onFocus();
            });
          } else {
            this.toolboxWorker.showSSHUpload(zsession, this.term, () => {
              this.onFocus();
            });
          }
        },
        on_retract: () => {},
        sender: (octets) => {
          this.socket.send(new Uint8Array(octets));
        },
      });
      // this.attachAddon = new AttachAddon(this.socket);
      // this.term.loadAddon(this.attachAddon);
    },
    async onKeydown(e) {
      // let key = arg.key;
      // console.log(key);

      if (this.sshSessionClosed) {
        this.tool.stopEvent(e);
        if (this.startIng) {
          this.tool.warn("SSH会话连接中，请稍后");
          return;
        }
        if (this.tool.keyIsEnter(e)) {
          this.reStart();
          this.tool.warn("SSH会话连接中，请稍后");
          return;
        }
        this.tool.warn("SSH会话已关闭，请重新连接！");
        return;
      }
      if (this.tool.keyIsCtrlC(e)) {
        this.doEventCopy();
      } else if (this.tool.keyIsCtrlV(e)) {
        this.doEventPaste();
      }
    },
    async doEventCopy() {
      let copiedText = this.term.getSelection();
      if (this.tool.isNotEmpty(copiedText)) {
        this.tool.stopEvent();
        let res = await this.tool.clipboardWrite(copiedText);
        if (res.success) {
          this.tool.success("复制成功");
        } else {
          this.tool.warn("复制失败，请允许访问剪贴板！");
        }
      }
    },
    async doEventPaste() {
      let readResult = await this.tool.readClipboardText();
      if (readResult.success) {
        if (this.tool.isNotEmpty(readResult.text)) {
          this.tool.stopEvent();
          this.toPaste(readResult.text);
        }
      } else {
        this.tool.warn("粘贴失败，请允许访问剪贴板！");
      }
    },
    toPaste(text) {
      if (this.tool.isNotEmpty(text)) {
        if (text.indexOf("\n") >= 0) {
          let showText = text;
          let div = this.tool.jQuery("<div/>");

          let textarea = this.tool.jQuery(
            `<textarea readonly="readonly" style="width: 100%;height: 200px;overflow: auto;color: #a15656;margin-top: 15px;outline: 0px;border: 1px solid #ddd;padding: 5px;"/>`
          );
          textarea.append(showText);

          div.append("<div>确认粘贴以下内容<div/>");
          div.append(textarea);
          this.tool
            .confirm(div.html())
            .then(() => {
              this.writeData(showText);
              this.tool.success("粘贴成功");
            })
            .catch(() => {});
        } else {
          this.writeData(text);
          this.tool.success("粘贴成功");
        }
      }
    },
    onKeyup(e) {},
    async onMousedown(e) {
      // let event = e || window.event;
      // this.tool.stopEvent(e);
      // console.log(event);
    },
    async onContextmenu() {
      let copiedText = this.term.getSelection();
      if (this.tool.isNotEmpty(copiedText)) {
        this.doEventCopy();
      } else {
        this.doEventPaste();
      }
    },
    async onMouseup(e) {},
    async initTerminal() {
      if (this.term != null) {
        this.term.dispose();
      }
      await this.initToken(this);
      this.initSocket(this);
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
        windowsMode: true,
        scrollback: 100000000,
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
      this.$refs.terminal.addEventListener("keydown", this.onKeydown, true);
      this.$refs.terminal.addEventListener("mouseup", this.onMouseup, true);
      this.$refs.terminal.addEventListener("mousedown", this.onMousedown, true);
      this.$refs.terminal.addEventListener(
        "contextmenu",
        this.onContextmenu,
        true
      );
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
      if (this.isDestroyed) {
        return;
      }
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
      if (this.socket != null) {
        this.socket.close();
      }
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
  destroyed() {
    this.isDestroyed = true;
  },
};
</script>

<style>
.toolbox-ssh-editor {
  width: 100%;
  height: 100%;
  position: relative;
}

.toolbox-ssh-quickCommand-box {
  background-color: #363636;
  color: #ddd;
}
.toolbox-ssh-quickCommand-header {
  display: flex;
  margin: 5px 0px;
}
.toolbox-ssh-quickCommand-search-box {
  margin: 0px 10px;
  width: 200px;
}
.toolbox-ssh-quickCommand-search {
  width: 100%;
  height: 26px;
  line-height: 26px;
  border: 1px solid #ddd;
  font-size: 12px;
  outline: 0px;
  background-color: transparent;
}
.toolbox-ssh-quickCommand-list {
  height: 300px;
  width: 600px;
}
.toolbox-ssh-quickCommand-one {
  display: flex;
  padding: 2px 10px;
}
.toolbox-ssh-quickCommand-name {
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.toolbox-ssh-editor .terminal-box {
  width: 100%;
  height: calc(100% - 30px);
  position: relative;
  background-color: black;
}
.terminal-box-back {
  width: calc(100% - 20px) !important;
  height: calc(100% - 30px);
  position: absolute;
  left: 0px;
  top: 0px;
  z-index: -1;
}

.toolbox-ssh-editor-ftp-box {
  position: absolute;
  right: 20px;
  bottom: 40px;
  background: #172029;
  transition: all 0s;
  transform: scale(0);
  z-index: -1;
}

.toolbox-ssh-editor-ftp-box.toolbox-ssh-editor-ftp-box-show {
  transform: scale(1);
  z-index: 10;
}
.toolbox-ssh-editor-ftp-box-top-line {
  position: absolute;
  top: 0px;
  left: 0px;
  width: 100%;
  margin-top: -2px;
  height: 4px;
  background: #4e4e4e;
  cursor: row-resize;
  z-index: 1;
}
.toolbox-ssh-editor-ftp-box-left-line {
  position: absolute;
  top: 0px;
  left: 0px;
  height: 100%;
  width: 4px;
  margin-left: -2px;
  background: #4e4e4e;
  cursor: col-resize;
  z-index: 1;
}

.toolbox-ssh-editor-ftp-box-header {
  height: 29px;
  border-bottom: 1px solid #2f2f2f;
}
.toolbox-ssh-editor-ftp-box-body {
  height: calc(100% - 30px);
}
.toolbox-ssh-editor .terminal-box .terminal {
  width: 100% !important;
}
.toolbox-ssh-editor .terminal-box .xterm-viewport {
  width: 100% !important;
  background-color: transparent !important;
}
.toolbox-ssh-editor .terminal-box .xterm-screen {
  width: calc(100% - 20px) !important;
  margin: 0px 5px;
}
.toolbox-ssh-editor .terminal-box .xterm-text-layer {
  width: 100% !important;
}
.toolbox-ssh-editor .terminal-box .xterm-selection-layer {
  width: 100% !important;
}
.toolbox-ssh-editor .terminal-box .xterm-link-layer {
  width: 100% !important;
}
.toolbox-ssh-editor .terminal-box .xterm-cursor-layer {
  width: 100% !important;
}

.toolbox-ssh-editor .terminal-box .xterm .xterm-viewport::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.toolbox-ssh-editor
  .terminal-box
  .xterm
  .xterm-viewport:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.toolbox-ssh-editor
  .terminal-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-thumb {
  border-radius: 0px;
  background: #6b6b6b;
}
.toolbox-ssh-editor
  .terminal-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-track {
  border-radius: 0;
  background: #383838;
}
.toolbox-ssh-editor
  .terminal-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-corner {
  background: #ddd;
}
</style>
