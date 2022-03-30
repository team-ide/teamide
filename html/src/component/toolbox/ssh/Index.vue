<template>
  <div class="toolbox-ssh-editor">
    <template v-if="ready">
      <div ref="terminal" style="width: 100%; height: 100%" />
    </template>
  </div>
</template>


<script>
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";

export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
      token: null,
    };
  },
  computed: {},
  watch: {
    token() {
      this.$nextTick(() => {
        if (this.socket != null) {
          this.socket.close();
        }
        if (this.term != null) {
          this.term.dispose();
        }
        this.$nextTick(() => {
          this.initTerminal();
        });
      });
    },
  },
  methods: {
    async init() {
      this.ready = true;
      await this.initToken();
    },
    async initToken() {
      let param = {};
      let res = await this.wrap.work("createToken", param);
      res.data = res.data || {};
      this.token = res.data.token;
    },
    initTerminal() {
      if (this.socket != null) {
        this.socket.close();
      }
      if (this.term != null) {
        this.term.dispose();
      }
      var term = new Terminal();
      var fitAddon = new FitAddon();
      let url = this.source.api;
      url = url.substring(url.indexOf(":"));
      url = "ws" + url + "ws/toolbox/ssh/connection?token=" + this.token;
      var socket = new WebSocket(url);
      var attachAddon = new AttachAddon(socket);

      term.loadAddon(attachAddon);
      term.loadAddon(fitAddon);
      term.open(this.$refs.terminal);
      fitAddon.fit();
      term.focus();
      console.log(term)
      socket.onopen = () => {
        if (this.tool.isNotEmpty(this.message)) {
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
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    if (this.socket != null) {
      this.socket.close();
    }
    if (this.term != null) {
      this.term.dispose();
    }
  },
};
</script>

<style>
.toolbox-ssh-editor {
  width: 100%;
  height: 100%;
}
</style>
