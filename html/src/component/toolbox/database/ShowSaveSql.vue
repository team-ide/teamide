<template>
  <el-dialog
    ref="modal"
    :title="
      '导出：[' +
      database +
      '].[' +
      (tableDetail == null ? '' : tableDetail.name) +
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
  props: ["source", "toolboxWorker"],
  data() {
    return {
      showDialog: false,
      showSQL: null,
      database: null,
      tableDetail: null,
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
        appendDatabase: true,
        databasePackingCharacter: "`",
        tablePackingCharacter: "`",
        columnPackingCharacter: "`",
        stringPackingCharacter: "'",
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(database, tableDetail, params) {
      this.database = database;
      this.tableDetail = tableDetail;
      this.insertList = params.insertList;
      this.updateList = params.updateList;
      this.updateWhereList = params.updateWhereList;
      this.deleteList = params.deleteList;
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

      data.appendSqlValue = true;
      data.database = this.database;
      data.table = this.tableDetail.name;
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
    this.toolboxWorker.showSaveSql = this.show;
    this.init();
  },
};
</script>

<style>
.toolbox-database-save-sql textarea {
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
