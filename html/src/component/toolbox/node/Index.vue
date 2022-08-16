<template>
  <div class="toolbox-node-editor">
    <template v-if="type == 'node-context'">
      <NodeContext :source="source"></NodeContext>
    </template>
    <template v-if="type == 'net-proxy-context'">
      <NetProxyContext :source="source"></NetProxyContext>
    </template>
    <template v-else-if="type == 'node-info'">
      <Info :source="source" :serverId="serverId"></Info>
    </template>
  </div>
</template>


<script>
import Info from "./Info.vue";
import NodeContext from "./NodeContext.vue";
import NetProxyContext from "./NetProxyContext.vue";

export default {
  components: { Info, NodeContext, NetProxyContext },
  props: ["source", "toolboxWorker", "extend"],
  data() {
    return {
      type: null,
      serverId: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      if (this.extend) {
        this.type = this.extend.type;
        this.serverId = this.extend.serverId;
      }
    },
    refresh() {},
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-node-editor {
  width: 100%;
  height: 100%;
}
</style>
