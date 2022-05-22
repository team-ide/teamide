<template>
  <div class="toolbox-database-import">
    <div class="pd-10">
      <el-form ref="form" :model="form" size="mini" inline>
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
        <div class="color-grey ft-12" style="user-select: text">
          <div>
            <span class="color-orange pdr-10">表达式</span>
            <span>表达式，如：'aa' + 'c'，返回“aac”；1 + 2，返回“3”</span>
          </div>
          <div>
            <span class="color-orange pdr-10">_$index</span>
            <span>索引，每个策略数据从0开始，最大为当前策略数据数量-1</span>
          </div>
          <div>
            <span class="color-orange pdr-10">_$now()</span>
            <span>当前时间对象</span>
          </div>
          <div>
            <span class="color-orange pdr-10">_$nowTime()</span>
            <span>当前时间戳</span>
          </div>
          <div>
            <span class="color-orange pdr-10">_$uuid()</span>
            <span>生成UUID</span>
          </div>
          <div>
            <span class="color-orange pdr-10">
              _$randomString(minLength, maxLength)
            </span>
            <span>随机字符串</span>
          </div>
          <div>
            <span class="color-orange pdr-10">_$randomInt(min, max)</span>
            <span>随机数字</span>
          </div>
          <div>
            <span class="color-orange pdr-10">
              _$randomUserName(minLength, maxLength)
            </span>
            <span>随机用户姓名</span>
          </div>
          <div>
            <span class="color-orange pdr-10"> _$toPinYin(str) </span>
            <span>转为拼音</span>
          </div>
        </div>
        <div
          v-if="tableDetail != null"
          class="
            mgt-20
            toolbox-database-table-data toolbox-database-table-data-table
          "
        >
          <div class="mgb-10">
            <div class="tm-link color-grey" @click="addStrategyData">添加</div>
          </div>
          <el-table
            :data="strategyDataList"
            borde
            style="width: 100%"
            size="mini"
          >
            <el-table-column width="110" label="导入数量">
              <template slot-scope="scope">
                <input v-model="scope.row._$importCount" style="width: 100%" />
              </template>
            </el-table-column>
            <template v-for="(column, index) in tableDetail.columnList">
              <el-table-column
                :key="index"
                :prop="column.name"
                :label="column.name"
                width="150"
              >
                <template slot-scope="scope">
                  <div class="">
                    <input
                      v-model="scope.row[column.name]"
                      :placeholder="
                        scope.row[column.name] == null ? 'NULL' : ''
                      "
                      type="text"
                      style="width: 100%"
                    />
                  </div>
                </template>
              </el-table-column>
            </template>
          </el-table>
        </div>
      </template>
      <div class="mgt-10" style="user-select: text">
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
                  tool.formatDate(
                    new Date(task.startTime),
                    "yyyy-MM-dd hh:mm:ss"
                  )
                }}
              </span>
            </span>
            <template v-if="task.isEnd">
              <span class="color-grey pdr-10">
                结束：
                <span>
                  {{
                    tool.formatDate(
                      new Date(task.endTime),
                      "yyyy-MM-dd hh:mm:ss"
                    )
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
      <div class="mgt-20" v-if="taskKey == null">
        <div class="tm-btn bg-green" @click="toImport">导入</div>
      </div>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap", "extend", "database", "table"],
  data() {
    return {
      ready: false,
      importTypes: [
        { text: "策略函数", value: "strategy" },
        { text: "SQL", value: "sql", disabled: true },
        { text: "Excel", value: "excel", disabled: true },
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
        this.tableDetail = await this.wrap.getTableDetail(
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
      data._$importCount = 1;

      let keys = [];
      this.tableDetail.columnList.forEach((column) => {
        if (column.primaryKey) {
          keys.push(column.name);
          if (
            column.type == "int" ||
            column.type == "bigint" ||
            column.type == "number"
          ) {
            data[column.name] = "0 + _$index";
          } else {
            data[column.name] = "_$uuid()";
          }
        } else if (column.notNull) {
          if (
            column.type == "int" ||
            column.type == "bigint" ||
            column.type == "number"
          ) {
            data[column.name] = 0;
          } else if (
            column.type == "date" ||
            column.type == "time" ||
            column.type == "datetime"
          ) {
            data[column.name] = "_$now()";
          } else {
            if (keys.length > 0) {
              data[column.name] =
                "'" + column.name + "' + " + keys.join(" + ") + "";
            } else {
              data[column.name] = "_$randomString(1, 5)";
            }
          }
        } else {
          data[column.name] = null;
        }
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
        one._$importCount = Number(one._$importCount);
      });

      let param = Object.assign({}, this.form);
      param.database = this.database;
      param.table = this.table;
      param.columnList = this.tableDetail.columnList;
      param.strategyDataList = this.strategyDataList;

      let res = await this.wrap.work("import", param);
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
      let res = await this.wrap.work("importStatus", param);
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
      await this.wrap.work("importStop", param);
    },
    async cleanTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.wrap.work("importClean", param);
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
