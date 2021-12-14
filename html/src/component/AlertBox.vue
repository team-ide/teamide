<template>
  <div>
    <template v-for="(one, index) in alerts">
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
      alerts: [],
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    del(alert) {
      this.$nextTick(() => {
        this.alerts.splice(this.alerts.indexOf(alert), 1);
      });
    },
    show(title, msg, hasCancel) {
      return new Promise((resolve, reject) => {
        let alert = {
          id: "alert-" + this.tool.getRandom(),
          title,
          msg,
          show: () => {},
          hide: () => {
            this.del(alert);
          },
          okOnly: true,
          okTitle: "确认",
          ok: () => {
            if (resolve) {
              resolve();
            }
          },
          cancelTitle: "取消",
          cancel: () => {
            if (reject) {
              reject();
            }
          },
        };
        if (hasCancel) {
          alert.okOnly = false;
        }
        this.alerts.push(alert);
        this.$nextTick(() => {
          this.$bvModal.show(alert.id);
        });
      });
    },
  },
  // 在实例创建完成后被立即调用
  created() {
    this.tool.alert = (msg) => {
      return this.show("提示", msg, false);
    };
    this.tool.confirm = (msg) => {
      return this.show("提示", msg, true);
    };
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {},
};
</script>

<style>
</style>
