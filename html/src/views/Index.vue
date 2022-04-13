<template>
  <div>
    11
    <div style="width: 500px; height: 500px">
      <TabEditor ref="TabEditor" :source="source">
        <template v-slot:body="{ tab }">
          <template v-if="tab.type == '1'">
            <div class="">type-1</div>
          </template>
          <template v-if="tab.type == '2'">
            <div class="">type-2</div>
          </template>
          <template v-if="tab.type == '3'">
            <div class="">type-3{{ tabData }}</div>
          </template>
        </template>
      </TabEditor>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      tabData: {
        a: "xxx",
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    initTabEditor() {
      this.addTab("xxx-1", "1");
      this.addTab("xxx-2", "2");
      this.addTab("xxx-3", "3");
      this.addTab("xxx-4", "1");
      this.addTab("xxx-5", "3");
    },
    addTab(key, type) {
      let tab = {};
      tab.key = key;
      tab.title = key;
      tab.name = key;
      tab.type = type;
      this.$refs.TabEditor.addTab(tab);
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.initTabEditor();
  },
};
</script>

<style>
</style>
