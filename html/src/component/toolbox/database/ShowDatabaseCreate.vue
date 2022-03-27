<template>
  <b-modal
    ref="modal"
    :title="database == null ? '' : database.name"
    :hide-header-close="false"
    :no-close-on-backdrop="true"
    :no-close-on-esc="true"
    :hide-backdrop="true"
    hide-footer
    size="lg"
  >
    <div class="ft-15">
      <b-form-textarea
        size="sm"
        rows="10"
        max-rows="30"
        v-model="createSql"
        readonly
      >
      </b-form-textarea>
    </div>
  </b-modal>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      createSql: null,
      database: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(database) {
      this.database = database;
      this.createSql = await this.loadCreateDatabase(database.name);
      this.$refs["modal"].show();
    },
    hide() {
      this.$refs["modal"].hide();
    },
    async loadCreateDatabase(database) {
      let param = {
        database: database,
      };
      let res = await this.wrap.work("showCreateDatabase", param);
      res.data = res.data || {};
      return res.data.create;
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.wrap.showDatabaseCreate = this.show;
    this.init();
  },
};
</script>

<style>
</style>
