<template>
  <ul class :class="{ 'contextmenu-submenu': isSub, contextmenu: !isSub }">
    <template v-if="menus != null">
      <template v-for="(menu, index) in menus">
        <li
          v-if="menu.divider"
          class="divider"
          :class="menu.class"
          :key="index"
        ></li>
        <li
          v-else-if="menu.header != null"
          class="header"
          :class="menu.class"
          :key="index"
          v-html="menu.header"
        ></li>
        <li
          v-else
          class
          :class="{
            'has-sub': menu.menus != null && menu.menus.length > 0,
            disabled: menu.disabled,
          }"
          :key="index"
          @click="onClick(menu)"
        >
          <a
            class
            :class="menu.class"
            :target="menu.target"
            :href="menu.href"
            v-html="menu.text"
          ></a>
          <template v-if="menu.menus != null">
            <ContextmenuMenus
              :menus="menu.menus"
              :contextmenu="contextmenu"
              :isSub="true"
            ></ContextmenuMenus>
          </template>
        </li>
      </template>
    </template>
  </ul>
</template>

<script>
import ContextmenuMenus from "@/components/ContextmenuMenus";
export default {
  name: "ContextmenuMenus",
  components: { ContextmenuMenus },
  props: ["contextmenu", "menus", "isSub"],
  data() {
    return {};
  },
  beforeMount() {},
  watch: {},
  methods: {
    onClick(menu) {
      if (menu) {
        if (menu.disabled) {
          return;
        }
        if (menu.onClick) {
          this.contextmenu.hide();
          menu.onClick();
        }
      }
    },
  },
  mounted() {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style  scoped>
.divider {
  *width: 100%;
  height: 2px;
  margin: 9px 1px;
  *margin: -5px 0 5px;
  overflow: hidden;
  background-color: #e5e5e5;
  border-bottom: 1px solid #fff;
}

.header {
  display: block;
  padding: 3px 30px 3px 15px;
  font-size: 12px;
  font-weight: 700;
  line-height: 24px;
  color: #999;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  text-decoration: none;
}
li {
  position: relative;
}

li a {
  display: block;
  padding: 3px 20px;
  clear: both;
  font-weight: 400;
  line-height: 24px;
  color: #333;
  white-space: nowrap;
  text-decoration: none;
  cursor: pointer;
}

li > a:hover,
li > a:focus,
.has-sub:hover > a {
  color: #fff !important;
  text-decoration: none;
  background-color: #009688;
  background-image: -moz-linear-gradient(top, #00897b, #00796b);
  background-image: -webkit-gradient(
    linear,
    0 0,
    0 100%,
    from(#00897b),
    to(#00796b)
  );
  background-image: -webkit-linear-gradient(top, #00897b, #00796b);
  background-image: -o-linear-gradient(top, #00897b, #00796b);
  background-image: linear-gradient(to bottom, #00897b, #00796b);
  background-repeat: repeat-x;
  filter: progid:dximagetransform.microsoft.gradient(startColorstr='#ff0088cc',
		endColorstr='#ff0077b3', GradientType=0);
}
.disabled {
  cursor: no-drop;
  pointer-events: none;
  opacity: 0.65;
  filter: alpha(opacity = 65);
}
.disabled > a,
.disabled > a:hover {
  color: #999;
}

.disabled > a:hover {
  text-decoration: none;
  cursor: default;
  background-color: transparent;
}

.contextmenu-submenu {
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

li:hover > .contextmenu-submenu {
  display: block;
}

.has-sub > a {
  padding-right: 30px;
}
.has-sub > a:after {
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

.has-sub > a:hover > a:after {
  border-left-color: #fff;
}
.contextmenu .header {
  cursor: default;
}
</style>