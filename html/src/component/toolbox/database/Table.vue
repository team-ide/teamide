<template>
  <div class="toolbox-database-table">
    <el-form class="pdlr-10" size="mini" inline>
      <el-form-item label="库名">
        <el-input v-model="form.ownerName" style="width: 120px"> </el-input>
      </el-form-item>
      <el-form-item label="新建表">
        <el-switch v-model="isInsert" :readonly="tableDetail == null">
        </el-switch>
      </el-form-item>
    </el-form>
    <div style="height: calc(100% - 130px)">
      <TableDetail
        class="pd-10"
        ref="TableDetail"
        :source="source"
        :toolboxWorker="toolboxWorker"
        :onChange="onTableDetailChange"
        :columnTypeInfoList="columnTypeInfoList"
        :indexTypeInfoList="indexTypeInfoList"
        :formData="form"
        :getFormData="getFormData"
        :isInsert="isInsert"
        :onError="onError"
      ></TableDetail>
    </div>
    <div class="" v-if="error != null">
      <div class="bg-red ft-12 mglr-10 pd-5">{{ error }}</div>
    </div>
    <div class="pd-10">
      <div
        class="tm-btn bg-green ft-13"
        @click="toExecuteSql"
        :class="{ 'tm-disabled': executeSqlIng }"
      >
        执行
      </div>
    </div>
  </div>
</template>

<script>
import TableDetail from "./TableDetail";

export default {
  components: { TableDetail },
  props: [
    "source",
    "toolboxWorker",
    "actived",
    "ownerName",
    "tableName",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      showSQL: null,
      characterSets: [
        { value: "utf8", text: "utf8" },
        { value: "utf8mb4", text: "utf8mb4" },
      ],
      form: {
        ownerName: null,
        tableName: "TB_XXX",
        tableComment: "",
        columnList: [],
        indexList: [],
        targetDatabaseType: "",
        appendOwnerName: true,
        ownerNamePackChar: "",
        tableNamePackChar: "",
        columnNamePackChar: "",
        sqlValuePackChar: "",
      },
      tableDetail: null,
      isInsert: false,
      error: null,
      executeSqlIng: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    onFocus() {
      if (this.inited) {
        return;
      }
      this.$nextTick(async () => {
        this.init();
      });
    },
    async init() {
      this.inited = true;
      let tableDetail = null;
      if (this.tool.isNotEmpty(this.tableName)) {
        tableDetail = await this.toolboxWorker.getTableDetail(
          this.ownerName,
          this.tableName
        );
      }

      this.initForm(this.ownerName, tableDetail);
    },
    async initForm(ownerName, tableDetail) {
      this.form.ownerName = ownerName;
      this.tableDetail = tableDetail;
      this.isInsert = tableDetail == null;
      this.error = null;
      this.executeSqlIng = false;

      this.form.tableName = "TB_XXX";
      this.form.tableComment = "";
      this.form.columnList.splice(0, this.form.columnList.length);
      this.form.indexList.splice(0, this.form.indexList.length);
      if (tableDetail != null) {
        this.form.tableName = tableDetail.tableName;
        this.form.tableComment = tableDetail.tableComment;
        this.form.oldTable = tableDetail;
        if (tableDetail.columnList) {
          let lastColumn = null;
          tableDetail.columnList.forEach((one) => {
            let column = Object.assign({}, one);
            column.oldColumn = one;
            column.deleted = false;
            column.oldBeforeColumn = lastColumn;
            this.form.columnList.push(column);
            lastColumn = column;
          });
        }
        if (tableDetail.indexList) {
          tableDetail.indexList.forEach((one) => {
            let index = Object.assign({}, one);
            index.oldIndex = one;
            index.deleted = false;
            this.form.indexList.push(index);
          });
        }
      }
      this.$nextTick(() => {
        this.$refs.TableDetail.init(this.form);
      });
    },
    getFormData() {
      this.form.columnList.forEach((one) => {
        one.columnLength = Number(one.columnLength);
        one.columnDecimal = Number(one.columnDecimal);
      });
      let data = Object.assign({}, this.form);
      data.columnList = [];
      this.form.columnList.forEach((one) => {
        let column = Object.assign({}, one);
        delete column.beforeColumn_;
        delete column.oldBeforeColumn;
        data.columnList.push(column);
      });
      this.toolboxWorker.formatParam(data);
      return data;
    },
    async toExecuteSql() {
      let data = this.getFormData();

      this.executeSqlIng = true;
      let res = null;
      if (this.isInsert) {
        res = await this.toolboxWorker.work("tableCreate", data);
      } else {
        res = await this.toolboxWorker.work("tableUpdate", data);
      }
      this.executeSqlIng = false;
      this.error = null;
      if (res.code != 0) {
        this.error = res.msg;
        return;
      }
      this.tool.success("执行成功");
      let tableDetail = await this.toolboxWorker.getTableDetail(
        data.ownerName,
        data.tableName
      );

      this.initForm(data.ownerName, tableDetail);

      return res.data || {};
    },
    async onTableDetailChange() {},
    onError(error) {
      this.error = error;
    },
  },
  created() {},
  updated() {},
  mounted() {
    if (this.actived) {
      this.init();
    }
  },
};
</script>

<style>
.toolbox-database-table {
  width: 100%;
  height: 100%;
}
</style>
