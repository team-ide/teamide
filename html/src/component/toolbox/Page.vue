<template>
  <div class="page-box">
    <template v-if="extend != null">
      <template v-if="extend.page == 'userCenter'">
        <UserCenter
          :source="source"
          :toolboxWorker="toolboxWorker"
          :extend="extend"
        ></UserCenter>
      </template>
    </template>
  </div>
</template>


<script>
import UserCenter from "./user/Center";

export default {
  components: { UserCenter },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {},
    refresh() {
      this.$nextTick(() => {});
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.page-box {
  width: 100%;
  height: 100%;
}
</style>
