<template>
  <div class="login-box">
    <b-container fluid class="bv-example-row">
      <b-row cols="1" cols-sm="1" cols-md="2">
        <b-col>
          <div class="login-left">
            <h1>Team IDE</h1>
            <hr />
            <p>团队在线工作 · 高效 · 安全 · 可靠</p>
          </div>
        </b-col>
        <b-col class="login-right">
          <b-form>
            <b-form-group
              label="账号"
            >
              <b-form-input
                v-model="form.account"
                placeholder="账号/邮箱"
                required
              ></b-form-input>
            </b-form-group>

            <b-form-group
              label="密码"
            >
              <b-form-input
                v-model="form.password"
                placeholder="登录密码"
                required
              ></b-form-input>
            </b-form-group>

            <b-form-group id="input-group-4" v-slot="{ ariaDescribedby }">
              <b-form-checkbox-group
                v-model="form.checked"
                id="checkboxes-4"
                :aria-describedby="ariaDescribedby"
              >
                <b-form-checkbox value="me">Check me out</b-form-checkbox>
                <b-form-checkbox value="that">Check that out</b-form-checkbox>
              </b-form-checkbox-group>
            </b-form-group>
            <b-button type="submit" variant="primary">Submit</b-button>
            <b-button type="reset" variant="danger">Reset</b-button>
          </b-form>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      form: {
        loginname: "",
        password: "",
      },
      ariaDescribedby: [],
    };
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
};
</script>

<style>
.login-box {
  width: 100%;
  height: 100%;
  position: fixed;
  left: 0px;
  top: 0px;
  z-index: 1000000;
  background: #fff;
}
</style>
