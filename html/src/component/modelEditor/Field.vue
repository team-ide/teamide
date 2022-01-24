<template>
  <li v-if="ready">
    <template v-if="tool.isNotEmpty(field.text)">
      <div class="text">{{ field.text }}</div>
    </template>

    <template v-if="field.fields == null">
      <ModelEditorFieldInput
        :source="source"
        :context="context"
        :bean="bean"
        :field="field"
        :wrap="wrap"
      ></ModelEditorFieldInput>
    </template>

    <template v-if="tool.isNotEmpty(field.comment)">
      <div class="comment">{{ field.comment }}</div>
    </template>
  </li>
</template>


<script>
export default {
  components: {},
  props: ["source", "context", "wrap", "bean", "field"],
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
