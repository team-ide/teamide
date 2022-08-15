<template>
  <div class="echarts-view-mem-box">
    <div ref="view" style="width: 100%; height: 300px"></div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "timeList", "virtualMemoryStatList"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    dispose() {
      if (this.$refs.view.echarts) {
        this.$refs.view.echarts.dispose();
        this.$refs.view.echarts = null;
      }
    },
    clear() {
      if (this.$refs.view.echarts) {
        this.$refs.view.echarts.clear();
      }
    },
    initView() {
      if (!this.$refs.view) {
        return;
      }
      let myChart = this.$refs.view.echarts;
      if (myChart) {
        // myChart.clear();
      } else {
        myChart = this.$echarts.init(this.$refs.view, "dark");
        this.$refs.view.echarts = myChart;
      }

      let legendData = [{ name: "总" }, { name: "已用" }, { name: "空闲" }];
      let xAxisData = [];
      let totalS = { name: "总", type: "line", smooth: true, data: [] };
      let usedS = { name: "已用", type: "line", smooth: true, data: [] };
      let freeS = { name: "空闲", type: "line", smooth: true, data: [] };
      let seriesList = [totalS, usedS, freeS];

      let timeList = this.timeList || [];
      let list = this.virtualMemoryStatList || [];

      list.forEach((one, i) => {
        let time = timeList[i];
        xAxisData.push(time);
        totalS.data.push(one.total);
        usedS.data.push(one.used);
        freeS.data.push(one.free);
      });

      let option = {
        backgroundColor: "#0f1b26",
        title: {
          text: "内存",
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
          name: "大小/G",
          type: "value",
        },
        series: seriesList,
      };
      myChart.setOption(option);
    },
    init() {},
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.dispose();
  },
};
</script>

<style>
.echarts-view-mem-box {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
</style>
