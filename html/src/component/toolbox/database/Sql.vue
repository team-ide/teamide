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
                    <div>
                      SQL:
                      <span class="pdlr-5">{{ one.sql }}</span>
                    </div>
                    <template v-if="one.error">
                      <div class="">
                        执行异常:
                        <span class="color-red pdlr-5">{{ one.error }}</span>
                      </div>
                    </template>
                    <template v-else>
                      <div>
                        <span class="color-green pdr-5"> 执行成功 </span>
                        <template
                          v-if="one.rowsAffected == 0 || one.rowsAffected > 0"
                        >
                          <span class=""
                            >受影响行数:
                            <span class="color-green pdlr-5">
                              {{ one.rowsAffected }}
                            </span>
                          </span>
                        </template>
                        <template v-if="one.dataList != null">
                          <span class=""
                            >查询行数:
                            <span class="color-green pdlr-5">
                              {{ one.dataList.length }}
                            </span>
                          </span>
                        </template>
                        <span>
                          开始时间:
                          <span class="pdlr-5">{{ one.startTime }} </span>
                        </span>
                        <span>
                          结束时间:
                          <span class="pdlr-5">{{ one.endTime }}</span>
                        </span>
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
                <div
                  class="
                    toolbox-database-table-data
                    toolbox-database-table-data-table
                  "
                >
                  <el-table
                    :data="tab.dataList"
                    :border="true"
                    height="100%"
                    style="width: 100%"
                    size="mini"
                  >
                    <el-table-column width="70" label="序号">
                      <template slot-scope="scope">
                        <span class="mgl-5">{{ scope.$index + 1 }}</span>
                      </template>
                    </el-table-column>
                    <template v-for="(column, index) in tab.columnList">
                      <el-table-column
                        :key="index"
                        :prop="column.name"
                        :label="column.name"
                        width="120"
                      >
                        <template slot-scope="scope">
                          <div class="">
                            <input
                              v-model="scope.row[column.name]"
                              :placeholder="
                                scope.row[column.name] == null ? 'null' : ''
                              "
                              type="text"
                            />
                          </div>
                        </template>
                      </el-table-column>
                    </template>
                  </el-table>
                </div>
              </div>
            </template>
          </template>
        </TabEditor>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "wrap", "extend", "databases", "tab"],
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
    async autoSaveSql() {
      if (this.lastSavedExecuteSQL == this.executeSQL) {
        setTimeout(this.autoSaveSql, 300);
        return;
      }
      this.lastSavedExecuteSQL = this.executeSQL;
      await this.toSaveSql();
      setTimeout(this.autoSaveSql, 300);
    },
    init() {
      if (this.extend) {
        this.executeSQL = this.extend.executeSQL;
      }
      this.autoSaveSql();
      this.ready = true;
    },
    async toSaveSql() {
      await this.wrap.updateOpenTabExtend(
        this.tab.tabId,
        ["executeSQL"],
        this.executeSQL
      );
    },
    async toExecuteSql() {
      this.toSaveSql();
      let task = await this.doExecuteSql();
      if (task.error) {
        this.tool.error(task.error);
      }
      this.executeList = task.executeList || [];
      this.initExecuteList();
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
      tab.key = "执行结果";
      tab.title = "执行结果";
      tab.name = "执行结果";
      tab.isExecuteList = true;
      this.addTab(tab);
      this.doActiveTab(tab);
    },
    addExecuteSelectTab(executeData) {
      let tab = {};
      tab.key = "查询结果" + executeData.selectIndex;
      tab.title = "查询结果" + executeData.selectIndex;
      tab.name = "查询结果" + executeData.selectIndex;
      tab.isSelect = true;
      tab.dataList = executeData.dataList || [];
      tab.columnList = executeData.columnList || [];
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
.sql-execute-list {
  font-size: 12px;
}
.sql-execute-select {
  width: 100%;
  height: 100%;
}
</style>
