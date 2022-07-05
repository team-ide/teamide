<template>
  <Workspace :source="source" ref="Workspace">
    <el-dropdown
      slot="mainTabLeftExtend"
      trigger="click"
      class="workspace-tabs-nav-dropdown"
    >
      <div class="workspace-tabs-nav tm-pointer pdlr-5 ft-12">
        全部
        <i class="mdi mdi-menu-down"></i>
      </div>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item> 全部 </el-dropdown-item>
        <el-dropdown-item> 分组1 </el-dropdown-item>
        <el-dropdown-item> 分组2 </el-dropdown-item>
        <el-dropdown-item> 分组3 </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <div
      slot="mainTabRightExtend"
      class="workspace-tabs-nav tm-pointer color-green pdlr-2"
      @click="toolbox.showSwitchToolboxType()"
    >
      <i class="mdi mdi-plus"></i>
    </div>
    <ToolboxType :source="source"></ToolboxType>
  </Workspace>
</template>

<script>
import ToolboxType from "./toolbox/ToolboxType";
export default {
  components: { ToolboxType },
  props: ["source"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    async init() {
      await this.source.initToolboxData();
      this.initOpens();
      let iconFonts = [
        "teamide-database",
        "teamide-redis",
        "teamide-elasticsearch",
        "teamide-zookeeper",
        "teamide-kafka",
        "teamide-ftp",
        "teamide-ssh",
      ];
      for (let i = 0; i < 20; i++) {
        let item = {};
        item.key = "key:item:" + i;
        item.name = "name:item:" + i;
        item.title = "title:item:" + i;
        item.show = true;
        item.isToolbox = true;
        item.toolboxType = "other";
        item.iconFont = iconFonts[i % iconFonts.length];
        item.extend = {
          type: "format",
        };
        this.addMainItem(item);
      }
    },
    addMainItem(item) {
      this.$refs.Workspace.mainItemsWorker.addItem(item);
    },
    toMainActiveItem(item) {
      this.$refs.Workspace.mainItemsWorker.toActiveItem(item);
    },
    addMainItemByOpen(open) {
      let item = {};
      item.key = open.openId;
      item.name = open.toolboxName;
      item.title = open.toolboxName;
      item.show = true;
      item.isToolbox = true;
      item.toolboxType = open.toolboxType;
      item.extend = open.extend || {};
      item.openId = open.openId;
      switch (item.toolboxType) {
        case "database":
          item.iconFont = "teamide-database";
          break;
        case "redis":
          item.iconFont = "teamide-redis";
          break;
        case "elasticsearch":
          item.iconFont = "teamide-elasticsearch";
          break;
        case "zookeeper":
          item.iconFont = "teamide-zookeeper";
          break;
        case "kafka":
          item.iconFont = "teamide-kafka";
          break;
        case "ssh":
          item.iconFont = "teamide-ssh";
          if (item.extend && item.extend.isFTP) {
            item.iconFont = "teamide-ftp";
          }
          break;
        case "other":
          break;
      }
      this.addMainItem(item);
      open.item = item;
      return item;
    },
    async initOpens() {
      let opens = [];
      let res = await this.server.toolbox.queryOpens({});
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        opens = res.data.opens || [];
      }
      opens.forEach((one) => {
        one.extend = this.tool.getOptionJSON(one.extend);
        this.addMainItemByOpen(one);
      });
      // 激活最后
      let activeOpen = null;
      opens.forEach(async (one) => {
        if (activeOpen == null) {
          activeOpen = one;
        } else {
          if (
            new Date(one.openTime).getTime() >
            new Date(activeOpen.openTime).getTime()
          ) {
            activeOpen = one;
          }
        }
      });
      if (activeOpen != null) {
        this.toMainActiveItem(activeOpen.item);
      } else {
        // this.source.toolbox.showToolboxType();
      }
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>
