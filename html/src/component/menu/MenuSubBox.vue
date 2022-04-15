<template>
  <div class="menu-sub-box" :class="{ 'menu-sub-left': subLeft }">
    <slot></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: ["subLeft"],
  data() {
    return {
      isMenuSubBox: true,
    };
  },
  beforeMount() {},
  watch: {},
  methods: {},
  updated() {},
  mounted() {},
};
</script>

<style  >
.menu-sub-box {
  top: 0;
  left: 100%;
  margin-top: -5px;
  margin-left: 0px;
  -webkit-border-radius: 0px 6px 6px 6px;
  -moz-border-radius: 0px 6px 6px 6px;
  border-radius: 0px 6px 6px 0px;
  display: none;
  position: absolute;
  background-color: #fff;
  padding: 5px 0px;
  border-left: 1px solid #ddd;
}
.menu-sub-left .menu-sub-box {
  right: 100%;
  left: auto;
  -webkit-border-radius: 6px 6px 0px 6px;
  -moz-border-radius: 6px 6px 0px 6px;
  border-radius: 6px 6px 0px 0px;
  border-right: 1px solid #ddd;
  border-left: 0px;
}
</style>