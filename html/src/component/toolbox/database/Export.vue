<template>
  <div class="toolbox-database-export">
    <div class="app-scroll-bar pd-10" style="height: calc(100% - 120px)">
      <el-form size="mini" inline>
        <el-form-item label="类型">
          <el-select v-model="formData.exportType" style="width: 100px">
            <el-option
              v-for="(one, index) in exportTypes"
              :key="index"
              :value="one.value"
              :label="one.text"
              :disabled="one.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <template
          v-if="formData.exportType == 'csv' || formData.exportType == 'txt'"
        >
          <el-form-item label="列分割字符">
            <el-input v-model="formData.separator" style="width: 100px">
            </el-input>
          </el-form-item>
          <el-form-item label="换行符转换">
            <el-input v-model="formData.linefeed" style="width: 100px">
            </el-input>
          </el-form-item>
        </template>
        <template v-if="formData.exportType == 'sql'">
          <Pack
            :source="source"
            :toolboxWorker="toolboxWorker"
            :form="formData"
            :change="packChange"
          >
          </Pack>
          <el-checkbox v-model="formData.exportStruct">导出结构体 </el-checkbox>
          <el-checkbox v-model="formData.exportData">导出数据 </el-checkbox>
          <template v-if="formData.exportData">
            <el-checkbox v-model="formData.exportBatchSql"
              >导出批量插入语句
            </el-checkbox>
            <template v-if="formData.exportBatchSql">
              <el-form-item label="批量插入数量">
                <el-input v-model="formData.batchNumber" style="width: 60px">
                </el-input>
              </el-form-item>
            </template>
          </template>
        </template>
        <el-checkbox v-model="formData.errorContinue"> 有错继续</el-checkbox>
      </el-form>
      <el-form size="mini" inline>
        <template v-if="ownerList == null || owners.length == 0">
          <el-form-item label="导出所有库">
            <div class="tm-link color-green mgr-5" @click="initOwners()">
              自定义导出库
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
            <el-form size="mini" inline>
              <el-form-item label="库名称">
                <el-input
                  v-model="owner.sourceName"
                  style="width: 150px"
                  readonly=""
                >
                </el-input>
              </el-form-item>
              <el-form-item :label="ownerExportNameLabel">
                <el-input v-model="owner.targetName" style="width: 150px">
                </el-input>
              </el-form-item>
            </el-form>
          </div>
          <div class="pdl-20" v-loading="owner.tableListLoading">
            <el-form size="mini" inline>
              <template
                v-if="owner.tableList == null || owner.tables.length == 0"
              >
                <el-form-item label="导出所有表">
                  <div
                    class="tm-link color-green mgr-5"
                    @click="initOwnerTables(owner)"
                  >
                    自定义导出表
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
                  <el-form size="mini" inline>
                    <el-form-item label="表名称">
                      <el-input
                        v-model="table.sourceName"
                        style="width: 150px"
                        readonly=""
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item :label="tableExportNameLabel">
                      <el-input v-model="table.targetName" style="width: 150px">
                      </el-input>
                    </el-form-item>
                    <template
                      v-if="
                        table.columnList == null || table.columnList.length == 0
                      "
                    >
                      <el-form-item label="导出所有字段">
                        <div
                          class="tm-link color-green mgr-5"
                          @click="initOwnerTableColumns(owner, table)"
                        >
                          自定义导出字段
                        </div>
                      </el-form-item>
                    </template>
                    <template v-else>
                      <div
                        class="tm-link color-green mgr-5"
                        @click="addExportColumn(table, {})"
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
                    <el-table-column label="字段">
                      <template slot-scope="scope">
                        <div class="">
                          <el-input
                            v-model="scope.row.sourceName"
                            type="text"
                          />
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="导出名称（列名，字段名）">
                      <template slot-scope="scope">
                        <div class="">
                          <el-input
                            v-model="scope.row.targetName"
                            type="text"
                          />
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      label="导出固定值（函数脚本，默认为查询出的值）"
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
                          @click="upExportColumn(table, scope.row)"
                        >
                          上移
                        </div>
                        <div
                          class="tm-link color-grey mglr-5"
                          @click="downExportColumn(table, scope.row)"
                        >
                          下移
                        </div>
                        <div
                          class="tm-link color-grey mglr-5"
                          @click="addExportColumn(table, {}, scope.row)"
                        >
                          插入
                        </div>
                        <div
                          class="tm-link color-red mglr-5"
                          @click="removeExportColumn(table, scope.row)"
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
              <template
                v-if="task.extend && tool.isNotEmpty(task.extend.downloadPath)"
              >
                <div @click="toDownload()" class="color-green tm-link mgr-10">
                  下载
                </div>
              </template>
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
    "actived",
    "extend",
    "ownerName",
    "tableName",
    "columnTypeInfoList",
    "indexTypeInfoList",
  ],
  data() {
    return {
      ready: false,
      exportTypes: [
        { text: "SQL", value: "sql" },
        { text: "Excel", value: "excel" },
        { text: "CSV", value: "csv" },
        { text: "Txt", value: "txt" },
      ],
      formData: {
        exportType: "sql",
        exportStruct: true,
        exportData: true,
        exportBatchSql: true,
        errorContinue: true,
        batchNumber: 200,

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
      exportColumnList: null,
      ownerExportNameLabel: "导出后库名称",
      tableExportNameLabel: "导出后表名称",
      taskId: null,
      task: null,
    };
  },
  computed: {},
  watch: {
    "formData.exportType"() {
      if (this.formData.exportType == "sql") {
        this.exportData = true;
        this.exportStruct = true;
      } else {
        this.exportData = true;
        this.exportStruct = false;
      }

      if (this.formData.exportType == "txt") {
        this.formData.separator = "|:-:|";
      } else if (this.formData.exportType == "csv") {
        this.formData.separator = ",";
      }
    },
  },
  methods: {
    onFocus() {
      if (this.inited) {
        return;
      }
      this.$nextTick(async () => {
        this.init();
      });
    },
    packChange() {},
    async init() {
      this.inited = true;
      let ownerList = [];
      this.ownersReadonly = false;
      this.tablesReadonly = false;
      if (this.tool.isNotEmpty(this.ownerName)) {
        this.ownersReadonly = true;
        let owner = {
          ownerName: this.ownerName,
        };
        this.initOwnerData(owner);
        if (this.tool.isNotEmpty(this.tableName)) {
          this.tablesReadonly = true;
          let table = {
            tableName: this.tableName,
          };
          this.initTableData(table);
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
    initOwnerData(owner) {
      owner.sourceName = owner.ownerName;
      owner.targetName = owner.ownerName;
      owner.tableListLoading = false;
      owner.tableList = null;
      owner.tables = [];
    },
    initTableData(table) {
      table.sourceName = table.tableName;
      table.targetName = table.tableName;
      table.columnListLoading = false;
      table.columnList = null;
      table.openColumnList = false;
    },
    initColumnData(column) {
      column.sourceName = column.columnName;
      column.targetName = column.columnName;
      column.value = null;
    },
    async initOwners() {
      if (this.ownerList == null) {
        let ownerList = await this.toolboxWorker.loadOwners();
        ownerList.forEach((owner) => {
          this.initOwnerData(owner);
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
          this.initTableData(table);
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
          this.initColumnData(column);
        });
        table.columnList = columnList;
        table.columnListLoading = false;
        table.openColumnList = true;
      }
    },

    upExportColumn(table, exportColumn) {
      this.tool.up(table, "columnList", exportColumn);
    },
    downExportColumn(table, exportColumn) {
      this.tool.down(table, "columnList", exportColumn);
    },
    addExportColumn(table, exportColumn, after) {
      exportColumn = exportColumn || {};
      exportColumn.sourceName = exportColumn.sourceName || "";
      exportColumn.targetName = exportColumn.targetName || "";
      exportColumn.value = exportColumn.value || "";

      let appendIndex = table.columnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = table.columnList.length;
      } else {
        appendIndex++;
      }
      table.columnList.splice(appendIndex, 0, exportColumn);
    },
    removeExportColumn(table, exportColumn) {
      let findIndex = table.columnList.indexOf(exportColumn);
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
      let param = this.toolboxWorker.getWorkParam(
        Object.assign({}, this.formData)
      );
      this.toolboxWorker.formatParam(param);

      param.batchNumber = Number(param.batchNumber);
      param.owners = [];
      this.owners.forEach((owner) => {
        let exportOwner = {
          sourceName: owner.sourceName,
          targetName: owner.targetName,
          tables: [],
        };
        owner.tables.forEach((table) => {
          let exportTable = {
            sourceName: table.sourceName,
            targetName: table.targetName,
            columns: [],
          };
          if (table.columnList) {
            table.columnList.forEach((column) => {
              let exportColumn = {
                sourceName: column.sourceName,
                targetName: column.targetName,
                value: column.value,
              };
              exportTable.columns.push(exportColumn);
            });
          }
          exportOwner.tables.push(exportTable);
        });
        param.owners.push(exportOwner);
      });

      let res = await this.server.database.export(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res.data;
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
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.taskId,
      });
      let res = await this.server.database.taskStatus(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      this.task = res.data;
      setTimeout(this.loadStatus, 100);
    },
    async stopTask() {
      if (this.task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.task.taskId,
      });
      await this.server.database.taskStop(param);
    },
    async taskClean() {
      if (this.task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: this.task.taskId,
      });
      this.task = null;
      this.taskId = null;
      await this.server.database.taskClean(param);
    },
    toDownload() {
      if (this.task == null) {
        this.tool.error("任务数据丢失");
        return;
      }
      let url =
        this.source.api +
        "api/database/exportDownload?taskId=" +
        encodeURIComponent(this.task.taskId) +
        "&jwt=" +
        encodeURIComponent(this.tool.getJWT());
      window.location.href = url;
    },
  },
  created() {},
  mounted() {
    if (this.actived) {
      this.init();
    }
  },
  beforeDestroy() {
    this.isDestroyed = true;
    this.taskClean();
  },
};
</script>

<style>
.toolbox-database-export {
  width: 100%;
  height: 100%;
}
</style>
