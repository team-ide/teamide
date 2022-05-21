<template>
  <el-dialog
    ref="modal"
    :title="'导出：' + (tableDetail == null ? '' : tableDetail.name)"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="mgt--20 toolbox-database-export-sql">
      <el-form ref="form" :model="form" size="mini" :inline="true">
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
              style="width: 90px"
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
            style="width: 90px"
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
            style="width: 90px"
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
            style="width: 60px"
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
      <div>
        <textarea v-model="showSQL"> </textarea>
      </div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      showDialog: false,
      showSQL: null,
      tableDetail: null,
      sqlTypes: [
        { value: "insert", text: "Insert" },
        { value: "update", text: "Update" },
        { value: "delete", text: "Delete" },
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
        appendDatabase: true,
        databasePackingCharacter: "`",
        tablePackingCharacter: "`",
        columnPackingCharacter: "`",
        stringPackingCharacter: "'",
        dateFunction: "",
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(database, tableDetail, dataList) {
      this.database = database;
      this.dataList = dataList || [];
      this.tableDetail = tableDetail;
      await this.toLoad();
      this.showDialog = true;
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
    },
    async loadSqls() {
      let data = Object.assign({}, this.form);
      let insertList = [];
      let updateList = [];
      let updateWhereList = [];
      let deleteList = [];

      let keys = [];
      this.tableDetail.columnList.forEach((column) => {
        if (column.primaryKey) {
          keys.push(column.name);
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
      data.database = this.database;
      data.table = this.tableDetail.name;
      data.columnList = this.tableDetail.columnList;

      data.insertList = insertList;
      data.updateList = updateList;
      data.updateWhereList = updateWhereList;
      data.deleteList = deleteList;

      let res = await this.wrap.work("dataListSql", data);
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
    this.wrap.showExportSql = this.show;
    this.init();
  },
};
</script>

<style>
.toolbox-database-export-sql textarea {
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
