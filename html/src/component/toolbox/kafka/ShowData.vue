<template>
  <el-dialog
    ref="modal"
    title="数据查看"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="700px"
  >
    <div class="ft-15">
      <el-input
        type="textarea"
        v-model="showData"
        :autosize="{ minRows: 5, maxRows: 10 }"
      >
      </el-input>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {
      showDialog: false,
      showData: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    show(data) {
      data = data || {};
      this.showData = JSON.stringify(data, null, "    ");

      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolboxWorker.showData = this.show;
    this.toolboxWorker.hideData = this.hide;
    this.init();
  },
};
</script>

<style>
</style>
