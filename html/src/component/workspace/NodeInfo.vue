<template>
  <div class="node-info-box">
    <template v-if="node != null">
      <div
        class="node-info"
        :class="{
          'node-info-started': node.status == 1,
          'node-info-stopped': node.status == 2,
          'node-info-error': node.status == 3,
          'node-info-root': node.isLocal,
        }"
      >
        <Icon class="node-info-icon mdi-checkbox-blank-circle"></Icon>
        <div class="node-info-text">{{ node.text }}</div>
      </div>
    </template>
  </div>
</template>

<script>
export default {
  name: "node-info",
  components: {},
  props: [],
  inject: ["getGraph", "getNode"],
  data() {
    return {
      source: null,
      node: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      let gNode = this.getNode();
      this.node = gNode.data;
      this.source = gNode.source;
      // 监听数据改变事件
      gNode.on("change:data", ({ current }) => {
        Object.assign(this.node, current);
      });
    },
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.node-info-box {
  position: relative;
  width: 100%;
  height: 100%;
  font-size: 12px;
}
.node-info {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  border: 2px solid #505050;
  color: #adadad;
}
.node-info-started.node-info {
  color: #4caf50;
  border-color: #4caf50;
  background: #163217;
}
.node-info-stopped.node-info {
  color: #f44336;
  border-color: #f44336;
  background: #381613;
}
.node-info-error.node-info {
  color: #ff9800;
  border-color: #ff9800;
  background: #48361a;
}
.node-info-root.node-info .node-info-text {
  color: wheat;
}
.node-info-icon {
  position: absolute;
  right: 5px;
  top: 5px;
  line-height: 12px;
  padding: 0px;
  margin: 0px;
}
.node-info-text {
  width: 100%;
  text-align: center;
  padding: 5px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 15px;
  font-weight: 600;
}
</style>
