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
        <b-form @submit="onSubmit" @reset="onReset" v-if="show">
          <b-form-group
            id="input-group-1"
            label="Email address:"
            label-for="input-1"
            description="We'll never share your email with anyone else."
          >
            <b-form-input
              id="input-1"
              v-model="form.email"
              type="email"
              placeholder="Enter email"
              required
            ></b-form-input>
          </b-form-group>

          <b-form-group
            id="input-group-2"
            label="Your Name:"
            label-for="input-2"
          >
            <b-form-input
              id="input-2"
              v-model="form.name"
              placeholder="Enter name"
              required
            ></b-form-input>
          </b-form-group>

          <b-form-group id="input-group-3" label="Food:" label-for="input-3">
            <b-form-select
              id="input-3"
              v-model="form.food"
              :options="foods"
              required
            ></b-form-select>
          </b-form-group>

          <b-form-group id="input-group-4" v-slot="{ ariaDescribedby }">
            <b-form-checkbox-group
              v-model="form.checked"
              id="checkboxes-4"
              :aria-describedby="ariaDescribedby"
            >
              <b-form-checkbox value="me">Check me out</b-form-checkbox>
              <b-form-checkbox value="that">Check that out</b-form-checkbox>
            </b-form-checkbox-group>
          </b-form-group>

          <b-button type="submit" variant="primary">Submit</b-button>
          <b-button type="reset" variant="danger">Reset</b-button>
        </b-form>
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
      let ok = form.ok;
      let validate = form.validate;
      form.ok = () => {
        if (validate) {
          if (!validate()) {
            return;
          }
        }
        ok && ok();
      };
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
