<template>
  <div class="toolbox-kafka-tabs">
    <template v-if="ready">
      <div class="toolbox-kafka-tabs-header">
        <div class="toolbox-tab-box">
          <template v-for="(one, index) in tabs">
            <div
              :key="'tab-' + index"
              class="toolbox-tab"
              :title="one.title"
              :class="{ active: one.active }"
            >
              <span class="text" @click="toSelectTab(one)">
                {{ one.name }}
              </span>
              <span
                class="delete-btn tm-pointer color-orange"
                @click="toDeleteTab(one)"
                title="关闭"
              >
                <b-icon icon="x"></b-icon>
              </span>
            </div>
          </template>
        </div>
      </div>
      <div class="toolbox-kafka-tabs-body">
        <div class="toolbox-tab-span-box">
          <template v-for="(one, index) in tabs">
            <div
              :key="'span-' + index"
              class="toolbox-tab-span"
              :class="{ active: one.active }"
            >
              <ToolboxKafkaTopicData
                :source="source"
                :toolbox="toolbox"
                :toolboxType="toolboxType"
                :data="data"
                :wrap="wrap"
                :topic="one.data"
              >
              </ToolboxKafkaTopicData>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "data", "toolboxType", "toolbox", "option", "wrap"],
  data() {
    return {
      ready: false,
      tabs: [],
      activeTab: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.wrap.doActiveTab = this.doActiveTab;
      this.wrap.createTab = this.createTab;
      this.wrap.createTabByData = this.createTabByData;
      this.wrap.getTabByData = this.getTabByData;
      this.wrap.addTab = this.addTab;
      this.ready = true;
    },
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
      });
    },
    getTabKeyByData(data) {
      let key = "" + data.name;

      return key;
    },
    getTabByData(data) {
      let key = this.getTabKeyByData(data);
      let tab = this.getTab(key);
      return tab;
    },
    createTabByData(data) {
      let key = this.getTabKeyByData(data);

      let tab = this.getTab(key);
      if (tab == null) {
        let title = data.name;
        let name = data.name;
        tab = {
          key,
          data,
          title,
          name,
        };
        tab.active = false;
        tab.changed = false;
      }
      return tab;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-kafka-tabs {
  width: 100%;
  height: 100%;
}
.toolbox-kafka-tabs-header {
  width: 100%;
  height: 25px;
  line-height: 25px;
  font-size: 14px;
  position: relative;
}
.toolbox-kafka-tabs-body {
  width: 100%;
  height: calc(100% - 25px);
  border-bottom: 1px solid #4e4e4e;
  position: relative;
}
</style>
