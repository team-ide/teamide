<template>
  <div class="frame-body-box">
    <template v-if="!source.hasPower($route.path)">
      <NoPower :source="source"></NoPower>
    </template>
    <template v-else>
      <template
        v-if="tool.isManagePage($route.path) || tool.isUserPage($route.path)"
      >
        <div class="frame-body">
          <div class="frame-body-left">
            <template v-if="tool.isManagePage($route.path)">
              <ManageNav :source="source"></ManageNav>
            </template>
            <template v-if="tool.isUserPage($route.path)">
              <UserNav :source="source"></UserNav>
            </template>
          </div>
          <div class="frame-body-center">
            <router-view :source="source"></router-view>
          </div>
          <div class="frame-body-right"></div>
        </div>
      </template>
      <template v-else>
        <router-view :source="source"></router-view>
      </template>
    </template>
  </div>
</template>

<script>
import ManageNav from "@/views/manage/Nav.vue";
import UserNav from "@/views/user/Nav.vue";
export default {
  components: { ManageNav, UserNav },
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
.frame-body-box {
  width: 100%;
  height: calc(100% - 56px);
  background: #ddd;
  position: relative;
  overflow: auto;
  padding: 0px 0px;
}
.frame-body {
  min-width: 1024px;
  max-width: 1200px;
  margin: 0px auto;
  min-height: calc(100% - 0px);
  background: #fff;
  display: flex;
}
.frame-body-left {
  width: 280px;
  border-right: 1px solid #ddd;
}
.frame-body-center {
  flex: 1;
}
.frame-body-right {
  width: 0px;
  border-left: 0px solid #ddd;
}
.app-nav-box {
  background: #343a40;
  color: #979696;
  height: 100%;
}
.app-nav-header {
  padding: 10px 30px;
  border-bottom: 1px solid #595959;
  font-weight: 600;
  font-size: 22px;
}
.app-nav-body {
  padding: 0px 0px;
  min-height: calc(100% - 54px);
}
.app-nav-body ul {
  font-size: 16px;
}
.app-nav-body ul ul {
  font-size: 15px;
}
.app-nav-body ul,
.app-nav-body li {
  list-style: none;
  padding: 0px;
  margin: 0px;
}
.app-nav-body > ul > li {
  border-bottom: 1px solid #595959;
  padding-bottom: 10px;
  padding-top: 10px;
}
.app-nav-body ul .app-nav-text {
  padding-left: 30px;
  line-height: 30px;
}
.app-nav-body ul ul .app-nav-text {
  padding-left: 60px;
}
.app-nav-body ul .app-nav-link {
  display: block;
  cursor: pointer;
  border-right: 3px solid transparent;
}
.app-nav-body ul .app-nav-link.tm-active {
  border-color: #cfcfcf;
  color: #cfcfcf;
}
.app-nav-body ul .app-nav-link:hover:not(.tm-active) {
  border-color: #bdbdbd;
  color: #bdbdbd;
}
.page-box {
  width: 100%;
  height: 100%;
}
.page-header {
  padding: 10px 20px;
  border-bottom: 1px solid #ddd;
  font-weight: 600;
  font-size: 22px;
}
.page-body {
  padding: 10px 20px;
}
</style>
