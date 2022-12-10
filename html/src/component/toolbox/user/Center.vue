<template>
  <div class="user-page">
    <div class="mglr-30 pdt-20">
      <el-form size="mini" @submit.native.prevent>
        <el-form-item label="用户名称">
          <el-input type="input" v-model="userForm.name"> </el-input>
        </el-form-item>
        <el-form-item label="登录账号">
          <el-input type="input" v-model="userForm.account"> </el-input>
        </el-form-item>
        <el-form-item label="用户邮箱">
          <el-input type="input" v-model="userForm.email"> </el-input>
        </el-form-item>
      </el-form>
      <div class="">
        <div class="tm-btn bg-teal-8 ft-18 pdtb-5" @click="doUpdate">修改</div>
      </div>
    </div>
    <div class="mglr-30 pdt-20">
      <el-form size="mini" @submit.native.prevent>
        <el-form-item label="原密码">
          <el-input type="password" v-model="passwordForm.oldPassword">
          </el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input type="password" v-model="passwordForm.password"> </el-input>
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input type="password" v-model="passwordForm.repassword">
          </el-input>
        </el-form-item>
      </el-form>
      <div class="">
        <div class="tm-btn bg-teal-8 ft-18 pdtb-5" @click="doUpdatePassword">
          修改密码
        </div>
      </div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      userForm: {
        name: null,
        avatar: null,
        account: null,
        email: null,
      },
      passwordForm: {
        oldPassword: null,
        password: null,
        repassword: null,
      },
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.loadUser();
    },
    refresh() {
      this.$nextTick(() => {});
    },
    async loadUser() {
      let res = await this.server.user.get();
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let data = res.data || {};
      let user = data.user || {};
      this.userForm.name = user.name;
      this.userForm.avatar = user.avatar;
      this.userForm.account = user.account;
      this.userForm.email = user.email;
    },
    async doUpdate() {
      if (this.tool.isEmpty(this.userForm.name)) {
        this.tool.error("请输入用户名称");
        return;
      }
      if (this.tool.isEmpty(this.userForm.account)) {
        this.tool.error("请输入登录账号");
        return;
      }
      if (this.tool.isEmpty(this.userForm.email)) {
        this.tool.error("请输入用户邮箱");
        return;
      }
      let res = await this.server.user.update(this.userForm);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        this.tool.success("修改成功");
        let data = res.data || {};
        if (data.user != null) {
          this.source.login.user.name = data.user.name;
          this.source.login.user.avatar = data.user.avatar;
          this.source.login.user.account = data.user.account;
          this.source.login.user.email = data.user.email;
        }
      }
    },
    async doUpdatePassword() {
      if (this.tool.isEmpty(this.passwordForm.oldPassword)) {
        this.tool.error("请输入原密码");
        return;
      }
      if (this.tool.isEmpty(this.passwordForm.password)) {
        this.tool.error("请输入登录密码");
        return;
      }
      if (this.passwordForm.password != this.passwordForm.repassword) {
        this.tool.error("两次密码输入不一致");
        return;
      }
      let param = {};
      Object.assign(param, this.passwordForm);
      param.oldPassword = this.tool.aesEncrypt(param.oldPassword);
      param.password = this.tool.aesEncrypt(param.password);
      let res = await this.server.user.updatePassword(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        this.tool.success("修改成功");
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.user-page {
  width: 100%;
  height: 100%;
}
</style>
