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
        <b-img
          v-bind="{
            blank: true,
            blankColor: '#777',
            width: 20,
            height: 20,
          }"
          rounded="circle"
          class="mgt-5 mgr-5"
          :src="source.login.user.avatarUrl"
        ></b-img>
        <el-dropdown
          trigger="click"
          size="mini"
          placement="left"
          ref="dropdown"
        >
          <div class="tm-link user-name">{{ source.login.user.name }}</div>
          <el-dropdown-menu slot="dropdown">
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
                        <template
                          v-if="
                            source.login.groups != null &&
                            source.login.groups.length > 0
                          "
                        >
                          <div class="mgt-5 ft-12">
                            <template
                              v-for="(one, index) in source.login.groups"
                            >
                              <span
                                :key="'group-' + index"
                                class="mglr-2"
                                :class="[`color-${one.color}`]"
                              >
                                {{ one.name }}
                                <b-icon
                                  v-if="one.leader"
                                  icon="gear"
                                  title="组长"
                                ></b-icon>
                              </span>
                            </template>
                          </div>
                        </template>
                      </b-card-text>
                    </b-card>
                    <hr />
                  </b-col>
                </b-row>
                <b-row
                  v-if="source.frame.userNavs.length > 0"
                  cols="3"
                  class="body-row"
                >
                  <template v-for="(one, index) in source.frame.userNavs">
                    <template v-if="index < 6">
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
                              @click="$router.push(`${one.link}`)"
                              :class="[
                                `color-${one.color}`,
                                $route.path == `${one.link}` ||
                                (one.match && one.match($route.path))
                                  ? 'tm-active'
                                  : '',
                              ]"
                            >
                              {{ one.name }}
                            </div>
                          </b-card-text>
                        </b-card>
                      </b-col>
                    </template>
                  </template>
                  <b-col>
                    <b-card style="border: 0px" align="center">
                      <b-card-text style="padding-top: 20px">
                        <div
                          @click="$router.push(`/user`)"
                          class="tm-link color-grey"
                        >
                          更多功能
                        </div>
                      </b-card-text>
                    </b-card>
                  </b-col>
                </b-row>
                <b-row
                  v-if="source.hasPower('logout')"
                  cols="1"
                  class="footer-row"
                >
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
          </el-dropdown-menu>
        </el-dropdown>
      </template>
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
