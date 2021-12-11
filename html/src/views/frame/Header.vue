<template>
  <div class="frame-header-box">
    <b-navbar
      :toggleable="source.header.toggleable"
      :type="source.header.type"
      :variant="source.header.variant"
    >
      <b-navbar-brand :href="source.url" tag="h1" class="mb-0">
        {{ source.header.title }}
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item @click="$router.push(`/`)" :active="$route.path == `/`">
            首页
          </b-nav-item>
          <b-nav-item
            @click="$router.push(`/workspace`)"
            :active="$route.path == `/workspace`"
          >
            工作台
          </b-nav-item>
          <b-nav-item
            @click="$router.push(`/toolbox`)"
            :active="$route.path == `/toolbox`"
          >
            工具箱
          </b-nav-item>
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-form>
            <b-form-input
              size="sm"
              class="mr-sm-2"
              placeholder="Search"
            ></b-form-input>
            <b-button size="sm" class="my-2 my-sm-0" type="submit">
              Search
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

          <b-nav-item-dropdown right class="user-dropdown">
            <template #button-content>
              <b-img
                v-bind="{
                  blank: true,
                  blankColor: '#777',
                  width: 30,
                  height: 30,
                  class: 'm1',
                }"
                rounded="circle"
                alt="Circle image"
              ></b-img>
              <span class="user-name">张三</span>
            </template>
            <div class="user-nav-box">
              <b-container class="bv-example-row mb-3">
                <b-row cols="3">
                  <b-col>
                    <b-card>
                      <b-card-title>
                        <b-icon icon="camera"></b-icon>
                      </b-card-title>
                      <b-card-text>个人首页</b-card-text>
                    </b-card>
                  </b-col>
                  <b-col>
                    <b-card>
                      <b-card-title>
                        <b-icon icon="camera"></b-icon>
                      </b-card-title>
                      <b-card-text>个人中心</b-card-text>
                    </b-card>
                  </b-col>
                  <b-col>
                    <b-card>
                      <b-card-title>
                        <b-icon icon="camera"></b-icon>
                      </b-card-title>
                      <b-card-text>修改密码</b-card-text>
                    </b-card>
                  </b-col>
                  <b-col>
                    <b-card>
                      <b-card-title>
                        <b-icon icon="camera"></b-icon>
                      </b-card-title>
                      <b-card-text>消息</b-card-text>
                    </b-card>
                  </b-col>
                  <b-col>
                    <b-card>
                      <b-card-title>
                        <b-icon icon="camera"></b-icon>
                      </b-card-title>
                      <b-card-text>帮助中心</b-card-text>
                    </b-card>
                  </b-col>
                </b-row>
              </b-container>
            </div>
          </b-nav-item-dropdown>
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
  width: 400px;
  margin: 5px 5px;
}
.user-dropdown .user-nav-box .container {
  margin-bottom: 0px !important;
  padding: 5px 15px;
}
.user-dropdown .user-nav-box .col {
  padding-right: 5px;
  padding-left: 5px;
}
.user-dropdown .user-nav-box .card {
  text-align: center;
  margin-bottom: 10px;
}
.user-dropdown .user-nav-box .card-text {
  font-size: 12px;
}
</style>
