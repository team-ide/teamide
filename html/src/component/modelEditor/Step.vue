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
        <span class="text">// if</span>
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
        <b-icon icon="x" class="ft-12"></b-icon>
      </span>
      <span
        class="tm-pointer color-green mgl-5"
        @click="wrap.push(bean, 'steps', { name: 'step' })"
        title="添加子步骤"
      >
        <b-icon icon="plus"></b-icon>
      </span>
    </div>
    <div :class="{ 'pdl-20': openIf }">
      <template v-if="bean.variables != null">
        <div>
          <div>
            <span class="text">变量</span>
            <span
              class="tm-pointer mgl-5"
              @click="
                wrap.push(bean, 'variables', {
                  name: 'name',
                  value: 'value',
                  dataType: 'string',
                })
              "
              title="添加变量"
            >
              <b-icon icon="plus"></b-icon>
            </span>
          </div>
          <template v-for="(one, index) in bean.variables">
            <div :key="'variable-' + index" class="pdl-20">
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="name"
                title="变量名称"
                :wrap="wrap"
              ></Input_>
              <span class="text">:</span>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="value"
                title="值"
                :wrap="wrap"
              ></Input_>
              <span class="text">类型</span>
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
                @click="wrap.del(bean, 'variables', one)"
                title="删除"
              >
                <b-icon icon="x" class="ft-12"></b-icon>
              </span>
            </div>
          </template>
        </div>
      </template>
      <template v-if="bean.validates != null">
        <div>
          <div>
            <span class="text">验证</span>
            <span
              class="tm-pointer mgl-5"
              @click="
                wrap.push(bean, 'validates', {
                  name: 'name',
                  pattern: '',
                })
              "
              title="添加验证"
            >
              <b-icon icon="plus"></b-icon>
            </span>
          </div>
          <template v-for="(one, index) in bean.validates">
            <div :key="'validate-' + index" class="pdl-20">
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="name"
                :wrap="wrap"
                title="变量"
              ></Input_>
              <span class="text">:</span>
              <Input_
                :source="source"
                :context="context"
                :bean="one"
                name="value"
                :wrap="wrap"
                title="pattern"
              ></Input_>
            </div>
          </template>
        </div>
      </template>
      <div>
        <Select_
          :source="source"
          :context="context"
          :bean="this"
          name="stepType"
          :options="stepTypes"
          placeholder="普通"
          class="pdl-0"
          @change="stepTypeChange"
        ></Select_>
        <span class="text">步骤</span>
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
    stepTypeChange(value) {
      this.stepType = value;
    },
    init() {
      if (this.bean.variables == null) {
        this.bean.variables = [];
        this.wrap.refresh();
        return;
      }
      if (this.bean.validates == null) {
        this.bean.validates = [];
        this.wrap.refresh();
        return;
      }
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
