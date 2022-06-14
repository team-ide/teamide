<template>
  <div class="frame-header-box">
    <div class="frame-header-left">
      <div class="ft-15 ft-700 color-white">
        {{ source.header.title }}
      </div>
    </div>
    <div class="frame-header-center">
      <div>
        <template v-for="(one, index) in source.frame.headerNavs">
          <div
            class="tm-link mgr-10"
            :key="index"
            @click="$router.push(`${one.link}`)"
            :class="{
              'tm-active':
                $route.path == `${one.link}` ||
                (one.match && one.match($route.path)),
            }"
          >
            {{ one.name }}
          </div>
        </template>
      </div>
    </div>
    <div class="frame-header-right">
      <template v-if="source.login.user == null">
        <div
          class="tm-link mgr-10"
          v-if="source.hasPower('login')"
          @click="tool.toLogin()"
        >
          登录
        </div>
      </template>
      <template v-else>
        <div class="tm-link mgr-10">
          {{ source.login.user.name }}
        </div>
      </template>
      <div class="tm-link mgr-10 ft-12" @click="tool.showUpdateCheck()">
        <template v-if="source.hasNewVersion"> 有新版本 </template>
        <template v-else> 检测新版本 </template>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {};
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
.frame-header-box {
  display: flex;
  background: #4e6064;
  font-size: 14px;
}

.frame-box-only-body .frame-header-box {
  display: none;
}
.frame-header-left {
  display: flex;
  padding: 0px 10px;
  line-height: 28px;
}
.frame-header-center {
  display: flex;
  flex: 1;
  padding: 0px 10px;
  line-height: 28px;
}
.frame-header-right {
  display: flex;
  padding: 0px 10px;
  line-height: 28px;
}
.frame-header-box .tm-link {
  color: #bababa;
}
.frame-header-box img {
  margin-right: 0px;
}
.frame-header-box .user-nav-box {
  margin: 5px 5px;
  min-width: 400px;
}
.frame-header-box .user-nav-box .container {
  margin-bottom: 0px !important;
  padding: 5px 15px;
}
.frame-header-box .user-nav-box .row {
  position: relative;
}
.frame-header-box .user-nav-box .col {
  padding-right: 5px;
  padding-left: 5px;
}
.frame-header-box .user-nav-box .header-row {
  margin: 0px auto;
  width: 100%;
}
.frame-header-box .user-nav-box .header-row .card-body {
  padding-bottom: 0px;
}
.frame-header-box .user-nav-box .body-row {
  margin: 0px auto;
  width: 100%;
}
.frame-header-box .user-nav-box .body-row .col {
  width: 33.33333%;
}
.frame-header-box .user-nav-box .footer-row {
  margin: 0px auto;
  width: 100%;
}
.frame-header-box .user-nav-box .footer-row .card-body {
  padding: 0px;
}
.frame-header-box .user-nav-box .body-row .card {
  margin-bottom: 10px;
}
.frame-header-box .user-nav-box .card-text {
  font-size: 12px;
}
</style>
