<template>
  <div class="toolbox-quick-command-box"></div>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.toolboxWorker.toInsertSSHCommand = this.toInsertSSHCommand;
      this.toolboxWorker.toUpdateSSHCommand = this.toUpdateSSHCommand;
      this.toolboxWorker.toDeleteSSHCommand = this.toDeleteSSHCommand;

      // source.quickCommands = null;
      // source.quickCommandSSHCommands = null;
      // source.getQuickCommandType = (name) => {
      //   if (source.quickCommandTypes == null) {
      //     return null;
      //   }
      //   let res = null;
      //   source.quickCommandTypes.forEach((one) => {
      //     if (one.name == name) {
      //       res = one;
      //     }
      //   });
      //   return res;
      // };
      // source.initToolboxQuickCommands = async () => {
      //   let data = {};
      //   if (source.login.user != null) {
      //     let res = await server.toolbox.quickCommand.query({});
      //     if (res.code != 0) {
      //       tool.error(res.msg);
      //     }
      //     data = res.data || {};
      //   }
      //   let quickCommands = data.quickCommands || [];

      //   let quickCommandSSHCommands = [];
      //   let quickCommandTypeSSHCommand =
      //     source.getQuickCommandType("SSH Command");

      //   quickCommands.forEach((one) => {
      //     if (quickCommandTypeSSHCommand) {
      //       if (one.quickCommandType == quickCommandTypeSSHCommand.value) {
      //         quickCommandSSHCommands.push(one);
      //       }
      //     }
      //   });
      //   source.quickCommands = quickCommands;
      //   source.quickCommandSSHCommands = quickCommandSSHCommands;
      // };
    },
    toInsertSSHCommand() {
      this.toolboxWorker.showQuickCommandSSHCommandForm({}, async (data) => {
        return await this.doInsert(data);
      });
    },
    toUpdateSSHCommand(param) {
      this.toolboxWorker.showQuickCommandSSHCommandForm(param, async (data) => {
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
      this.source.initToolboxQuickCommands();
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
      this.source.initToolboxQuickCommands();
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
      this.source.initToolboxQuickCommands();
      return true;
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
