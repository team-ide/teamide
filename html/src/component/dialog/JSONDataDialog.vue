<template>
  <el-dialog
    ref="modal"
    :title="title || '数据查看'"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="1200px"
  >
    <div class="mgt--20">
      <div style="height: 560px !important">
        <Editor
          ref="Editor"
          :source="source"
          :value="text"
          language="json"
        ></Editor>
      </div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      showDialog: false,
      text: null,
      title: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    async show(data) {
      this.showDialog = true;
      this.$nextTick(() => {
        try {
          data = data || {};
          this.text = JSON.stringify(data, null, "    ");
        } catch (e) {
          this.text = e;
        }
        this.$refs.Editor.setValue(this.text);
      });
    },
    hide() {
      this.showDialog = false;
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>
