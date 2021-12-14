<template>
  <div class="login-page">
    <div class="login-box bg-cyan-6 pd-20">
      <div class="login-left ">
        <h1 class="pdtb-10 pdlr-20">Team IDE</h1>
        <hr />
        <p class="ft-16 pdtb-10 pdlr-20">
          <span class="pdlr-5">团队协作</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">工作报告</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">高效</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">安全</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">可靠</span>
        </p>
        <h3 class="pdtb-10 pdlr-20">Toolbox</h3>
        <p class="ft-16 pdtb-10 pdlr-20">
          <span class="pdlr-5">Redis</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">Mysql</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">Zookeeper</span>
          <br/>
          <span class="pdlr-5">Elasticsearch</span>
          <span class="pdlr-5">·</span>
          <span class="pdlr-5">Kafka</span>
        </p>
      </div>
      <div class="login-right">
        <Form v-if="loginForm != null" :form="loginForm" class="pd-10"></Form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      loginForm: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {
      this.loginForm = this.form.build(this.form.login);
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
.login-page {
  width: 100%;
  height: 100%;
  position: fixed;
  left: 0px;
  top: 0px;
  z-index: 1000000;
  background: #fff;
}
.login-box {
  position: absolute;
  width: 860px;
  height: 400px;
  left: 50%;
  top: 50%;
  margin-left: -420px;
  margin-top: -260px;
}
.login-left {
  width: 400px;
  height: 100%;
  float: left;
}
.login-right {
  width: 400px;
  height: 100%;
  float: right;
}
</style>
