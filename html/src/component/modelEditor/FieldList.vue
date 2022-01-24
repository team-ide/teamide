<template>
  <li v-if="ready">
    <table class="model-table">
      <thead>
        <tr>
          <template v-for="(one, index) in field.fields">
            <th :key="'field-table-th-' + index">
              {{ one.text }}
            </th>
          </template>
          <th width="100">
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
      <template v-if="list != null">
        <tbody>
          <template v-for="(oneBean, oneBeanIndex) in list">
            <tr :key="'field-table-tr-' + oneBeanIndex">
              <template v-for="(one, index) in field.fields">
                <td :key="'field-table-td-' + index">
                  <ModelEditorFieldInput
                    :source="source"
                    :context="context"
                    :bean="oneBean"
                    :field="one"
                    :wrap="wrap"
                  ></ModelEditorFieldInput>
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
export default {
  components: {},
  props: ["source", "context", "wrap", "bean", "field"],
  data() {
    return {
      ready: false,
      list: null,
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
      if (this.bean[this.field.name] == null) {
        this.bean[this.field.name] = [];
        this.wrap.refresh();
        return;
      }
      this.initList();
      this.ready = true;
    },
    initList() {
      this.list = this.bean[this.field.name];
    },
  },
  created() {},
  mounted() {
    this.init();
  },
};
</script>

<style>
.model-table {
  border: 1px dashed #4e4e4e;
}
.model-table thead {
  border-bottom: 1px dashed #4e4e4e;
}
.model-table th {
  text-align: center;
  line-height: 25px;
  padding: 0px 10px;
  min-width: 75px;
}
.model-table td {
  text-align: center;
}
</style>
