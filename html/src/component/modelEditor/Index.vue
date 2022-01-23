<template>
  <div class="model-editor" v-if="group != null && group.fields != null">
    <ul v-if="data != null">
      <template v-for="(one, index) in group.fields">
        <ModelEditorField
          :key="'field-' + index"
          :source="source"
          :field="one"
          :bean="data"
          :wrap="wrap"
        >
        </ModelEditorField>
        <template v-if="one.fields != null && !one.isList">
          <ModelEditorFieldBean
            :key="'field-bean-' + index"
            :source="source"
            :field="one"
            :bean="data"
            :wrap="wrap"
          >
          </ModelEditorFieldBean>
        </template>
        <template v-if="one.fields != null && one.isList">
          <ModelEditorFieldList
            :key="'field-list-' + index"
            :source="source"
            :field="one"
            :bean="data"
            :wrap="wrap"
          >
          </ModelEditorFieldList>
        </template>
      </template>
    </ul>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "model", "group"],
  data() {
    return {
      key: this.tool.getNumber(),
      wrap: {},
      data: null,
      ready: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {
    model(newModel, oldModel) {
      this.initData();
    },
  },
  methods: {
    init() {
      this.set(JSON.parse(JSON.stringify(this.model || {})));
    },
    initData() {
      this.set(JSON.parse(JSON.stringify(this.data || {})));
    },
    get() {
      return this.data;
    },
    set(data) {
      this.data = data;
    },
    refresh() {
      this.initData();
    },
    onChange(bean, name, value) {
      bean = bean || {};
      bean[name] = value;
      this.toCallChange();
    },
    push(bean, name, value) {
      bean = bean || {};
      bean[name] = bean[name] || [];
      value = value || {};
      bean[name].push(value);
      this.toCallChange();
    },
    del(bean, name, value) {
      bean = bean || {};
      bean[name] = bean[name] || [];
      value = value || {};
      let index = bean[name].indexOf(value);
      if (index >= 0) {
        bean[name].splice(index, 1);
        this.toCallChange();
      }
    },
    up(bean, name, value) {
      bean = bean || {};
      bean[name] = bean[name] || [];
      value = value || {};
      let index = bean[name].indexOf(value);
      if (index > 0) {
        bean[name].splice(index, 1);
        bean[name].splice(index - 1, 0, value);
        this.toCallChange();
      }
    },
    down(bean, name, value) {
      bean = bean || {};
      bean[name] = bean[name] || [];
      value = value || {};
      let index = bean[name].indexOf(value);
      if (index >= 0 && index < bean[name].length - 1) {
        bean[name].splice(index, 1);
        bean[name].splice(index + 1, 0, value);
        this.toCallChange();
      }
    },
    toCallChange() {
      let change_model = JSON.parse(JSON.stringify(this.data));
      // console.log(this.data);
      this.last_change_model = change_model;
      this.$emit("change", this.group, change_model);
    },
  },
  // 在实例创建完成后被立即调用
  created() {
    this.wrap.push = this.push;
    this.wrap.del = this.del;
    this.wrap.up = this.up;
    this.wrap.down = this.down;
    this.wrap.onChange = this.onChange;
    this.wrap.refresh = this.refresh;
  },
  // el 被新创建的 vm.$el 替换，并挂载到实例上去之后调用
  mounted() {
    this.init();
  },
};
</script>

<style>
.model-editor {
  width: 100%;
}
.model-editor ul {
  margin-top: 10px;
}
.model-editor ul,
.model-editor li {
  list-style: none;
  padding: 0px;
  font-size: 12px;
}
.model-editor li {
  display: flex;
}
.model-editor .text {
  min-width: 80px;
}
.model-editor .text,
.model-editor .input,
.model-editor .comment {
  padding: 2px 5px;
}
.model-editor input {
  padding: 1px 5px;
  border: 0px;
  outline: none;
  background: transparent;
  color: #f9f9f9;
  border-bottom: 1px dashed #f9f9f9;
  min-width: 40px;
}
.model-editor table {
  border: 1px dashed #4e4e4e;
}
.model-editor thead {
  border-bottom: 1px dashed #4e4e4e;
}
.model-editor table th {
  text-align: center;
  line-height: 25px;
}
</style>
