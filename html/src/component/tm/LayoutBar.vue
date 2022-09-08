<template>
  <div
    class="tm-layout-bar"
    :class="classObject"
    :style="{ width: w, height: h }"
  >
    <div
      class="tm-layout-bar-part"
      :class="partClassObject"
      @mousedown="mousedown"
    ></div>
  </div>
</template>

<script>
import co from "./base/index.js";
export default {
  name: "tm-layout-bar",
  props: ["left", "right", "top", "bottom", "color"],
  data() {
    return {
      co,
      isLayoutBar: true,
      w: null,
      h: null,
    };
  },
  computed: {
    classObject() {
      let option = this;
      let res = {};
      if (co.isNotEmpty(option.color)) {
        let color = co.addColor(option.color);
        res["bg-" + color] = true;
      }
      if (this.left != undefined) {
        res["tm-layout-bar-left"] = true;
      }
      if (this.right != undefined) {
        res["tm-layout-bar-right"] = true;
      }
      if (this.top != undefined) {
        res["tm-layout-bar-top"] = true;
      }
      if (this.bottom != undefined) {
        res["tm-layout-bar-bottom"] = true;
      }
      return res;
    },
    partClassObject() {
      let option = this;
      let res = {};
      if (co.isNotEmpty(option.color)) {
        let color = co.addColor(option.color);
        res["bg-" + color] = true;
      }
      return res;
    },
  },
  watch: {},
  methods: {
    mousedown(event) {
      if (this.controlLayout == null) {
        return;
      }
      // 获取event对象，兼容性写法
      event = event || window.event;
      // 鼠标按下时的位置
      this.mouseDownX = event.clientX;
      this.mouseDownY = event.clientY;
      this.downMoveSize = this.lastMoveSize;
      this.mouseDowned = true;
    },
    mouseup(event) {
      this.mouseDowned = false;
    },
    mousemove(event) {
      if (!this.mouseDowned) {
        return;
      }
      // 获取event对象，兼容性写法
      event = event || window.event;
      // 鼠标按下时的位置

      if (this.left != undefined || this.right != undefined) {
        let moveSize = event.clientX - this.mouseDownX;
        if (this.left != undefined) {
          moveSize = -moveSize;
        }
        if (co.isNotEmpty(this.downMoveSize)) {
          moveSize += this.downMoveSize;
        }
        this.lastMoveSize = moveSize;
        this.controlLayout.onMoveX(moveSize);
      }
      if (this.top != undefined || this.bottom != undefined) {
        let moveSize = event.clientY - this.mouseDownY;
        if (this.top != undefined) {
          moveSize = -moveSize;
        }
        if (co.isNotEmpty(this.downMoveSize)) {
          moveSize += this.downMoveSize;
        }
        this.lastMoveSize = moveSize;
        this.controlLayout.onMoveY(moveSize);
      }
    },
    getControlLayout() {
      let parentChildren = this.$parent.$children;
      let layout = null;
      if (parentChildren != null) {
        let index = parentChildren.indexOf(this);
        if (this.left != undefined || this.top != undefined) {
          layout = parentChildren[index + 1];
        }
        if (this.right != undefined || this.bottom != undefined) {
          layout = parentChildren[index - 1];
        }
      }
      return layout;
    },
    init() {
      this.controlLayout = this.getControlLayout();
      if (this.left != undefined || this.right != undefined) {
        this.h = "100%";
        this.w = this.$parent.barSize;
      }
      if (this.top != undefined || this.bottom != undefined) {
        this.w = "100%";
        this.h = this.$parent.barSize;
      }
      document.addEventListener("mouseup", (event) => {
        this.mouseup(event);
      });
      document.addEventListener("mousemove", (event) => {
        this.mousemove(event);
      });
    },
  },
  mounted() {
    this.init();
  },
  beforeDestroy() {},
};
</script>
