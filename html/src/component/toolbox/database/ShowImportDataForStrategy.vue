<template>
  <el-dialog
    ref="modal"
    :title="'导入（策略）：' + (tableDetail == null ? '' : tableDetail.name)"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="taskKey == null"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="ft-12 mgt--20" v-if="tableDetail != null && datas != null">
      <div class="color-grey ft-12">
        <div>
          <span class="color-orange pdr-10">表达式</span>
          <span>表达式，如：'aa' + 'c'，返回“aac”；1 + 2，返回“3”</span>
        </div>
        <div>
          <span class="color-orange pdr-10">_$index</span>
          <span>索引</span>
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
      <div class="mgt-20" style="user-select: none">
        <el-table
          :data="datas"
          :border="true"
          height="100%"
          style="width: 100%"
          size="mini"
        >
          <el-table-column width="110" label="导入数量">
            <template slot-scope="scope">
              <input v-model="scope.row._$importCount" style="width: 100%" />
            </template>
          </el-table-column>
          <template v-for="(column, index) in tableDetail.columns">
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
                    :placeholder="scope.row[column.name] == null ? 'null' : ''"
                    type="text"
                    style="width: 100%"
                  />
                </div>
              </template>
            </el-table-column>
          </template>
        </el-table>
      </div>
      <div class="mgt-10">
        <div class="ft-12">
          <span class="color-grey">导入状态：</span>
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
                成功： <span>{{ task.importSuccess }}</span>
              </span>
              <span class="color-error pdr-10">
                失败： <span>{{ task.importError }}</span>
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
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "wrap"],
  data() {
    return {
      showDialog: false,
      sql: null,
      database: null,
      tableDetail: null,
      datas: null,
      task: null,
      taskKey: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(database, tableDetail) {
      this.database = database;
      this.tableDetail = tableDetail;
      this.task = null;
      this.taskKey = null;
      this.datas = [];
      await this.addData();
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
      this.taskKey = null;
      this.task = null;
    },
    async toImport() {
      this.task = null;
      this.taskKey = null;
      let res = await this.doImport();
      this.taskKey = res.taskKey;
      this.loadStatus();
    },
    async doImport() {
      this.datas.forEach((one) => {
        one._$importCount = Number(one._$importCount);
      });
      let param = {
        database: this.database,
        table: this.tableDetail.name,
        columns: this.tableDetail.columns,
        importDatas: this.datas,
      };
      let res = await this.wrap.work("importDataForStrategy", param);
      res.data = res.data || {};
      return res.data;
    },
    async addData() {
      if (this.tableDetail == null) {
        return;
      }
      let data = {};
      data._$importCount = 1;

      let keys = [];
      this.tableDetail.columns.forEach((column) => {
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
      this.datas.push(data);
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
      let res = await this.wrap.work("importDataForStrategyStatus", param);
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
      await this.wrap.work("importDataForStrategyStop", param);
    },
    async cleanTask() {
      if (this.taskKey == null) {
        return;
      }
      let param = {
        taskKey: this.taskKey,
      };
      await this.wrap.work("importDataForStrategyClean", param);
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.wrap.showImportDataForStrategy = this.show;
    this.init();
  },
};
</script>

<style>
</style>
