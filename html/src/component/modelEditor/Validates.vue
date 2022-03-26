<template>
  <div v-if="ready">
    <table v-if="bean.validates != null" class="model-table">
      <thead>
        <tr>
          <th>验证</th>
          <th>正则</th>
          <th>必填</th>
          <th>错误</th>
          <th>错误码</th>
          <th>错误信息</th>
          <th width="100">
            <span
              class="tm-pointer mgr-5"
              @click="wrap.push(bean, 'validates', {})"
              title="添加"
            >
              <b-icon icon="plus"></b-icon>
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <template v-for="(one, index) in bean.validates">
          <tr :key="'validate-' + index">
            <td>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="name"
                title="变量名称"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="pattern"
                title="正则"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <Switch_
                :source="source"
                :context="context"
                :bean="one"
                name="required"
                title="必填"
                :wrap="wrap"
              ></Switch_>
            </td>
            <td>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="error"
                title="错误"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="errorCode"
                title="错误码"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="errorMsg"
                title="错误信息"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.up(bean, 'validates', one)"
                title="上移"
              >
                <b-icon icon="caret-up" class="ft-13"></b-icon>
              </span>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.down(bean, 'validates', one)"
                title="下移"
              >
                <b-icon icon="caret-down" class="ft-13"></b-icon>
              </span>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.del(bean, 'validates', one)"
                title="删除"
              >
                <b-icon icon="backspace" class="ft-13"></b-icon>
              </span>
            </td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template>


<script>
import Input_ from "./Input.vue";
import Switch_ from "./Switch.vue";

export default {
  components: { Input_, Switch_ },
  props: ["source", "context", "wrap", "bean"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {
    bean() {
      this.init();
    },
  },
  methods: {
    init() {
      this.ready = true;
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
</style>
