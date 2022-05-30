<template>
  <el-dialog
    ref="modal"
    :title="`SSH快速指令`"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="true"
    :append-to-body="true"
    :visible="showDialog"
    :before-close="hide"
    width="800px"
  >
    <div class="">
      <el-form ref="form" size="mini" @submit.native.prevent>
        <el-form-item label="名称">
          <el-input type="input" v-model="name"> </el-input>
        </el-form-item>
        <el-form-item label="说明">
          <el-input type="input" v-model="comment"> </el-input>
        </el-form-item>

        <el-form-item label="指令">
          <el-input type="textarea" v-model="command"> </el-input>
        </el-form-item>
      </el-form>
    </div>
    <div class="">
      <div
        class="tm-btn bg-teal-8 ft-18 pdtb-5"
        :class="{ 'tm-disabled': saveBtnDisabled }"
        @click="doSave"
      >
        保存
      </div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  components: {},
  props: ["source", "toolbox"],
  data() {
    return {
      showDialog: false,
      name: null,
      comment: null,
      command: null,
      saveBtnDisabled: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    show(data, callback) {
      data = data || {};

      this.quickCommandType = this.toolbox.getQuickCommandType("SSH Command");
      if (this.quickCommandType == null) {
        this.tool.error("SSH Command类型不存在");
        return;
      }

      this.quickCommandId = data.quickCommandId;
      this.name = data.name;
      this.comment = data.comment;
      this.option = this.toolbox.getOptionJSON(data.option);
      this.option = this.option || {};
      this.command = this.option.command;

      this.callback = callback;
      this.showDialog = true;
    },
    hide() {
      this.showDialog = false;
    },
    async doSave() {
      let name = this.name;
      if (this.tool.isEmpty(name)) {
        this.tool.error("请输入名称");
        return;
      }
      let command = this.command;
      if (this.tool.isEmpty(command)) {
        this.tool.error("请输入指令");
        return;
      }

      this.saveBtnDisabled = true;
      let option = this.option || {};
      option.command = command;
      let param = {
        quickCommandId: this.quickCommandId,
        name: this.name,
        comment: this.comment,
        option: option,
        quickCommandType: this.quickCommandType.value,
      };
      let flag = await this.callback(param);
      this.saveBtnDisabled = false;
      if (flag) {
        this.hide();
      }
    },
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.toolbox.showQuickCommandSSHCommandForm = this.show;
    this.toolbox.hideQuickCommandSSHCommandForm = this.hide;
    this.init();
  },
};
</script>

<style>
</style>
