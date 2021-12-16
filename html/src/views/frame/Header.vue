<template>
  <div class="frame-header-box">
    <b-navbar
      :toggleable="source.header.toggleable"
      :type="source.header.type"
      :variant="source.header.variant"
    >
      <b-navbar-brand :href="source.url" class="ft-20 ft-900 color-grey">
        {{ source.header.title }}
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <template v-for="(one, index) in source.header.navs">
            <b-nav-item
              :key="index"
              @click="$router.push(`${one.link}`)"
              :active="
                $route.path == `${one.link}` ||
                (one.match && one.match($route.path))
              "
            >
              {{ one.name }}
            </b-nav-item>
          </template>
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-form>
            <b-form-input
              size="sm"
              class="mr-sm-2"
              placeholder="请输入搜索"
            ></b-form-input>
            <b-button size="sm">
              <b-icon icon="search"></b-icon>
            </b-button>
          </b-nav-form>
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right>
            <template #button-content>
              <span>语言:中文</span>
            </template>
            <b-dropdown-item href="#">中文</b-dropdown-item>
            <b-dropdown-item href="#" disabled>英文</b-dropdown-item>
          </b-nav-item-dropdown>

          <template v-if="source.login.user == null">
            <b-nav-item v-if="source.hasPower('login')" @click="tool.toLogin()">
              登录
            </b-nav-item>
          </template>
          <template v-else>
            <b-nav-item-dropdown right class="user-dropdown">
              <template #button-content>
                <b-img
                  v-bind="{
                    blank: true,
                    blankColor: '#777',
                    width: 30,
                    height: 30,
                  }"
                  rounded="circle"
                  :src="source.login.user.avatarUrl"
                ></b-img>
                <span class="user-name">{{ source.login.user.name }}</span>
              </template>
              <li class="user-nav-box">
                <b-container class="bv-example-row">
                  <b-row cols="1" class="header-row">
                    <b-col>
                      <b-card class="bd-0 text-center">
                        <b-card-title>
                          <b-img
                            v-bind="{
                              blank: true,
                              blankColor: '#777',
                              width: 80,
                              height: 80,
                            }"
                            rounded="circle"
                            :src="source.login.user.avatarUrl"
                          ></b-img>
                        </b-card-title>
                        <b-card-text>
                          <div class="ft-15">
                            {{ source.login.user.name }}
                          </div>
                          <div class="mgt-5 ft-12">
                            <span class="mglr-2 color-green"> A组 </span>
                            <span class="mglr-2 color-orange">
                              B组
                              <b-icon icon="gear" title="组长"></b-icon>
                            </span>
                            <span class="mglr-2 color-blue"> E组 </span>
                          </div>
                          <div class="mgt-5 ft-12">
                            <span class="mglr-2 color-green"> A部门 </span>
                            <span class="mglr-2 color-orange"> B部门 </span>
                            <span class="mglr-2 color-blue"> E部门 </span>
                          </div>
                        </b-card-text>
                      </b-card>
                      <hr />
                    </b-col>
                  </b-row>
                  <b-row cols="3" class="body-row">
                    <template v-for="(one, index) in source.login.headerNavs">
                      <b-col :key="index">
                        <b-card align="center" :class="`bd-${one.color}`">
                          <b-card-title>
                            <b-icon
                              :icon="one.icon"
                              :class="`color-${one.color}`"
                            ></b-icon>
                          </b-card-title>
                          <b-card-text>
                            <div
                              class="tm-link"
                              @click="$router.push(`/${one.link}`)"
                              :class="`color-${one.color}`"
                            >
                              {{ one.name }}
                            </div>
                          </b-card-text>
                        </b-card>
                      </b-col>
                    </template>
                    <b-col>
                      <b-card style="border: 0px" align="center">
                        <b-card-text style="padding-top: 20px">
                          <div href="#" class="tm-link color-grey">
                            更多功能
                          </div>
                        </b-card-text>
                      </b-card>
                    </b-col>
                  </b-row>
                  <b-row cols="1" class="footer-row">
                    <b-col>
                      <hr />
                      <b-card class="bd-0 text-center">
                        <b-card-text>
                          <div
                            v-if="source.hasPower('logout')"
                            @click="tool.toLogout()"
                            class="tm-link color-orange"
                          >
                            登出
                          </div>
                        </b-card-text>
                      </b-card>
                    </b-col>
                  </b-row>
                </b-container>
              </li>
            </b-nav-item-dropdown>
          </template>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
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
  mounted() {
    console.log(this);
  },
};
</script>

<style>
.user-dropdown .dropdown-toggle {
  padding-top: 5px;
  padding-bottom: 5px;
}
.user-dropdown .dropdown-toggle img {
  margin-right: 0px;
}
.user-dropdown .dropdown-toggle .user-name {
  vertical-align: middle;
  font-size: 14px;
  margin-left: 8px;
}
.user-dropdown .dropdown-toggle::after {
  vertical-align: middle;
  margin-left: 8px;
}
.user-dropdown .dropdown-menu {
  padding: 0px;
}
.user-dropdown .user-nav-box {
  margin: 5px 5px;
  min-width: 400px;
}
.user-dropdown .user-nav-box .container {
  margin-bottom: 0px !important;
  padding: 5px 15px;
}
.user-dropdown .user-nav-box .row {
}
.user-dropdown .user-nav-box .col {
  padding-right: 5px;
  padding-left: 5px;
}
.user-dropdown .user-nav-box .header-row {
  margin: 0px auto;
  width: 100%;
}
.user-dropdown .user-nav-box .header-row .card-body {
  padding-bottom: 0px;
}
.user-dropdown .user-nav-box .body-row {
  margin: 0px auto;
  width: 100%;
}
.user-dropdown .user-nav-box .body-row .col {
  width: 33.33333%;
}
.user-dropdown .user-nav-box .footer-row {
  margin: 0px auto;
  width: 100%;
}
.user-dropdown .user-nav-box .footer-row .card-body {
  padding: 0px;
}
.user-dropdown .user-nav-box .body-row .card {
  margin-bottom: 10px;
}
.user-dropdown .user-nav-box .card-text {
  font-size: 12px;
}
</style>
