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
        <div class="node-info-dialog-body scrollbar">
          <div class="pd-10">
            <div
              class="tm-btn tm-btn-xs bg-grey-6"
              @click="cleanNodeMonitorData"
            >
              清理节点监控数据
            </div>
          </div>
          <div>
            <CPU
              ref="cpu"
              :source="source"
              :timeList="timeList"
              :cpuPercentsList="cpuPercentsList"
            ></CPU>
          </div>
          <div>
            <Mem
              ref="mem"
              :source="source"
              :timeList="timeList"
              :virtualMemoryStatList="virtualMemoryStatList"
            ></Mem>
          </div>
          <div>
            <Disk
              ref="disk"
              :source="source"
              :timeList="timeList"
              :diskUsageStatList="diskUsageStatList"
            ></Disk>
          </div>
          <div>
            <Net
              ref="net"
              :source="source"
              :timeList="timeList"
              :netIOCountersStatsList="netIOCountersStatsList"
            ></Net>
          </div>
        </div>
      </tm-layout>
    </tm-layout>
  </el-dialog>
</template>

<script>
import CPU from "./echarts/CPU.vue";
import Mem from "./echarts/Mem.vue";
import Net from "./echarts/Net.vue";
import Disk from "./echarts/Disk.vue";

export default {
  components: { CPU, Mem, Net, Disk },
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
      viewList: ["cpu", "disk", "mem", "net"],
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
      this.viewList.forEach((one) => {
        if (this.$refs[one] && this.$refs[one].dispose) {
          this.$refs[one].dispose();
        }
      });

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

      this.viewList.forEach((one) => {
        if (this.$refs[one] && this.$refs[one].clear) {
          this.$refs[one].clear();
        }
      });

      this.lastTimestamp = null;
      this.timeList = [];
      this.virtualMemoryStatList = [];
      this.cpuPercentsList = [];
      this.diskUsageStatList = [];
      this.netIOCountersStatsList = [];
      this.loadMonitorData();
    },
    initView() {
      this.viewList.forEach((one) => {
        if (this.$refs[one] && this.$refs[one].initView) {
          this.$refs[one].initView();
        }
      });
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

          one.virtualMemoryStat.total = Number(
            one.virtualMemoryStat.total / 1024 / 1024 / 1024
          ).toFixed(2);
          one.virtualMemoryStat.free = Number(
            one.virtualMemoryStat.free / 1024 / 1024 / 1024
          ).toFixed(2);
          one.virtualMemoryStat.used = Number(
            one.virtualMemoryStat.used / 1024 / 1024 / 1024
          ).toFixed(2);
          this.virtualMemoryStatList.push(one.virtualMemoryStat);

          one.cpuPercents.forEach((d, i) => {
            one.cpuPercents[i] = Number(d).toFixed(2);
          });
          this.cpuPercentsList.push(one.cpuPercents);

          one.diskUsageStat.total = Number(
            one.diskUsageStat.total / 1024 / 1024 / 1024
          ).toFixed(2);
          one.diskUsageStat.free = Number(
            one.diskUsageStat.free / 1024 / 1024 / 1024
          ).toFixed(2);
          one.diskUsageStat.used = Number(
            one.diskUsageStat.used / 1024 / 1024 / 1024
          ).toFixed(2);
          this.diskUsageStatList.push(one.diskUsageStat);

          one.netIOCountersStats.forEach((d, i) => {
            one.netIOCountersStats[i].speedRecv =
              one.netIOCountersStats[i].speedRecv || 0;
            one.netIOCountersStats[i].speedRecv = Number(
              d.speedRecv / 1024 / 1024
            ).toFixed(2);
            one.netIOCountersStats[i].speedSent =
              one.netIOCountersStats[i].speedSent || 0;
            one.netIOCountersStats[i].speedSent = Number(
              d.speedSent / 1024 / 1024
            ).toFixed(2);
          });
          this.netIOCountersStatsList.push(one.netIOCountersStats);
        });
        if (data.monitorDataList.length == size) {
          this.loadMonitorDataIng = false;
          this.loadMonitorData();
          return;
        }
      }
      this.initView();
      window.setTimeout(() => {
        this.loadMonitorDataIng = false;
        this.loadMonitorData();
      }, 5000);
    },
    async loadInfo() {
      try {
        let res = await this.server.node.system.info({ nodeId: this.nodeId });
        if (res.code == 0) {
          return res.data;
        } else {
          this.tool.error(res.msg);
          return null;
        }
      } catch (error) {
        return null;
      }
    },
    async queryMonitorData(timestamp, size) {
      if (this.tool.isEmpty(timestamp) || timestamp <= 0) {
        timestamp = new Date().getTime() - 60 * 1000 * 5;
      }
      size = size || 30;
      try {
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
      } catch (error) {
        this.tool.error(error.message);
        return null;
      }
    },
    async cleanMonitorData() {
      try {
        let res = await this.server.node.system.cleanMonitorData({
          nodeId: this.nodeId,
        });
        if (res.code == 0) {
          return res.data;
        } else {
          this.tool.error(res.msg);
          return null;
        }
      } catch (error) {}
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
  destroyed() {
    this.showBox = false;
  },
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
