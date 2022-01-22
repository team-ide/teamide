<template>
  <div class="" style="display: flex">
    <template v-if="application.apps == null"> </template>
    <template v-else-if="application.apps.length == 0">
      <div
        class="mgt-10 mgl-20 ft-21 ft-800"
        style="float: left; line-height: 30px"
      >
        暂无应用
      </div>
    </template>
    <template v-else-if="app != null">
      <div
        class="mgt-10 mgl-20 ft-21 ft-800"
        style="float: left; line-height: 30px"
      >
        {{ app.name }}
      </div>
      <b-dropdown class="mgt-10 mgl-20 pd-0" size="sm">
        <template v-for="(one, index) in application.apps">
          <b-dropdown-item
            :disabled="one.name == app.name"
            :key="index"
            @click="changeApp(one)"
          >
            {{ one.name }}
          </b-dropdown-item>
        </template>
      </b-dropdown>
    </template>
    <div
      v-if="source.hasPower('application_insert')"
      class="tm-link color-green mgt-10 mgl-15 ft-19"
      title="新建应用"
      @click="toInsert()"
    >
      <b-icon icon="plus-square"></b-icon>
    </div>
    <template v-if="app != null">
      <div
        v-if="source.hasPower('application_delete')"
        class="tm-link color-orange mgt-10 mgl-15 ft-19"
        title="新建应用"
        @click="toDelete(app)"
      >
        <b-icon icon="x-square"></b-icon>
      </div>
    </template>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "application", "app"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    changeApp(app) {
      this.application.app = app;
    },
    async init() {
      await this.loadList();
      await this.initApp();
    },
    async loadList() {
      let param = {};
      let res = await this.server.application.list(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let data = res.data || {};
      let list = data.list || [];
      this.application.apps = list;
    },
    async initApp() {
      if (this.application.app != null && this.application.apps != null) {
        let find = null;
        this.application.apps.forEach((one) => {
          if (one.name == this.application.app.name) {
            find = one;
          }
        });
        if (find == null) {
          this.application.app = null;
        }
      }
      if (this.application.app != null) {
        return;
      }
      if (this.application.apps.length > 0) {
        this.application.app = this.application.apps[0];
      }
    },
    toInsert() {
      let data = {};
      this.application.showAppForm(data, this.doInsert);
    },
    toUpdate(app) {
      let data = {
        name: app.name,
      };
      this.updateData = app;
      this.application.showAppForm(data, this.doUpdate);
    },
    toDelete(app) {
      this.tool
        .confirm("删除应用[" + app.name + "]将无法恢复，确定删除？")
        .then(() => {
          this.doDelete(app);
        })
        .catch((e) => {});
    },
    async doInsert(app) {
      let find;
      this.application.apps.forEach((one) => {
        if (one.name == app.name) {
          find = one;
        }
      });
      if (find != null) {
        this.tool.error("应用[" + app.name + "]已存在");
        return false;
      }
      let res = await this.server.application.insert(app);
      if (res.code == 0) {
        this.init();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(app) {
      let find;
      this.application.apps.forEach((one) => {
        if (one.name == app.name) {
          find = one;
        }
      });
      if (find != null && find != this.updateData) {
        this.tool.error("应用[" + app.name + "]已存在");
        return false;
      }
      let res = await this.server.application.rename({
        name: this.updateData.name,
        rename: app.name,
      });
      if (res.code == 0) {
        this.init();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doDelete(app) {
      let res = await this.server.application.delete({
        name: app.name,
      });
      if (res.code == 0) {
        this.init();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
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
</style>
