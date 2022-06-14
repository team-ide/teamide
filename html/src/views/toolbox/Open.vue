<template>
  <div class="toolbox-open-box" v-if="source.toolbox.context != null">
    <template v-if="openTab != null">
      <div class="toolbox-box">
        <ToolboxEditor
          :source="source"
          :toolbox="source.toolbox"
          :tab="openTab"
          :toolboxType="openTab.toolboxType"
          :openId="openTab.openId"
          :toolboxData="openTab.toolboxData"
          :extend="openTab.extend"
          :active="openTab.active"
          :updateOpenExtend="updateOpenExtend"
          :updateOpenComment="updateOpenComment"
        >
        </ToolboxEditor>
      </div>
      <QuickCommand
        v-if="source.toolbox.context != null"
        :source="source"
        :toolbox="source.toolbox"
      >
      </QuickCommand>
      <QuickCommandSSHCommandForm
        v-if="source.toolbox.context != null"
        :source="source"
        :toolbox="source.toolbox"
      >
      </QuickCommandSSHCommandForm>
    </template>
  </div>
</template>

<script>
import QuickCommand from "./QuickCommand";
import QuickCommandSSHCommandForm from "./QuickCommandSSHCommandForm";

export default {
  components: { QuickCommand, QuickCommandSSHCommandForm },
  props: ["source", "data"],
  data() {
    return {
      openTab: null,
    };
  },
  computed: {},
  watch: {
    openTab() {
      this.initTitle();
    },
  },
  methods: {
    initTitle() {
      let tab = this.openTab;
      if (tab == null) {
        return;
      }
      let title = "Team · IDE - " + tab.title;
      if (this.tool.isNotEmpty(tab.comment)) {
        title += ` (${tab.comment})`;
      }
      document.title = title;
    },
    async init() {
      if (this.source.toolbox.context == null) {
        await this.source.toolbox.initContext();
      }
      let openData = await this.getOpen();
      if (openData == null) {
        return;
      }
      this.openTab = await this.getTabByOpenData(openData);
      this.openTab.active = true;
      this.initTitle();
    },
    async updateOpenExtend(openId, keyValueMap) {
      let tab = this.openTab;
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
        openId: openId,
        extend: JSON.stringify(tab.extend),
      };
      let res = await this.server.toolbox.updateOpenExtend(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    updateOpenComment(openId, comment) {
      let tab = this.openTab;
      if (tab == null) {
        return;
      }
      tab.comment = comment;
    },
    async getTabByOpenData(openData) {
      let toolboxData = this.source.toolbox.getToolboxData(openData.toolboxId);
      if (toolboxData == null) {
        return;
      }
      let toolboxType = this.source.toolbox.getToolboxType(
        toolboxData.toolboxType
      );
      if (toolboxType == null) {
        return;
      }
      openData.toolboxData = toolboxData;
      openData.toolboxType = toolboxType;
      if (this.tool.isNotEmpty(openData.extend)) {
        openData.extend = JSON.parse(openData.extend);
      } else {
        openData.extend = {};
      }
      let tab = await this.source.toolbox.createToolboxDataTab(openData);
      return tab;
    },
    async getOpen() {
      if (this.tool.isEmpty(this.data)) {
        this.tool.error("OpenID等数据丢失");
      }
      let res = await this.server.toolbox.getOpen({
        openId: Number(this.data),
      });
      if (res.code != 0) {
        this.tool.error(res.msg);
        return null;
      }
      let openData = res.data.open;
      return openData;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-open-box {
  width: 100%;
  height: 100%;
}
</style>
