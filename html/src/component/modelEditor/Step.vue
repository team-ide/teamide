<template>
  <div v-if="ready" class="model-step">
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
        <span class="text">添加条件</span>
      </template>
      <input
        type="checkbox"
        class="model-switch"
        title="去除或设置条件"
        v-model="openIf"
      />

      <span
        class="tm-pointer color-red mgl-5"
        @click="wrap.del(parent, 'steps', bean)"
        title="删除步骤"
      >
        删除步骤
      </span>
      <AddStep
        :source="source"
        :context="context"
        :wrap="wrap"
        :bean="bean"
        :isChild="true"
      >
      </AddStep>
    </div>
    <div :class="{ 'pdl-20': openIf }">
      <div class="mgt-10" v-if="bean.variables != null">
        <span class="text">定义变量</span>
        <Variables
          class="mgl-60 mgt--20"
          :source="source"
          :context="context"
          :bean="bean"
          :wrap="wrap"
        >
        </Variables>
      </div>
      <div class="mgt-10" v-if="bean.validates != null">
        <span class="text">变量验证</span>
        <Validates
          class="mgl-60 mgt--20"
          :source="source"
          :context="context"
          :bean="bean"
          :wrap="wrap"
        >
        </Validates>
      </div>
      <template v-if="tool.isEmpty(stepType)">
        <div class="pdl-20">base</div>
      </template>
      <template v-else-if="stepType == 'sqlSelect'">
        <div class="pdl-20">sqlSelect</div>
      </template>

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
import Variables from "./Variables.vue";
import Validates from "./Validates.vue";
import AddStep from "./AddStep.vue";

export default {
  components: { Input_, Variables, Validates, AddStep },
  props: ["source", "context", "wrap", "bean", "parent"],
  data() {
    return {
      ready: false,
      openIf: false,
      stepType: null,
    };
  },
  computed: {},
  watch: {
    bean() {
      this.init();
    },
  },
  methods: {
    stepTypeChange(value) {
      this.stepType = value;
    },
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
.model-step {
  margin: 2px;
  border: 1px solid #3a3a3a;
}
</style>
