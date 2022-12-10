<template>
  <div class="log-page">
    <tm-layout>
      <tm-layout height="100px">
        <el-form
          class="pdt-10 pdlr-10"
          size="mini"
          :inline="true"
          @submit.native.prevent
        >
          <el-form-item label="操作">
            <el-select
              placeholder="请选择操作"
              v-model="formData.action"
              style="width: 180px"
              clearable
            >
              <template v-for="(one, index) in powers">
                <el-option :key="index" :label="one.text" :value="one.action">
                </el-option>
              </template>
            </el-select>
          </el-form-item>
          <el-form-item label="时间区间">
            <el-date-picker
              v-model="formData.startTime"
              type="datetime"
              placeholder="选择日期时间"
            >
            </el-date-picker>
            <span>-</span>
            <el-date-picker
              v-model="formData.endTime"
              type="datetime"
              placeholder="选择日期时间"
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item label="">
            <div class="">
              <div class="tm-btn tm-btn-sm bg-grey ft-13" @click="doSearch()">
                搜索
              </div>
              <div class="tm-btn tm-btn-sm bg-orange ft-13" @click="toClean()">
                清理
              </div>
            </div>
          </el-form-item>
        </el-form>
      </tm-layout>
      <tm-layout height="auto">
        <div style="height: 100%">
          <el-table
            :data="dataList"
            :border="true"
            height="100%"
            style="width: 100%"
            size="mini"
            @row-contextmenu="rowContextmenu"
          >
            <el-table-column width="60" label="序号" fixed>
              <template slot-scope="scope">
                <span class="mgl-5">{{ scope.$index + 1 }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="userName" label="姓名" width="100">
            </el-table-column>
            <el-table-column prop="ip" label="IP" width="80"> </el-table-column>
            <el-table-column prop="actionText" label="操作" width="180">
            </el-table-column>
            <el-table-column prop="paramStr" label="参数"> </el-table-column>
            <el-table-column prop="useTime" label="耗时(ms)" width="80">
            </el-table-column>
            <el-table-column label="开始时间" width="160">
              <template slot-scope="scope">
                {{ tool.formatDateByTime(scope.row.startTime) }}
              </template>
            </el-table-column>
            <el-table-column label="结束时间" width="160">
              <template slot-scope="scope">
                {{ tool.formatDateByTime(scope.row.endTime) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template slot-scope="scope">
                <div
                  class="tm-link color-grey"
                  @click="toolboxWorker.showJSONData(scope.row)"
                >
                  查看
                </div>
                <div
                  class="tm-link color-grey mgl-5"
                  @click="toolboxWorker.showJSONData(scope.row.paramStr)"
                >
                  查看参数
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </tm-layout>
      <tm-layout height="30px">
        <div class="text-center">
          <el-pagination
            small
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="formData.pageNo"
            :page-sizes="[10, 50, 100, 200, 500]"
            :page-size="formData.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="totalCount"
            :disabled="totalCount <= 0"
          >
          </el-pagination>
        </div>
      </tm-layout>
    </tm-layout>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      formData: {
        action: "",
        startTime: null,
        endTime: null,
        pageSize: 50,
        pageNo: 1,
      },
      totalCount: 0,
      dataList: null,
      powerMap: null,
      powers: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      await this.initPowerData();
      this.doSearch();
    },
    refresh() {
      this.$nextTick(() => {});
    },
    rowContextmenu() {},
    handleSizeChange(pageSize) {
      this.formData.pageSize = pageSize;
      this.doSearch();
    },
    handleCurrentChange(pageNo) {
      this.formData.pageNo = pageNo;
      this.doSearch();
    },
    async doSearch() {
      let param = Object.assign({}, this.formData);
      if (this.tool.isNotEmpty(param.startTime)) {
        param.startTime = new Date(param.startTime).getTime() / 1000;
      }
      if (this.tool.isNotEmpty(param.endTime)) {
        param.endTime = new Date(param.endTime).getTime() / 1000;
      }
      let res = await this.server.log.queryPage(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let data = res.data || {};
      this.totalCount = data.totalCount;
      let dataList = data.dataList || [];
      let powerMap = this.powerMap || {};
      dataList.forEach((one) => {
        let paramStr = "";
        if (this.tool.isNotEmpty(one.param) && this.tool.isNotEmpty(one.data)) {
          paramStr = JSON.stringify({ param: one.param, data: one.data });
        } else if (this.tool.isNotEmpty(one.param)) {
          paramStr = one.param;
        } else if (this.tool.isNotEmpty(one.data)) {
          paramStr = one.data;
        }
        one.paramStr = paramStr;
        if (powerMap[one.action]) {
          one.actionText = powerMap[one.action].text;
        } else {
          one.actionText = one.action;
        }
      });
      this.dataList = dataList;
    },
    async initPowerData() {
      let res = await this.server.power.data({});
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      let data = res.data || {};
      let powers = data.powers || [];
      let powerMap = {};
      powers.forEach((one) => {
        powerMap[one.action] = one;
      });
      this.powers = powers;
      this.powerMap = powerMap;
    },
    toClean() {
      let msg = "日志清理将删除用户所有日志数据，且无法恢复，确认清理？";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doClean();
        })
        .catch((e) => {});
    },
    async doClean() {
      let res = await this.server.log.clean({});
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        this.tool.success("清理成功");
        this.doSearch();
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.log-page {
  width: 100%;
  height: 100%;
}
</style>
