<template>
  <div class="toolbox-database-sql">
    <tm-layout height="100%">
      <tm-layout height="50px" class="" style="overflow: hidden">
        <el-form
          class="pdt-10 pdl-10"
          ref="form"
          :model="form"
          size="mini"
          inline
        >
          <el-form-item label="数据库">
            <el-select v-model="form.database" style="width: 150px">
              <el-option
                v-for="(one, index) in databases"
                :key="index"
                :value="one.name"
                :label="one.name"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="开启事务">
            <el-switch v-model="form.openTransaction"> </el-switch>
          </el-form-item>
          <el-form-item label="忽略异常继续">
            <el-switch v-model="form.errorContinue"> </el-switch>
          </el-form-item>

          <div class="tm-btn tm-btn-sm bg-green ft-13" @click="toExecuteSql">
            执行
          </div>
          <div
            class="tm-btn tm-btn-sm bg-green ft-13"
            @click="toExecuteSelectSql"
            @mousedown="toSelectSql"
          >
            执行选中
          </div>
        </el-form>
      </tm-layout>
      <tm-layout height="300px" class="" style="overflow: hidden">
        <textarea ref="sqlTextarea" v-model="executeSQL"> </textarea>
      </tm-layout>
      <tm-layout-bar bottom></tm-layout-bar>
      <tm-layout height="auto">
        <TabEditor
          ref="TabEditor"
          :source="source"
          :onRemoveTab="onRemoveTab"
          :onActiveTab="onActiveTab"
          class="sql-execute-tabs"
        >
          <template v-slot:body="{ tab }">
            <template v-if="tab.isExecuteList != null">
              <div class="sql-execute-list">
                <template v-for="(one, index) in executeList">
                  <div :key="index" class="sql-execute-one mgb-10">
                    <div class="color-grey">
                      <span class="pdr-5">{{ one.startTime }}</span>
                      执行
                    </div>
                    <div>
                      <span class="">{{ one.sql }}</span>
                    </div>
                    <template v-if="one.error">
                      <div class="color-orange">
                        执行异常
                        <span class="color-orange pdlr-5">{{ one.error }}</span>
                      </div>
                    </template>
                    <template v-else>
                      <div>
                        <span class="color-green pdr-5"> 执行成功 </span>
                        <template
                          v-if="one.rowsAffected == 0 || one.rowsAffected > 0"
                        >
                          <span class="">
                            受影响行数:
                            <span class="color-green pdlr-5">
                              {{ one.rowsAffected }}
                            </span>
                          </span>
                        </template>
                        <template v-if="one.dataList != null">
                          <span class="">
                            查询行数:
                            <span class="color-green pdlr-5">
                              {{ one.dataList.length }}
                            </span>
                          </span>
                        </template>
                        <span>
                          耗时:
                          <span class="pdlr-5">{{ one.useTime }}毫秒</span>
                        </span>
                      </div>
                    </template>
                  </div>
                </template>
              </div>
            </template>
            <template v-else-if="tab.isSelect">
              <div class="sql-execute-select">
                <SqlSelectDataList :source="source" :wrap="wrap" :tab="tab">
                </SqlSelectDataList>
              </div>
            </template>
          </template>
        </TabEditor>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
import SqlSelectDataList from "./SqlSelectDataList.vue";

export default {
  components: { SqlSelectDataList },
  props: ["source", "wrap", "extend", "databases", "tab"],
  data() {
    return {
      ready: false,
      executeSQL: null,
      form: {
        database: null,
        openTransaction: true,
        errorContinue: false,
      },
      executeList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    async autoSaveSql() {
      let keyValueMap = {};
      if (this.lastSavedExecuteSQL != this.executeSQL) {
        this.lastSavedExecuteSQL = this.executeSQL;
        keyValueMap.executeSQL = this.executeSQL;
      }
      if (this.lastSavedDatabase != this.form.database) {
        this.lastSavedDatabase = this.form.database;
        keyValueMap.database = this.form.database;
      }
      await this.wrap.updateOpenTabExtend(this.tab.tabId, keyValueMap);
      setTimeout(this.autoSaveSql, 300);
    },
    init() {
      if (this.extend) {
        this.executeSQL = this.extend.executeSQL;
        this.form.database = this.extend.database;
      }
      this.autoSaveSql();
      this.ready = true;
    },
    toSelectSql() {
      this.tool.stopEvent(window.event);
      let startIndex = this.$refs.sqlTextarea.selectionStart || 0;
      let endIndex = this.$refs.sqlTextarea.selectionEnd || 0;
      if (endIndex <= startIndex) {
        return;
      }
      this.$refs.sqlTextarea.setSelectionRange(startIndex, endIndex); //将光标定位在textarea的开头，需要定位到其他位置的请自行修改
      this.$refs.sqlTextarea.focus();
    },
    async toExecuteSelectSql() {
      let startIndex = this.$refs.sqlTextarea.selectionStart || 0;
      let endIndex = this.$refs.sqlTextarea.selectionEnd || 0;
      if (endIndex <= startIndex) {
        this.tool.warn("没有SQL被选中");
        return;
      }
      let sql = this.executeSQL.substring(startIndex, endIndex);

      await this.doExecuteSql(sql);
    },
    async toExecuteSql() {
      await this.doExecuteSql(this.executeSQL);
    },
    async doExecuteSql(executeSQL) {
      let data = Object.assign({}, this.form);

      data.executeSQL = executeSQL;
      let res = await this.wrap.work("executeSQL", data);
      if (res.code != 0) {
        return;
      }
      res.data = res.data || {};
      let task = res.data.task;
      if (task.error) {
        this.tool.error(task.error);
      }
      this.executeList = task.executeList || [];
      this.initExecuteList();
    },
    initExecuteList() {
      this.cleanTab();
      this.addExecuteListTab();
      let selectIndex = 0;
      this.executeList.forEach((one) => {
        if (one.isSelect && one.error == null) {
          one.selectIndex = selectIndex;
          this.addExecuteSelectTab(one);
          selectIndex++;
        }
      });
    },
    addExecuteListTab() {
      let tab = {};
      tab.key = "执行结果-" + this.tool.getNumber();
      tab.title = "执行结果";
      tab.name = "执行结果";
      tab.isExecuteList = true;
      this.addTab(tab);
      this.doActiveTab(tab);
    },
    addExecuteSelectTab(executeData) {
      executeData.dataList = executeData.dataList || [];
      let title = `第${executeData.selectIndex + 1}条查询结果（${
        executeData.dataList.length
      }条记录）`;
      let tab = {};
      tab.key = "查询结果-" + this.tool.getNumber();
      tab.title = title;
      tab.name = title;
      tab.isSelect = true;
      tab.columnList = executeData.columnList || [];
      tab.dataList = executeData.dataList;
      this.addTab(tab);
      // this.doActiveTab(tab);
    },
    getTab(tab) {
      return this.$refs.TabEditor.getTab(tab);
    },
    onRemoveTab(tab) {},
    onActiveTab(tab) {},
    addTab(tab) {
      return this.$refs.TabEditor.addTab(tab);
    },
    cleanTab() {
      return this.$refs.TabEditor.toDeleteAll();
    },
    doActiveTab(tab) {
      return this.$refs.TabEditor.doActiveTab(tab);
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
  border-left-color: transparent;
  border-right-color: transparent;
  border-bottom-color: transparent;
}
.toolbox-database-sql textarea::selection {
  color: #494949;
  background: lightblue;
}
.sql-execute-list {
  font-size: 12px;
  user-select: text;
}
.sql-execute-one {
  padding: 0px 5px;
}
.sql-execute-select {
  width: 100%;
  height: 100%;
}
</style>
