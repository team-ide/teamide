<template>
  <div ref="terminal" style="width: 100%; height: 100%" />
</template>
<script>
import server from "@/server";
import tool from "@/tool";
import source from "@/source";

import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";

export default {
  name: "Xterm",
  props: ["action", "message"],

  watch: {
    action() {
      this.$nextTick(() => {
        if (this.socket != null) {
          this.socket.close();
        }
        if (this.term != null) {
          this.term.dispose();
        }
        setTimeout(() => {
          this.init();
        }, 1000);
      });
    },
  },
  methods: {
    init() {
      if (this.socket != null) {
        this.socket.close();
      }
      if (this.term != null) {
        this.term.dispose();
      }
      var term = new Terminal();
      var fitAddon = new FitAddon();
      let url = source.ROOT_URL;
      url = url.substring(url.indexOf(":"));
      url = "ws" + url + this.action;
      var socket = new WebSocket(url);
      var attachAddon = new AttachAddon(socket);

      term.loadAddon(attachAddon);
      term.loadAddon(fitAddon);
      term.open(this.$refs.terminal);
      fitAddon.fit();
      term.focus();
      socket.onopen = () => {
        if (tool.isNotEmpty(this.message)) {
          socket.send(this.message);
        }
      }; // 当连接建立时向终端发送一个换行符，不这么做的话最初终端是没有内容的，输入换行符可让终端显示当前用户的工作路径

      window.onresize = function () {
        // 窗口尺寸变化时，终端尺寸自适应
        fitAddon.fit();
      };

      this.term = term;
      this.socket = socket;
    },
  },
  mounted() {
    this.$nextTick(() => {
      setTimeout(() => {
        this.init();
      }, 1000);
    });
  },
  beforeDestroy() {
    this.socket.close();
    this.term.dispose();
  },
};
</script>
<style >
#terminal {
  height: 100%;
}
</style>