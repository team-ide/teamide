<template>
  <div class="echarts-view-net-list-box">
    <template v-for="(one, index) in netList">
      <NetOne
        :key="index"
        ref="view"
        :source="source"
        :timeList="timeList"
        :name="one.name"
        :netStatList="one.netStatList"
      ></NetOne>
    </template>
  </div>
</template>

<script>
import NetOne from "./NetOne.vue";

export default {
  components: { NetOne },
  props: ["source", "timeList", "netIOCountersStatsList"],
  data() {
    return {
      netList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    dispose() {
      let list = this.$refs.view || [];
      list.forEach((one) => {
        one.dispose && one.dispose();
      });
    },
    clear() {
      let list = this.$refs.view || [];
      list.forEach((one) => {
        one.clear && one.clear();
      });
    },
    toInitView() {
      let list = this.$refs.view || [];
      list.forEach((one) => {
        one.initView && one.initView();
      });
    },
    initView() {
      let list = this.netIOCountersStatsList || [];
      let netList = [];
      let netMap = {};
      list.forEach((netS) => {
        netS.forEach((one) => {
          let find = netMap[one.name];
          if (find == null) {
            find = {
              name: one.name,
              netStatList: [],
            };
            netMap[one.name] = find;
            netList.push(find);
          }
          find.netStatList.push({
            speedSent: one.speedSent,
            speedRecv: one.speedRecv,
          });
        });
      });
      this.netList = netList;
      this.$nextTick(() => {
        this.toInitView();
      });
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
.echarts-view-net-list-box {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0f1b26;
}
</style>
