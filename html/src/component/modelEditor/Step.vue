<template>
  <div v-if="ready">
    <div>
      <template v-if="openIf">
        <span class="text">if (</span>
        <Input_
          :source="source"
          :context="context"
          :bean="bean"
          name="if"
          :wrap="wrap"
        ></Input_>
        <span class="text">) {</span>
      </template>
      <template v-else>
        <span class="text">if</span>
      </template>
      <input
        type="checkbox"
        class="model-switch"
        title="去除或设置条件"
        v-model="openIf"
      />

      <Select_
        :source="source"
        :context="context"
        :bean="this"
        name="stepType"
        :options="stepTypes"
        placeholder="基础步骤"
      ></Select_>

      <span
        class="tm-pointer color-red mgl-5"
        @click="wrap.del(parent, 'steps', bean)"
        title="删除步骤"
      >
        <b-icon icon="x" class="ft-12"></b-icon>
      </span>
      <span
        class="tm-pointer color-green mgl-5"
        @click="wrap.push(bean, 'steps', { name: 'arg', dataType: '' })"
        title="添加子步骤"
      >
        <b-icon icon="plus"></b-icon>
      </span>
    </div>
    <div :class="{ 'pdl-20': openIf }">
      <div>
        {{bean}}
      </div>

      <ModelEditorSteps
        class="pdl-20"
        :source="source"
        :context="context"
        :wrap="wrap"
        :bean="bean"
      >
      </ModelEditorSteps>
    </div>
    <div>
      <template v-if="openIf">
        <span class="text">}</span>
      </template>
    </div>
  </div>
</template>


<script>
import Input_ from "./Input.vue";
import Select_ from "./Select.vue";

export default {
  components: { Input_, Select_ },
  props: ["source", "context", "wrap", "bean", "parent"],
  data() {
    return {
      ready: false,
      openIf: false,
      stepType: null,
      stepTypes: [
        { value: "sqlSelect", text: "Sql Select" },
        { value: "sqlInsert", text: "Sql Insert" },
        { value: "sqlUpdate", text: "Sql Update" },
        { value: "sqlDelete", text: "Sql Delete" },
        { value: "redisGet", text: "Redis Get" },
        { value: "redisSet", text: "Redis Set" },
        { value: "redisDel", text: "Redis Del" },
      ],
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
      if (this.tool.isNotEmpty(this.bean.if)) {
        this.openIf = true;
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
