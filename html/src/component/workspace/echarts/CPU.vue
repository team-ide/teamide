<template>
  <div class="echarts-view-cpu-box">
    <div class="pd-10">
      <el-switch
        class="color-grey"
        v-model="showCpuOther"
        active-text="显示CPU所有核心"
        inactive-text="显示CPU总核心"
      >
      </el-switch>
    </div>
    <div ref="view" style="width: 100%; height: 300px"></div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "timeList", "cpuPercentsList"],
  data() {
    return {
      showCpuOther: false,
    };
  },
  computed: {},
  watch: {
    showCpuOther() {
      if (this.$refs.view.echarts) {
        this.$refs.view.echarts.clear();
      }
      this.initView();
    },
  },
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

      let legendData = [{ name: "CPU" }];
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
          if (this.showCpuOther) {
            let ser = seriesList[dIndex + 1];
            if (ser == null) {
              let name = "CPU-" + dIndex;
              legendData.push({ name: name });
              ser = {
                name: name,
                type: "line",
                smooth: true,
                data: [],
              };
              seriesList.push(ser);
            }
            ser.data.push(Number(d));
          }
          cP += Number(d);
        });
        seriesList[0].data.push(cP.toFixed(2));
      });

      let option = {
        backgroundColor: "#0f1b26",
        title: {
          text: "CPU",
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
.echarts-view-cpu-box {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
.echarts-view-cpu-box .el-switch__label {
  color: #ffffff;
}
.echarts-view-cpu-box .el-switch__label.is-active {
  color: #409eff;
}
</style>
