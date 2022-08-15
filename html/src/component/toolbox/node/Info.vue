<template>
  <div class="toolbox-node-info-editor">
    <tm-layout height="100%">
      <tm-layout height="150px">
        <div class="toolbox-node-info-editor-header pd-10 color-grey">
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
        <div class="toolbox-node-info-editor-body scrollbar">
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
  </div>
</template>

<script>
import CPU from "./echarts/CPU.vue";
import Mem from "./echarts/Mem.vue";
import Net from "./echarts/Net.vue";
import Disk from "./echarts/Disk.vue";

export default {
  components: { CPU, Mem, Net, Disk },
  props: ["source", "serverId"],
  data() {
    return {
      isDestroyed: false,
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
    async init() {
      let systemInfo = await this.loadInfo();
      if (systemInfo != null) {
        this.hostInfoStat = systemInfo.hostInfoStat;
      }
      this.lastTimestamp = null;
      this.loadMonitorData();
    },
    dispose() {
      this.viewList.forEach((one) => {
        if (this.$refs[one] && this.$refs[one].dispose) {
          this.$refs[one].dispose();
        }
      });
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
      if (this.isDestroyed || this.tool.isEmpty(this.serverId)) {
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
        let res = await this.server.node.system.info({ nodeId: this.serverId });
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
          nodeId: this.serverId,
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
          nodeId: this.serverId,
        });
        if (res.code == 0) {
          return res.data;
        } else {
          this.tool.error(res.msg);
          return null;
        }
      } catch (error) {}
    },
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
  destroyed() {
    this.isDestroyed = true;
    this.dispose();
  },
};
</script>

<style>
.toolbox-node-info-editor {
  width: 100%;
  height: 100%;
  user-select: text;
}
.toolbox-node-info-editor-header {
  position: relative;
}
.toolbox-node-info-editor-body {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
</style>
