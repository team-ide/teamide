<template>
  <div id="app">
    <template v-if="source.ready">
      <InfoBox :source="source"></InfoBox>
      <SystemInfoBox :source="source"></SystemInfoBox>
      <AlertBox :source="source"></AlertBox>
      <Frame
        :source="source"
        v-if="!source.frame.remove"
        v-show="source.frame.show"
      >
      </Frame>
      <Login
        :source="source"
        v-if="!source.login.remove"
        v-show="source.login.show"
      >
      </Login>
      <Register
        :source="source"
        v-if="!source.register.remove"
        v-show="source.register.show"
      >
      </Register>
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
  </div>
</template>

<script>
import source from "@/source";

import Frame from "@/views/frame/Index.vue";

export default {
  components: { Frame },
  props: [],
  data() {
    return { source, contextmenu: { menus: [] } };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    showContextmenu(menus) {
      let e = window.event;
      e.preventDefault();
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
        let tags = ["input", "textarea"];
        if (tags.indexOf(("" + e.target.tagName).toLowerCase()) >= 0) {
          return;
        }
        this.tool.stopEvent(e);
      });
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.tool.showContextmenu = this.showContextmenu;
    this.bindEvent();
  },
  destroyed() {},
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
.el-message,
.el-dialog__title {
  user-select: text;
}

/* 滚动条样式*/
.scrollbar {
  overflow: scroll;
}

.scrollbar:hover::-webkit-scrollbar-thumb {
  box-shadow: inset 0 0 10px #333333;
  background: #333333;
}
.scrollbar:hover::-webkit-scrollbar-track {
  box-shadow: inset 0 0 10px #262626;
  background: #262626;
}
.scrollbar:hover::-webkit-scrollbar-corner {
  background: #262626;
}

.scrollbar::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.scrollbar:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.scrollbar::-webkit-scrollbar-thumb {
  border-radius: 0px;
}
.scrollbar::-webkit-scrollbar-track {
  border-radius: 0;
}
.scrollbar::-webkit-scrollbar-corner {
  background: transparent;
}

.scrollbar-textarea textarea::-webkit-scrollbar-thumb {
  box-shadow: inset 0 0 10px #333333;
  background: #333333;
}
.scrollbar-textarea textarea:hover::-webkit-scrollbar-track {
  box-shadow: inset 0 0 10px #262626;
  background: #262626;
}
.scrollbar-textarea textarea:hover::-webkit-scrollbar-corner {
  background: #262626;
}

.scrollbar-textarea textarea::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.scrollbar-textarea textarea:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.scrollbar-textarea textarea::-webkit-scrollbar-thumb {
  border-radius: 0px;
}
.scrollbar-textarea textarea::-webkit-scrollbar-track {
  border-radius: 0;
}
.scrollbar-textarea textarea::-webkit-scrollbar-corner {
  background: transparent;
}

.xterm .xterm-viewport::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.xterm .xterm-viewport:hover::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
.xterm .xterm-viewport::-webkit-scrollbar-thumb {
  border-radius: 0px;
  box-shadow: inset 0 0 10px #333333;
  background: #333333;
}
.xterm .xterm-viewport::-webkit-scrollbar-track {
  border-radius: 0;
  box-shadow: inset 0 0 10px #ddd;
  background: #262626;
}
.xterm .xterm-viewport::-webkit-scrollbar-corner {
  background: #262626;
}
.tm-link {
  text-decoration: none !important; /* 去除默认的下划线 */
}
.mdi {
  vertical-align: middle;
}
</style>
