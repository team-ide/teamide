<template>
  <div class="echarts-view-net-box">
    <div ref="view" style="width: 100%; height: 300px"></div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "timeList", "netStatList", "name"],
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

      let legendData = [{ name: "发送" }, { name: "接收" }];
      let xAxisData = [];
      let sentS = { name: "发送", type: "line", smooth: true, data: [] };
      let recvS = { name: "接收", type: "line", smooth: true, data: [] };
      let seriesList = [sentS, recvS];

      let timeList = this.timeList || [];
      let list = this.netStatList || [];

      list.forEach((one, i) => {
        let time = timeList[i];
        xAxisData.push(time);
        recvS.data.push(one.speedRecv);
        sentS.data.push(one.speedSent);
      });

      let option = {
        backgroundColor: "#0f1b26",
        title: {
          text: "网络（" + this.name + "）",
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
          name: "速度(M/秒)",
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
.echarts-view-net-box {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
</style>
