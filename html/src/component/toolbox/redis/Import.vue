<template>
  <div class="toolbox-redis-import">
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
      <el-form-item label="数据库">
        <el-input v-model="form.database"> </el-input>
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
              <el-form size="mini" @submit.native.prevent label-width="80px">
                <el-form-item label="导入数量" class="mgb-5">
                  <el-input v-model="strategyData.count"> </el-input>
                </el-form-item>
                <el-form-item label="值类型">
                  <el-select
                    placeholder="请选择类型"
                    v-model="strategyData.valueType"
                    style="width: 100%"
                  >
                    <el-option label="string" value="string"> </el-option>
                    <el-option label="list" value="list"></el-option>
                    <el-option label="set" value="set"></el-option>
                    <el-option label="hash" value="hash"></el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="Key" class="mgb-5">
                  <el-input v-model="strategyData.key"> </el-input>
                </el-form-item>
                <template v-if="strategyData.valueType == 'string'">
                  <div class="ft-12 mgb-5">String 值</div>
                  <el-form-item label="Value" class="mgb-5">
                    <el-input
                      type="textarea"
                      v-model="strategyData.value"
                      :autosize="{ minRows: 5, maxRows: 10 }"
                    >
                    </el-input>
                  </el-form-item>
                </template>
                <template v-else-if="strategyData.valueType == 'list'">
                  <div class="ft-12 mgb-5">List 值</div>
                  <el-form-item label="值数量" class="mgb-5">
                    <el-input v-model="strategyData.valueCount"> </el-input>
                  </el-form-item>
                  <el-form-item label="List Value" class="mgb-5">
                    <el-input
                      type="textarea"
                      v-model="strategyData.listValue"
                      :autosize="{ minRows: 5, maxRows: 10 }"
                    >
                    </el-input>
                  </el-form-item>
                </template>
                <template v-else-if="strategyData.valueType == 'set'">
                  <div class="ft-12 mgb-5">Set 值</div>
                  <el-form-item label="值数量" class="mgb-5">
                    <el-input v-model="strategyData.valueCount"> </el-input>
                  </el-form-item>
                  <el-form-item label="Set Value" class="mgb-5">
                    <el-input
                      type="textarea"
                      v-model="strategyData.setValue"
                      :autosize="{ minRows: 5, maxRows: 10 }"
                    >
                    </el-input>
                  </el-form-item>
                </template>
                <template v-else-if="strategyData.valueType == 'hash'">
                  <div class="">Set 值</div>
                  <el-form-item label="值数量" class="mgb-5">
                    <el-input v-model="strategyData.valueCount"> </el-input>
                  </el-form-item>
                  <el-form-item label="Hash Key" class="mgb-5">
                    <el-input v-model="strategyData.hashKey"> </el-input>
                  </el-form-item>
                  <el-form-item label="Hash Value" class="mgb-5">
                    <el-input
                      type="textarea"
                      v-model="strategyData.hashValue"
                      :autosize="{ minRows: 5, maxRows: 10 }"
                    >
                    </el-input>
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
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      ready: false,
      importTypes: [
        { text: "策略函数", value: "strategy" },
        { text: "Excel", value: "excel", disabled: true },
        { text: "文本", value: "text", disabled: true },
      ],
      form: {
        database: 0,
        importType: "strategy",
      },
      strategyDataList: null,
      taskKey: null,
      task: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      let extend = this.extend || {};
      let database = extend.database || 0;
      this.form.database = Number(database);
      this.strategyDataList = [];
      await this.addStrategyData();
      this.ready = true;
    },
    async addStrategyData() {
      let data = {};
      data.count = 1;
      data.key = "'xx:key:' + (_$index + 0)";
      data.value = "'xx:value:' + (_$index + 0)";
      data.valueType = "string";
      data.valueCount = 1;
      data.listValue = "'list:value:' + (_$value_index + 0)";
      data.setValue = "'set:value:' + (_$value_index + 0)";
      data.hashKey = "'hash:key:' + (_$value_index + 0)";
      data.hashValue = "'hash:value:' + (_$value_index + 0)";

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
        one.valueCount = Number(one.valueCount);
      });
      this.form.database = Number(this.form.database);

      let param = Object.assign({}, this.form);
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
.toolbox-redis-import {
  width: 100%;
  height: 100%;
}
</style>
