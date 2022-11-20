<template>
  <el-dialog
    ref="modal"
    title="SQL"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1000px"
  >
    <div class="mgt--20">
      <el-form class="" size="mini" inline>
        <Pack
          :source="source"
          :toolboxWorker="toolboxWorker"
          :form="form"
          :change="toLoad"
        >
        </Pack>
      </el-form>
      <div class="">
        <div style="height: 480px !important">
          <Editor
            ref="Editor"
            :source="source"
            :value="showSQL"
            language="sql"
          ></Editor>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script>
import Pack from "./Pack";

export default {
  components: { Pack },
  props: [
    "source",
    "toolboxWorker",
    "isInsert",
    "form",
    "getFormData",
    "onError",
  ],
  data() {
    return {
      showSQL: null,
      showDialog: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "form.targetDatabaseType"() {
      this.toLoad();
    },
  },
  methods: {
    async show() {
      this.showDialog = true;
      await this.toLoad();
    },
    hide() {
      this.showDialog = false;
    },
    async toLoad() {
      this.showSQL = "";
      let res = await this.loadSqls();
      let sqlList = res.sqlList || [];
      sqlList.forEach((sql) => {
        this.showSQL += sql + ";\n\n";
      });
      this.$refs.Editor.setValue(this.showSQL);
    },
    async loadSqls() {
      let data = this.getFormData();
      let res = null;
      if (this.isInsert) {
        res = await this.toolboxWorker.work("tableCreateSql", data);
      } else {
        res = await this.toolboxWorker.work("tableUpdateSql", data);
      }
      this.error = null;
      if (res.code != 0) {
        this.onError(res.msg);
        return;
      }
      return res.data || {};
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>
