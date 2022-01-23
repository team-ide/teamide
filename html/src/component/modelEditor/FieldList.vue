<template>
  <li>
    <table>
      <thead>
        <tr>
          <template v-for="(one, index) in field.fields">
            <th
              :key="'field-table-th-' + index"
              :width="tool.isEmpty(one.width) ? 60 : one.width"
            >
              {{ one.text }}
            </th>
          </template>
          <th>
            <div class="btn-group">
              <div
                class="tm-link color-green mgr-5"
                @click="wrap.push(bean, field.name, {})"
                title="添加"
              >
                <b-icon icon="plus-circle-fill"></b-icon>
              </div>
            </div>
          </th>
        </tr>
      </thead>
      <template v-if="bean[field.name] != null">
        <tbody>
          <template v-for="(oneBean, oneBeanIndex) in bean[field.name]">
            <tr :key="'field-table-tr-' + oneBeanIndex">
              <template v-for="(one, index) in field.fields">
                <td :key="'field-table-td-' + index">
                  <template v-if="tool.isEmpty(one.type) || one.type == 'text'">
                    <div class="input">
                      <Input
                        :bean="oneBean"
                        :name="one.name"
                        :readonly="one.readonly"
                        :wrap="wrap"
                      ></Input>
                    </div>
                  </template>
                </td>
              </template>
              <td>
                <div class="btn-group">
                  <div
                    class="tm-link mgr-5"
                    @click="wrap.up(bean, field.name, oneBean)"
                    title="上移"
                  >
                    <b-icon icon="caret-up-fill" class="ft-13"></b-icon>
                  </div>
                  <div
                    class="tm-link mgr-5"
                    @click="wrap.down(bean, field.name, oneBean)"
                    title="下移"
                  >
                    <b-icon icon="caret-down-fill" class="ft-13"></b-icon>
                  </div>
                  <div
                    class="tm-link mgr-5"
                    @click="wrap.del(bean, field.name, oneBean)"
                    title="删除"
                  >
                    <b-icon icon="backspace-fill" class="ft-13"></b-icon>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </template>
    </table>
  </li>
</template>


<script>
import Input from "./Input.vue";

export default {
  components: { Input },
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
