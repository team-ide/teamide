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
      <div v-if="source.status == 'connecting'">
        <!-- <b-alert show>
          <h1>服务器连接中，请稍后！</h1>
        </b-alert> -->
      </div>
      <div v-if="source.status == 'connected'">
        <!-- <b-alert show variant="success">
          <h1>服务器连接成功，请使用！</h1>
        </b-alert> -->
      </div>
      <div
        v-if="source.status == 'error'"
        class="alert alert-danger"
        style="position: fixed; width: 100%; height: 100%; top: 0px"
      >
        <b-alert show variant="danger">
          <h4 class="alert-heading">服务器连接异常！</h4>
          <hr />
          <p class="mb-0">我们很遗憾的通知您：</p>
          <p class="mb-0" style="text-indent: 33px; margin-top: 10px">
            服务器暂时无法正常访问，请您不要着急，我们正在紧急修复中！！！
          </p>
        </b-alert>
      </div>
    </template>
  </div>
</template>

<script>
import source from "@/source";

import Frame from "@/views/frame/Index.vue";

export default {
  components: { Frame },
  props: [],
  data() {
    return { source };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {},
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
  destroyed() {},
};
</script>

<style>
html,
body {
  height: 100%;
  width: 100%;
}
#app {
  height: 100%;
  width: 100%;
}
*,
:after,
:before {
  box-sizing: border-box;
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
</style>
