<template>
  <el-dialog
    ref="modal"
    :title="
      '导出：[' +
      ownerName +
      '].[' +
      (tableDetail == null ? '' : tableDetail.tableName) +
      '] 数据为SQL'
    "
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="mgt--20 toolbox-database-save-sql">
      <el-form ref="form" :model="form" size="mini" :inline="true">
        <Pack
          :source="source"
          :toolboxWorker="toolboxWorker"
          :form="form"
          :change="toLoad"
        >
        </Pack>
      </el-form>
      <div style="height: 480px !important">
        <Editor
          ref="Editor"
          :source="source"
          :value="showSQL"
          language="sql"
        ></Editor>
      </div>
    </div>
  </el-dialog>
</template>

<script>
import Pack from "./Pack";

export default {
  components: { Pack },
  props: ["source", "toolboxWorker"],
  data() {
    return {
      showDialog: false,
      showSQL: null,
      ownerName: null,
      tableDetail: null,
      form: {
        targetDatabaseType: "",
        appendOwnerName: true,
        ownerNamePackChar: "",
        tableNamePackChar: "",
        columnNamePackChar: "",
        sqlValuePackChar: "",
      },
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
    async show(ownerName, tableDetail, params) {
      this.ownerName = ownerName;
      this.tableDetail = tableDetail;
      this.insertList = params.insertList;
      this.updateList = params.updateList;
      this.updateWhereList = params.updateWhereList;
      this.deleteList = params.deleteList;
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
      let valuesList = res.valuesList || [];
      sqlList.forEach((sql) => {
        this.showSQL += sql + ";\n\n";
      });
      this.$refs.Editor.setValue(this.showSQL);
    },
    async loadSqls() {
      let data = Object.assign({}, this.form);
      this.toolboxWorker.formatParam(data);

      data.appendSqlValue = true;
      data.ownerName = this.ownerName;
      data.tableName = this.tableDetail.tableName;
      data.columnList = this.tableDetail.columnList;

      data.insertList = this.insertList;
      data.updateList = this.updateList;
      data.updateWhereList = this.updateWhereList;
      data.deleteList = this.deleteList;

      let res = await this.toolboxWorker.work("dataListSql", data);
      if (res.code != 0) {
        return;
      }
      return res.data || {};
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolboxWorker.showTableDataListExecSql = this.show;
    this.init();
  },
};
</script>

<style>
.toolbox-database-save-sql-textarea {
  width: 100%;
  height: 400px;
  letter-spacing: 1px;
  word-spacing: 5px;
  word-break: break-all;
  font-size: 12px;
  border: 1px solid #ddd;
  padding: 0px 5px;
  outline: none;
  user-select: none;
  resize: none;
}
</style>
