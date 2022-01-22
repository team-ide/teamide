<template>
  <div class="application-box" :style="boxStyleObject">
    <tm-layout height="100%">
      <tm-layout
        class="application-layout-header"
        :height="style.header.height"
      >
        <Header
          ref="header"
          :source="source"
          :application="source.application"
          :app="source.application.app"
          :style="headerStyleObject"
        ></Header>
      </tm-layout>
      <tm-layout height="auto">
        <tm-layout height="100%">
          <tm-layout :width="style.left.width">
            <Left
              ref="left"
              v-if="source.application.context != null"
              :source="source"
              :application="source.application"
              :app="source.application.app"
              :context="source.application.context"
              :style="leftStyleObject"
            ></Left>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="auto">
            <Main
              ref="main"
              v-if="source.application.context != null"
              :source="source"
              :application="source.application"
              :app="source.application.app"
              :context="source.application.context"
              :style="mainStyleObject"
            ></Main>
          </tm-layout>
        </tm-layout>
      </tm-layout>
    </tm-layout>
    <AppForm :source="source" :application="source.application"></AppForm>
    <ModelForm :source="source" :application="source.application"></ModelForm>
  </div>
</template>

<script>
import Header from "./Header";
import Left from "./Left";
import Main from "./Main";
import AppForm from "./AppForm";
import ModelForm from "./ModelForm";

export default {
  components: { Header, Left, Main, AppForm, ModelForm },
  props: ["source"],
  data() {
    return {
      style: {
        backgroundColor: "#2d2d2d",
        color: "#adadad",
        header: {
          height: "50px",
        },
        left: {
          width: "260px",
        },
        main: {},
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {
    boxStyleObject: function () {
      return {
        backgroundColor: this.style.backgroundColor,
        color: this.style.color,
      };
    },
    headerStyleObject: function () {
      return {};
    },
    leftStyleObject: function () {
      return {};
    },
    mainStyleObject: function () {
      return {};
    },
  },
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "source.application.app"(newApp, oldApp) {
      if (newApp == null || oldApp == null) {
        this.initContext();
      } else if (newApp.name != oldApp.name) {
        this.initContext();
      }
    },
  },
  methods: {
    init() {},
    initContext() {
      if (this.source.application.app == null) {
        this.source.application.context = null;
      } else {
        this.loadContext();
      }
    },
    async loadContext() {
      let param = {
        appName: this.source.application.app.name,
      };
      let res = await this.server.application.context.load(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let context = res.data;
      this.source.application.context = context;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
  },
};
</script>

<style>
.application-box {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
.application-box .tm-layout {
  overflow: unset;
}
.application-box .tm-layout > .tm-layout-bar {
  background-color: #4e4e4e;
}
.application-box .tm-layout-bar > .tm-layout-bar-part {
  background-color: #4e4e4e;
}
.application-box .application-layout-header {
  border-bottom: 1px solid #4e4e4e;
}
</style>
