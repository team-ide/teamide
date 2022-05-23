<template>
  <div
    class="toolbox-dropdown-box scrollbar"
    :class="{ 'toolbox-dropdown-box-show': showBox }"
  >
    <div class="toolbox-dropdown-box-header">
      <div style="text-align: right">
        <span title="关闭" class="tm-link color-write mgr-10" @click="hide">
          <i class="mdi mdi-close ft-22"></i>
        </span>
      </div>
    </div>
    <div class="toolbox-type-box" v-if="searchMap != null">
      <template v-for="toolboxType in toolbox.types">
        <div :key="toolboxType.name" class="toolbox-type-one">
          <div class="toolbox-type-title">
            <div class="toolbox-type-title-text">
              <template v-if="toolboxType.name == 'database'">
                <IconFont class="teamide-database"> </IconFont>
              </template>
              <template v-else-if="toolboxType.name == 'redis'">
                <IconFont class="teamide-redis"> </IconFont>
              </template>
              <template v-else-if="toolboxType.name == 'elasticsearch'">
                <IconFont class="teamide-elasticsearch"> </IconFont>
              </template>
              <template v-else-if="toolboxType.name == 'kafka'">
                <IconFont class="teamide-kafka"> </IconFont>
              </template>
              <template v-else-if="toolboxType.name == 'zookeeper'">
                <IconFont class="teamide-zookeeper"> </IconFont>
              </template>
              <template v-else-if="toolboxType.name == 'ssh'">
                <IconFont class="teamide-ssh"> </IconFont>
                <IconFont class="teamide-ftp"> </IconFont>
              </template>
              <span class="toolbox-type-text">
                {{ toolboxType.text || toolboxType.name }}
              </span>
              <span
                class="tm-link color-green mgl-10"
                title="新增"
                @click="toolbox.toInsert(toolboxType)"
              >
                <i class="mdi mdi-plus ft-14"></i>
              </span>
            </div>
            <div class="toolbox-type-data-search-box">
              <input
                class="toolbox-type-data-search"
                v-model="searchMap[toolboxType.name]"
                placeholder="输入过滤"
              />
            </div>
          </div>
          <div class="toolbox-type-data-box scrollbar">
            <template
              v-if="
                context[toolboxType.name] == null ||
                context[toolboxType.name].length == 0
              "
            >
              <span
                class="tm-link color-green"
                title="新增"
                @click="toolbox.toInsert(toolboxType)"
              >
                新增
              </span>
            </template>
            <template v-else>
              <template v-for="toolboxData in context[toolboxType.name]">
                <div
                  :key="toolboxData.toolboxId"
                  v-if="
                    tool.isEmpty(searchMap[toolboxType.name]) ||
                    toolboxData.name
                      .toLowerCase()
                      .indexOf(searchMap[toolboxType.name].toLowerCase()) >= 0
                  "
                  class="toolbox-type-data"
                >
                  <span
                    class="toolbox-type-data-text tm-link color-grey"
                    title="打开"
                    @click="toolbox.toolboxDataOpen(toolboxData)"
                  >
                    {{ toolboxData.name }}
                  </span>
                  <div class="toolbox-type-data-btn-box">
                    <span
                      title="打开FTP"
                      v-if="toolboxType.name == 'ssh'"
                      class="tm-link color-green"
                      @click="toolbox.toolboxDataOpenSfpt(toolboxData)"
                    >
                      <IconFont
                        class="teamide-ftp ft-12"
                        style="vertical-align: -2px"
                      >
                      </IconFont>
                    </span>
                    <span
                      title="编辑"
                      class="tm-link color-grey"
                      @click="toolbox.toUpdate(toolboxType, toolboxData)"
                    >
                      <i class="mdi mdi-square-edit-outline ft-13"></i>
                    </span>
                    <span
                      title="复制"
                      class="tm-link color-grey"
                      @click="toolbox.toCopy(toolboxType, toolboxData)"
                    >
                      <i class="mdi mdi-content-copy ft-12"></i>
                    </span>
                    <span
                      title="删除"
                      class="tm-link color-orange-8"
                      @click="toolbox.toDelete(toolboxType, toolboxData)"
                    >
                      <i class="mdi mdi-delete-outline ft-14"></i>
                    </span>
                  </div>
                </div>
              </template>
            </template>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context"],
  data() {
    return {
      showBox: false,
      searchMap: null,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {
      let searchMap = {};
      this.toolbox.types.forEach((one) => {
        searchMap[one.name] = "";
      });
      this.searchMap = searchMap;
    },
    show() {
      this.showBox = true;
    },
    showSwitch() {
      this.showBox = !this.showBox;
    },
    hide() {
      this.showBox = false;
    },
  },
  created() {},
  updated() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-dropdown-box {
  position: absolute;
  top: 25px;
  width: 100%;
  background: #0f1b26;
  transition: all 0s;
  transform: scale(0);
  height: calc(100% - 25px);
}
.toolbox-dropdown-box.toolbox-dropdown-box-show {
  transform: scale(1);
}

.toolbox-dropdown-box .toolbox-type-box {
  width: 100%;
  font-size: 12px;
}
.toolbox-dropdown-box .toolbox-type-box:after {
  content: "";
  display: table;
  clear: both;
}
.toolbox-dropdown-box .toolbox-type-one {
  /* width: calc(25% - 12.5px); */
  width: 270px;
  float: left;
  margin: 0px 0px 10px 10px;
}
.toolbox-dropdown-box .toolbox-type-title {
  padding: 0px 10px;
  background: #2b3f51;
  color: #ffffff;
  line-height: 23px;
  display: flex;
}
.toolbox-dropdown-box .toolbox-type-title-text {
  flex: 1;
}
.toolbox-dropdown-box .toolbox-type-title .icon {
  margin-right: 5px;
}
.toolbox-dropdown-box .toolbox-type-title .tm-link {
  padding: 0px;
}
.toolbox-dropdown-box .toolbox-type-title .toolbox-type-data-search-box {
  width: 100px;
}
.toolbox-dropdown-box
  .toolbox-type-title
  .toolbox-type-data-search-box
  .toolbox-type-data-search {
  width: 100%;
  height: 20px;
  line-height: 20px;
  margin-top: 4px;
  border: 1px solid #767676;
  font-size: 12px;
}
.toolbox-dropdown-box .toolbox-type-data-box {
  background: #1b2a38;
  padding: 5px 10px;
  padding-right: 0px;
  height: 250px;
}
.toolbox-dropdown-box .toolbox-type-data {
  display: flex;
  overflow: hidden;
  padding: 2px 0px;
}
.toolbox-dropdown-box .toolbox-type-data-text {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
  flex: 1;
}
.toolbox-dropdown-box .toolbox-type-data .tm-link {
  padding: 0px;
}
.toolbox-dropdown-box .toolbox-type-data-btn-box {
  display: inline-block;
  text-align: right;
  width: 85px;
}
.toolbox-dropdown-box .toolbox-type-data-btn-box .tm-link {
  margin-right: 5px;
}
</style>
