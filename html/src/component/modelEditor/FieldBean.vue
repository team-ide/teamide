<template>
  <li>
    <ul>
      <template v-for="(one, index) in field.fields">
        <ModelEditorField
          :key="'field-' + index"
          :source="source"
          :field="one"
          :bean="bean[field.name]"
          :wrap="wrap"
        >
        </ModelEditorField>
        <template v-if="one.fields != null && !one.isList">
          <ModelEditorFieldBean
            :key="'field-bean-' + index"
            :source="source"
            :field="one"
            :bean="bean"
            :wrap="wrap"
          >
          </ModelEditorFieldBean>
        </template>
        <template v-if="one.fields != null && one.isList">
          <ModelEditorFieldList
            :key="'field-list-' + index"
            :source="source"
            :field="one"
            :bean="bean"
            :wrap="wrap"
          >
          </ModelEditorFieldList>
        </template>
      </template>
    </ul>
  </li>
</template>


<script>
export default {
  components: {},
  props: ["source", "wrap", "bean", "field"],
  data() {
    return {};
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {},
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
