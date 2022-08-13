<template>
  <el-dialog
    ref="modal"
    :title="`节点查看`"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    :append-to-body="true"
    :visible="showBox"
    :before-close="hide"
    :fullscreen="true"
    width="100%"
    class="node-info-dialog"
  >
    <tm-layout height="100%">
      <tm-layout height="150px">
        <div class="node-info-dialog-header pd-10 color-grey">
          <template v-if="node != null && node.model != null">
            <div>
              节点:
              <span class="pdlr-5 color-grey-3"> {{ node.model.name }}</span>
              ID:
              <span class="pdlr-5 color-grey-3">
                {{ node.model.serverId }}</span
              >
            </div>
          </template>
          <template v-if="hostInfoStat != null">
            <div class="pdtb-3">
              Host Name:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.hostname }}
              </span>
              Boot Time:
              <span class="pdlr-5 color-grey-3">
                {{ tool.formatDateByTime(hostInfoStat.bootTime * 1000) }}
              </span>
            </div>
            <div class="pdtb-3">
              OS:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.os }}
              </span>
            </div>
            <div class="pdtb-3">
              Platform:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.platform }}
              </span>
              Platform Family:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.platformFamily }}
              </span>
              Platform Version:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.platformVersion }}
              </span>
            </div>
            <div class="pdtb-3">
              Kernel Version:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.kernelVersion }}
              </span>
              Kernel Arch:
              <span class="pdlr-5 color-grey-3">
                {{ hostInfoStat.kernelArch }}
              </span>
            </div>
          </template>
        </div>
      </tm-layout>
      <tm-layout height="auto" class="" style="overflow: auto">
        <div class="node-info-dialog-body">
          <div class="pd-10">
            <div
              class="tm-btn tm-btn-xs bg-grey-6"
              @click="cleanNodeMonitorData"
            >
              清理节点监控数据
            </div>
          </div>
          <div ref="cpuView" style="width: 100%; height: 300px"></div>
        </div>
      </tm-layout>
    </tm-layout>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showBox: false,
      node: null,
      hostInfoStat: null,
      loading: false,
      timeList: [],
      virtualMemoryStatList: [],
      cpuPercentsList: [],
      diskUsageStatList: [],
      netIOCountersStatsList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    async show(node) {
      this.node = node;
      this.nodeId = node.info.id;
      this.showBox = true;
      let systemInfo = await this.loadInfo();
      if (systemInfo != null) {
        this.hostInfoStat = systemInfo.hostInfoStat;
      }
      this.lastTimestamp = null;
      this.loadMonitorData();
    },
    hide() {
      if (this.$refs.cpuView.echarts) {
        this.$refs.cpuView.echarts.dispose();
        this.$refs.cpuView.echarts = null;
      }
      this.showBox = false;
      this.cleanCacheData();
    },
    cleanCacheData() {
      this.node = null;
      this.nodeId = null;
      this.hostInfoStat = null;
      this.lastTimestamp = null;
      this.timeList = [];
      this.virtualMemoryStatList = [];
      this.cpuPercentsList = [];
      this.diskUsageStatList = [];
      this.netIOCountersStatsList = [];
    },
    async cleanNodeMonitorData() {
      await this.cleanMonitorData();

      if (this.$refs.cpuView.echarts) {
        this.$refs.cpuView.echarts.clear();
      }
      this.lastTimestamp = null;
      this.timeList = [];
      this.virtualMemoryStatList = [];
      this.cpuPercentsList = [];
      this.diskUsageStatList = [];
      this.netIOCountersStatsList = [];
      this.loadMonitorData();
    },
    initView() {
      this.initCpuView();
    },
    initCpuView() {
      if (!this.$refs.cpuView) {
        return;
      }
      let myChart = this.$refs.cpuView.echarts;
      if (myChart) {
        // myChart.clear();
      } else {
        myChart = this.$echarts.init(this.$refs.cpuView);
        this.$refs.cpuView.echarts = myChart;
      }

      let legendData = [{ name: "CPU", textStyle: { color: "#ffffff" } }];
      let xAxisData = [];
      let seriesList = [
        {
          name: "CPU",
          type: "line",
          smooth: true,
          data: [],
        },
      ];

      let timeList = this.timeList || [];
      let list = this.cpuPercentsList || [];

      list.forEach((one, i) => {
        let time = timeList[i];
        xAxisData.push(time);
        let cpuPercents = one || [];
        let cP = 0;
        cpuPercents.forEach((d, dIndex) => {
          let ser = seriesList[dIndex + 1];
          if (ser == null) {
            let name = "CPU-" + dIndex;
            legendData.push({ name: name, textStyle: { color: "#ffffff" } });
            ser = {
              name: name,
              type: "line",
              smooth: true,
              data: [],
            };
            seriesList.push(ser);
          }
          ser.data.push(Number(d));
          cP += Number(d);
        });
        seriesList[0].data.push(cP.toFixed(2));
      });

      let option = {
        title: {
          text: "CPU",
          textStyle: { color: "#ffffff" },
        },
        tooltip: {
          trigger: "axis",
          confine: true,
        },
        legend: {
          data: legendData,
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "3%",
          containLabel: true,
        },
        toolbox: {
          feature: {
            saveAsImage: {},
          },
        },
        xAxis: {
          name: "时间",
          type: "category",
          boundaryGap: false,
          data: xAxisData,
        },
        yAxis: {
          name: "消耗/%",
          type: "value",
        },
        series: seriesList,
      };
      myChart.setOption(option);
    },
    async loadMonitorData() {
      if (!this.showBox || this.tool.isEmpty(this.nodeId)) {
        return;
      }
      if (this.loadMonitorDataIng) {
        return;
      }
      this.loadMonitorDataIng = true;
      let size = 30;
      let data = await this.queryMonitorData(this.lastTimestamp, size);
      if (data && data.monitorDataList && data.monitorDataList.length > 0) {
        this.lastTimestamp = data.lastTimestamp;
        let list = data.monitorDataList;
        list.forEach((one) => {
          this.timeList.push(
            this.tool.formatDateByTime(one.startTime, "HH:mm:ss")
          );
          this.virtualMemoryStatList.push(one.virtualMemoryStat);
          one.cpuPercents.forEach((d, i) => {
            one.cpuPercents[i] = Number(d).toFixed(2);
          });
          this.cpuPercentsList.push(one.cpuPercents);
          this.diskUsageStatList.push(one.diskUsageStat);
          this.netIOCountersStatsList.push(one.netIOCountersStats);
        });
        if (data.monitorDataList.length == size) {
          if (this.loadMonitorDataIng) {
            return;
          }
          this.loadMonitorDataIng = false;
          this.loadMonitorData();
          return;
        }
      }
      this.initView();
      window.setTimeout(() => {
        if (this.loadMonitorDataIng) {
          return;
        }
        this.loadMonitorDataIng = false;
        this.loadMonitorData();
      }, 5000);
    },
    async loadInfo() {
      let res = await this.server.node.system.info({ nodeId: this.nodeId });
      if (res.code == 0) {
        return res.data;
      } else {
        this.tool.error(res.msg);
        return null;
      }
    },
    async queryMonitorData(timestamp, size) {
      if (this.tool.isEmpty(timestamp) || timestamp <= 0) {
        timestamp = new Date().getTime() - 60 * 1000 * 5;
      }
      size = size || 30;
      let res = await this.server.node.system.queryMonitorData({
        nodeId: this.nodeId,
        timestamp: Number(timestamp),
        size: Number(size),
      });
      if (res.code == 0) {
        return res.data;
      } else {
        this.tool.error(res.msg);
        return null;
      }
    },
    async cleanMonitorData() {
      let res = await this.server.node.system.cleanMonitorData({
        nodeId: this.nodeId,
      });
      if (res.code == 0) {
        return res.data;
      } else {
        this.tool.error(res.msg);
        return null;
      }
    },
    init() {},
  },
  created() {},
  updated() {
    this.tool.showNodeInfo = this.show;
    this.tool.hideNodeInfo = this.hide;
  },
  mounted() {
    this.init();
    this.tool.showNodeInfo = this.show;
    this.tool.hideNodeInfo = this.hide;
  },
  destroyed() {},
};
</script>

<style>
.node-info-dialog {
  width: 100%;
  height: 100%;
  user-select: text;
}
.node-info-dialog .el-dialog {
  background: #0f1b26;
  color: #ffffff;
  position: absolute;
  top: 30px;
  bottom: 30px;
  left: 30px;
  right: 30px;
  width: auto;
  height: auto;
}
.node-info-dialog .el-dialog__title {
  color: #ffffff;
}
.node-info-dialog .el-dialog__body {
  position: relative;
  width: 100%;
  height: calc(100% - 55px);
  padding: 0px;
}
.node-info-dialog-header {
  position: relative;
}
.node-info-dialog-body {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
</style>
