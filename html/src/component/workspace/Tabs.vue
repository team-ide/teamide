<template>
  <div class="workspace-tabs">
    <div class="workspace-tabs-left">
      <slot name="leftExtend"></slot>
      <el-dropdown
        v-if="leftTabs.length > 0"
        trigger="click"
        @command="handleCommand"
        class="workspace-tabs-nav-dropdown"
      >
        <div class="workspace-tabs-nav tm-pointer">
          <i class="mdi mdi-menu-down"></i>
        </div>
        <el-dropdown-menu slot="dropdown">
          <template v-for="(one, index) in leftTabs">
            <el-dropdown-item :key="index" :command="one">
              {{ one.name || one.title }}
            </el-dropdown-item>
          </template>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <div class="workspace-tabs-body" ref="workspaceTabsBody">
      <template v-for="(one, index) in mainTabs">
        <div
          :key="index"
          class="workspace-tabs-one"
          :title="one.title"
          :class="{ active: one == itemsWorker.activeItem }"
          @contextmenu.prevent="tabContextmenu(one)"
          @mouseup="tabMouseup(one)"
        >
          <span class="workspace-tabs-one-text" @click="toSelectTab(one)">
            <template v-if="tool.isNotEmpty(one.iconFont)">
              <IconFont class="mgr-5" :class="one.iconFont"> </IconFont>
            </template>
            {{ one.name || one.title }}
          </span>
          <span
            class="workspace-tabs-delete-btn tm-pointer color-orange"
            @click="toDeleteTab(one)"
            title="关闭"
          >
            <Icon class="mdi-close"></Icon>
          </span>
        </div>
      </template>
    </div>
    <div class="workspace-tabs-right">
      <el-dropdown
        v-if="rightTabs.length > 0"
        trigger="click"
        @command="handleCommand"
        placement="bottom-start"
        class="workspace-tabs-nav-dropdown"
      >
        <div class="workspace-tabs-nav tm-pointer">
          <i class="mdi mdi-menu-down"></i>
        </div>
        <el-dropdown-menu slot="dropdown">
          <template v-for="(one, index) in rightTabs">
            <el-dropdown-item :key="index" :command="one">
              {{ one.name || one.title }}
            </el-dropdown-item>
          </template>
        </el-dropdown-menu>
      </el-dropdown>
      <slot name="rightExtend"></slot>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "itemsWorker"],
  data() {
    return {
      mainTabs: [],
      leftTabs: [],
      rightTabs: [],
      tabsWidth: 0,
      tabHeaderWidthCount: 0,
      tabEditorHeaderWidth: 0,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "itemsWorker.items"() {
      this.$nextTick(() => {
        this.initTabs();
      });
    },
    "itemsWorker.activeItem"() {
      this.$nextTick(() => {
        this.initTabs();
      });
    },
  },
  methods: {
    init() {
      this.$nextTick(() => {
        this.initTabs();
      });
    },
    initTabs() {
      let leftTabs = [];
      let rightTabs = [];
      let mainTabs = [];

      this.itemsWorker.items.forEach((one) => {
        if (!one.show) {
          return;
        }
        mainTabs.push(one);
      });
      this.mainTabs = mainTabs;

      this.$nextTick(() => {
        let tabsWidth = this.tool.jQuery(this.$el).width();
        let workspaceTabsBody = this.tool.jQuery(this.$refs.workspaceTabsBody);
        let tabsBodyWidth = workspaceTabsBody.width();
        let tabWidthCount = 0;
        let children = workspaceTabsBody.children();
        this.mainTabs.forEach((one, index) => {
          let tabWidth = this.tool.jQuery(children[index]).width();
          one._tabWidth = tabWidth;
          tabWidthCount += Number(tabWidth);
        });
        let scrollLeft = 0;
        let activeIndex = this.getActiveIndex();
        let showWidth = 0;
        if (tabsBodyWidth < tabWidthCount) {
          let showIndex = 0;
          this.mainTabs.forEach((one) => {
            if (showIndex < activeIndex - 2) {
              scrollLeft += Number(one._tabWidth);
              leftTabs.push(one);
            } else {
              showWidth += Number(one._tabWidth);
            }
            if (showWidth > tabsBodyWidth) {
              rightTabs.push(one);
            }
            showIndex++;
          });
        } else {
          this.mainTabs.forEach((one, index) => {
            showWidth += Number(one._tabWidth);
          });
        }
        this.leftTabs = leftTabs;
        this.rightTabs = rightTabs;
        if (scrollLeft > 0) {
          scrollLeft += 2;
          if (scrollLeft > tabWidthCount - showWidth) {
            scrollLeft = tabWidthCount - showWidth;
          }
          scrollLeft = parseInt(scrollLeft);
        }
        workspaceTabsBody.scrollLeft(scrollLeft);
      });
    },
    resize() {
      this.$nextTick(() => {
        this.initTabs();
      });
    },
    getActiveIndex() {
      return this.mainTabs.indexOf(this.itemsWorker.activeItem);
    },
    handleCommand(tab) {
      this.toSelectTab(tab);
    },
    tabMouseup(tab) {
      if (window.event && window.event.button != 1) {
        return;
      }
      if (tab == null) {
        return;
      }
      this.itemsWorker.toRemoveItem(tab);
    },
    tabContextmenu(tab) {
      if (tab == null) {
        return;
      }
      let menus = [];

      menus.push({
        header: tab.name,
      });
      if (this.hasOpenNewWindow) {
        menus.push({
          text: "新窗口打开",
          onClick: () => {
            this.openNewWindow && this.openNewWindow(tab);
          },
        });
      }
      menus.push({
        text: "关闭",
        onClick: () => {
          this.toDeleteTab(tab);
        },
      });
      menus.push({
        text: "打开新标签",
        onClick: () => {
          this.itemsWorker.toCopyItem(tab);
        },
      });
      menus.push({
        text: "选择",
        onClick: () => {
          this.toSelectTab(tab);
        },
      });
      menus.push({
        text: "关闭其它",
        onClick: () => {
          this.itemsWorker.toDeleteOtherItem(tab);
        },
      });
      menus.push({
        text: "关闭所有",
        onClick: () => {
          this.itemsWorker.toRemoveAll();
        },
      });

      this.tool.showContextmenu(menus);
    },
    toSelectTab(tab) {
      this.itemsWorker.toActiveItem(tab);
    },
    toDeleteTab(tab) {
      this.itemsWorker.toRemoveItem(tab);
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
    window.addEventListener("resize", this.resize);
  },
  destroyed() {
    window.removeEventListener("resize", this.resize);
  },
};
</script>

<style >
.workspace-tabs {
  width: 100%;
  height: 100%;
  font-size: 14px;
  position: relative;
  user-select: none;
  display: flex;
  border-bottom: 1px solid #4e4e4e;
}

.workspace-tabs-body {
  display: flex;
  position: relative;
  overflow: hidden;
}
.workspace-tabs-left,
.workspace-tabs-right {
  height: 100%;
  position: relative;
  padding: 0px 0px;
  background: #2a2a2a;
  display: flex;
}
.workspace-tabs-nav-dropdown {
  height: 100%;
  align-items: center;
}
.workspace-tabs-nav {
  height: 100%;
  align-items: center;
  display: flex;
  white-space: nowrap;
}
.workspace-tabs-left .workspace-tabs-nav {
  border-right: 1px solid #404040;
}
.workspace-tabs-right .workspace-tabs-nav {
  border-left: 1px solid #404040;
}
.workspace-tabs-nav .mdi {
  font-size: 16px;
}
.workspace-tabs-one {
  display: flex;
  border-right: 1px solid #4e4e4e;
  position: relative;
  border-top-left-radius: 0px;
  border-top-right-radius: 10px;
  align-items: center;
}
.workspace-tabs-one.active {
  background-color: #172029;
}
.workspace-tabs-one .workspace-tabs-one-text {
  display: flex;
  height: 100%;
  align-items: center;
  padding: 0px 5px 0px 10px;
  cursor: pointer;
  white-space: nowrap;
}
.workspace-tabs-one .workspace-tabs-delete-btn {
  transition: all 0.1s;
  transform: scale(0);
  margin: 0px 5px;
}
.workspace-tabs-one:hover .workspace-tabs-delete-btn {
  transform: scale(1);
}
</style>
