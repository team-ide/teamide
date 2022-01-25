<template>
  <div v-if="ready">
    <ul>
      <li>
        <div class="text">名称</div>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="name"
          :readonly="true"
          :wrap="wrap"
        ></Input_>
        <span class="pdlr-10">注释</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="comment"
          :wrap="wrap"
        ></Input_>
      </li>
      <li>
        <div class="text">详细描述</div>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="description"
          :wrap="wrap"
        ></Input_>
      </li>
      <li>
        <div class="text">Web Api</div>
        <input type="checkbox" class="model-switch" v-model="webApiOpen" />
        <template v-if="bean.api != null && bean.api.request != null">
          <span class="pdlr-10">请求地址</span>
          <Input_
            :source="source"
            :context="context"
            :bean="bean.api.request"
            name="path"
            :wrap="wrap"
          ></Input_>
          <span class="pdlr-10">请求方法</span>
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
        <div class="text">入参</div>
        (
        <template v-for="(one, index) in bean.inVariables">
          <div :key="'inVariable-' + index" class="inVariable-box">
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
            <div
              class="tm-link color-red mgl-5"
              @click="wrap.del(bean, 'inVariables', one)"
              title="删除"
            >
              <b-icon icon="x" class="ft-12"></b-icon>
            </div>
            <template v-if="index < bean.inVariables.length - 1">
              <span class="pdl-5 pdr-10">,</span>
            </template>
          </div>
        </template>
        )
        <div class="btn-group">
          <div
            class="tm-link color-green mgl-5"
            @click="
              wrap.push(bean, 'inVariables', { name: 'arg', dataType: '' })
            "
            title="添加入参"
          >
            <b-icon icon="plus-circle-fill"></b-icon>
          </div>
        </div>
      </li>
    </ul>
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
.inVariable-box {
  display: flex;
}
</style>
