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
            <el-select v-model="form.ownerName" style="width: 150px">
              <el-option
                v-for="(one, index) in owners"
                :key="index"
                :value="one.ownerName"
                :label="one.ownerName"
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
          >
            执行选中
          </div>
        </el-form>
      </tm-layout>
      <tm-layout height="300px" class="" style="overflow: hidden">
        <Editor
          ref="Editor"
          :source="source"
          :value="executeSQL"
          language="sql"
          :change="executeSQLChange"
          :onContextMenu="editorContextmenu"
        ></Editor>
      </tm-layout>
      <tm-layout-bar bottom></tm-layout-bar>
      <tm-layout height="auto">
        <div class="default-tabs-container">
          <WorkspaceTabs :source="source" :itemsWorker="sqlItemsWorker">
          </WorkspaceTabs>
        </div>
        <div class="default-spans-container">
          <WorkspaceSpans :source="source" :itemsWorker="sqlItemsWorker">
            <template v-slot:span="{ item }">
              <template v-if="item.isExecuteList != null">
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
                          <span class="color-orange pdlr-5">{{
                            one.error
                          }}</span>
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
              <template v-else-if="item.isSelect">
                <div class="sql-execute-select">
                  <SqlSelectDataList
                    :source="source"
                    :toolboxWorker="toolboxWorker"
                    :item="item"
                  >
                  </SqlSelectDataList>
                </div>
              </template>
            </template>
          </WorkspaceSpans>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
import SqlSelectDataList from "./SqlSelectDataList.vue";

export default {
  components: { SqlSelectDataList },
  props: ["source", "toolboxWorker", "extend", "owners", "tabId"],
  data() {
    let sqlItemsWorker = this.tool.newItemsWorker({
      async onRemoveItem(item) {},
      async onActiveItem(item) {},
    });

    return {
      ready: false,
      executeSQL: null,
      sqlItemsWorker: sqlItemsWorker,
      form: {
        ownerName: null,
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
      if (this.isDestroyed) {
        return;
      }
      let keyValueMap = {};
      if (this.lastSavedExecuteSQL != this.executeSQL) {
        this.lastSavedExecuteSQL = this.executeSQL;
        keyValueMap.executeSQL = this.executeSQL;
      }
      if (this.lastSavedDatabase != this.form.ownerName) {
        this.lastSavedDatabase = this.form.ownerName;
        keyValueMap.ownerName = this.form.ownerName;
      }
      await this.toolboxWorker.updateOpenTabExtend(this.tabId, keyValueMap);
      setTimeout(this.autoSaveSql, 300);
    },
    init() {
      if (this.extend) {
        this.executeSQL = this.extend.executeSQL;
        this.form.ownerName = this.extend.ownerName;
      }
      this.autoSaveSql();
      this.ready = true;
      this.$refs.Editor.setValue(this.executeSQL);
    },
    executeSQLChange(value) {
      this.executeSQL = value;
    },
    editorContextmenu() {
      let menus = [];
      let sql = this.$refs.Editor.getSelection();
      menus.push({
        text: "执行选中",
        disabled: this.tool.isEmpty(sql),
        onClick: () => {
          this.toExecuteSelectSql();
        },
      });
      menus.push({
        text: "执行全部",
        onClick: () => {
          this.toExecuteSql();
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    async toExecuteSelectSql() {
      let sql = this.$refs.Editor.getSelection();
      if (this.tool.isEmpty(sql)) {
        this.tool.warn("没有SQL被选中");
        return;
      }

      await this.doExecuteSql(sql);
    },
    async toExecuteSql() {
      await this.doExecuteSql(this.executeSQL);
    },
    async doExecuteSql(executeSQL) {
      let data = Object.assign({}, this.form);

      data.executeSQL = executeSQL;
      let res = await this.toolboxWorker.work("executeSQL", data);
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
      this.doActiveTab(tab);
    },
    getTab(tab) {
      return this.sqlItemsWorker.getItem(tab);
    },
    addTab(tab) {
      this.sqlItemsWorker.addItem(tab);
    },
    cleanTab() {
      return this.sqlItemsWorker.toRemoveAll();
    },
    doActiveTab(tab) {
      return this.sqlItemsWorker.toActiveItem(tab);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
  },
};
</script>

<style>
.toolbox-database-sql {
  width: 100%;
  height: 100%;
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
