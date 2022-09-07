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
      <tm-layout height="30px"> </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
import _worker from "./worker.js";
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";

export default {
  components: {},
  props: ["source", "toolboxWorker", "place", "placeId"],
  data() {
    let worker = _worker.newWorker({
      workerId: this.toolboxWorker.workerId,
      place: this.place,
      placeId: this.placeId,
      onSocketOpen: this.onSocketOpen,
      onSocketClose: this.onSocketClose,
      onSocketError: this.onSocketError,
    });
    return {
      rows: 40,
      cols: 100,
      worker: worker,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.worker.init();
      this.initTerm();
    },
    refresh() {},
    onFocus() {
      this.term && this.term.focus();
    },
    initTerm() {
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
      this.term.open(this.$refs.terminalXtermBox, true);

      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.fitAddon.fit();

      this.term.focus();
    },
    onSocketOpen() {},
    onSocketClose() {
      if (this.isDestroyed) {
        return;
      }
      this.worker.refresh();
    },
    onSocketError() {},
  },
  created() {},
  mounted() {
    this.init();
  },
  destroyed() {
    this.isDestroyed = true;
    this.worker.close();
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
}
.toolbox-terminal-box .terminal-xterm-box .xterm-viewport {
  width: 100% !important;
  background-color: transparent !important;
}
.toolbox-terminal-box .terminal-xterm-box .xterm-screen {
  width: calc(100% - 20px) !important;
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
