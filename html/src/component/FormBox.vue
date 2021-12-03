<template>
  <div>
    <template v-for="(one, index) in forms">
      <b-modal
        :key="index"
        :id="one.id"
        :title="one.title"
        @show="one.show"
        @cancel="one.cancel"
        @hide="one.hide"
        @ok="one.ok"
        :ok-title="one.okTitle"
        :cancel-title="one.cancelTitle"
        :ok-only="one.okOnly"
        :hide-header-close="true"
        :no-close-on-backdrop="true"
        :no-close-on-esc="true"
        auto-focus-button="ok"
        :hide-backdrop="true"
      >
        <div v-html="one.msg"></div>
      </b-modal>
    </template>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source"],
  data() {
    return {
      forms: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    del(form) {
      this.$nextTick(() => {
        this.forms.splice(this.forms.indexOf(form), 1);
      });
    },
    show(config) {
      let form = config || {};
      form.id = "form-" + this.tool.getRandom();
      form.inputs = form.inputs || [];

      form.hide = () => {
        this.del(form);
      };
      form.okTitle = form.okTitle || "确认";
      form.cancelTitle = form.cancelTitle || "取消";
      this.forms.push(form);
      this.$nextTick(() => {
        this.$bvModal.show(form.id);
      });
    },
  },
  // 在实例创建完成后被立即调用
  created() {
    this.tool.form = (config) => {
      this.show(config);
    };
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
.login-box {
  width: 100%;
  height: 100%;
}
</style>
