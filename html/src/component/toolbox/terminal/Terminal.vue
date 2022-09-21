<template>
  <div class="toolbox-terminal-box">
    <tm-layout height="100%">
      <tm-layout height="auto">
        <div
          class="terminal-xterm-box"
          ref="terminalXtermBox"
          style="
            padding-top: 0px;
            padding-right: 0px;
            padding-bottom: 0px;
            padding-left: 0px;
          "
        />
        <div class="terminal-xterm-box-back" ref="terminalXtermBoxBack" />
      </tm-layout>
      <tm-layout height="30px">
        <div class="pdt-2 pdlr-10">
          <div class="ft-12 tm-link color-grey mglr-5" @click="openFtpWindow()">
            FTP
          </div>
        </div>
      </tm-layout>
    </tm-layout>
    <template v-if="isOpenFTP">
      <div
        class="toolbox-terminal-file-manager-box"
        :class="{ 'toolbox-terminal-file-manager-box-show': isShowFTP }"
        :style="{
          width: `${ftpWidth}px`,
          height: `${ftpHeight}px`,
        }"
      >
        <div
          ref="ftpBoxTopLine"
          class="toolbox-terminal-file-manager-box-top-line"
        ></div>
        <div
          ref="ftpBoxLeftLine"
          class="toolbox-terminal-file-manager-box-left-line"
        ></div>
        <div class="toolbox-terminal-file-manager-box-header">
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
        <div class="toolbox-terminal-file-manager-box-body">
          <tm-layout height="100%">
            <tm-layout height="auto">
              <FileManager
                :source="source"
                :toolboxWorker="toolboxWorker"
                :place="place"
                :placeId="placeId"
                :openDir="extend == null ? '' : extend.openDir"
                :onChangeOpenDir="onChangeOpenDir"
              ></FileManager>
            </tm-layout>
            <tm-layout-bar top></tm-layout-bar>
            <tm-layout height="200px">
              <Progress
                :source="source"
                :toolboxWorker="toolboxWorker"
              ></Progress>
            </tm-layout>
          </tm-layout>
        </div>
      </div>
    </template>
    <Download :source="source" :toolboxWorker="toolboxWorker"></Download>
    <Upload :source="source" :toolboxWorker="toolboxWorker"></Upload>
    <ConfirmPaste
      :source="source"
      :toolboxWorker="toolboxWorker"
    ></ConfirmPaste>
  </div>
</template>


<script>
import _worker from "./worker.js";
import "xterm/css/xterm.css";
import Zmodem from "zmodem.js";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
// import { AttachAddon } from "xterm-addon-attach";

// https://juejin.cn/post/6918911964009725959

import FileManager from "../file-manager/FileManager.vue";
import Progress from "../file-manager/Progress.vue";
import Download from "./Download.vue";
import Upload from "./Upload.vue";
import ConfirmPaste from "./ConfirmPaste.vue";

export default {
  components: { FileManager, Progress, Download, Upload, ConfirmPaste },
  props: ["source", "toolboxWorker", "place", "placeId", "extend"],
  data() {
    let worker = _worker.newWorker({
      workerId: this.toolboxWorker.workerId,
      place: this.place,
      placeId: this.placeId,
      onSocketOpen: this.onSocketOpen,
      onSocketClose: this.onSocketClose,
      onSocketError: this.onSocketError,
      onSocketData: this.onSocketData,
    });
    return {
      worker: worker,
      ftpWidth: 900,
      ftpHeight: 600,

      isOpenFTP: false,
      isShowFTP: true,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.$nextTick(() => {
        this.initTerm();
        this.worker.init();
      });
    },
    onFocus() {
      this.term && this.term.focus();
    },
    refresh() {},
    onChangeOpenDir(openDir) {
      let data = this.extend || {};
      if (this.tool.isEmpty(openDir)) {
        openDir = "";
      }
      if (this.tool.isEmpty(openDir) && this.tool.isEmpty(data.openDir)) {
        return;
      }
      if (openDir == data.openDir) {
        return;
      }
      data.openDir = openDir;
      var keyValueMap = {};
      keyValueMap["openDir"] = openDir;
      this.toolboxWorker.updateExtend(keyValueMap);
    },
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
    onSocketData(data) {
      this.zsentry.consume(data);
      // if (typeof data === "string") {
      //   this.term.write(data);
      // } else {
      //   this.term.write(new Uint8Array(data));
      // }
    },
    initTerm() {
      if (this.term != null) {
        this.term.dispose();
      }

      this.term = new Terminal({
        useStyle: true,
        cursorBlink: true, //光标闪烁
        cursorStyle: "underline", // 光标样式 'block' | 'underline' | 'bar'
        rendererType: "canvas", //渲染类型
        width: 500,
        height: 400,
        windowsMode: true,
        scrollback: 100000000, //终端中的回滚量
        // rows: this.rows, //行数
        // cols: this.cols, // 不指定行数，自动回车后光标从下一行开始
        convertEol: true, //启用时，光标将设置为下一行的开头
        // disableStdin: false, //是否应禁用输入
        // // cursorStyle: "underline", //光标样式
        // theme: {
        //   foreground: "#ECECEC", //字体
        //   background: "#000000", //背景色
        //   cursor: "help", //设置光标
        //   lineHeight: 20,
        // },
      });
      this.term.open(this.$refs.terminalXtermBox, true);

      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.fitAddon.fit();

      this.term.focus();
      this.worker.cols = this.term.cols;
      this.worker.rows = this.term.rows;

      this.term.onData((data) => {
        if (this.checkIsExit(data)) {
          this.worker.sendDataToWS(data);
        }
      });
      this.term.onBinary((data) => {
        if (this.checkIsExit(data)) {
          this.worker.sendDataToWS(data);
        }
      });

      this.zsentry = new Zmodem.Sentry({
        //发送终端
        to_terminal: (octets) => {
          this.term.write(new Uint8Array(octets));
        },
        // 属于 Zmodem 相关流
        on_detect: (detection) => {
          let zsession = detection.confirm();
          if (zsession.type === "receive") {
            this.toolboxWorker.showDownload(zsession, this.term, () => {
              this.onFocus();
            });
          } else {
            this.toolboxWorker.showUpload(zsession, this.term, () => {
              this.onFocus();
            });
          }
        },
        // 撤回
        on_retract: () => {},
        sender: (octets) => {
          this.worker.socket.send(new Uint8Array(octets));
        },
      });

      this.terminal_back_width = this.tool
        .jQuery(this.$refs.terminalXtermBoxBack)
        .width();
      this.terminal_back_height = this.tool
        .jQuery(this.$refs.terminalXtermBoxBack)
        .height();

      this.$refs.terminalXtermBox.addEventListener(
        "keydown",
        this.onKeydown,
        true
      );
      this.$refs.terminalXtermBox.addEventListener(
        "mouseup",
        this.onMouseup,
        true
      );
      this.$refs.terminalXtermBox.addEventListener(
        "mousedown",
        this.onMousedown,
        true
      );
      this.$refs.terminalXtermBox.addEventListener(
        "contextmenu",
        this.onContextmenu,
        true
      );

      this.changeSizeTimer();
    },
    checkIsExit(data) {
      if (this.worker.socket == null) {
        return false;
      }
      return true;
    },
    onSocketOpen() {
      // const attachAddon = new AttachAddon(this.worker.socket);
      // this.term.loadAddon(attachAddon);
    },
    onSocketClose() {
      if (this.isDestroyed) {
        return;
      }
      this.term.write("\r\n终端会话已关闭，输入回车重新连接！\r\n");
      // this.worker.refresh();
    },
    onSocketError() {},
    async onKeydown(e) {
      // let key = arg.key;
      // console.log(key);

      if (this.worker.socket == null) {
        this.tool.stopEvent(e);

        if (this.worker.building) {
          return;
        }
        if (this.tool.keyIsEnter(e)) {
          this.term.write("\r\n终端会话连接中，请稍后！\r\n");
          this.worker.refresh();
          return;
        }
        this.term.write("\r\n终端会话已关闭，输入回车重新连接！\r\n");

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
          // this.tool.success("复制成功");
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
      if (this.tool.isEmpty(text)) {
        return;
      }
      if (text.indexOf("\n") >= 0) {
        text = text.replace(/(\r\n|\n|\r|↵)/g, `\n`);
        this.toolboxWorker.showConfirmPaste(
          text,
          () => {
            this.worker.sendDataToWS(text);
            this.onFocus();
            // this.tool.success("粘贴成功");
          },
          () => {
            this.onFocus();
          }
        );
      } else {
        this.worker.sendDataToWS(text);
        // this.tool.success("粘贴成功");
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
    changeSizeTimer() {
      if (this.isDestroyed) {
        return;
      }
      if (
        this.tool.jQuery(this.$refs.terminalXtermBoxBack).width() !=
          this.terminal_back_width ||
        this.tool.jQuery(this.$refs.terminalXtermBoxBack).height() !=
          this.terminal_back_height
      ) {
        this.terminal_back_width = this.tool
          .jQuery(this.$refs.terminalXtermBoxBack)
          .width();
        this.terminal_back_height = this.tool
          .jQuery(this.$refs.terminalXtermBoxBack)
          .height();
        this.tool.jQuery(this.term.element).css({
          width: parseInt(this.terminal_back_width),
          height: parseInt(this.terminal_back_height),
        });
        this.fitAddon.fit();
        // console.log(this.term.element);

        if (
          this.term.cols != this.worker.cols ||
          this.term.rows != this.worker.rows
        ) {
          this.worker.cols = this.term.cols;
          this.worker.rows = this.term.rows;
          this.worker.changeSize();
        }
      }

      // window.setTimeout(() => {
      // 窗口尺寸变化时，终端尺寸自适应

      window.setTimeout(() => {
        this.changeSizeTimer();
      }, 200);
      // }, 100);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
    this.worker.close();
    if (this.term != null) {
      this.term.dispose();
    }
  },
};
</script>

<style>
.toolbox-terminal-box {
  width: 100%;
  height: 100%;
}

.toolbox-terminal-box .terminal-xterm-box {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: black;
  overflow: hidden;
}
.toolbox-terminal-box .terminal-xterm-box-back {
  width: calc(100% - 20px) !important;
  height: 100%;
  position: absolute;
  left: 0px;
  top: 0px;
  z-index: -1;
}

.toolbox-terminal-box .terminal-xterm-box .terminal {
  width: 100% !important;
  height: 100% !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-viewport {
  width: 100% !important;
  background-color: transparent !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-screen {
  width: calc(100% - 20px) !important;
  height: 100% !important;
  margin: 0px 5px;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-text-layer {
  width: 100% !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-selection-layer {
  width: 100% !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-link-layer {
  width: 100% !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-cursor-layer {
  width: 100% !important;
}

.toolbox-terminal-box
  .terminal-xterm-box
  .xterm
  .xterm-viewport::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.toolbox-terminal-box
  .terminal-xterm-box
  .xterm
  .xterm-viewport:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.toolbox-terminal-box
  .terminal-xterm-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-thumb {
  border-radius: 0px;
  background: #6b6b6b;
}
.toolbox-terminal-box
  .terminal-xterm-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-track {
  border-radius: 0;
  background: #383838;
}
.toolbox-terminal-box
  .terminal-xterm-box
  .xterm
  .xterm-viewport::-webkit-scrollbar-corner {
  background: #ddd;
}

.toolbox-terminal-file-manager-box {
  position: absolute;
  right: 20px;
  bottom: 40px;
  background: #172029;
  transition: all 0s;
  transform: scale(0);
  z-index: -1;
}

.toolbox-terminal-file-manager-box.toolbox-terminal-file-manager-box-show {
  transform: scale(1);
  z-index: 10;
}
.toolbox-terminal-file-manager-box-top-line {
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
.toolbox-terminal-file-manager-box-left-line {
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

.toolbox-terminal-file-manager-box-header {
  height: 29px;
  border-bottom: 1px solid #2f2f2f;
}
.toolbox-terminal-file-manager-box-body {
  height: calc(100% - 30px);
}
</style>
