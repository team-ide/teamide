<template>
  <div
    class="toolbox-context-box"
    :class="{ 'toolbox-context-box-show': showBox }"
  >
    <div class="toolbox-context-box-header">
      <div style="text-align: right">
        <span title="关闭" class="tm-link color-write mgr-10" @click="hide">
          <i class="mdi mdi-close ft-21"></i>
        </span>
      </div>
    </div>
    <div class="toolbox-context-body" v-if="selectGroup != null">
      <div class="toolbox-group-box scrollbar">
        <template v-for="group in groupList">
          <div
            :key="group.groupId"
            class="toolbox-group-one"
            :class="{ active: group.groupId == selectGroup.groupId }"
            @click="toSelectGroup(group)"
            @contextmenu="groupContextmenu(group)"
          >
            <div class="toolbox-group-title">
              <div class="toolbox-group-title-text">
                {{ group.name }}
              </div>
            </div>
          </div>
        </template>
        <div class="pd-10">
          <span
            class="tm-link color-green"
            title="新增分组"
            @click="toolbox.toInsertGroup()"
          >
            新增分组
          </span>
        </div>
      </div>

      <div class="toolbox-type-box scrollbar" v-if="searchMap != null">
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
                  selectGroup.context[toolboxType.name] == null ||
                  selectGroup.context[toolboxType.name].length == 0
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
                <template
                  v-for="toolboxData in selectGroup.context[toolboxType.name]"
                >
                  <div
                    :key="toolboxData.toolboxId"
                    v-if="
                      tool.isEmpty(searchMap[toolboxType.name]) ||
                      toolboxData.name
                        .toLowerCase()
                        .indexOf(searchMap[toolboxType.name].toLowerCase()) >= 0
                    "
                    class="toolbox-type-data"
                    @contextmenu="dataContextmenu(toolboxType, toolboxData)"
                  >
                    <span
                      class="toolbox-type-data-text tm-link color-grey"
                      title="打开"
                      @click="toolbox.toolboxDataOpen(toolboxData)"
                    >
                      {{ toolboxData.name }}
                    </span>
                    <div class="toolbox-type-data-btn-box"></div>
                  </div>
                </template>
              </template>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox", "context", "groups"],
  data() {
    return {
      showBox: false,
      searchMap: null,
      selectGroup: null,
      groupList: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    context() {
      this.initGroup();
    },
    groups() {
      this.initGroup();
    },
  },
  methods: {
    init() {
      this.initGroup();
      let searchMap = {};
      this.toolbox.types.forEach((one) => {
        searchMap[one.name] = "";
      });
      this.searchMap = searchMap;
      this.initGroup();
    },
    toSelectGroup(group) {
      if (group == null) {
        group = this.groupList[0];
      }
      this.selectGroup = group;
    },
    initGroup() {
      let groupList = [];
      let context = this.context || {};
      let groups = this.groups || [];
      groupList.push({
        groupId: null,
        name: "未分组",
        context: {},
      });
      groups.forEach((one) => {
        groupList.push({
          groupId: one.groupId,
          name: one.name,
          context: {},
        });
      });
      let selectGroup = groupList[0];
      if (this.selectGroup && this.selectGroup.groupId != null) {
        groupList.forEach((one) => {
          if (one.groupId == this.selectGroup.groupId) {
            selectGroup = one;
          }
        });
      }
      this.toolbox.types.forEach((type) => {
        let list = context[type.name];
        groupList.forEach((one) => {
          let groupToolboxList = [];
          list.forEach((tOne) => {
            if (
              this.tool.isEmpty(one.groupId) &&
              this.tool.isEmpty(tOne.groupId)
            ) {
              groupToolboxList.push(tOne);
            } else if (one.groupId == tOne.groupId) {
              groupToolboxList.push(tOne);
            }
          });
          one.context[type.name] = groupToolboxList;
        });
      });
      this.groupList = groupList;
      this.selectGroup = selectGroup;
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
    groupContextmenu(group) {
      let menus = [];
      menus.push({
        header: group.name,
      });
      menus.push({
        text: "修改",
        onClick: () => {
          this.toolbox.toUpdateGroup(group);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toolbox.toDeleteGroup(group);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    dataContextmenu(toolboxType, toolboxData) {
      let menus = [];
      menus.push({
        header: toolboxType.text + ":" + toolboxData.name,
      });
      menus.push({
        text: "打开",
        onClick: () => {
          this.toolbox.toolboxDataOpen(toolboxData);
        },
      });
      if (toolboxType.name == "ssh") {
        menus.push({
          text: "打开FTP",
          onClick: () => {
            this.toolbox.toolboxDataOpenSfpt(toolboxData);
          },
        });
      }
      if (this.groupList.length > 0) {
        let moveGroupMenu = {
          text: "移动分组",
          menus: [],
        };
        menus.push(moveGroupMenu);
        this.groupList.forEach((one) => {
          moveGroupMenu.menus.push({
            text: one.name,
            onClick: () => {
              this.toolbox.moveGroup(toolboxData.toolboxId, one.groupId);
            },
          });
        });
      }
      menus.push({
        text: "修改",
        onClick: () => {
          this.toolbox.toUpdate(toolboxType, toolboxData);
        },
      });
      menus.push({
        text: "复制",
        onClick: () => {
          this.toolbox.toCopy(toolboxType, toolboxData);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toolbox.toDelete(toolboxType, toolboxData);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
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
.toolbox-context-box {
  position: absolute;
  top: 25px;
  width: 100%;
  background: #0f1b26;
  transition: all 0s;
  transform: scale(0);
  height: calc(100% - 25px);
}
.toolbox-context-box.toolbox-context-box-show {
  transform: scale(1);
}
.toolbox-context-box-header {
  height: 30px;
}
.toolbox-context-body {
  width: 100%;
  height: calc(100% - 30px);
  display: flex;
}
.toolbox-context-box .toolbox-group-box {
  width: 200px;
  font-size: 12px;
  height: 100%;
}

.toolbox-context-box .toolbox-group-one {
  /* width: calc(25% - 12.5px); */
  width: 100%;
  cursor: pointer;
  margin-top: 10px;
}
.toolbox-context-box .toolbox-group-title {
  padding: 0px 10px;
  background: #15222d;
  color: #ffffff;
  line-height: 30px;
  display: flex;
}
.toolbox-context-box .toolbox-group-one.active .toolbox-group-title {
  background: #3f4e5d;
}
.toolbox-context-box .toolbox-group-one:hover .toolbox-group-title {
  background: #2f3f4f;
}
.toolbox-context-box .toolbox-group-title-text {
  flex: 1;
}
.toolbox-context-box .toolbox-group-title .icon {
  margin-right: 5px;
}
.toolbox-context-box .toolbox-group-title .tm-link {
  padding: 0px;
}

.toolbox-context-box .toolbox-type-box {
  flex: 1;
  font-size: 12px;
  height: 100%;
}
.toolbox-context-box .toolbox-type-box:after {
  content: "";
  display: table;
  clear: both;
}
.toolbox-context-box .toolbox-type-one {
  /* width: calc(25% - 12.5px); */
  width: 290px;
  float: left;
  margin: 0px 0px 10px 10px;
}
.toolbox-context-box .toolbox-type-title {
  padding: 0px 10px;
  background: #2b3f51;
  color: #ffffff;
  line-height: 23px;
  display: flex;
}
.toolbox-context-box .toolbox-type-title-text {
  flex: 1;
}
.toolbox-context-box .toolbox-type-title .icon {
  margin-right: 5px;
}
.toolbox-context-box .toolbox-type-title .tm-link {
  padding: 0px;
}
.toolbox-context-box .toolbox-type-title .toolbox-type-data-search-box {
  width: 100px;
}
.toolbox-context-box
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
.toolbox-context-box .toolbox-type-data-box {
  background: #1b2a38;
  padding: 5px 10px;
  padding-right: 0px;
  height: 250px;
}
.toolbox-context-box .toolbox-type-data {
  display: flex;
  overflow: hidden;
  padding: 2px 0px;
}
.toolbox-context-box .toolbox-type-data-text {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
  flex: 1;
}
.toolbox-context-box .toolbox-type-data .tm-link {
  padding: 0px;
}
.toolbox-context-box .toolbox-type-data-btn-box {
  display: inline-block;
  text-align: right;
  width: 85px;
}
.toolbox-context-box .toolbox-type-data-btn-box .tm-link {
  margin-right: 5px;
}
</style>
