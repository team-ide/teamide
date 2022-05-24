<template>
  <div class="toolbox-box scrollbar-textarea" :style="boxStyleObject">
    <template v-if="ready">
      <tm-layout height="100%">
        <Main
          ref="main"
          v-if="source.toolbox.context != null"
          :source="source"
          :toolbox="source.toolbox"
          :context="source.toolbox.context"
          :style="mainStyleObject"
        >
        </Main>
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
        <ToolboxType
          ref="ToolboxType"
          v-if="source.toolbox.context != null"
          :source="source"
          :toolbox="source.toolbox"
          :context="source.toolbox.context"
          :groups="source.toolbox.groups"
        >
        </ToolboxType>
      </tm-layout>
    </template>
  </div>
</template>

<script>
import Main from "./Main";
import ToolboxType from "./ToolboxType";

export default {
  components: { Main, ToolboxType },
  props: ["source"],
  data() {
    return {
      ready: false,
      style: {
        backgroundColor: "#2d2d2d",
        color: "#adadad",
        header: {},
        left: {
          width: "260px",
        },
        main: {},
      },
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {
    boxStyleObject: function () {
      return {
        backgroundColor: this.style.backgroundColor,
        color: this.style.color,
      };
    },
    leftStyleObject: function () {
      return {};
    },
    mainStyleObject: function () {
      return {};
    },
  },
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    "$route.path"() {
      this.init();
    },
  },
  methods: {
    async init() {
      this.source.toolbox.toInsert = this.toInsert;
      this.source.toolbox.toUpdate = this.toUpdate;
      this.source.toolbox.toCopy = this.toCopy;
      this.source.toolbox.toDelete = this.toDelete;
      this.source.toolbox.toInsertGroup = this.toInsertGroup;
      this.source.toolbox.toUpdateGroup = this.toUpdateGroup;
      this.source.toolbox.toDeleteGroup = this.toDeleteGroup;
      this.source.toolbox.showSwitchToolboxType = () => {
        this.$refs.ToolboxType.showSwitch();
      };
      this.source.toolbox.showToolboxType = () => {
        this.$refs.ToolboxType.show();
      };
      this.source.toolbox.hideToolboxType = () => {
        this.$refs.ToolboxType.hide();
      };
      if (this.ready) {
        return;
      }
      if (!this.tool.isToolboxPage(this.$route.path)) {
        return;
      }
      this.source.toolbox.initContext = this.initContext;
      if (this.source.toolbox.context == null) {
        await this.initContext();
      }
      this.ready = true;
    },
    async initContext() {
      let res = await this.server.toolbox.data();
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        let data = res.data || {};

        data.mysqlColumnTypeInfos.forEach((one) => {
          one.name = one.name.toLowerCase();
        });
        this.source.toolbox.mysqlColumnTypeInfos = data.mysqlColumnTypeInfos;
        this.source.toolbox.databaseTypes = data.databaseTypes;
        this.source.toolbox.types = data.types;
        data.types.forEach((one) => {
          this.form.toolbox[one.name] = one.configForm;
          if (one.otherForm) {
            for (let formName in one.otherForm) {
              this.form.toolbox[one.name][formName] = one.otherForm[formName];
            }
          }
        });
        this.source.toolbox.sqlConditionalOperations =
          data.sqlConditionalOperations;
      }

      await this.loadContext();
    },
    async loadContext() {
      let param = {};
      let res = await this.server.toolbox.context(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        let context = res.data.context || {};
        let groups = res.data.groups || [];
        this.source.toolbox.groups = groups;
        this.source.toolbox.context = context;
      }
    },
    toInsert(toolboxType) {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
      let toolboxData = {};
      let optionsJSON = {};

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    toCopy(toolboxType, copy) {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
      let toolboxData = {};
      Object.assign(toolboxData, copy);
      delete toolboxData.toolboxId;
      toolboxData.name = toolboxData.name + " Copy";

      let optionsJSON = this.getOptionJSON(toolboxData.option);

      this.$refs.InsertToolbox.show({
        title: `新增[${toolboxType.text}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    toUpdate(toolboxType, toolboxData) {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
      this.updateData = toolboxData;

      let optionsJSON = this.getOptionJSON(toolboxData.option);

      this.$refs.UpdateToolbox.show({
        title: `编辑[${toolboxType.text}][${toolboxData.name}]工具`,
        form: [this.form.toolbox, toolboxType.configForm],
        data: [toolboxData, optionsJSON],
        toolboxType,
      });
    },
    getOptionJSON(option) {
      let json = {};
      if (this.tool.isNotEmpty(option)) {
        json = JSON.parse(option);
      }
      return json;
    },
    toDelete(toolboxType, toolboxData) {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
      this.tool
        .confirm(
          "删除[" +
            toolboxType.text +
            "]工具[" +
            toolboxData.name +
            "]将无法回复，确定删除？"
        )
        .then(async () => {
          return this.doDelete(toolboxType, toolboxData);
        })
        .catch((e) => {});
    },
    async doDelete(toolboxType, toolboxData) {
      let res = await this.server.toolbox.delete(toolboxData);
      if (res.code == 0) {
        this.source.toolbox.initContext();
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
        this.source.toolbox.initContext();
        this.tool.success("修改成功");
        let tab = this.source.toolbox.getTabByData(toolboxData);
        if (tab != null) {
          Object.assign(tab.toolboxData, toolboxData);
        }
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
      toolboxData.option = JSON.stringify(optionJSON);
      let res = await this.server.toolbox.insert(toolboxData);
      if (res.code == 0) {
        this.source.toolbox.initContext();
        this.tool.success("新增成功");
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },

    toInsertGroup() {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
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
      // this.source.toolbox.hideToolboxType();
      this.updateGroupData = data;

      let optionsJSON = this.getOptionJSON(data.option);

      this.$refs.UpdateToolboxGroup.show({
        title: `编辑[${data.name}]工具分组`,
        form: [this.form.toolbox.group, this.form.toolbox.group.option],
        data: [data, optionsJSON],
      });
    },
    toDeleteGroup(data) {
      this.tool.stopEvent();
      // this.source.toolbox.hideToolboxType();
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
        this.source.toolbox.initContext();
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
        this.source.toolbox.initContext();
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
        this.source.toolbox.initContext();
        return true;
      } else {
        this.tool.error(res.msg);
        return false;
      }
    },
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
  },
};
</script>

<style>
.toolbox-box {
  width: 100%;
  height: 100%;
  margin: 0px;
  padding: 0px;
  position: relative;
}
</style>
