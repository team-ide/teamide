<template>
  <div v-if="ready">
    <table v-if="bean.variables != null" class="model-table">
      <thead>
        <tr>
          <th>变量</th>
          <th>值</th>
          <th>类型</th>
          <th v-if="bean.api != null && bean.api.request != null">位置</th>
          <th width="100">
            <span
              class="tm-pointer mgr-5"
              @click="wrap.push(bean, 'variables', {})"
              title="添加"
            >
              <b-icon icon="plus"></b-icon>
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <template v-for="(one, index) in bean.variables">
          <tr :key="'variable-' + index">
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
                name="value"
                title="变量值"
                :wrap="wrap"
              ></Input_>
            </td>
            <td>
              <Select_
                :source="source"
                :context="context"
                :bean="one"
                name="dataType"
                :isDataTypeOption="true"
                :wrap="wrap"
              ></Select_>
            </td>
            <td v-if="bean.api != null && bean.api.request != null">
              <Select_
                :source="source"
                :context="context"
                :bean="one"
                name="dataType"
                :isDataTypeOption="true"
                :wrap="wrap"
              ></Select_>
            </td>
            <td>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.up(bean, 'variables', one)"
                title="上移"
              >
                <b-icon icon="caret-up" class="ft-13"></b-icon>
              </span>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.down(bean, 'variables', one)"
                title="下移"
              >
                <b-icon icon="caret-down" class="ft-13"></b-icon>
              </span>
              <span
                class="tm-pointer mgr-5"
                @click="wrap.del(bean, 'variables', one)"
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
import Select_ from "./Select.vue";

export default {
  components: { Input_, Select_ },
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
