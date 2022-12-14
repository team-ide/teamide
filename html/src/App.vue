<template>
  <div id="app">
    <template v-if="source.ready">
      <InfoBox :source="source"></InfoBox>
      <SystemInfoBox :source="source"></SystemInfoBox>
      <AlertBox :source="source"></AlertBox>
      <router-view :source="source"></router-view>
      <Login v-show="source.login.show" :source="source"></Login>
      <Register v-show="source.register.show" :source="source"></Register>
      <UpdateCheck :source="source"></UpdateCheck>
    </template>
    <template v-else>
      <div v-if="source.status == 'connecting'"></div>
      <div v-if="source.status == 'connected'"></div>
      <div
        v-if="source.status == 'error'"
        style="
          position: fixed;
          width: 100%;
          height: 100%;
          top: 0px;
          background: #454a4e;
          color: #e61d1d;
        "
      >
        <h4 style="padding: 20px 20px; font-size: 25px">服务器连接异常！</h4>
        <hr />
        <p style="padding: 20px 20px; font-size: 20px">我们很遗憾的通知您：</p>
        <p style="text-indent: 60px; margin-top: 10px; font-size: 20px">
          服务器暂时无法正常访问，请您不要着急，我们正在紧急修复中！！！
        </p>
      </div>
    </template>
    <Contextmenu :contextmenu="contextmenu" ref="Contextmenu"></Contextmenu>
    <JSONDataDialog ref="JSONDataDialog" :source="source"></JSONDataDialog>
    <TextDialog ref="TextDialog" :source="source"></TextDialog>
    <div class="editor-for-copy" ref="editorForCopy"></div>
  </div>
</template>

<script>
import source from "@/source";

import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";
import UpdateCheck from "@/views/UpdateCheck.vue";

export default {
  components: { Login, Register, UpdateCheck },
  props: [],
  data() {
    return { source, contextmenu: { menus: [] } };
  },
  computed: {},
  watch: {
    "source.login.user"() {
      if (this.source.login.user == null) {
        if (this.source.login.userId != null) {
          this.source.login.userId = null;
        }
      } else {
        if (this.source.login.userId != this.source.login.user.userId) {
          this.source.login.userId = this.source.login.user.userId;
        }
      }
    },
    "source.login.userId"() {
      var info = "change userId:";
      if (this.source.login.userId != null) {
        info += this.source.login.userId;
      }
      this.server.websocketSendText(info);
    },
    "source.ready"() {
      if (this.source.ready) {
        this.server.openWebsocket();
      }
    },
  },
  methods: {
    init() {
      this.tool.showJSONData = this.$refs.JSONDataDialog.show;
      this.tool.showText = this.$refs.TextDialog.show;

      window.onMonacoLoad(() => {
        let monaco = window.monaco;
        this.monacoInstance = monaco.editor.create(this.$refs.editorForCopy, {
          theme: "vs-dark", //官方自带三种主题vs, hc-black, or vs-dark
          minimap: { enabled: false }, // 缩略导航
          value: "", //编辑器初始显示文字
          language: "html",
          selectOnLineNumbers: false, //显示行号
          roundedSelection: false,
          cursorStyle: "line", //光标样式
          automaticLayout: false, //自动布局
          glyphMargin: false, //字形边缘
          useTabStops: false,
          fontSize: 13, //字体大小
          autoIndent: false, //自动布局
          // quickSuggestionsDelay: 500, //代码提示延时
          contextmenu: false,
        });
        this.tool.copyByEditor = (text) => {
          this.monacoInstance.setValue(text || "");
          this.monacoInstance.setSelection({
            startLineNumber: 0,
            startColumn: 0,
            endLineNumber: 999999,
            endColumn: 999999,
          });
          let action = this.monacoInstance.getAction(
            "editor.action.clipboardCopyWithSyntaxHighlightingAction"
          );
          if (action) {
            action.run();
            return { success: true };
          } else {
            return { success: false };
          }
        };
      });
    },
    showContextmenu(menus) {
      let e = window.event;
      this.tool.stopEvent(e || window.event);
      this.contextmenu.menus = menus;
      this.$refs.Contextmenu.show(e);
    },
    onKeyDown(event) {
      event = event || window.event;
      if (this.tool.keyIsCtrlS(event)) {
        this.tool.stopEvent(event);
      }
    },
    onKeyUp(event) {
      event = event || window.event;
      if (this.tool.keyIsCtrlS(event)) {
        this.tool.stopEvent(event);
      }
    },
    bindEvent() {
      if (this.bindEvented) {
        return;
      }
      this.bindEvented = true;
      window.addEventListener("keydown", (e) => {
        this.onKeyDown(e);
      });
      window.addEventListener("keyup", (e) => {
        this.onKeyUp(e);
      });
      window.document.body.addEventListener("contextmenu", (e) => {
        let tags = ["input", "textarea", "canvas"];
        if (tags.indexOf(("" + e.target.tagName).toLowerCase()) >= 0) {
          return;
        }
        this.tool.stopEvent(e || window.event);
        // this.showContextmenu([{text:"测试菜单"}])
      });
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.tool.showContextmenu = this.showContextmenu;
    this.bindEvent();
    this.init();
  },
  beforeDestroy() {},
};
</script>

<style>
html,
body {
  height: 100%;
  width: 100%;
  user-select: none;
  padding: 0px;
  margin: 0px;
}
#app {
  user-select: none;
  height: 100%;
  width: 100%;
  padding: 0px;
  margin: 0px;
}
*,
:after,
:before {
  box-sizing: border-box;
}
.el-message {
  user-select: text;
}
.el-message-box {
  user-select: text;
}

/* 滚动条样式*/
.app-scroll-bar {
  overflow: scroll;
}

.app-scroll-bar:hover::-webkit-scrollbar-thumb {
  box-shadow: inset 0 0 10px #333333;
  background: #333333;
}
.app-scroll-bar:hover::-webkit-scrollbar-track {
  box-shadow: inset 0 0 10px #262626;
  background: #262626;
}
.app-scroll-bar:hover::-webkit-scrollbar-corner {
  background: #262626;
}

.app-scroll-bar::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.app-scroll-bar:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.app-scroll-bar::-webkit-scrollbar-thumb {
  border-radius: 0px;
}
.app-scroll-bar::-webkit-scrollbar-track {
  border-radius: 0;
}
.app-scroll-bar::-webkit-scrollbar-corner {
  background: transparent;
}

.app-scroll-bar-textarea textarea::-webkit-scrollbar-thumb {
  box-shadow: inset 0 0 10px #333333;
  background: #333333;
}
.app-scroll-bar-textarea textarea:hover::-webkit-scrollbar-track {
  box-shadow: inset 0 0 10px #262626;
  background: #262626;
}
.app-scroll-bar-textarea textarea:hover::-webkit-scrollbar-corner {
  background: #262626;
}

.app-scroll-bar-textarea textarea::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.app-scroll-bar-textarea textarea:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.app-scroll-bar-textarea textarea::-webkit-scrollbar-thumb {
  border-radius: 0px;
}
.app-scroll-bar-textarea textarea::-webkit-scrollbar-track {
  border-radius: 0;
}
.app-scroll-bar-textarea textarea::-webkit-scrollbar-corner {
  background: transparent;
}

.tm-link {
  text-decoration: none !important; /* 去除默认的下划线 */
}
.mdi {
  vertical-align: middle;
}
.editor-for-copy {
  transform: scale(0);
  width: 100px;
  height: 100px;
}
</style>
