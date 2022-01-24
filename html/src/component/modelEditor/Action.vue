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
        <div class="text">注释</div>
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
          <div class="text">请求地址</div>
          <Input_
            :source="source"
            :context="context"
            :bean="bean.api.request"
            name="path"
            :wrap="wrap"
          ></Input_>
          <div class="text">请求方法</div>
          <Input_
            :source="source"
            :context="context"
            :bean="bean.api.request"
            name="method"
            :wrap="wrap"
          ></Input_>
        </template>
      </li>
    </ul>
  </div>
</template>


<script>
import Input_ from "./Input.vue";

export default {
  components: { Input_ },
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
      if (this.isInitWebApiOpen) {
        return;
      }
      let api = null;
      if (open) {
        api = this.bean.api || {};
        api.request = api.request || {
          path: "/",
        };
        api.response = api.response || {};
      }
      this.wrap && this.wrap.onChange(this.bean, "api", api);
    },
  },
  methods: {
    init() {
      if (this.bean.api != null && this.bean.api.request != null) {
        this.isInitWebApiOpen = true;
        this.webApiOpen = true;
        this.$nextTick(() => {
          this.isInitWebApiOpen = false;
        });
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
