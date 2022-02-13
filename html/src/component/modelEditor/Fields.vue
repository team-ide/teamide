<template>
  <ul class="mgt-40">
    <template v-for="(one, index) in fields">
      <ModelEditorField
        :key="'field-' + index"
        :source="source"
        :context="context"
        :field="one"
        :bean="bean"
        :wrap="wrap"
      >
      </ModelEditorField>
      <template v-if="one.fields != null && !one.isList">
        <ModelEditorFields
          :key="'field-bean-' + index"
          :source="source"
          :context="context"
          :fields="one.fields"
          :bean="bean"
          :wrap="wrap"
        >
        </ModelEditorFields>
      </template>
      <template v-if="one.fields != null && one.isList">
        <ModelEditorFieldList
          :key="'field-list-' + index"
          :source="source"
          :context="context"
          :field="one"
          :bean="bean"
          :wrap="wrap"
        >
        </ModelEditorFieldList>
      </template>
    </template>
  </ul>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "wrap", "bean", "fields"],
  data() {
    return {
      ready: false,
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
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
