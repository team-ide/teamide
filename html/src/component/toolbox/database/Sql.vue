<template>
  <div class="toolbox-database-sql">
    <tm-layout height="100%">
      <tm-layout height="50px" class="" style="overflow: hidden">
        <el-form
          class="mgt-10"
          ref="form"
          :model="form"
          label-width="60px"
          size="mini"
          :inline="true"
        >
          <el-form-item label="数据库">
            <el-select v-model="form.database">
              <el-option
                v-for="(one, index) in databases"
                :key="index"
                :value="one.name"
              >
                {{ one.name }}
              </el-option>
            </el-select>
          </el-form-item>

          <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toExecuteSql">
            执行
          </div>
        </el-form>
      </tm-layout>
      <tm-layout height="300px" class="" style="overflow: hidden">
        <textarea v-model="executeSQL"> </textarea>
      </tm-layout>
      <tm-layout-bar bottom></tm-layout-bar>
      <tm-layout height="auto">
        <div class="sql-execute-list">
          <template v-for="(one, index) in executeList">
            <div :key="index" class="sql-execute-one">
              <div>SQL:{{ one.sql }}</div>
              <div>{{ one }}</div>
            </div>
          </template>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "wrap", "extend", "databases"],
  data() {
    return {
      ready: false,
      executeSQL: null,
      form: {
        database: null,
      },
      executeList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.ready = true;
    },
    async toExecuteSql() {
      let task = await this.doExecuteSql();
      this.executeList = task.executeList || [];
    },
    async doExecuteSql() {
      let data = Object.assign({}, this.form);

      data.executeSQL = this.executeSQL;
      let res = await this.wrap.work("executeSQL", data);
      if (res.code != 0) {
        return;
      }
      res.data = res.data || {};
      return res.data.task;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-sql {
  width: 100%;
  height: 100%;
}
.toolbox-database-sql {
  width: 100%;
  height: 100%;
}
.toolbox-database-sql textarea {
  width: 100%;
  height: 100%;
  letter-spacing: 1px;
  word-spacing: 5px;
  word-break: break-all;
  font-size: 12px;
  padding: 5px 5px;
  outline: none;
  user-select: none;
  resize: none;
}
</style>
