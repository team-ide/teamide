<template>
  <li v-if="ready">
    <template v-if="tool.isNotEmpty(field.text)">
      <div class="text">{{ field.text }}</div>
    </template>

    <template v-if="field.fields == null">
      <template v-if="tool.isEmpty(field.type) || field.type == 'text'">
        <div class="input">
          <Input
            :bean="bean"
            :name="field.name"
            :wrap="wrap"
            :readonly="field.readonly"
          ></Input>
        </div>
      </template>
    </template>

    <template v-if="tool.isNotEmpty(field.comment)">
      <div class="comment">{{ field.comment }}</div>
    </template>
  </li>
</template>


<script>
import Input from "./Input.vue";

export default {
  components: { Input },
  props: ["source", "wrap", "bean", "field"],
  data() {
    return {
      ready: false,
    };
  },
  // 计算属性 只有依赖数据发生改变，才会重新进行计算
  computed: {},
  // 计算属性 数据变，直接会触发相应的操作
  watch: {},
  methods: {
    init() {},
  },
  // 在实例创建完成后被立即调用
  created() {
    if (this.bean[this.field.name] == undefined) {
      this.bean[this.field.name] = null;
    }
    if (this.field.fields != null) {
      if (this.bean[this.field.name] == null) {
        if (this.field.isList) {
          this.bean[this.field.name] = [];
        } else {
          this.bean[this.field.name] = {};
        }
        this.wrap.refresh();
        return;
      }
    }
    this.ready = true;
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
