<template>
  <div class="toolbox-database-import">
    <div class="app-scroll-bar pd-10" style="height: calc(100% - 120px)">
      <el-form ref="form" :model="form" size="mini" inline>
        <el-form-item label="类型">
          <el-select v-model="form.importType" style="width: 100px">
            <el-option
              v-for="(one, index) in importTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
              :disabled="one.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <template v-if="form.importType == 'csv' || form.importType == 'txt'">
          <el-form-item label="列分割字符">
            <el-input v-model="form.separator" style="width: 100px"> </el-input>
          </el-form-item>
          <el-form-item label="换行符转换">
            <el-input v-model="form.linefeed" style="width: 100px"> </el-input>
          </el-form-item>
        </template>
        <template v-if="form.importType == 'sql'">
          <Pack
            :source="source"
            :toolboxWorker="toolboxWorker"
            :form="form"
            :change="packChange"
          >
          </Pack>
        </template>

        <el-checkbox v-model="form.errorContinue"> 有错继续</el-checkbox>
      </el-form>
      <el-form ref="form" :model="form" size="mini" inline>
        <template v-if="ownerList == null">
          <el-form-item label="库未加载">
            <div class="tm-link color-green mgr-5" @click="initOwners()">
              加载库
            </div>
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item label="选择库">
            <el-checkbox-group v-model="owners">
              <el-checkbox
                v-for="(owner, index) in ownerList"
                :label="owner"
                :key="index"
                :disabled="ownersReadonly"
              >
                {{ owner.ownerName }}
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>
      </el-form>

      <template v-for="(owner, ownerIndex) in owners">
        <div :key="ownerIndex">
          <div>
            <el-form :model="form" size="mini" inline>
              <el-form-item label="名称">
                <el-input v-model="owner.name" style="width: 150px" readonly="">
                </el-input>
              </el-form-item>
              <el-form-item label="文件路径">
                <el-input v-model="owner.path" style="width: 150px"> </el-input>
              </el-form-item>
            </el-form>
          </div>
          <div class="pdl-20" v-loading="owner.tableListLoading">
            <el-form :model="form" size="mini" inline>
              <template v-if="owner.tableList == null">
                <el-form-item label="表未加载">
                  <div
                    class="tm-link color-green mgr-5"
                    @click="initOwnerTables(owner)"
                  >
                    加载表
                  </div>
                </el-form-item>
              </template>
              <template v-else>
                <el-form-item label="选择表">
                  <el-checkbox-group v-model="owner.tables">
                    <el-checkbox
                      v-for="(table, index) in owner.tableList"
                      :label="table"
                      :key="index"
                      :disabled="tablesReadonly"
                    >
                      {{ table.tableName }}
                    </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              </template>
            </el-form>
            <template v-for="(table, tableIndex) in owner.tables">
              <div :key="tableIndex">
                <div v-loading="table.columnListLoading">
                  <el-form ref="form" :model="form" size="mini" inline>
                    <el-form-item label="表名称">
                      <el-input
                        v-model="table.name"
                        style="width: 150px"
                        readonly=""
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item label="文件路径">
                      <el-input v-model="table.path" style="width: 150px">
                      </el-input>
                    </el-form-item>
                    <template v-if="table.columnList == null">
                      <el-form-item label="字段未加载">
                        <div
                          class="tm-link color-green mgr-5"
                          @click="initOwnerTableColumns(owner, table)"
                        >
                          加载字段
                        </div>
                      </el-form-item>
                    </template>
                    <template v-else>
                      <div
                        class="tm-link color-green mgr-5"
                        @click="addImportColumn(table, {})"
                      >
                        添加字段
                      </div>
                      <div
                        class="tm-link color-orange mgr-5"
                        v-if="table.openColumnList"
                        @click="table.openColumnList = false"
                      >
                        收起字段
                      </div>
                      <div
                        class="tm-link color-orange mgr-5"
                        v-if="!table.openColumnList"
                        @click="table.openColumnList = true"
                      >
                        展开字段
                      </div>
                    </template>
                  </el-form>
                </div>
                <div v-if="table.openColumnList">
                  <el-table
                    :data="table.columnList"
                    border
                    style="width: 100%"
                    size="mini"
                  >
                    <el-table-column label="字段名称">
                      <template slot-scope="scope">
                        <div class="">
                          <el-input v-model="scope.row.name" type="text" />
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      label="导入固定值（函数脚本，默认为查询出的值）"
                    >
                      <template slot-scope="scope">
                        <div class="">
                          <el-input v-model="scope.row.value" type="text" />
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="操作" width="200px">
                      <template slot-scope="scope">
                        <div
                          class="tm-link color-grey mglr-5"
                          @click="upImportColumn(table, scope.row)"
                        >
                          上移
                        </div>
                        <div
                          class="tm-link color-grey mglr-5"
                          @click="downImportColumn(table, scope.row)"
                        >
                          下移
                        </div>
                        <div
                          class="tm-link color-grey mglr-5"
                          @click="addImportColumn(table, {}, scope.row)"
                        >
                          插入
                        </div>
                        <div
                          class="tm-link color-red mglr-5"
                          @click="removeImportColumn(table, scope.row)"
                        >
                          删除
                        </div>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </template>
          </div>
        </div>
      </template>
    </div>
    <div class="mglr-10 mgt-10" style="user-select: text">
      <div class="ft-12">
        <span class="color-grey">任务状态：</span>
        <template v-if="task == null">
          <span class="color-orange pdr-10">暂未开始</span>
        </template>
        <template v-else>
          <template v-if="!task.isEnd">
            <span class="color-orange pdr-10"> 处理中 </span>
          </template>
          <template v-else-if="task.isStop">
            <span class="color-red pdr-10"> 已停止 </span>
          </template>
          <template v-else>
            <span class="color-green pdr-10"> 执行完成 </span>
          </template>
          <span class="color-grey pdr-10">
            开始：
            <span>
              {{
                tool.formatDate(new Date(task.startTime), "yyyy-MM-dd hh:mm:ss")
              }}
            </span>
          </span>
          <template v-if="task.isEnd">
            <span class="color-grey pdr-10">
              结束：
              <span>
                {{
                  tool.formatDate(new Date(task.endTime), "yyyy-MM-dd hh:mm:ss")
                }}
              </span>
            </span>
            <span class="color-grey pdr-10">
              耗时： <span>{{ task.useTime }} 毫秒</span>
            </span>
          </template>
          <template v-if="!task.isEnd">
            <div @click="stopTask()" class="color-red tm-link mgr-10">
              停止执行
            </div>
          </template>
          <div class="mgt-5">
            <span class="color-grey pdr-10">
              库 总数/成功/失败：
              <span>
                {{ task.ownerCount }}
                /
                <span class="color-green">
                  {{ task.ownerSuccessCount }}
                </span>
                /
                <span class="color-red">
                  {{ task.ownerErrorCount }}
                </span>
              </span>
            </span>
            <span class="color-grey pdr-10">
              表 总数/成功/失败：
              <span>
                {{ task.tableCount }}
                /
                <span class="color-green">
                  {{ task.tableSuccessCount }}
                </span>
                /
                <span class="color-red">
                  {{ task.tableErrorCount }}
                </span>
              </span>
            </span>
            <span class="color-grey pdr-10">
              数据 总数/成功/失败：
              <span>
                {{ task.dataCount }}
                /
                <span class="color-green">
                  {{ task.dataSuccessCount }}
                </span>
                /
                <span class="color-red">
                  {{ task.dataErrorCount }}
                </span>
              </span>
            </span>
            <template v-if="task.isEnd">
              <div @click="taskClean()" class="color-orange tm-link mgr-10">
                删除
              </div>
            </template>
          </div>
          <template v-if="tool.isNotEmpty(task.error)">
            <div class="mgt-5 color-error pdr-10">
              异常： <span>{{ task.error }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>
    <div class="pdlr-10 mgt-10" v-if="taskId == null">
      <div class="tm-btn bg-green" @click="toDo">开始</div>
    </div>
  </div>
</template>


<script>
import Pack from "./Pack";

export default {
  components: { Pack },
  props: [
    "source",
    "toolboxWorker",
    "extend",
    "ownerName",
    "tableName",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      ready: false,
      importTypes: [
        { text: "SQL", value: "sql" },
        { text: "Excel", value: "excel" },
        { text: "CSV", value: "csv" },
        { text: "Txt", value: "txt" },
      ],
      form: {
        importType: "sql",
        errorContinue: true,

        separator: "|:-:|",
        linefeed: "|:-n-:|",
        targetDatabaseType: "",
        appendOwnerName: true,
        ownerNamePackChar: "",
        tableNamePackChar: "",
        columnNamePackChar: "",
        sqlValuePackChar: "",
      },
      ownersReadonly: false,
      tablesReadonly: false,
      owners: [],
      ownerList: [],
      importColumnList: null,
      ownerImportNameLabel: "导入后库名称",
      tableImportNameLabel: "导入后表名称",
      taskId: null,
      task: null,
    };
  },
  computed: {},
  watch: {
    "form.importType"() {
      if (this.form.importType == "sql") {
      } else if (this.form.importType == "txt") {
        this.form.separator = "|:-:|";
      } else if (this.form.importType == "csv") {
        this.form.separator = ",";
      }
    },
  },
  methods: {
    packChange() {},
    async init() {
      let ownerList = [];
      this.ownersReadonly = false;
      this.tablesReadonly = false;
      if (this.tool.isNotEmpty(this.ownerName)) {
        this.ownersReadonly = true;
        let owner = {
          ownerName: this.ownerName,
          name: this.ownerName,
          path: null,
          tableListLoading: false,
          tableList: null,
          tables: [],
        };

        if (this.tool.isNotEmpty(this.tableName)) {
          this.tablesReadonly = true;
          let table = {
            tableName: this.tableName,
            name: this.tableName,
            path: null,
            columnList: null,
            columnListLoading: false,
            openColumnList: false,
          };
          owner.tableList = [];
          owner.tableList.push(table);
          owner.tables.push(table);
        }
        ownerList.push(owner);
        this.owners.push(owner);
        this.ownerList = ownerList;
      } else {
        this.ownerList = null;
      }

      this.ready = true;
    },
    async initOwners() {
      if (this.ownerList == null) {
        let ownerList = await this.toolboxWorker.loadOwners();
        ownerList.forEach((owner) => {
          owner.tableList = null;
          owner.tables = [];
          owner.tableListLoading = false;
          owner.name = owner.ownerName;
          owner.path = "";
        });
        this.ownerList = ownerList;
        ownerList.forEach(async (owner) => {
          this.owners.push(owner);
        });
      }
    },
    async initOwnerTables(owner) {
      if (owner.tableList == null) {
        owner.tableListLoading = true;
        let tableList = await this.toolboxWorker.loadTables(owner.ownerName);
        tableList.forEach((table) => {
          table.name = table.tableName;
          table.path = "";
          table.openColumnList = false;
          table.columnList = null;
          table.columnListLoading = false;
        });
        owner.tableList = tableList;
        owner.tableListLoading = false;
      }
      owner.tableList.forEach(async (table) => {
        owner.tables.push(table);
      });
    },
    async initOwnerTableColumns(owner, table) {
      if (table.columnList == null) {
        table.columnListLoading = true;
        let detail = await this.toolboxWorker.getTableDetail(
          owner.ownerName,
          table.tableName
        );
        let columnList = [];
        if (detail) {
          columnList = detail.columnList || [];
        }
        columnList.forEach((column) => {
          column.name = column.columnName;
          column.path = "";
          column.value = null;
        });
        table.columnList = columnList;
        table.columnListLoading = false;
        table.openColumnList = true;
      }
    },

    upImportColumn(table, importColumn) {
      this.tool.up(table, "columnList", importColumn);
    },
    downImportColumn(table, importColumn) {
      this.tool.down(table, "columnList", importColumn);
    },
    addImportColumn(table, importColumn, after) {
      importColumn = importColumn || {};
      importColumn.columnName = importColumn.columnName || "";
      importColumn.importName = importColumn.importName || "";
      importColumn.value = importColumn.value || "";

      let appendIndex = table.columnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = table.columnList.length;
      } else {
        appendIndex++;
      }
      table.columnList.splice(appendIndex, 0, importColumn);
    },
    removeImportColumn(table, importColumn) {
      let findIndex = table.columnList.indexOf(importColumn);
      if (findIndex >= 0) {
        table.columnList.splice(findIndex, 1);
      }
    },
    async toDo() {
      if (this.task != null) {
        this.taskClean();
      }
      this.task = null;
      this.taskId = null;
      let res = await this.start();
      if (res) {
        this.taskId = res.taskId;
        this.loadStatus();
      }
    },
    async start() {
      let param = Object.assign({ owners: [] }, this.form);
      this.toolboxWorker.formatParam(param);

      param.owners = [];
      this.owners.forEach((owner) => {
        let importOwner = {
          name: owner.name,
          path: owner.path,
          tables: [],
        };
        owner.tables.forEach((table) => {
          let importTable = {
            name: table.name,
            path: table.path,
            columns: [],
          };
          if (table.columnList) {
            table.columnList.forEach((column) => {
              let importColumn = {
                name: column.name,
                path: column.path,
                value: column.value,
              };
              importTable.columns.push(importColumn);
            });
          }
          importOwner.tables.push(importTable);
        });
        param.owners.push(importOwner);
      });

      let res = await this.toolboxWorker.work("import", param);
      res.data = res.data || {};
      return res.data.task;
    },
    async loadStatus() {
      if (this.taskId == null) {
        return;
      }
      if (this.task != null && this.task.isEnd) {
        this.taskId = null;
        return;
      }
      if (this.isDestroyed) {
        return;
      }
      let param = {
        taskId: this.taskId,
      };
      let res = await this.toolboxWorker.work("taskStatus", param);
      res.data = res.data || {};
      this.task = res.data.task;
      setTimeout(this.loadStatus, 100);
    },
    async stopTask() {
      if (this.task == null) {
        return;
      }
      let param = {
        taskId: this.task.taskId,
      };
      await this.toolboxWorker.work("taskStop", param);
    },
    async taskClean() {
      if (this.task == null) {
        return;
      }
      let param = {
        taskId: this.task.taskId,
      };
      this.task = null;
      this.taskId = null;
      await this.toolboxWorker.work("taskClean", param);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
    this.taskClean();
  },
};
</script>

<style>
.toolbox-database-import {
  width: 100%;
  height: 100%;
}
</style>
