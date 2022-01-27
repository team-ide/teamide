<template>
  <div class="menu-sub-box">
    <slot></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: [],
  data() {
    return {
      isMenuSubBox: true,
    };
  },
  beforeMount() {},
  watch: {},
  methods: {},
  mounted() {},
};
</script>

<style  >
.menu-sub-box {
  top: 0;
  left: 100%;
  margin-top: -5px;
  margin-left: 0px;
  -webkit-border-radius: 0 6px 6px 6px;
  -moz-border-radius: 0 6px 6px 6px;
  border-radius: 0px 6px 6px 0px;
  display: none;
  position: absolute;
  background-color: #fff;
  padding: 5px 0px;
  border-left: 1px solid #ddd;
}
</style>