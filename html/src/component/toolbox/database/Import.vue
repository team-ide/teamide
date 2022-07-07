<template>
  <div class="toolbox-database-import">
    <el-form
      class="pdt-10 pdlr-10"
      size="mini"
      @submit.native.prevent
      label-width="80px"
      inline
    >
      <el-form-item label="导入类型">
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
    </el-form>

    <template v-if="form.importType == 'strategy'">
      <ScriptInfo></ScriptInfo>
      <div class="mglr-10 mgt-10">
        <div class="tm-link color-grey" @click="addStrategyData">添加</div>
      </div>
      <div class="mglr-10 mgt-10 scrollbar" style="height: calc(100% - 380px)">
        <template v-for="(strategyData, index) in strategyDataList">
          <div :key="index" class="mgb-10">
            <div class="ft-12 mgb-5">数据策略[ {{ index + 1 }} ]</div>
            <div>
              <el-form size="mini" @submit.native.prevent label-width="200px">
                <el-form-item label="导入数量" class="mgb-5">
                  <el-input v-model="strategyData.count"> </el-input>
                </el-form-item>
                <el-form-item label="批量保存数量" class="mgb-5">
                  <el-input v-model="strategyData.batchNumber"> </el-input>
                </el-form-item>
                <template v-for="(column, index) in strategyData.columnList">
                  <el-form-item :key="index" :label="column.name" class="mgb-5">
                    <el-input v-model="column.value"> </el-input>
                  </el-form-item>
                </template>
              </el-form>
            </div>
          </div>
        </template>
      </div>
    </template>
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
          <template v-if="task.isStop">
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
            <div @click="stopTask" class="color-red tm-link mgr-10">
              停止执行
            </div>
          </template>
          <div class="mgt-5">
            <span class="color-grey pdr-10">
              预计导入： <span>{{ task.dataCount }}</span>
            </span>
            <span class="color-grey pdr-10">
              已准备数据： <span>{{ task.readyDataCount }}</span>
            </span>
            <span class="color-success pdr-10">
              成功： <span>{{ task.successCount }}</span>
            </span>
            <span class="color-error pdr-10">
              异常： <span>{{ task.errorCount }}</span>
            </span>
          </div>
          <template v-if="task.error != null">
            <div class="mgt-5 color-error pdr-10">
              异常： <span>{{ task.error }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>
    <div class="mglr-10 mgt-20" v-if="taskKey == null">
      <div class="tm-btn bg-green" @click="toImport">导入</div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend", "database", "table"],
  data() {
    return {
      ready: false,
      importTypes: [
        { text: "策略函数", value: "strategy" },
        { text: "SQL", value: "sql", disabled: true },
        { text: "Excel", value: "excel", disabled: true },
        { text: "文本", value: "text", disabled: true },
      ],
      form: {
        importType: "strategy",
      },
      strategyDataList: null,
      tableDetail: null,
      taskKey: null,
      task: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      if (this.tool.isNotEmpty(this.table)) {
        this.tableDetail = await this.toolboxWorker.getTableDetail(
          this.database,
          this.table
        );
      }
      this.strategyDataList = [];
      await this.addStrategyData();
      this.ready = true;
    },
    async addStrategyData() {
      if (this.tableDetail == null) {
        return;
      }
      let data = {};
      data.count = 1;
      data.batchNumber = 100;
      data.columnList = [];

      let keys = [];
      this.tableDetail.columnList.forEach((column) => {
        let value = null;
        if (column.primaryKey) {
          keys.push(column.name);
          if (
            column.type == "int" ||
            column.type == "bigint" ||
            column.type == "number"
          ) {
            value = "0 + _$index";
          } else {
            value = "_$uuid()";
          }
        } else if (column.notNull) {
          if (
            column.type == "int" ||
            column.type == "bigint" ||
            column.type == "number"
          ) {
            value = "0";
          } else if (
            column.type == "date" ||
            column.type == "time" ||
            column.type == "datetime"
          ) {
            value = "_$now()";
          } else {
            if (keys.length > 0) {
              value = "'" + column.name + "' + " + keys.join(" + ") + "";
            } else {
              value = "_$randomString(1, 5)";
            }
          }
        }

        data.columnList.push({
          name: column.name,
          value: value,
        });
      });
      this.strategyDataList.push(data);
    },
    async toImport() {
      this.task = null;
      this.taskKey = null;
      let res = await this.doImport();
      this.taskKey = res.taskKey;
      this.loadStatus();
    },
    async doImport() {
      this.strategyDataList.forEach((one) => {
        one.count = Number(one.count);
        one.batchNumber = Number(one.batchNumber);
      });

      let param = Object.assign({}, this.form);
      param.database = this.database;
      param.table = this.table;
      param.columnList = this.tableDetail.columnList;
      param.strategyDataList = this.strategyDataList;

      let res = await this.toolboxWorker.work("import", param);
      res.data = res.data || {};
      return res.data;
    },
    async loadStatus() {
      if (this.taskKey == null) {
        return;
      }
      if (this.task != null && this.task.isEnd) {
        this.taskKey = null;
        this.cleanTask();
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      let res = await this.toolboxWorker.work("importStatus", param);
      res.data = res.data || {};
      this.task = res.data.task;
      setTimeout(this.loadStatus, 100);
    },
    async stopTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.toolboxWorker.work("importStop", param);
    },
    async cleanTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.toolboxWorker.work("importClean", param);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-database-import {
  width: 100%;
  height: 100%;
}
</style>
