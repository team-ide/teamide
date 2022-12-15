<template>
  <div class="toolbox-elasticsearch-import">
    <tm-layout height="100%">
      <tm-layout height="140px">
        <div class="pd-10">
          <el-form size="mini" inline>
            <el-form-item label="类型">
              <el-select v-model="formData.importType" style="width: 100px">
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
            <el-form-item label="索引">
              <el-input v-model="formData.indexName" style="width: 250px">
              </el-input>
            </el-form-item>
            <el-form-item label="导入数量">
              <el-input v-model="formData.dataNumber" style="width: 100px">
              </el-input>
            </el-form-item>
            <el-form-item label="批量导入">
              <el-input v-model="formData.batchNumber" style="width: 100px">
              </el-input>
            </el-form-item>
            <el-form-item label="线程数量">
              <el-input v-model="formData.threadNumber" style="width: 100px">
              </el-input>
            </el-form-item>
            <el-form-item label="ID">
              <el-input v-model="formData.id" style="width: 300px"> </el-input>
            </el-form-item>
            <el-checkbox v-model="formData.errorContinue">
              有错继续
            </el-checkbox>
          </el-form>
        </div>
        <div class="pdlr-10 mgt--10">
          <div class="tm-link color-green mgr-5" @click="addImportColumn({})">
            添加字段
          </div>
          <div class="tm-link color-grey mgr-5" @click="showColumnListCode()">
            字段模板
          </div>
        </div>
      </tm-layout>
      <tm-layout height="400px">
        <div class="pdlr-10" style="height: 100%">
          <el-table
            :data="columnList"
            border
            size="mini"
            style="width: 100%"
            height="100%"
          >
            <el-table-column label="字段名称">
              <template slot-scope="scope">
                <div class="">
                  <el-input v-model="scope.row.name" type="text" />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="导入固定值（函数脚本，默认为查询出的值）">
              <template slot-scope="scope">
                <div class="">
                  <el-input v-model="scope.row.value" type="text" />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="重复使用次数">
              <template slot-scope="scope">
                <div class="">
                  <el-input v-model="scope.row.reuseNumber" type="text" />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200px">
              <template slot-scope="scope">
                <div
                  class="tm-link color-grey mglr-5"
                  @click="upImportColumn(scope.row)"
                >
                  上移
                </div>
                <div
                  class="tm-link color-grey mglr-5"
                  @click="downImportColumn(scope.row)"
                >
                  下移
                </div>
                <div
                  class="tm-link color-grey mglr-5"
                  @click="addImportColumn({}, scope.row)"
                >
                  插入
                </div>
                <div
                  class="tm-link color-red mglr-5"
                  @click="removeImportColumn(scope.row)"
                >
                  删除
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </tm-layout>
      <tm-layout height="auto" class="app-scroll-bar">
        <template v-for="(task, index) in taskList">
          <div class="mglr-10 mgt-10" :key="index" style="user-select: text">
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
                </template>
                <template v-if="!task.isEnd">
                  <span class="color-grey pdr-10">
                    当前：
                    <span>
                      {{
                        tool.formatDate(
                          new Date(task.nowTime),
                          "yyyy-MM-dd hh:mm:ss"
                        )
                      }}
                    </span>
                  </span>
                </template>
                <span class="color-grey pdr-10">
                  耗时：
                  <span class="color-green">
                    {{ task.useTimeValue }}
                  </span>
                  {{ task.useTimeUnit }}
                </span>
                <template v-if="!task.isEnd">
                  <div @click="stopTask(task)" class="color-red tm-link mgr-10">
                    停止执行
                  </div>
                </template>
                <div class="mgt-5">
                  <span class="color-grey pdr-10">
                    数据准备 需要 / 成功 / 失败：
                    <span>
                      {{ task.dataNumber }}
                      /
                      <span class="color-green">
                        {{ task.readyDataStatistics.dataSuccessCount }}
                      </span>
                      /
                      <span class="color-red">
                        {{ task.readyDataStatistics.dataErrorCount }}
                      </span>
                    </span>
                  </span>
                  <template v-if="task.readyDataStatistics.useTime > 0">
                    <span class="color-grey pdr-10">
                      耗时：
                      <span class="color-green">
                        {{ task.readyDataStatistics.useTimeValue }}
                      </span>
                      {{ task.readyDataStatistics.useTimeUnit }}
                    </span>
                    <span class="color-grey pdr-10">
                      平均：
                      <span class="color-green">
                        {{ task.readyDataStatistics.dataAverage }}
                      </span>
                      个 / 秒
                    </span>
                  </template>
                </div>
                <div class="mgt-5">
                  <span class="color-grey pdr-10">
                    数据导入 总数 / 成功 / 失败：
                    <span>
                      {{ task.doDataStatistics.dataCount }}
                      /
                      <span class="color-green">
                        {{ task.doDataStatistics.dataSuccessCount }}
                      </span>
                      /
                      <span class="color-red">
                        {{ task.doDataStatistics.dataErrorCount }}
                      </span>
                    </span>
                  </span>
                  <template v-if="task.doDataStatistics.useTime > 0">
                    <span class="color-grey pdr-10">
                      耗时：
                      <span class="color-green">
                        {{ task.doDataStatistics.useTimeValue }}
                      </span>
                      {{ task.doDataStatistics.useTimeUnit }}
                    </span>
                    <span class="color-grey pdr-10">
                      平均：
                      <span class="color-green">
                        {{ task.doDataStatistics.dataAverage }}
                      </span>
                      个 / 秒
                    </span>
                  </template>
                  <template v-if="task.isEnd">
                    <div
                      @click="taskClean(task)"
                      class="color-orange tm-link mgr-10"
                    >
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
        </template>
      </tm-layout>
      <tm-layout height="40px">
        <div class="pdlr-10 mgt-8">
          <div class="tm-btn tm-btn-sm bg-green" @click="toDo">
            开始一个任务
          </div>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "tabId", "actived", "extend", "indexName"],
  data() {
    return {
      ready: false,
      importTypes: [{ text: "数据策略", value: "strategy" }],
      formData: {
        importType: "strategy",
        indexName: null,
        dataNumber: 1,
        batchNumber: 200,
        threadNumber: 1,
        errorContinue: true,
        id: "",
      },
      columnList: [],
      taskList: [],
    };
  },
  computed: {},
  watch: {
    "formData.importType"() {},
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
      this.formData.indexName = this.indexName;
      this.formData.id = "'id-' + (_$index + 0)";
      this.ready = true;
      if (this.extend) {
        if (this.extend.columnList) {
          this.columnList = this.extend.columnList;
        }
        if (this.extend.formData) {
          Object.assign(this.formData, this.extend.formData);
        }
      }
      this.autoSaveSql();
      this.loadTaskList();
    },
    upImportColumn(importColumn) {
      this.tool.up(this, "columnList", importColumn);
    },
    downImportColumn(importColumn) {
      this.tool.down(this, "columnList", importColumn);
    },
    addImportColumn(importColumn, after) {
      importColumn = importColumn || {};
      importColumn.name = importColumn.name || "name";
      importColumn.value = importColumn.value || "'name-' + (_$index + 0)";
      importColumn.reuseNumber = 1;
      let appendIndex = this.columnList.indexOf(after);
      if (appendIndex < 0) {
        appendIndex = this.columnList.length;
      } else {
        appendIndex++;
      }
      this.columnList.splice(appendIndex, 0, importColumn);
    },
    removeImportColumn(importColumn) {
      let findIndex = this.columnList.indexOf(importColumn);
      if (findIndex >= 0) {
        this.columnList.splice(findIndex, 1);
      }
    },
    showColumnListCode() {
      this.tool.showJSONData(this.columnList, {
        onSave: (res) => {
          if (res.jsonError) {
            this.tool.error(res.jsonError);
            return;
          }
          res.jsonData = res.jsonData || [];
          if (!res.jsonData.forEach) {
            this.tool.error("非有效JSON");
            return;
          }
          this.columnList.splice(0, this.columnList.length);
          res.jsonData.forEach((one) => {
            one.name = one.name || "";
            one.value = one.value || "";
            one.reuseNumber = one.reuseNumber || 1;
            this.columnList.push(one);
          });
        },
      });
    },
    async toDo() {
      await this.start();
    },
    async start() {
      let param = this.toolboxWorker.getWorkParam(Object.assign(this.formData));

      param.batchNumber = Number(param.batchNumber);
      param.threadNumber = Number(param.threadNumber);
      param.dataNumber = Number(param.dataNumber);
      if (this.tool.isEmpty(param.id)) {
        this.tool.error("请输入id值策略");
        return;
      }
      param.columnList = [];
      this.columnList.forEach((one) => {
        one.reuseNumber = Number(one.reuseNumber);
        param.columnList.push(one);
      });

      let res = await this.server.elasticsearch.import(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      this.loadTaskList();
    },
    async stopTask(task) {
      if (task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: task.taskId,
      });
      await this.server.elasticsearch.taskStop(param);
    },
    async taskClean(task) {
      if (task == null) {
        return;
      }
      let param = this.toolboxWorker.getWorkParam({
        taskId: task.taskId,
      });
      await this.server.elasticsearch.taskClean(param);
      this.loadTaskList();
    },
    async loadTaskList() {
      let param = this.toolboxWorker.getWorkParam({});
      let res = await this.server.elasticsearch.taskList(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let taskList = res.data || [];
      this.taskList = taskList;
      taskList.forEach((one) => {
        this.tool.formatTime(one, "useTime", "useTimeValue", "useTimeUnit");
        this.tool.formatTime(
          one.readyDataStatistics,
          "useTime",
          "useTimeValue",
          "useTimeUnit"
        );
        this.tool.formatTime(
          one.doDataStatistics,
          "useTime",
          "useTimeValue",
          "useTimeUnit"
        );

        this.loadTaskStatus(one);
      });
    },
    async loadTaskStatus(task) {
      if (
        this.isDestroyed ||
        task == null ||
        task.isEnd ||
        task.taskStatusIng
      ) {
        return;
      }
      if (this.taskList.indexOf(task) < 0) {
        return;
      }
      task.taskStatusIng = true;

      let param = this.toolboxWorker.getWorkParam({
        taskId: task.taskId,
      });
      let res = await this.server.elasticsearch.taskStatus(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      Object.assign(task, res.data || {});
      delete task.taskStatusIng;

      this.tool.formatTime(task, "useTime", "useTimeValue", "useTimeUnit");
      this.tool.formatTime(
        task.readyDataStatistics,
        "useTime",
        "useTimeValue",
        "useTimeUnit"
      );
      this.tool.formatTime(
        task.doDataStatistics,
        "useTime",
        "useTimeValue",
        "useTimeUnit"
      );
      setTimeout(() => {
        this.loadTaskStatus(task);
      }, 100);
    },
    async autoSaveSql() {
      if (this.isDestroyed) {
        return;
      }
      if (this.autoSaveSqlIng) {
        return;
      }
      this.autoSaveSqlIng = true;
      let keyValueMap = {};
      let columnListStr = JSON.stringify(this.columnList);
      if (this.columnListStr != columnListStr) {
        this.columnListStr = columnListStr;
        keyValueMap.columnList = this.columnList;
      }
      let formDataStr = JSON.stringify(this.formData);
      if (this.formDataStr != formDataStr) {
        this.formDataStr = formDataStr;
        keyValueMap.formData = this.formData;
      }
      if (Object.keys(keyValueMap).length > 0) {
        await this.toolboxWorker.updateOpenTabExtend(this.tabId, keyValueMap);
      }
      this.autoSaveSqlIng = false;
      setTimeout(this.autoSaveSql, 300);
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
.toolbox-elasticsearch-import {
  width: 100%;
  height: 100%;
}
</style>
