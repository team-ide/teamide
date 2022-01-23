<template>
  <div class="model-editor" v-if="config != null && config.fields != null">
    <ul>
      <template v-for="(one, index) in config.fields">
        <li :key="'field-' + index">
          <template v-if="tool.isNotEmpty(one.text)">
            <div class="text">{{ one.text }}</div>
          </template>

          <template v-if="one.fields == null">
            <template v-if="tool.isEmpty(one.type) || one.type == 'text'">
              <div class="input">
                <Input :bean="bean" :name="one.name"></Input>
              </div>
            </template>
          </template>

          <template v-if="tool.isNotEmpty(one.comment)">
            <div class="comment">{{ one.comment }}</div>
          </template>
        </li>
        <template v-if="one.fields != null && one.isList">
          <li :key="'field-table-' + index" class="pdl-10 mgtb-10">
            <table>
              <thead>
                <tr>
                  <template v-for="(sub, subIndex) in one.fields">
                    <th
                      :key="'field-table-th-' + subIndex"
                      :width="tool.isEmpty(sub.width) ? 60 : sub.width"
                    >
                      {{ sub.text }}
                    </th>
                  </template>
                  <th>
                    <div class="btn-group">
                      <div
                        class="tm-link color-green mgr-5"
                        @click="toUpdate(group, model)"
                        title="添加"
                      >
                        <b-icon icon="plus-circle-fill"></b-icon>
                      </div>
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <template v-for="(sub, subIndex) in one.fields">
                    <td :key="'field-table-td-' + subIndex">
                      <template
                        v-if="tool.isEmpty(sub.type) || sub.type == 'text'"
                      >
                        <div class="input">
                          <Input :bean="bean" :name="sub.name"></Input>
                        </div>
                      </template>
                    </td>
                  </template>
                  <td>
                    <div class="btn-group">
                      <div
                        class="tm-link mgr-5"
                        @click="toUpdate(group, model)"
                        title="上移"
                      >
                        <b-icon icon="caret-up-fill" class="ft-13"></b-icon>
                      </div>
                      <div
                        class="tm-link mgr-5"
                        @click="toUpdate(group, model)"
                        title="下移"
                      >
                        <b-icon icon="caret-down-fill" class="ft-13"></b-icon>
                      </div>
                      <div
                        class="tm-link mgr-5"
                        @click="toUpdate(group, model)"
                        title="删除"
                      >
                        <b-icon icon="backspace-fill" class="ft-13"></b-icon>
                      </div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </li>
        </template>
      </template>
    </ul>
  </div>
</template>


<script>
import Input from "./Input.vue";

export default {
  components: { Input },
  props: ["source", "bean", "config"],
  data() {
    return {
      key: this.tool.getNumber(),
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
