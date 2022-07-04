<template>
  <div class="workspace-tabs">
    <div class="workspace-tabs-left" v-if="leftTabs.length > 0">
      <el-dropdown trigger="click" @command="handleCommand">
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
      <template v-for="one in mainTabs">
        <div
          :key="one.key"
          :tab-key="one.key"
          class="workspace-tabs-one"
          :title="one.title"
          :class="{ active: one == activeTab }"
          @contextmenu.prevent="tabContextmenu"
          @mouseup="tabMouseup"
        >
          <span
            class="text"
            @click="toSelectTab(one)"
            @dblclick="toCopyTab(one)"
          >
            <template v-if="slotTab">
              <slot name="tab" :tab="one"></slot>
            </template>
            <template v-else>
              {{ one.name || one.title }}
            </template>
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
    <div class="workspace-tabs-right" v-if="rightTabs.length > 0">
      <el-dropdown
        trigger="click"
        @command="handleCommand"
        size="mini"
        placement="bottom-start"
      >
        <div class="workspace-tabs-nav tm-pointer">
          <i class="mdi mdi-menu-down"></i>
        </div>
        <el-dropdown-menu slot="dropdown" class="workspace-tabs-right-dropdown">
          <template v-for="(one, index) in rightTabs">
            <el-dropdown-item :key="index" :command="one">
              {{ one.name || one.title }}
            </el-dropdown-item>
          </template>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <slot name="extend" ref="extend"></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: [
    "source",
    "onActive",
    "onRemoveTab",
    "onActiveTab",
    "slotTab",
    "copyTab",
    "hasOpenNewWindow",
    "openNewWindow",
    "tabs",
    "activeTab",
  ],
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
    tabs() {
      this.$nextTick(() => {
        this.initTabs();
      });
    },
    activeTab() {
      this.$nextTick(() => {
        this.initTabs();
      });
      this.onActiveTabFocue();
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

      this.tabs.forEach((one) => {
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
      return this.mainTabs.indexOf(this.activeTab);
    },
    onFocus() {
      this.onActiveTabFocue();
    },
    handleCommand(tab) {
      this.toSelectTab(tab);
    },
    tabMouseup(e) {
      e = e || window.event;
      if (e.button != 1) {
        return;
      }
      let tab = this.getTabByTarget(e.target);
      if (tab == null) {
        return;
      }
      this.toDeleteTab(tab);
    },
    toSelectTab(tab) {
      this.onActiveTab(tab);
    },
    toDeleteTab(tab) {
      this.onRemoveTab(tab);
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
  font-size: 12px;
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
.workspace-tabs-right-dropdown {
  left: auto !important;
  right: 0px;
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
.workspace-tabs-one .text {
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
