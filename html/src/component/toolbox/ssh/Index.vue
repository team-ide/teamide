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
      rows: 40,
      cols: 100,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      this.ready = true;
      await this.initToken();
      this.$nextTick(() => {
        this.initSocket();
      });
    },
    async initToken() {
      if (this.tool.isEmpty(this.token)) {
        let param = {};
        let res = await this.wrap.work("createToken", param);
        res.data = res.data || {};
        this.token = res.data.token;
      }
    },
    initSocket() {
      if (this.socket != null) {
        this.socket.close();
      }

      this.initTerminal();
      let url = this.source.api;
      url = url.substring(url.indexOf(":"));
      url = "ws" + url + "ws/toolbox/ssh/connection?token=" + this.token;
      url += "&cols=" + this.cols;
      url += "&rows=" + this.rows;
      this.socket = new WebSocket(url);

      // this.socket.onopen = () => {
      //   this.initTerminal();
      //   setTimeout(() => {
      //     // this.socket.send("\r");
      //   }, 1000);
      // };
      // // 当连接建立时向终端发送一个换行符，不这么做的话最初终端是没有内容的，输入换行符可让终端显示当前用户的工作路径
      this.socket.onmessage = (event) => {
        // 接收推送的消息
        let data = event.data.toString();
        if (data == '{"event":"ready"}') {
          this.socket.send("TeamIDE:event:start");
        } else if (data == '{"event":"shell created"}') {
          this.initAttachAddon();
        } else {
          console.log("data");
        }
      };
      // this.socket.onclose = () => {
      //   console.log("close socket");
      // };
      // this.socket.onerror = () => {
      //   console.log("socket error");
      // };
    },
    initAttachAddon() {
      var attachAddon = new AttachAddon(this.socket);

      this.term.loadAddon(attachAddon);
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
      this.term.open(this.$refs.terminal);

      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.fitAddon.fit();

      this.term.focus();
      console.log(this.fitAddon);
      this.cols = this.fitAddon._terminal.cols;
      this.rows = this.fitAddon._terminal.rows;
      window.onresize = function () {
        // 窗口尺寸变化时，终端尺寸自适应
        this.fitAddon.fit();
      };
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
