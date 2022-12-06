<template>
  <div class="register-page bg-teal-9">
    <div
      class="tm-link"
      style="
        position: absolute;
        right: 10px;
        top: 10px;
        font-size: 30px;
        color: #ffffff;
        line-height: 30px;
      "
      @click="tool.hideRegister()"
    >
      <Icon class="mdi-close"></Icon>
    </div>
    <div class="register-box bg-teal-5 pd-20">
      <div class="register-left">
        <div class="ft-25 pdtb-10 pdlr-20">Team IDE</div>
        <p class="ft-15 pdtb-5 pdlr-20">
          <span class="pdlr-5">Redis</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Zookeeper</span>
          <br />
          <span class="pdlr-5">Elasticsearch</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Kafka</span>
          <br />
          <span class="pdlr-5">Mysql</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Oracle</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">达梦</span>
          <br />
          <span class="pdlr-5">神通</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">金仓</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Sqlite</span>
        </p>
      </div>
      <div class="register-right" v-if="registerForm != null">
        <Form
          :formBuild="registerForm"
          :formData="registerData"
          :saveShow="false"
          class="pd-10"
        >
          <div class="tm-row pdtb-10">
            <div
              v-if="source.hasPower('register')"
              class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
              :class="{ 'tm-disabled': registerBtnDisabled }"
              @click="doRegister"
            >
              注&nbsp;&nbsp;册
            </div>
          </div>
          <div v-if="source.hasPower('login')" class="pdtb-10 text-right ft-13">
            已有账号？
            <div class="tm-link color-orange mgt--1" @click="tool.toLogin()">
              立即登录
            </div>
          </div>
        </Form>
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
      registerForm: null,
      registerData: null,
      registerBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    doRegister() {
      this.registerBtnDisabled = true;
      this.registerForm.validate(this.registerData).then((res) => {
        if (res.valid) {
          let param = {};
          Object.assign(param, this.registerData);
          let aesPassword = this.tool.aesEncrypt(param.password);
          param.password = aesPassword;
          this.server
            .register(param)
            .then((res) => {
              this.registerBtnDisabled = false;
              if (res.code == 0) {
                this.tool.success("注册成功！");
                setTimeout(() => {
                  this.tool.toLogin();
                }, 300);
              } else {
                this.tool.error(res.msg);
              }
            })
            .catch((e) => {
              this.registerBtnDisabled = false;
            });
        } else {
          this.registerBtnDisabled = false;
        }
      });
    },
    init() {
      this.registerForm = this.form.build(this.form.register);
      let registerData = this.registerForm.newDefaultData();
      this.registerData = registerData;
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
.register-page {
  width: 100%;
  height: 100%;
  position: fixed;
  left: 0px;
  top: 0px;
  z-index: 100;
  background: #fff;
}
.register-box {
  position: absolute;
  width: 860px;
  height: 500px;
  left: 50%;
  top: 50%;
  margin-left: -420px;
  margin-top: -300px;
}
.register-left {
  width: 400px;
  height: 100%;
  float: left;
  font-weight: 700;
}
.register-right {
  width: 400px;
  height: 100%;
  float: right;
}
.register-page .el-form-item__label {
  color: #ffffff;
}
</style>
