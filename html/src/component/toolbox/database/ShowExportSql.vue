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
    <div class="mgt--20 toolbox-database-export-sql">
      <el-form ref="form" :model="form" size="mini" inline>
        <el-form-item label="SQL类型">
          <el-select
            v-model="form.sqlType"
            @change="toLoad"
            style="width: 100px"
          >
            <el-option
              v-for="(one, index) in sqlTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <Pack
          :source="source"
          :toolboxWorker="toolboxWorker"
          :form="form"
          :change="toLoad"
        >
        </Pack>
        <el-form-item label="日期函数">
          <el-select
            v-model="form.dateFunction"
            @change="toLoad"
            style="width: 120px"
          >
            <el-option
              v-for="(one, index) in dateFunctions"
              :key="index"
              :value="one.value"
              :label="one.text"
            >
            </el-option>
          </el-select>
        </el-form-item>
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
      sqlTypes: [
        { value: "insert", text: "Insert" },
        { value: "update", text: "Update" },
        { value: "delete", text: "Delete" },
      ],
      dateFunctions: [
        {
          value: "",
          text: "无函数",
        },
        {
          value: "to_date('$value','yyyy-mm-dd hh24:mi:ss')",
          text: "to_date('$value','yyyy-mm-dd hh24:mi:ss')",
        },
        {
          value: "to_timestamp('$value','yyyy-mm-dd hh24:mi:ss')",
          text: "to_timestamp('$value','yyyy-mm-dd hh24:mi:ss')",
        },
      ],
      form: {
        sqlType: "insert",
        targetDatabaseType: "",
        appendOwnerName: true,
        ownerNamePackChar: "",
        tableNamePackChar: "",
        columnNamePackChar: "",
        sqlValuePackChar: "",
        dateFunction: "",
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
    async show(ownerName, tableDetail, dataList) {
      this.ownerName = ownerName;
      this.dataList = dataList || [];
      this.tableDetail = tableDetail;
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

      let insertList = [];
      let updateList = [];
      let updateWhereList = [];
      let deleteList = [];

      let keys = [];
      this.tableDetail.columnList.forEach((column) => {
        if (column.primaryKey) {
          keys.push(column.columnName);
        }
      });

      this.dataList.forEach((data) => {
        switch (this.form.sqlType) {
          case "insert":
            insertList.push(data);
            break;
          case "update":
            if (keys.length > 0) {
              let whereData = {};
              keys.forEach((key) => {
                whereData[key] = data[key];
              });
              updateWhereList.push(whereData);
            } else {
              updateWhereList.push(data);
            }

            if (keys.length > 0) {
              let updateData = {};
              for (let name in data) {
                if (keys.indexOf(name) < 0) {
                  updateData[name] = data[name];
                }
              }
              updateList.push(updateData);
            } else {
              updateList.push(data);
            }
            break;
          case "delete":
            if (keys.length > 0) {
              let whereData = {};
              keys.forEach((key) => {
                whereData[key] = data[key];
              });
              deleteList.push(whereData);
            } else {
              deleteList.push(data);
            }
            break;
        }
      });

      data.appendSqlValue = true;
      data.ownerName = this.ownerName;
      data.tableName = this.tableDetail.tableName;
      data.columnList = this.tableDetail.columnList;

      data.insertList = insertList;
      data.updateList = updateList;
      data.updateWhereList = updateWhereList;
      data.deleteList = deleteList;

      let res = await this.toolboxWorker.work("dataListSql", data);
      if (res.code != 0) {
        return {};
      }
      return res.data || {};
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolboxWorker.showExportSql = this.show;
    this.init();
  },
};
</script>

<style>
</style>
