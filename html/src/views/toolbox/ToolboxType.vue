<template>
  <div
    class="toolbox-context-box"
    :class="{ 'toolbox-context-box-show': showBox }"
  >
    <div class="toolbox-context-box-header">
      <div class="color-white ft-16 pdl-10 pdt-5">
        工具箱
        <span class="color-orange mgl-20 ft-12"
          >请右击进行操作，点击某个分组，查看一分组的各类工具</span
        >
      </div>
      <div
        style="display: inline-block; position: absolute; top: 0px; right: 0px"
      >
        <span title="关闭" class="tm-link color-write mgr-0" @click="hide">
          <i class="mdi mdi-close ft-21"></i>
        </span>
      </div>
    </div>
    <div class="toolbox-context-body" v-if="selectGroup != null">
      <div class="toolbox-group-box">
        <div class="toolbox-group-header">
          <div class="toolbox-group-header-text">分组</div>
          <span
            class="tm-link color-green mgl-5 mgt-5"
            title="新增分组"
            @click="toInsertGroup()"
          >
            <i class="mdi mdi-plus ft-14"></i>
          </span>
          <div class="toolbox-group-header-search-box mgl-10">
            <input
              class="toolbox-group-header-search"
              v-model="searchGroup"
              placeholder="输入过滤"
            />
          </div>
        </div>
        <div class="toolbox-group-body scrollbar">
          <template v-for="group in groupList">
            <div
              :key="group.groupId"
              class="toolbox-group-one"
              :class="{ active: group.groupId == selectGroup.groupId }"
              v-if="
                tool.isEmpty(searchGroup) ||
                group.name.toLowerCase().indexOf(searchGroup.toLowerCase()) >= 0
              "
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
        </div>
      </div>

      <div class="toolbox-type-box scrollbar" v-if="searchMap != null">
        <template v-for="toolboxType in toolboxTypes">
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
                  v-if="toolboxType.name != 'other'"
                  @click="toInsert(toolboxType, selectGroup)"
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
                  v-if="toolboxType.name != 'other'"
                  @click="toInsert(toolboxType, selectGroup)"
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
                    @click="toolboxDataOpen(toolboxData)"
                  >
                    <span class="toolbox-type-data-text" title="打开">
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
    <FormDialog
      ref="InsertToolbox"
      :source="source"
      title="新增Toolbox"
      :onSave="doInsert"
    ></FormDialog>
    <FormDialog
      ref="UpdateToolbox"
      :source="source"
      title="编辑Toolbox"
      :onSave="doUpdate"
    ></FormDialog>
    <FormDialog
      ref="InsertToolboxGroup"
      :source="source"
      title="新增工具分组"
      :onSave="doInsertGroup"
    ></FormDialog>
    <FormDialog
      ref="UpdateToolboxGroup"
      :source="source"
      title="编辑工具分组"
      :onSave="doUpdateGroup"
    ></FormDialog>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["source", "openByToolboxId"],
  data() {
    return {
      showBox: false,
      searchMap: null,
      selectGroup: null,
      groupList: [],
      searchGroup: null,
      toolboxTypes: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    showBox() {
      if (this.showBox) {
        this.initData();
      }
    },
  },
  methods: {
    show() {
      this.showBox = true;
    },
    showSwitch() {
      this.showBox = !this.showBox;
    },
    hide() {
      this.showBox = false;
    },
    init() {},
    async initData() {
      await this.source.initToolboxGroups();
      if (this.toolboxTypes.length == 0) {
        this.toolboxTypes = this.source.toolboxTypes || [];
      }
      await this.initToolboxDataGroup();
      if (this.searchMap == null) {
        let searchMap = {};
        this.toolboxTypes.forEach((one) => {
          searchMap[one.name] = "";
        });
        this.searchMap = searchMap;
      }
    },
    toSelectGroup(group) {
      if (group == null) {
        group = this.groupList[0];
      }
      this.selectGroup = group;
    },
    async initToolboxDataGroup() {
      let groupList = [];
      let toolboxList = [];

      let res = await this.server.toolbox.list({});
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        let data = res.data || {};
        toolboxList = data.toolboxList || [];
      }

      let context = {};

      toolboxList.forEach((one) => {
        context[one.toolboxType] = context[one.toolboxType] || [];
        context[one.toolboxType].push(one);
      });

      let groups = this.source.toolboxGroups || [];
      groupList.push({
        groupId: null,
        name: "未分组",
        context: {
          other: context["other"] || [],
        },
      });
      groups.forEach((one) => {
        groupList.push({
          groupId: one.groupId,
          name: one.name,
          context: {
            other: context["other"] || [],
          },
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
      this.toolboxTypes.forEach((type) => {
        if (type.name == "other") {
          return;
        }
        let list = context[type.name] || [];
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
    groupContextmenu(group) {
      let menus = [];
      menus.push({
        header: group.name,
      });
      menus.push({
        text: "修改",
        onClick: () => {
          this.toUpdateGroup(group);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toDeleteGroup(group);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    async moveGroup(toolboxId, groupId) {
      let res = await this.server.toolbox.moveGroup({
        toolboxId: toolboxId,
        groupId: groupId,
      });
      if (res.code == 0) {
        this.tool.success("移动成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    toolboxDataOpen(toolboxData) {
      let extend = {};
      if (toolboxData.toolboxType == "other") {
        extend = this.tool.getOptionJSON(toolboxData.option);
      }
      this.openByToolboxId(toolboxData.toolboxId, extend);
      this.hide();
    },
    toolboxDataOpenSfpt(toolboxData) {
      this.openByToolboxId(toolboxData.toolboxId, { isFTP: true });
      this.hide();
    },
    dataContextmenu(toolboxType, toolboxData) {
      if (toolboxType.name == "other") {
        return;
      }
      let menus = [];
      menus.push({
        header: toolboxType.text + ":" + toolboxData.name,
      });
      menus.push({
        text: "打开",
        onClick: () => {
          this.toolboxDataOpen(toolboxData);
        },
      });
      if (toolboxType.name == "ssh") {
        menus.push({
          text: "打开FTP",
          onClick: () => {
            this.toolboxDataOpenSfpt(toolboxData);
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
              this.moveGroup(toolboxData.toolboxId, one.groupId);
            },
          });
        });
      }
      menus.push({
        text: "修改",
        onClick: () => {
          this.toUpdate(toolboxType, toolboxData);
        },
      });
      menus.push({
        text: "复制",
        onClick: () => {
          this.toCopy(toolboxType, toolboxData);
        },
      });
      menus.push({
        text: "删除",
        onClick: () => {
          this.toDelete(toolboxType, toolboxData);
        },
      });

      if (menus.length > 0) {
        this.tool.showContextmenu(menus);
      }
    },
    toInsert(toolboxType, selectGroup) {
      this.tool.stopEvent();
      let toolboxData = {};
      let optionsJSON = {};

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
        selectGroup,
      });
    },
    toCopy(toolboxType, copy) {
      this.tool.stopEvent();
      let toolboxData = {};
      Object.assign(toolboxData, copy);
      delete toolboxData.toolboxId;
      toolboxData.name = toolboxData.name + " Copy";

      let optionsJSON = this.tool.getOptionJSON(toolboxData.option);

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    toUpdate(toolboxType, toolboxData) {
      this.tool.stopEvent();
      this.updateData = toolboxData;

      let optionsJSON = this.tool.getOptionJSON(toolboxData.option);

      this.$refs.UpdateToolbox.show({
        title: `编辑[${toolboxType.text}][${toolboxData.name}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    toDelete(toolboxType, toolboxData) {
      this.tool.stopEvent();
      this.tool
        .confirm(
          "删除[" +
            toolboxType.text +
            "]工具[" +
            toolboxData.name +
            "]将无法恢复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(toolboxType, toolboxData);
        })
        .catch((e) => {});
    },
    async doDelete(toolboxType, toolboxData) {
      let res = await this.server.toolbox.delete(toolboxData);
      if (res.code == 0) {
        this.initData();
        this.tool.success("删除成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdate(dataList, config) {
      let toolboxData = dataList[0];
      let optionJSON = dataList[1];
      let toolboxType = config.toolboxType;
      toolboxData.toolboxType = toolboxType.name;
      toolboxData.toolboxId = this.updateData.toolboxId;
      toolboxData.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.update(toolboxData);
      if (res.code == 0) {
        this.initData();
        this.tool.success("修改成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsert(dataList, config) {
      let toolboxData = dataList[0];
      let optionJSON = dataList[1];
      let toolboxType = config.toolboxType;
      toolboxData.toolboxType = toolboxType.name;
      if (config.selectGroup) {
        toolboxData.groupId = config.selectGroup.groupId;
      }
      toolboxData.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.insert(toolboxData);
      if (res.code == 0) {
        this.initData();
        this.tool.success("新增成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },

    toInsertGroup() {
      this.tool.stopEvent();
      let data = {};
      let optionsJSON = {};

      this.$refs.InsertToolboxGroup.show({
        title: `新增工具分组`,
        form: [this.form.toolbox.group, this.form.toolbox.group.option],
        data: [data, optionsJSON],
      });
    },
    toUpdateGroup(data) {
      this.tool.stopEvent();
      this.updateGroupData = data;

      let optionsJSON = this.tool.getOptionJSON(data.option);

      this.$refs.UpdateToolboxGroup.show({
        title: `编辑[${data.name}]工具分组`,
        form: [this.form.toolbox.group, this.form.toolbox.group.option],
        data: [data, optionsJSON],
      });
    },
    toDeleteGroup(data) {
      this.tool.stopEvent();
      this.tool
        .confirm(
          "删除工具分组[" + data.name + "]并将该分组下工具移出该组，确定删除？"
        )
        .then(async () => {
          return this.doDeleteGroup(data);
        })
        .catch((e) => {});
    },
    async doDeleteGroup(data) {
      let res = await this.server.toolbox.group.delete(data);
      if (res.code == 0) {
        this.tool.success("删除分组成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doUpdateGroup(dataList, config) {
      let data = dataList[0];
      let optionJSON = dataList[1];
      data.groupId = this.updateGroupData.groupId;
      data.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.group.update(data);
      if (res.code == 0) {
        this.tool.success("修改分组成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
    async doInsertGroup(dataList, config) {
      let data = dataList[0];
      let optionJSON = dataList[1];
      data.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.group.insert(data);
      if (res.code == 0) {
        this.tool.success("新增分组成功");
        this.initData();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
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
  color: #ffffff;
}
.toolbox-context-box.toolbox-context-box-show {
  transform: scale(1);
}
.toolbox-context-box-header {
  height: 40px;
}
.toolbox-context-body {
  width: 100%;
  height: calc(100% - 40px);
  display: flex;
}
.toolbox-context-box .toolbox-group-box {
  width: 200px;
  font-size: 12px;
  height: 100%;
  background: #1b2a38;
}

.toolbox-context-box .toolbox-group-header {
  padding: 0px 10px;
  height: 30px;
  display: flex;
}
.toolbox-context-box .toolbox-group-header-text {
  flex: 1;
  line-height: 30px;
}
.toolbox-context-box .toolbox-group-header .icon {
  margin-right: 5px;
}
.toolbox-context-box .toolbox-group-header .tm-link {
  padding: 0px;
}

.toolbox-context-box .toolbox-group-header .toolbox-group-header-search-box {
  width: 100px;
}
.toolbox-context-box
  .toolbox-group-header-search-box
  .toolbox-group-header-search {
  width: 100%;
  height: 26px;
  line-height: 26px;
  margin-top: 4px;
  border: 1px solid #767676;
  font-size: 12px;
  background: transparent;
}

.toolbox-context-box .toolbox-group-body {
  /* width: calc(25% - 12.5px); */
  width: 100%;
  height: calc(100% - 30px);
  padding-left: 8px;
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
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
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
  background: transparent;
}
.toolbox-context-box .toolbox-type-data-box {
  background: #1b2a38;
  padding: 5px 0px;
  padding-right: 0px;
  height: 250px;
}
.toolbox-context-box .toolbox-type-data {
  display: flex;
  overflow: hidden;
  padding: 2px 0px;
  cursor: pointer;
}
.toolbox-context-box .toolbox-type-data:hover {
  background: #2f3f4f;
}
.toolbox-context-box .toolbox-type-data-text {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: left;
  flex: 1;
  padding: 3px 10px;
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
