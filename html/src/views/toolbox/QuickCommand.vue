<template>
  <div class="toolbox-quick-command-box"></div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.toolbox.toInsertSSHCommand = this.toInsertSSHCommand;
      this.toolbox.toUpdateSSHCommand = this.toUpdateSSHCommand;
      this.toolbox.toDeleteSSHCommand = this.toDeleteSSHCommand;
      this.load();
    },
    toInsertSSHCommand() {
      this.toolbox.showQuickCommandSSHCommandForm({}, async (data) => {
        return await this.doInsert(data);
      });
    },
    toUpdateSSHCommand(param) {
      this.toolbox.showQuickCommandSSHCommandForm(param, async (data) => {
        return await this.doUpdate(data);
      });
    },
    toDeleteSSHCommand(param) {
      let msg = "确认删除";
      msg += "快速指令[" + param.name + "]";
      msg += "?";
      this.tool
        .confirm(msg)
        .then(async () => {
          this.doDelete(param);
        })
        .catch((e) => {});
    },
    async doInsert(data) {
      let param = {};
      Object.assign(param, data);
      if (param.quickCommandId) {
        param.quickCommandId = Number(param.quickCommandId);
      }
      param.quickCommandType = Number(param.quickCommandType);
      param.option = JSON.stringify(param.option);
      let res = await this.server.toolbox.quickCommand.insert(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      this.tool.success("新增成功");
      this.load();
      return true;
    },
    async doUpdate(data) {
      let param = {};
      Object.assign(param, data);
      param.quickCommandId = Number(param.quickCommandId);
      param.quickCommandType = Number(param.quickCommandType);
      param.option = JSON.stringify(param.option);
      let res = await this.server.toolbox.quickCommand.update(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      this.tool.success("修改成功");
      this.load();
      return true;
    },
    async doDelete(data) {
      let param = {
        quickCommandId: Number(data.quickCommandId),
      };
      let res = await this.server.toolbox.quickCommand.delete(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
        return false;
      }
      this.tool.success("删除成功");
      this.load();
      return true;
    },
    async load() {
      let param = {};
      let res = await this.server.toolbox.quickCommand.query(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      } else {
        let quickCommands = res.data.quickCommands || [];

        let quickCommandSSHCommands = [];
        let quickCommandTypeSSHCommand =
          this.toolbox.getQuickCommandType("SSH Command");

        quickCommands.forEach((one) => {
          if (quickCommandTypeSSHCommand) {
            if (one.quickCommandType == quickCommandTypeSSHCommand.value) {
              quickCommandSSHCommands.push(one);
            }
          }
        });
        this.source.toolbox.quickCommands = quickCommands;
        this.source.toolbox.quickCommandSSHCommands = quickCommandSSHCommands;
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
.toolbox-quick-command-box {
  display: none;
}
</style>
