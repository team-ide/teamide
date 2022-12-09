<template>
  <div class="toolbox-editor" v-if="toolboxType != null" tabindex="-1">
    <template v-if="ready">
      <template v-if="toolboxType == 'redis'">
        <ToolboxRedisEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxRedisEditor>
      </template>
      <template v-else-if="toolboxType == 'database'">
        <ToolboxDatabaseEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxDatabaseEditor>
      </template>
      <template v-else-if="toolboxType == 'zookeeper'">
        <ToolboxZookeeperEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxZookeeperEditor>
      </template>
      <template v-else-if="toolboxType == 'elasticsearch'">
        <ToolboxElasticsearchEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxElasticsearchEditor>
      </template>
      <template v-else-if="toolboxType == 'kafka'">
        <ToolboxKafkaEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxKafkaEditor>
      </template>
      <template v-else-if="toolboxType == 'ssh'">
        <ToolboxSSHEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxSSHEditor>
      </template>
      <template v-else-if="toolboxType == 'other'">
        <ToolboxOtherEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </ToolboxOtherEditor>
      </template>
      <template v-else-if="toolboxType == 'node'">
        <NodeEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </NodeEditor>
      </template>
      <template v-else-if="toolboxType == 'file-manager'">
        <FileManagerEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </FileManagerEditor>
      </template>
      <template v-else-if="toolboxType == 'terminal'">
        <TerminalEditor
          :source="source"
          :extend="extend"
          :toolboxWorker="toolboxWorker"
        >
        </TerminalEditor>
      </template>
      <template v-else-if="toolboxType == 'page'">
        <Page :source="source" :extend="extend" :toolboxWorker="toolboxWorker">
        </Page>
      </template>
    </template>
    <JSONDataDialog ref="JSONDataDialog" :source="source"></JSONDataDialog>
  </div>
</template>


<script>
import "./toolbox.css";
import Page from "./Page.vue";
import NodeEditor from "./node/Index.vue";
import FileManagerEditor from "./file-manager/Index.vue";
import TerminalEditor from "./terminal/Index.vue";
import toolboxWorker_ from './toolboxWorker.js';

export default {
  components: {
    Page,
    NodeEditor,
    FileManagerEditor,
    TerminalEditor,
  },
  props: ["source", "extend", "toolboxType", "toolboxId", "openId"],
  data() {
    let toolboxWorker = toolboxWorker_.newToolboxWorker({
      toolboxId: this.toolboxId,
      openId: this.openId,
      toolboxType: this.toolboxType,
      extend: this.extend,
    });

    return {
      toolboxWorker: toolboxWorker,
      extendJSON: null,
      ready: false,
    };
  },
  computed: {},
  watch: {
    extend(newExtent, oldExtent) {
      if (newExtent == null || oldExtent == null) {
        return;
      }
      if (JSON.stringify(newExtent) == JSON.stringify(oldExtent)) {
        return;
      }
      this.toolboxWorker.updateOpenExtend(this.extend);
    },
  },
  methods: {
    async init() {
      if (this.inited) {
        return;
      }
      this.inited = true;
      this.toolboxWorker.showJSONData = (data) => {
        this.$refs.JSONDataDialog.show(data);
      };
      await this.toolboxWorker.init();
      this.ready = true;
    },
    async onFocus() {
      this.$nextTick(async () => {
        await this.init();
        this.$el && this.$el.focus && this.$el.focus();
        this.$children.forEach((one) => {
          one.onFocus && one.onFocus();
        });
      });
    },
    reload() {},
    onKeyDown() {
      if (this.tool.keyIsF5()) {
        this.tool.stopEvent();
        this.$children.forEach((one) => {
          one.refresh && one.refresh();
        });
      }
    },
    bindEvent() {
      if (this.bindEvented) {
        return;
      }
      this.bindEvented = true;
      this.$el.addEventListener("keydown", (e) => {
        this.onKeyDown(e);
      });
    },
  },
  created() {},
  mounted() {
    this.bindEvent();
  },
  beforeDestroy() {
    let param = {};
    this.toolboxWorker.work("close", param);
  },
  updated() {},
};
</script>

<style>
</style>
