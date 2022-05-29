<template>
  <div class="tm-layout" :style="{ width: w, height: h }">
    <slot></slot>
  </div>
</template>

<script>
import co from "./base/index.js";
export default {
  name: "tm-layout",
  props: ["width", "height", "bar"],
  data() {
    return {
      co,
      w: "100%",
      h: "100%",
      barSize: "2px",
      moveSize: null,
    };
  },
  computed: {},
  watch: {},
  methods: {
    getHW(type) {
      if (co.isNotEmpty(this[type])) {
        let res = "";
        if (this[type] == "auto") {
          res = "calc(100%";
          let parentChildren = this.$parent.$children;
          if (parentChildren != null) {
            parentChildren.forEach((one) => {
              if (one.isLayoutBar) {
                res += " - " + this.barSize;
              } else {
                if (one != this) {
                  if (co.isNotEmpty(one[type]) && one[type] != "auto") {
                    res += " - " + one[type];
                  }
                  let moveSize = one.getMoveSize();
                  if (co.isNotEmpty(moveSize)) {
                    if (moveSize < 0) {
                      res += " + " + -moveSize + "px";
                    } else {
                      res += " - " + moveSize + "px";
                    }
                  }
                }
              }
            });
          }
          if (co.isNotEmpty(this.moveSize)) {
            if (this.moveSize < 0) {
              res += " - " + -this.moveSize + "px";
            } else {
              res += "+" + this.moveSize + "px";
            }
          }
          res += ")";
        } else {
          if (co.isNotEmpty(this.moveSize)) {
            res = "calc(" + this[type];
            if (this.moveSize < 0) {
              res += " - " + -this.moveSize + "px";
            } else {
              res += " + " + this.moveSize + "px";
            }
            res += ")";
          } else {
            res = this[type];
          }
        }
        return res;
      }
    },
    getMoveSize() {
      return this.moveSize;
    },
    onMoveX(moveSize) {
      this.moveSize = moveSize;
      let parentChildren = this.$parent.$children;
      if (parentChildren != null) {
        parentChildren.forEach((one) => {
          one.initSize && one.initSize();
        });
      }
    },
    onMoveY(moveSize) {
      this.moveSize = moveSize;
      let parentChildren = this.$parent.$children;
      if (parentChildren != null) {
        parentChildren.forEach((one) => {
          one.initSize && one.initSize();
        });
      }
    },
    initSize() {
      let w = this.getHW("width");
      if (co.isNotEmpty(w)) {
        this.w = w;
      }
      let h = this.getHW("height");
      if (co.isNotEmpty(h)) {
        this.h = h;
      }
    },
    init() {
      this.initSize();
    },
  },
  mounted() {
    this.init();
  },
  destroyed() {},
};
</script>
