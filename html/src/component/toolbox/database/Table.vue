<template>
  <div class="toolbox-database-table">
    <div class="scrollbar" style="height: calc(100% - 90px)">
      <TableDetail
        class="pd-10"
        ref="TableDetail"
        :source="source"
        :wrap="wrap"
        :toolbox="toolbox"
        :onChange="onTableDetailChange"
      ></TableDetail>
      <el-form
        class="database-table-detail-form pdlr-10"
        ref="form"
        size="mini"
        inline
      >
        <el-form-item label="库名">
          <el-input
            v-model="form.database"
            @change="toLoad"
            style="width: 120px"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="新建表">
          <el-switch
            v-model="isInsert"
            @change="toLoad"
            :readonly="tableDetail == null"
          >
          </el-switch>
        </el-form-item>
        <el-form-item label="数据库类型">
          <el-select
            placeholder="当前库类型"
            v-model="form.databaseType"
            style="width: 120px"
          >
            <el-option label="MySql" value="mysql"> </el-option>
            <el-option label="Sqlite" value="sqlite"> </el-option>
            <el-option label="Oracle" value="oracle"> </el-option>
            <el-option label="达梦" value="dameng"> </el-option>
            <el-option label="神通" value="shentong"> </el-option>
            <el-option label="金仓" value="kingbase"> </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="追加库名">
          <el-switch v-model="form.appendDatabase" @change="toLoad">
          </el-switch>
        </el-form-item>
        <template v-if="form.appendDatabase">
          <el-form-item label="库名包装">
            <el-select
              placeholder="不包装"
              v-model="form.databasePackingCharacter"
              @change="toLoad"
              style="width: 120px"
            >
              <el-option
                v-for="(one, index) in packingCharacters"
                :key="index"
                :value="one.value"
                :label="one.text"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="表名包装">
          <el-select
            placeholder="不包装"
            v-model="form.tablePackingCharacter"
            @change="toLoad"
            style="width: 120px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字段包装">
          <el-select
            placeholder="不包装"
            v-model="form.columnPackingCharacter"
            @change="toLoad"
            style="width: 120px"
          >
            <el-option
              v-for="(one, index) in packingCharacters"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="字符值包装">
          <el-select
            v-model="form.stringPackingCharacter"
            @change="toLoad"
            style="width: 120px"
          >
            <el-option
              v-for="(one, index) in stringPackingCharacters"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="pdlr-10">
        <div class="pdb-5">SQL预览</div>
        <textarea v-model="showSQL" class="database-show-sql"> </textarea>
      </div>
    </div>
    <div class="" v-if="error != null">
      <div class="bg-red ft-12 pd-5">{{ error }}</div>
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
import TableDetail from "./TableDetail.vue";

export default {
  components: { TableDetail },
  props: ["source", "toolbox", "wrap", "database", "table"],
  data() {
    return {
      showSQL: null,
      characterSets: [
        { value: "utf8", text: "utf8" },
        { value: "utf8mb4", text: "utf8mb4" },
      ],
      packingCharacters: [
        { value: "", text: "不包装" },
        { value: "'", text: "'" },
        { value: '"', text: '"' },
        { value: "`", text: "`" },
      ],
      stringPackingCharacters: [
        { value: "'", text: "'" },
        { value: '"', text: '"' },
      ],
      form: {
        database: null,
        databaseType: null,
        name: "TB_XXX",
        comment: "",
        appendDatabase: false,
        databasePackingCharacter: "`",
        tablePackingCharacter: "`",
        columnPackingCharacter: "`",
        stringPackingCharacter: "'",
        columnList: [],
        indexList: [],
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
  watch: {
    "form.databaseType"() {
      if (
        this.tool.isEmpty(this.form.databaseType) ||
        this.form.databaseType == "mysql"
      ) {
        this.form.databasePackingCharacter = "`";
        this.form.tablePackingCharacter = "`";
        this.form.columnPackingCharacter = "`";
      } else {
        this.form.databasePackingCharacter = `"`;
        this.form.tablePackingCharacter = `"`;
        this.form.columnPackingCharacter = `"`;
      }
      this.toLoad();
    },
  },
  methods: {
    async init() {
      let tableDetail = null;
      if (this.tool.isNotEmpty(this.table)) {
        tableDetail = await this.wrap.getTableDetail(this.database, this.table);
      }

      this.show(this.database, tableDetail);
    },
    async show(database, tableDetail) {
      this.form.database = database;
      this.tableDetail = tableDetail;
      this.isInsert = tableDetail == null;
      this.error = null;
      this.executeSqlIng = false;

      this.form.name = "TB_XXX";
      this.form.comment = "";
      this.form.oldName = "";
      this.form.oldComment = "";
      this.form.columnList.splice(0, this.form.columnList.length);
      this.form.indexList.splice(0, this.form.indexList.length);
      if (tableDetail != null) {
        this.form.name = tableDetail.name;
        this.form.comment = tableDetail.comment;
        this.form.oldName = tableDetail.name;
        this.form.oldComment = tableDetail.comment;
        if (tableDetail.columnList) {
          let lastColumn = null;
          tableDetail.columnList.forEach((one) => {
            let column = Object.assign({}, one);
            column.oldName = column.name;
            column.oldComment = column.comment;
            column.oldType = column.type;
            column.oldLength = column.length;
            column.oldDecimal = column.decimal;
            column.oldPrimaryKey = column.primaryKey;
            column.oldNotNull = column.notNull;
            column.oldDefault = column.default;
            column.deleted = false;
            column.oldBeforeColumn = lastColumn;
            this.form.columnList.push(column);
            lastColumn = column;
          });
        }
        if (tableDetail.indexList) {
          tableDetail.indexList.forEach((one) => {
            let index = Object.assign({}, one);
            index.oldName = index.name;
            index.oldComment = index.comment;
            index.oldType = index.type;
            index.oldColumns = index.columns;
            index.deleted = false;
            this.form.indexList.push(index);
          });
        }
      }
      await this.toLoad();
      this.$nextTick(() => {
        this.$refs.TableDetail.init(this.form);
      });
    },
    getFormData() {
      this.form.columnList.forEach((one) => {
        one.length = Number(one.length);
        one.decimal = Number(one.decimal);
      });
      let data = Object.assign({}, this.form);
      data.columnList = [];
      this.form.columnList.forEach((one) => {
        let column = Object.assign({}, one);
        delete column.beforeColumn_;
        delete column.oldBeforeColumn;
        data.columnList.push(column);
      });
      return data;
    },
    async toExecuteSql() {
      let data = this.getFormData();

      this.executeSqlIng = true;
      let res = null;
      if (this.isInsert) {
        res = await this.wrap.work("createTable", data);
      } else {
        res = await this.wrap.work("updateTable", data);
      }
      this.executeSqlIng = false;
      this.error = null;
      if (res.code != 0) {
        this.error = res.msg;
        return;
      }
      this.tool.success("执行成功");
      this.init();
      return res.data || {};
    },
    async onTableDetailChange() {
      await this.toLoad();
    },
    async toLoad() {
      this.showSQL = "";
      let res = await this.loadSqls();
      let sqlList = res.sqlList || [];
      sqlList.forEach((sql) => {
        this.showSQL += sql + ";\n\n";
      });
    },
    async loadSqls() {
      let data = this.getFormData();
      let res = null;
      if (this.isInsert) {
        res = await this.wrap.work("createTableSql", data);
      } else {
        res = await this.wrap.work("updateTableSql", data);
      }
      this.error = null;
      if (res.code != 0) {
        this.error = res.msg;
        return;
      }
      return res.data || {};
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
.toolbox-database-table {
  width: 100%;
  height: 100%;
}
.database-show-sql {
  width: 100%;
  height: 300px;
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
