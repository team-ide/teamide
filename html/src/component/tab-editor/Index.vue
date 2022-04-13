<template>
  <div class="tab-editor">
    <div class="tab-editor-header">
      <div class="tab-header-box">
        <template v-for="one in tabs">
          <div
            :key="one.key"
            class="tab-header"
            :title="one.title"
            :class="{ active: one.active }"
          >
            <template v-if="$slots.tab">
              <slot name="tab" :tab="one"></slot>
            </template>
            <template v-else>
              <span class="text" @click="toSelectTab(one)">
                {{ one.name || one.title }}
              </span>
              <span
                class="tab-delete-btn tm-pointer color-orange"
                @click="toDeleteTab(one)"
                title="关闭"
              >
                <b-icon icon="x"></b-icon>
              </span>
            </template>
          </div>
        </template>
      </div>
    </div>
    <div class="tab-editor-body">
      <div class="tab-body-box">
        <template v-for="one in tabs">
          <div :key="one.key" class="tab-body" :class="{ active: one.active }">
            <slot name="body" :tab="one"></slot>
          </div>
        </template>
      </div>
    </div>
    <div class="tab-editor-footer">
      <template v-if="activeTab != null">
        <span class="ft-12 pdlr-10">{{ activeTab.title }}</span>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "onActive"],
  data() {
    return {
      activeTab: null,
      tabs: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    toSelectTab(tab) {
      this.doActiveTab(tab);
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
    addTab(tab) {
      let find = this.getTab(tab);
      if (find != null) {
        return;
      }
      this.tabs.push(tab);
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
  mounted() {},
};
</script>

<style >
.tab-editor {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #383838;
  color: #ffffff;
}
.tab-editor-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 14px;
  position: relative;
  user-select: none;
}
.tab-editor-body {
  width: 100%;
  height: calc(100% - 25px - 20px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
.tab-editor-footer {
  width: 100%;
  height: 20px;
  line-height: 20px;
  position: relative;
  display: flex;
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
  background-color: #2d2d2d;
}
.tab-header .text {
  padding: 0px 5px 0px 10px;
  cursor: pointer;
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
