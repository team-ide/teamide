<template>
  <div class="tab-editor">
    <div class="tab-editor-header" ref="tabEditorHeader">
      <div class="tab-header-left" v-if="leftTabs.length > 0">
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="tab-header-nav tm-pointer">
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
      <div class="tab-header-box" ref="headerBox">
        <template v-for="one in tabs">
          <div
            :ref="'tab:' + one.key"
            :key="one.key"
            :tab-key="one.key"
            class="tab-header"
            :title="one.title"
            :class="{ active: one.active }"
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
              class="tab-delete-btn tm-pointer color-orange"
              @click="toDeleteTab(one)"
              title="关闭"
            >
              <Icon class="mdi-close"></Icon>
            </span>
          </div>
        </template>
      </div>
      <div class="tab-header-right" v-if="rightTabs.length > 0">
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="tab-header-nav tm-pointer">
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
      </div>
      <slot name="extend" ref="extend"></slot>
    </div>
    <div class="tab-editor-body">
      <div class="tab-body-box">
        <template v-for="one in tabs">
          <div :key="one.key" class="tab-body" :class="{ active: one.active }">
            <slot name="body" :tab="one"> </slot>
          </div>
        </template>
      </div>
    </div>
    <!-- <div class="tab-editor-footer">
      <template v-if="activeTab != null">
        <span class="ft-12 pdlr-10">{{ activeTab.title }}</span>
      </template>
    </div> -->
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
    "onOffsetRightDistance",
    "slotTab",
    "copyTab",
  ],
  data() {
    return {
      activeTab: null,
      leftTabs: [],
      showLeftTabs: false,
      rightTabs: [],
      showRightTabs: false,
      tabs: [],
      headerBoxWidth: 0,
      tabHeaderWidthCount: 0,
      tabEditorHeaderWidth: 0,
      extendWidth: "0px",
      leftWidth: "0px",
      rightWidth: "0px",
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    tabs() {
      this.leftTabs = [];
      this.rightTabs = [];
      this.$nextTick(() => {
        this.initHeader();
      });
    },
    activeTab() {
      this.leftTabs = [];
      this.rightTabs = [];
      this.$nextTick(() => {
        this.initHeader();
      });
      this.onActiveTabFocue();
    },
  },
  methods: {
    resize() {
      this.leftTabs = [];
      this.rightTabs = [];
      this.$nextTick(() => {
        this.initHeader();
      });
    },
    onFocus() {
      this.onActiveTabFocue();
    },
    onActiveTabFocue() {
      if (this.activeTab) {
        let tabBodySlot = this.getTabBodySlot(this.activeTab);
        if (tabBodySlot == null) {
          return;
        }
        tabBodySlot.onFocus && tabBodySlot.onFocus();
      }
    },
    getTabs() {
      return this.tabs;
    },
    getTabByTarget(target) {
      let res = null;
      if (target) {
        let box = this.tool.jQuery(target).closest(".tab-header");
        if (box.length > 0) {
          let key = box.attr("tab-key");
          res = this.getTab(key);
        }
      }
      return res;
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
    tabContextmenu(e) {
      e = e || window.event;
      let tab = this.getTabByTarget(e.target);
      if (tab == null) {
        return;
      }
      let menus = [];

      menus.push({
        header: tab.name,
      });
      menus.push({
        text: "关闭",
        onClick: () => {
          this.toDeleteTab(tab);
        },
      });
      menus.push({
        text: "打开新标签",
        onClick: () => {
          this.toCopyTab(tab);
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
          this.toDeleteOtherTab(tab);
        },
      });
      menus.push({
        text: "关闭所有",
        onClick: () => {
          this.toDeleteAll();
        },
      });

      this.tool.showContextmenu(menus);
    },
    handleCommand(tab) {
      this.toSelectTab(tab);
    },
    initHeader() {
      let leftTabs = [];
      let rightTabs = [];
      let scrollLeft = 0;
      let offsetRightDistance = 0;
      this.initTabHeaderWidth();
      let activeIndex = this.tabs.indexOf(this.activeTab);
      let showWidth = 0;
      if (this.headerBoxWidth < this.tabHeaderWidthCount) {
        this.tabs.forEach((one, index) => {
          if (index < activeIndex - 2) {
            scrollLeft += Number(one.headerWidth);
            leftTabs.push(one);
          } else {
            showWidth += Number(one.headerWidth);
          }
          if (showWidth > this.headerBoxWidth) {
            rightTabs.push(one);
          }
        });
      } else {
        this.tabs.forEach((one, index) => {
          showWidth += Number(one.headerWidth);
        });
      }
      this.leftTabs = leftTabs;
      this.rightTabs = rightTabs;
      if (scrollLeft > 0) {
        scrollLeft += 2;
        if (scrollLeft > this.tabHeaderWidthCount - showWidth) {
          scrollLeft = this.tabHeaderWidthCount - showWidth;
        }
        scrollLeft = parseInt(scrollLeft);
      }
      this.tool.jQuery(this.$refs.headerBox).scrollLeft(scrollLeft);

      offsetRightDistance =
        this.tabEditorHeaderWidth - this.tabHeaderWidthCount;
      this.onOffsetRightDistance &&
        this.onOffsetRightDistance(offsetRightDistance);
    },
    initTabHeaderWidth() {
      this.tabEditorHeaderWidth = this.tool
        .jQuery(this.$refs.tabEditorHeader)
        .width();
      this.headerBoxWidth = this.tool.jQuery(this.$refs.headerBox).width();
      this.tabHeaderWidthCount = 0;
      this.tabs.forEach((one) => {
        if (one.headerWidth == null) {
          let el = this.$refs["tab:" + one.key];
          if (!el) {
            return;
          }
          if (el.length > 0) {
            el = el[0];
          }
          one.headerWidth = this.tool.jQuery(el).width();
        }
        this.tabHeaderWidthCount += Number(one.headerWidth);
      });
    },
    toSelectTab(tab) {
      this.doActiveTab(tab);
    },
    getTabBodySlot(tab) {
      let slot = null;
      this.$children.forEach((one) => {
        if (one.tab == tab) {
          slot = one;
        }
      });
      return slot;
    },
    getTab(tab) {
      let res = null;
      this.tabs.forEach((one) => {
        if (one == tab || one.key == tab || one.key == tab.key) {
          res = one;
        }
      });
      return res;
    },
    toDeleteTab(tab) {
      this.removeTab(tab);
    },
    addTab(tab, fromTab) {
      let find = this.getTab(tab);
      if (find != null) {
        return;
      }
      let fromIndex = this.tabs.indexOf(fromTab);
      if (fromIndex < 0) {
        this.tabs.push(tab);
      } else {
        this.tabs.splice(fromIndex + 1, 0, tab);
      }
    },
    toCopyTab(tab) {
      this.copyTab && this.copyTab(tab);
    },
    toDeleteOtherTab(tab) {
      this.toSelectTab(tab);
      let deleteTabs = [];
      this.tabs.forEach((one) => {
        if (one != tab) {
          deleteTabs.push(one);
        }
      });
      this.toDeleteTabs(deleteTabs);
    },
    toDeleteAll() {
      let deleteTabs = [];
      this.tabs.forEach((one) => {
        deleteTabs.push(one);
      });
      this.toDeleteTabs(deleteTabs);
    },
    toDeleteTabs(deleteTabs) {
      deleteTabs.forEach((one) => {
        let tabIndex = this.tabs.indexOf(one);
        this.tabs.splice(tabIndex, 1);
        this.onRemoveTab && this.onRemoveTab(one);
      });
    },
    removeTab(tab) {
      let find = this.getTab(tab);
      if (find == null) {
        return;
      }
      let tabIndex = this.tabs.indexOf(find);
      this.tabs.splice(tabIndex, 1);
      if (find.active) {
        let nextTabIndex = tabIndex - 1;
        if (nextTabIndex < 0) {
          nextTabIndex = 0;
        }
        this.doActiveTab(this.tabs[nextTabIndex]);
      }
      this.onRemoveTab && this.onRemoveTab(find);
    },
    doActiveTab(tab) {
      this.$nextTick(() => {
        tab = this.getTab(tab);
        this.tabs.forEach((one) => {
          if (one != tab) {
            one.active = false;
          }
        });
        this.tabs.forEach((one) => {
          if (one == tab) {
            one.active = true;
          }
        });
        this.activeTab = tab;
        if (tab != null) {
          this.onActiveTab && this.onActiveTab(tab);
        }
      });
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    window.addEventListener("resize", this.resize);
  },
  destroyed() {
    window.removeEventListener("resize", this.resize);
  },
};
</script>

<style >
.tab-editor {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #383838;
  color: #d9d9d9;
}
.tab-editor-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 12px;
  position: relative;
  user-select: none;
  display: flex;
  border-bottom: 1px solid #4e4e4e;
}
.tab-editor-body {
  width: 100%;
  height: calc(100% - 25px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.tab-editor-footer {
  width: 100%;
  height: 0px;
  line-height: 0px;
  position: relative;
  display: flex;
}
.tab-header-nav {
  font-size: 16px;
  color: #ffffff;
}
.tab-header-left {
  display: flex;
  position: relative;
  text-align: center;
  background-color: #2a3036;
  padding: 0px 5px;
}
.tab-header-extend {
  display: flex;
  position: relative;
  text-align: center;
  background-color: #2a3036;
  padding: 0px 5px;
}
.tab-header-right {
  display: flex;
  position: relative;
  text-align: center;
  background-color: #2a3036;
  padding: 0px 5px;
}
.tab-header-box {
  display: flex;
  position: relative;
  overflow: hidden;
}
.tab-header {
  display: flex;
  border-right: 1px solid #4e4e4e;
  position: relative;
  border-top-left-radius: 0px;
  border-top-right-radius: 10px;
}
.tab-header.active {
  background-color: #172029;
}
.tab-header .text {
  padding: 0px 5px 0px 10px;
  cursor: pointer;
  white-space: nowrap;
  line-height: 23px;
}
.tab-header .text .mdi {
  vertical-align: -1px;
}
.tab-header .tab-delete-btn {
  transition: all 0.1s;
  transform: scale(0);
  margin: 0px 5px;
}
.tab-header:hover .tab-delete-btn {
  transform: scale(1);
}
.tab-body-box {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #172029;
}
.tab-body {
  width: 100%;
  height: 100%;
  position: absolute;
  left: 0px;
  right: 0px;
  transition: all 0s;
  transform: scale(0);
}
.tab-body.active {
  transform: scale(1);
}
</style>
