<template>
  <el-dialog
    ref="modal"
    :title="'创建数据库'"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="700px"
  >
    <el-form
      class="mgt--20"
      ref="form"
      :model="form"
      label-width="120px"
      size="mini"
    >
      <el-form-item label="数据库名称">
        <el-input v-model="form.name" @change="toLoad"> </el-input>
      </el-form-item>
      <el-form-item label="字符集">
        <el-select v-model="form.characterSet" @change="toLoad">
          <el-option
            v-for="(one, index) in characterSets"
            :key="index"
            :value="one.value"
            :label="one.text"
          >
          </el-option>
        </el-select>
      </el-form-item>

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
      <el-form-item label="SQL预览">
        <textarea v-model="showSQL" class="database-show-sql"> </textarea>
      </el-form-item>
    </el-form>
    <div class="tb-5" v-if="error != null">
      <div class="bg-red ft-12">{{ error }}</div>
    </div>
    <div class="tb-5">
      <div
        class="tm-btn bg-green ft-13"
        @click="toExecuteSql"
        :class="{ 'tm-disabled': executeSqlIng }"
      >
        执行
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
      form: {
        name: "XXX_DB",
        characterSet: "utf8mb4",
        databasePackingCharacter: "`",
      },
      error: null,
      executeSqlIng: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(callback) {
      this.callback = callback;
      this.executeSqlIng = false;
      this.error = null;
      await this.toLoad();
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    async toExecuteSql() {
      this.executeSqlIng = true;
      let data = Object.assign({}, this.form);
      let res = await this.toolboxWorker.work("createDatabase", data);
      this.error = null;
      this.executeSqlIng = false;
      if (res.code != 0) {
        this.error = res.msg;
        return;
      }
      this.tool.success("创建成功");
      this.callback && this.callback(true);
      this.hide();
      return res.data || {};
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
      let data = Object.assign({}, this.form);
      let res = await this.toolboxWorker.work("createDatabaseSql", data);
      this.error = null;
      if (res.code != 0) {
        this.error = res.msg;
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
    this.toolboxWorker.showCreateDatabase = this.show;
    this.init();
  },
};
</script>

<style>
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
