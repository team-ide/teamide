<template>
  <div class="app-nav-box">
    <div class="app-nav-header">系统管理</div>
    <div class="app-nav-body">
      <ul>
        <template v-for="(one, index) in source.login.manageNavs">
          <li :key="index">
            <div class="app-nav">
              <template v-if="tool.isEmpty(one.link)">
                <div class="app-nav-text">{{ one.name }}</div>
              </template>
              <template v-else>
                <div
                  class="app-nav-link app-nav-text"
                  @click="$router.push(`${one.link}`)"
                  :class="{ 'app-nav-active': $route.path == `${one.link}` }"
                >
                  {{ one.name }}
                </div>
              </template>
            </div>
            <template v-if="one.navs != null && one.navs.length > 0">
              <ul>
                <template v-for="(sub, subIndex) in one.navs">
                  <li :key="subIndex">
                    <div class="app-nav">
                      <template v-if="tool.isEmpty(sub.link)">
                        <div class="app-nav-text">{{ sub.name }}</div>
                      </template>
                      <template v-else>
                        <div
                          class="app-nav-link app-nav-text"
                          @click="$router.push(`${sub.link}`)"
                          :class="{
                            'app-nav-active': $route.path == `${sub.link}`,
                          }"
                        >
                          {{ sub.name }}
                        </div>
                      </template>
                    </div>
                  </li>
                </template>
              </ul>
            </template>
          </li>
        </template>
      </ul>
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
</style>
