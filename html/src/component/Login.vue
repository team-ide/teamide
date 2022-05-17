<template>
  <div class="login-page bg-teal-9">
    <div class="login-box bg-teal-5 pd-20">
      <div class="login-left">
        <div class="ft-25 pdtb-10 pdlr-20">Team IDE</div>
        <p class="ft-15 ft-16 pdtb-5 pdlr-20">
          <span class="pdlr-5">团队协作</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">工作报告</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">高效</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">安全</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">可靠</span>
        </p>
        <hr />
        <div class="ft-25 pdtb-10 pdlr-20">Toolbox</div>
        <p class="ft-15 pdtb-5 pdlr-20">
          <span class="pdlr-5">Redis</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Mysql</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Zookeeper</span>
          <br />
          <span class="pdlr-5">Elasticsearch</span>
          <span class="pdlr-5 ft-20">·</span>
          <span class="pdlr-5">Kafka</span>
        </p>
      </div>
      <div class="login-right">
        <Form
          v-if="loginForm != null"
          :formBuild="loginForm"
          :formData="loginData"
          :saveShow="false"
          class="pd-10"
        >
          <!-- <b-form-group>
            <b-form-checkbox
              v-model="rememberPassword"
              :value="true"
              class="float-left mgr-20"
            >
              记住密码
            </b-form-checkbox>
            <b-form-checkbox
              v-model="autoLogin"
              :value="true"
              class="float-left mgr-20"
            >
              自动登录
            </b-form-checkbox>
          </b-form-group>
          <div class="pdtb-10">
            <div
              v-if="source.hasPower('login')"
              class="tm-btn bg-teal-8 ft-18 pdtb-5 tm-btn-block"
              :class="{ 'tm-disabled': loginBtnDisabled }"
              @click="doLogin"
            >
              登&nbsp;&nbsp;录
            </div>
          </div>
          <div
            v-if="source.hasPower('register')"
            class="pdtb-10 text-right ft-13"
          >
            没有账号？
            <div class="tm-link color-orange mgt--1" @click="tool.toRegister()">
              立即注册
            </div>
          </div> -->
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
      loginForm: null,
      loginData: null,
      rememberPassword: false,
      autoLogin: false,
      loginBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    rememberPassword() {
      if (!this.rememberPassword) {
        this.autoLogin = false;
      }
    },
    autoLogin() {
      if (this.autoLogin) {
        this.rememberPassword = true;
      }
    },
  },
  methods: {
    doLogin() {
      this.loginBtnDisabled = true;
      this.loginForm.validate(this.loginData).then((res) => {
        if (res.valid) {
          let param = {};
          Object.assign(param, this.loginData);
          let aesPassword = this.tool.aesEncrypt(param.password);
          param.password = aesPassword;
          this.server
            .login(param)
            .then((res) => {
              this.loginBtnDisabled = false;
              if (res.code == 0) {
                this.tool.setJWT(res.data);
                this.tool.success("登录成功！");
                this.tool.initSession();
                setTimeout(() => {
                  this.tool.hideLogin();
                }, 300);
              } else {
                this.tool.error(res.msg);
              }
            })
            .catch((e) => {
              this.loginBtnDisabled = false;
            });
        } else {
          this.loginBtnDisabled = false;
        }
      });
    },
    init() {
      this.loginForm = this.form.build(this.form.login);
      let loginData = this.loginForm.newDefaultData();
      this.loginData = loginData;
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
  z-index: 100;
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
  font-weight: 700;
}
.login-right {
  width: 400px;
  height: 100%;
  float: right;
}
</style>
