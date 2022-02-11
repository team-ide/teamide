<template>
  <li v-if="ready">
    <ul v-if="data != null">
      <template v-for="(one, index) in field.fields">
        <ModelEditorField
          :key="'field-' + index"
          :source="source"
          :context="context"
          :field="one"
          :bean="data[field.name]"
          :wrap="wrap"
          v-if="
            one.ifScript == null ||
            wrap.ifScript(one.ifScript, data[field.name])
          "
        >
        </ModelEditorField>
        <template v-if="one.fields != null && !one.isList">
          <ModelEditorFieldBean
            :key="'field-bean-' + index"
            :source="source"
            :context="context"
            :field="one"
            :bean="data"
            :wrap="wrap"
            v-if="one.ifScript == null || wrap.ifScript(one.ifScript, data)"
          >
          </ModelEditorFieldBean>
        </template>
        <template v-if="one.fields != null && one.isList">
          <ModelEditorFieldList
            :key="'field-list-' + index"
            :source="source"
            :context="context"
            :field="one"
            :bean="data"
            :wrap="wrap"
            v-if="one.ifScript == null || wrap.ifScript(one.ifScript, data)"
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
  props: ["source", "context", "wrap", "bean", "field"],
  data() {
    return {
      ready: false,
      data: null,
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
        this.bean[this.field.name] = {};
        this.wrap.refresh();
        return;
      }
      this.initData();
      this.ready = true;
    },
    initData() {
      this.data = this.bean[this.field.name];
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
