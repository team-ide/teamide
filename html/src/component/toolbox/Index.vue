<template>
  <div class="toolbox-editor" v-if="toolboxType != null" tabindex="-1">
    <template v-if="ready">
      <template v-if="toolboxType.name == 'redis'">
        <ToolboxRedisEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxRedisEditor>
      </template>
      <template v-else-if="toolboxType.name == 'database'">
        <ToolboxDatabaseEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxDatabaseEditor>
      </template>
      <template v-else-if="toolboxType.name == 'zookeeper'">
        <ToolboxZookeeperEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxZookeeperEditor>
      </template>
      <template v-else-if="toolboxType.name == 'elasticsearch'">
        <ToolboxElasticsearchEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxElasticsearchEditor>
      </template>
      <template v-else-if="toolboxType.name == 'kafka'">
        <ToolboxKafkaEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxKafkaEditor>
      </template>
      <template v-else-if="toolboxType.name == 'ssh'">
        <ToolboxSSHEditor
          :source="source"
          :toolbox="toolbox"
          :toolboxType="toolboxType"
          :extend="extend"
          :wrap="wrap"
        >
        </ToolboxSSHEditor>
      </template>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: [
    "source",
    "extend",
    "toolboxData",
    "toolboxType",
    "toolbox",
    "active",
    "openId",
    "updateOpenExtend",
    "updateOpenComment",
  ],
  data() {
    return {
      extendJSON: null,
      ready: false,
      wrap: {
        tabs: [],
      },
    };
  },
  computed: {},
  watch: {
    active() {
      this.init();
    },
    extend(newExtent, oldExtent) {
      if (newExtent == null || oldExtent == null) {
        return;
      }
      if (JSON.stringify(newExtent) == JSON.stringify(oldExtent)) {
        return;
      }
      this.wrap.updateOpenExtend(this.extend);
    },
  },
  methods: {
    init() {
      if (this.inited) {
        return;
      }
      if (!this.active) {
        return;
      }
      this.inited = true;
      this.wrap.work = this.work;
      this.wrap.openTabByExtend = this.openTabByExtend;
      this.wrap.onActiveTab = this.onActiveTab;
      this.wrap.onRemoveTab = this.onRemoveTab;
      this.wrap.updateComment = this.updateComment;
      this.wrap.updateOpenExtend = this.updateOpenExtend;
      this.wrap.updateOpenTabExtend = this.updateOpenTabExtend;
      this.wrap.updateExtend = this.updateExtend;
      this.ready = true;
      this.initOpenTabs();
    },
    async work(work, data) {
      let param = {
        toolboxId: this.toolboxData.toolboxId,
        work: work,
        data: data,
      };
      let res = await this.server.toolbox.work(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
      return res;
    },
    onFocus() {
      this.$el.focus();
      this.$children.forEach((one) => {
        one.onFocus && one.onFocus();
      });
    },
    reload() {},
    onRemoveTab(tab) {
      this.toolbox.closeOpenTab(tab.tabId);
    },
    onActiveTab(tab) {
      this.toolbox.activeOpenTab(tab.tabId);
    },
    doActiveTab(tab) {
      this.wrap.doActiveTab(tab);
    },
    updateExtend(keyValueMap) {
      this.updateOpenExtend(this.openId, keyValueMap);
    },
    updateComment(comment) {
      this.updateOpenComment(this.openId, comment);
    },
    async updateOpenTabExtend(tabId, keyValueMap) {
      let tab = this.wrap.getTab("" + tabId);
      if (tab == null) {
        return;
      }
      if (keyValueMap == null) {
        return;
      }
      if (Object.keys(keyValueMap) == 0) {
        return;
      }
      let obj = tab.extend;
      for (let key in keyValueMap) {
        let value = keyValueMap[key];
        let names = key.split(".");
        names.forEach((name, index) => {
          if (index < names.length - 1) {
            obj[name] = obj[name] || {};
            obj = obj[name];
          } else {
            obj[name] = value;
          }
        });
      }

      let param = {
        tabId: tabId,
        extend: JSON.stringify(tab.extend),
      };
      let res = await this.server.toolbox.updateOpenTabExtend(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    async openTabByExtend(extend, fromTab) {
      let data = {
        openId: this.openId,
        toolboxId: this.toolboxData.toolboxId,
      };

      let tabData = await this.toolbox.openTab(data, extend);
      if (tabData == null) {
        return;
      }
      let tab = await this.openByTabData(tabData, fromTab);
      if (tab != null) {
        this.doActiveTab(tab);
      }
    },
    async openByTabData(tabData, fromTab) {
      if (this.tool.isNotEmpty(tabData.extend)) {
        tabData.extend = JSON.parse(tabData.extend);
      } else {
        tabData.extend = null;
      }
      let tab = this.createTabByTabData(tabData);
      this.wrap.addTab(tab, fromTab);
      return tab;
    },
    createTabByTabData(tabData) {
      let key = tabData.tabId;

      let tab = this.wrap.getTab(key);
      if (tab == null) {
        tab = this.toolbox.createOpenTabTab(tabData);
        tab.key = key;
      }
      return tab;
    },
    async initOpenTabs() {
      let tabs = await this.toolbox.loadOpenTabs(this.openId);

      await tabs.forEach(async (tabData) => {
        await this.openByTabData(tabData);
      });

      // 激活最后
      let activeTabData = null;
      tabs.forEach(async (tabData) => {
        if (activeTabData == null) {
          activeTabData = tabData;
        } else {
          if (
            new Date(tabData.openTime).getTime() >
            new Date(activeTabData.openTime).getTime()
          ) {
            activeTabData = tabData;
          }
        }
      });
      if (activeTabData != null) {
        this.doActiveTab(activeTabData.tabId);
      }
    },
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
    this.init();
    this.bindEvent();
  },
  updated() {},
  beforeDestroy() {
    if (this.wrap.destroy != null) {
      this.wrap.destroy();
    }
  },
};
</script>

<style>
.toolbox-editor {
  width: 100%;
  height: 100%;
  overflow: auto;
  border: 0px;
  outline: 0px;
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
}
.toolbox-editor table thead {
  border: 1px solid #4e4e4e;
}
.toolbox-editor table th {
  text-align: center;
  line-height: 30px;
}
.toolbox-editor table td {
  border-right: 1px solid #4e4e4e;
  border-bottom: 1px solid #4e4e4e;
  padding: 3px 5px;
}
.toolbox-editor table tbody {
  border-left: 1px solid #4e4e4e;
}
.toolbox-editor table td .input {
  padding: 0px 0px;
}
.toolbox-editor table td .model-input {
  min-width: 80px;
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
