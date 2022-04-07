<template>
  <div class="toolbox-ssh-editor">
    <div ref="terminal" style="width: 100%; height: 100%" />
  </div>
</template>


<script>
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { AttachAddon } from "xterm-addon-attach";

export default {
  components: {},
  props: [
    "source",
    "data",
    "toolboxType",
    "toolbox",
    "option",
    "wrap",
    "token",
    "socket",
  ],
  data() {
    return {
      rows: 40,
      cols: 100,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.initTerminal();
    },
    onEvent(event) {
      if (event == "shell ready") {
        this.toStart();
      } else if (event == "shell created") {
        this.initAttachAddon();
      }
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
    initAttachAddon() {
      this.attachAddon = new AttachAddon(this.socket);
      this.term.loadAddon(this.attachAddon);
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
      this.cols = this.term.cols;
      this.rows = this.term.rows;
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
