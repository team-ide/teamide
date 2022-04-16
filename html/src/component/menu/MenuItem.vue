<template>
  <div
    class="menu-item"
    :class="{
      divider: divider,
      disabled: disabled,
      header: header,
      'has-sub-box': hasSubBox,
      'menu-sub-left': subLeft,
    }"
    @click="onClick"
  >
    <a class :class="aClass" :target="target" :href="href">
      <slot></slot>
    </a>
    <slot name="MenuSubBox"></slot>
  </div>
</template>

<script>
export default {
  components: {},
  props: [
    "divider",
    "disabled",
    "header",
    "target",
    "href",
    "aClass",
    "subLeft",
  ],
  data() {
    return {
      hasSubBox: false,
    };
  },
  beforeMount() {},
  watch: {},
  methods: {
    onClick() {
      if (this.disabled) {
        return;
      }
      this.$emit("click");
    },
  },
  mounted() {
    this.$children.forEach((one) => {
      if (one.isMenuSubBox) {
        this.hasSubBox = true;
      }
    });
  },
};
</script>

<style  >
.menu-item.divider {
  *width: 100%;
  height: 2px;
  margin: 9px 1px;
  *margin: -5px 0 5px;
  overflow: hidden;
  background-color: #e5e5e5;
  border-bottom: 1px solid #fff;
}

.menu-item.header {
  display: block;
  padding: 3px 30px 3px 5px;
  font-size: 12px;
  font-weight: 700;
  line-height: 24px;
  color: #999;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  text-decoration: none;
  cursor: default;
}
.menu-item {
  position: relative;
}

.menu-item > a {
  display: block;
  padding: 0px 10px;
  clear: both;
  font-weight: 400;
  line-height: 30px;
  color: #333 !important;
  white-space: nowrap;
  text-decoration: none;
  /* cursor: pointer; */
}

.menu-box.menu-mini .menu-item > a {
  line-height: 22px;
}
.menu-item:hover > a,
.menu-item:focus > a {
  text-decoration: none;
  background-color: #ddd;
}
.menu-item.disabled {
  cursor: no-drop;
  pointer-events: none;
  opacity: 0.65;
  filter: alpha(opacity = 65);
}
.menu-item.disabled > a,
.menu-item.disabled > a:hover {
  color: #999;
}

.menu-item.disabled > a:hover {
  text-decoration: none;
  cursor: default;
  background-color: transparent;
}

.menu-item:hover > .menu-sub-box {
  display: block;
}

.menu-item.has-sub-box > a {
  padding-right: 30px;
}
.menu-item.has-sub-box > a:after {
  display: block;
  float: right;
  width: 0;
  height: 0;
  margin-top: 5px;
  margin-right: -20px;
  border-color: transparent;
  border-left-color: #ccc;
  border-style: solid;
  border-width: 5px 0 5px 5px;
  content: " ";
}

.menu-item.has-sub-box > a:hover > a:after {
  border-left-color: #fff;
}

.menu-sub-left .menu-item.has-sub-box > a {
  padding-left: 30px;
}

.menu-sub-left .menu-item.has-sub-box > a:after {
  float: left;
  border-width: 5px 5px 5px 0;
  border-right-color: #ccc;
  margin-left: -20px;
}
.menu-sub-left .menu-item.has-sub-box > a:hover > a:after {
  border-right-color: #fff;
}
</style>