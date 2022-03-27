<template>
  <div class="toolbox-box" :style="boxStyleObject">
    <tm-layout height="100%">
      <tm-layout height="auto">
        <tm-layout height="100%">
          <tm-layout :width="style.left.width">
            <Left
              ref="left"
              v-if="source.toolbox.context != null"
              :source="source"
              :toolbox="source.toolbox"
              :context="source.toolbox.context"
              :style="leftStyleObject"
            ></Left>
          </tm-layout>
          <tm-layout-bar right></tm-layout-bar>
          <tm-layout width="auto">
            <Main
              ref="main"
              v-if="source.toolbox.context != null"
              :source="source"
              :toolbox="source.toolbox"
              :context="source.toolbox.context"
              :style="mainStyleObject"
            ></Main>
          </tm-layout>
        </tm-layout>
      </tm-layout>
    </tm-layout>
    <ToolboxForm :source="source" :toolbox="source.toolbox"></ToolboxForm>
  </div>
</template>

<script>
import Left from "./Left";
import Main from "./Main";
import ToolboxForm from "./ToolboxForm";

export default {
  components: { Left, Main, ToolboxForm },
  props: ["source"],
  data() {
    return {
      style: {
        backgroundColor: "#2d2d2d",
        color: "#adadad",
        header: {},
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
    leftStyleObject: function () {
      return {};
    },
    mainStyleObject: function () {
      return {};
    },
  },
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {
      this.source.toolbox.initContext = this.initContext;
      if (this.source.toolbox.context == null) {
        this.initContext();
      }
    },
    initContext() {
      this.loadContext();
    },
    async loadContext() {
      let param = {};
      let res = await this.server.toolbox.context(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let context = res.data.context || {};
      this.source.toolbox.context = context;
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
.toolbox-box {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
.toolbox-box .tm-layout > .tm-layout-bar {
  background-color: #4e4e4e;
}
.toolbox-box .tm-layout-bar > .tm-layout-bar-part {
  background-color: #4e4e4e;
}
.toolbox-box .toolbox-layout-header {
  border-bottom: 1px solid #4e4e4e;
}
</style>
