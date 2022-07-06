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
    </template>
  </div>
</template>


<script>
import "./toolbox.css";

export default {
  components: {},
  props: ["source", "extend", "toolboxType", "toolboxId", "openId"],
  data() {
    let toolboxWorker = this.tool.newToolboxWorker({
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
      await this.toolboxWorker.init();
      this.ready = true;
    },
    async onFocus() {
      await this.init();
      this.$nextTick(() => {
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
  updated() {},
  beforeDestroy() {},
};
</script>

<style>
.toolbox-editor {
  width: 100%;
  height: 100%;
  overflow: auto;
  border: 0px;
  outline: 0px;
  background-color: #383838;
  color: #d9d9d9;
}
/* 
.toolbox-editor ul {
  margin-top: 10px;
}
.toolbox-editor ul,
.toolbox-editor li {
  list-style: none;
  padding: 0px;
  font-size: 12px;
}
.toolbox-editor li {
  display: block;
  line-height: 22px;
  margin-bottom: 3px;
} */
.toolbox-editor .text {
  display: inline-block;
  min-width: 80px;
}
.toolbox-editor .text,
.toolbox-editor .input,
.toolbox-editor .comment {
  padding: 0px 5px;
}

.toolbox-editor table {
  padding: 0px 0px;
  width: 100%;
  border-collapse: collapse;
}
.toolbox-editor table th {
  border: 1px solid #4e4e4e;
  text-align: center;
  line-height: 30px;
}
.toolbox-editor table tr td {
  border: 1px solid #4e4e4e;
  padding: 3px 5px;
}
.toolbox-editor table td .input {
  padding: 0px 0px;
}
.toolbox-editor table td .model-input {
  min-width: 80px;
}

.toolbox-editor .el-tree-node__content {
  border-bottom: 1px dotted #696969;
}

.toolbox-editor .el-tree {
  background-color: transparent;
  color: unset;
  width: auto;
  font-size: 12px;
  user-select: none;
  border: 1px dotted #696969;
  border-bottom: 0px;
}
.toolbox-editor .el-tree .mdi {
  vertical-align: middle;
}
.toolbox-editor .el-tree-node__content {
  position: relative;
}

.toolbox-editor .el-tree-node__content .toolbox-editor-tree-btn-group {
  display: none;
}
.toolbox-editor .el-tree-node__content:hover .toolbox-editor-tree-btn-group {
  display: block;
}
.toolbox-editor .toolbox-editor-tree-span .toolbox-editor-tree-btn-group {
  position: absolute;
  right: 5px;
  top: -1px;
}
.toolbox-editor .toolbox-editor-tree-span {
  /* flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between; */
  font-size: 13px;
}
.toolbox-editor .el-tree__empty-block {
  display: none;
}
.toolbox-editor .el-tree-node__children {
  overflow: visible !important;
}
.toolbox-editor .el-tree-node__children.v-enter-active,
.toolbox-editor .el-tree-node__children.v-leave-active {
  overflow: hidden !important;
}
.toolbox-editor .el-tree .el-tree-node.is-current > .el-tree-node__content {
  background-color: #636363 !important;
}
.toolbox-editor .el-tree .el-tree-node:focus > .el-tree-node__content {
  background-color: #545454;
}
.toolbox-editor .el-tree .el-tree-node > .el-tree-node__content:hover {
  background-color: #545454;
}

.part-box {
  line-height: 20px;
  font-size: 12px;
  overflow: auto;
  width: 100%;
  height: 100%;
}
.part-box,
.part-box li {
  padding: 0px;
  margin: 0px;
  list-style: none;
}
.part-box li {
  text-overflow: ellipsis;
  white-space: nowrap;
  word-break: keep-all;
}

.part-box input,
.part-box select {
  color: #ffffff;
  width: 40px;
  min-width: 40px;
  border: 1px dashed transparent;
  background-color: transparent;
  height: 20px;
  max-width: 100%;
  padding: 0px;
  padding-left: 2px;
  padding-right: 2px;
  box-sizing: border-box;
  outline: none;
  font-size: 12px;
}

.part-box input {
  border-bottom: 1px dashed #636363;
}
.part-box select {
  -moz-appearance: auto;
  -webkit-appearance: auto;
}
.part-box option {
  background-color: #ffffff;
  color: #3e3e3e;
}
.part-box input[type="checkbox"] {
  width: 10px;
  min-width: 10px;
  height: 13px;
  vertical-align: -3px;
  margin-left: 6px;
}

.part-box textarea {
  color: #ffffff;
  height: 70px;
  border: 1px dashed #636363;
  text-align: left;
  padding: 5px;
  min-width: 500px;
  background-color: transparent;
  font-size: 12px;
  vertical-align: text-top;
}
</style>
