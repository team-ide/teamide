<template>
  <div v-if="ready">
    <ul>
      <li>
        <span class="text">名称</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="name"
          :readonly="true"
          :wrap="wrap"
        ></Input_>
        <span class="text">注释</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="comment"
          :wrap="wrap"
        ></Input_>
      </li>
      <li>
        <span class="text">详细描述</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="description"
          :wrap="wrap"
        ></Input_>
      </li>
      <li>
        <span class="text">Web Api</span>
        <input type="checkbox" class="model-switch" v-model="webApiOpen" />
        <template v-if="bean.api != null && bean.api.request != null">
          <span class="text">请求地址</span>
          <Input_
            :source="source"
            :context="context"
            :bean="bean.api.request"
            name="path"
            :wrap="wrap"
          ></Input_>
          <span class="text">请求方法</span>
          <Input_
            :source="source"
            :context="context"
            :bean="bean.api.request"
            name="method"
            :wrap="wrap"
          ></Input_>
        </template>
      </li>
      <li>
        <span class="text">入参</span>
        (
        <template v-for="(one, index) in bean.inVariables">
          <span :key="'inVariable-' + index">
            <span class="pdlr-5" style="width: auto">名称</span>
            <Input_
              :source="source"
              :context="context"
              :bean="one"
              name="name"
              :wrap="wrap"
            ></Input_>
            <span class="pdlr-5" style="width: auto">类型</span>
            <Select_
              :source="source"
              :context="context"
              :bean="one"
              name="dataType"
              :isDataTypeOption="true"
              :wrap="wrap"
            ></Select_>
            <span
              class="tm-pointer color-red mgl-5"
              @click="wrap.del(bean, 'inVariables', one)"
              title="删除"
            >
              <b-icon icon="x" class="ft-12"></b-icon>
            </span>
            <template v-if="index < bean.inVariables.length - 1">
              <span class="pdl-5 pdr-10">,</span>
            </template>
          </span>
        </template>
        )
        <span
          class="tm-pointer color-green mgl-5"
          @click="wrap.push(bean, 'inVariables', { name: 'arg', dataType: '' })"
          title="添加入参"
        >
          <b-icon icon="plus-circle-fill"></b-icon>
        </span>
      </li>
      <li>
        <span class="text">出参</span>
        <span class="pdlr-5" style="width: auto">名称</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean.outVariable"
          name="name"
          :wrap="wrap"
        ></Input_>
        <span class="pdlr-5" style="width: auto">类型</span>
        <Select_
          :source="source"
          :context="context"
          :bean="bean.outVariable"
          name="dataType"
          :isDataTypeOption="true"
          :wrap="wrap"
        ></Select_>
      </li>
      <li>
        <span class="text">步骤</span>
        <AddStep :source="source" :context="context" :wrap="wrap" :bean="bean">
        </AddStep>

        <ModelEditorSteps
          class="pdl-20"
          :source="source"
          :context="context"
          :wrap="wrap"
          :bean="bean"
        >
        </ModelEditorSteps>
      </li>
    </ul>
  </div>
</template>


<script>
import Input_ from "./Input.vue";
import Select_ from "./Select.vue";
import AddStep from "./AddStep.vue";

export default {
  components: { Input_, Select_, AddStep },
  props: ["source", "context", "wrap", "bean"],
  data() {
    return {
      ready: false,
      webApiOpen: false,
    };
  },
  computed: {},
  watch: {
    bean() {
      this.init();
    },
    webApiOpen(open) {
      let api = null;
      if (open) {
        if (this.bean.api != null && this.bean.api.request != null) {
          this.bean.api.response = this.bean.api.response || {};
          return;
        }
        api = this.bean.api || {};
        api.request = api.request || {
          path: "/",
        };
      } else {
        if (this.bean.api == null || this.bean.api.request == null) {
          return;
        }
      }
      this.wrap && this.wrap.onChange(this.bean, "api", api);
    },
  },
  methods: {
    init() {
      if (this.bean.inVariables == null) {
        this.bean.inVariables = [];
        this.wrap.refresh();
        return;
      }
      if (this.bean.outVariable == null) {
        this.bean.outVariable = {};
        this.wrap.refresh();
        return;
      }
      if (this.bean.steps == null) {
        this.bean.steps = [];
        this.wrap.refresh();
        return;
      }
      if (this.bean.api != null && this.bean.api.request != null) {
        this.webApiOpen = true;
      }
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
