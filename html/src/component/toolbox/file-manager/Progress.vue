<template>
  <div class="toolbox-file-manager-progress-box">
    <div
      class="progress-list-box app-scroll-bar"
      @contextmenu.prevent="workContextmenu"
    >
      <template v-for="one in progressList">
        <div :key="'progress-' + one.progressId" class="progress-box">
          <div class="progress-text">
            <span class="pdr-5">{{
              tool.formatDate(new Date(one.startTime))
            }}</span>
            <template v-if="one.place == 'local'">
              <span class="pdr-5 color-grey-2">本地</span>
            </template>
            <template v-else-if="one.place == 'ssh'">
              <span class="pdr-5 color-grey-2">SSH</span>
            </template>
            <template v-else-if="one.place == 'node'">
              <span class="pdr-5 color-grey-2">节点</span>
            </template>
            <template v-if="one.work == 'create'">
              <span class="pdr-5">新建</span>
              <template v-if="one.data.isDir">
                <span class="pdr-5">文件夹</span>
              </template>
              <template v-else>
                <span class="pdr-5">文件</span>
              </template>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
            <template v-if="one.work == 'rename'">
              <span class="pdr-5">重命名</span>

              <span class="pdr-5">修改：{{ one.data.oldPath }}</span>
              <span class="pdr-5">为：{{ one.data.newPath }}</span>
            </template>
            <template v-if="one.work == 'move'">
              <span class="pdr-5">移动</span>

              <span class="pdr-5">移动：{{ one.data.oldPath }}</span>
              <span class="pdr-5">到：{{ one.data.newPath }}</span>
            </template>
            <template v-else-if="one.work == 'copy'">
              <span class="pdr-5">复制</span>

              <template v-if="one.data.fromPlace == 'local'">
                <span class="pdr-5 color-grey-2">本地</span>
              </template>
              <template v-else-if="one.data.fromPlace == 'ssh'">
                <span class="pdr-5 color-grey-2">SSH</span>
              </template>
              <template v-else-if="one.data.fromPlace == 'node'">
                <span class="pdr-5 color-grey-2">节点</span>
              </template>
              <span class="pdr-5">{{ one.data.fromPath }}</span>
              <span class="pdr-5">到</span>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
            <template v-else-if="one.work == 'upload'">
              <span class="pdr-5">上传</span>
              <span class="pdr-5">文件：{{ one.data.path }}</span>
            </template>
            <template v-else-if="one.work == 'write'">
              <span class="pdr-5">写入</span>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
            <template v-else-if="one.work == 'read'">
              <span class="pdr-5">读取</span>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
            <template v-else-if="one.work == 'download'">
              <span class="pdr-5">下载</span>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
            <template v-else-if="one.work == 'remove'">
              <span class="pdr-5">删除</span>
              <span class="pdr-5">{{ one.data.path }}</span>
            </template>
          </div>
          <div class="progress-status">
            <template v-if="one.data != null">
              <template v-if="one.data.count != null">
                <span class="color-grey-2 mgr-5">
                  文件：
                  {{ one.data.count }}
                  /
                  {{ one.data.successCount }}
                </span>
              </template>
              <template v-if="one.data.unitSize != null">
                <span class="color-grey-2 mgr-5">
                  大小：
                  {{ one.data.unitSize }}
                  {{ one.data.unit }}
                  <template v-if="one.data.unitSuccessSize != null">
                    /
                    {{ one.data.unitSuccessSize }}
                    {{ one.data.unitSuccess }}
                  </template>
                </span>
              </template>
              <template v-if="one.data.unitSleepSize != null">
                <span class="color-grey-2 mgr-5">
                  速度：
                  {{ one.data.unitSleepSize }}
                  {{ one.data.unitSleep }}
                  / 秒
                </span>
              </template>
              <template v-if="one.data.percentage != null">
                <span class="color-grey-2 mgr-5">
                  进度：
                  {{ one.data.percentage }}
                </span>
              </template>
            </template>
            <template v-if="one.error">
              <span class="color-red">{{ one.error }}</span>
            </template>
            <template v-else-if="one.waitActionIng">
              <span class="color-orange">{{ one.waitActionMessage }}</span>
              <template v-for="(action, index) in one.waitActionList">
                <span
                  :key="index"
                  class="tm-link mgl-5"
                  :class="action.color"
                  @click="doCallAction(one, action)"
                >
                  {{ action.text }}
                </span>
              </template>
            </template>
            <template v-else-if="one.isEnd">
              <span class="color-green">完成</span>
            </template>
            <template v-else>
              <span class="tm-link mgl-5 color-red-4" @click="doCallStop(one)">
                终止
              </span>
              <span class="mgl-5 color-orange">执行中</span>
            </template>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>


<script>
export default {
  components: {},
  props: ["source", "toolboxWorker"],
  data() {
    return {
      progressList: [],
    };
  },
  computed: {},
  watch: {},
  methods: {
    init() {
      this.unbindFileWorkProgress();
      this.bindFileWorkProgress();
    },
    bindFileWorkProgress() {
      this.server.addListenOnEvent(
        "file-work-progress",
        this.onFileWorkProgress
      );
    },
    unbindFileWorkProgress() {
      this.server.removeListenOnEvent(
        "file-work-progress",
        this.onFileWorkProgress
      );
    },

    async doCallAction(progress, action) {
      let param = {
        progressId: progress.progressId,
        action: action.value,
      };
      let res = await this.server.fileManager.callAction(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    async doCallStop(progress) {
      let param = {
        progressId: progress.progressId,
      };
      let res = await this.server.fileManager.callStop(param);
      if (res.code != 0) {
        this.tool.error(res.msg);
      }
    },
    getProgress(arg) {
      let res = null;
      this.progressList.forEach((one) => {
        if (one == arg || one.progressId == arg) {
          res = one;
        }
      });
      return res;
    },
    onFileWorkProgress(data) {
      if (this.isDestroyed) {
        return;
      }
      try {
        if (data.workerId != this.toolboxWorker.workerId) {
          return;
        }
        this.onProgress(data);
      } catch (error) {}
    },
    onProgress(progress) {
      if (progress.data.size) {
        this.tool.formatSize(progress.data, "size", "unitSize", "unit");

        if (progress.data.successSize) {
          this.tool.formatSize(
            progress.data,
            "successSize",
            "unitSuccessSize",
            "unitSuccess"
          );

          if (progress.data.size > 0 && progress.data.successSize > 0) {
            let percentage = Number(
              (progress.data.successSize / progress.data.size) * 100
            ).toFixed(0);
            progress.data.percentage = percentage + "%";
          } else {
            if (progress.data.size == progress.data.successSize) {
              progress.data.percentage = "100%";
            } else {
              progress.data.percentage = "0%";
            }
          }

          if (progress.startTime && progress.timestamp) {
            let usetime = progress.timestamp - progress.startTime;
            if (progress.endTime) {
              usetime = progress.endTime - progress.startTime;
            }
            let sleepSize = Number(
              (progress.data.successSize / usetime) * 1000
            ).toFixed(2);
            progress.data.sleepSize = sleepSize;

            this.tool.formatSize(
              progress.data,
              "sleepSize",
              "unitSleepSize",
              "unitSleep"
            );
          }
        }
      }

      if (progress.isEnd && progress.work == "upload") {
        if (progress.data) {
          if (progress.data.fileInfo) {
            this.toolboxWorker.onUploadFileInfo &&
              this.toolboxWorker.onUploadFileInfo(progress.data.fileInfo);
          }
          if (progress.data.fileDir) {
            this.toolboxWorker.onUploadFileInfo &&
              this.toolboxWorker.onUploadFileInfo(progress.data.fileDir);
          }
        }
      }

      let find = this.getProgress(progress.progressId);
      if (find) {
        Object.assign(find, progress);
        return;
      }
      this.progressList.push(progress);
    },
    cleanEnd() {
      let list = [];
      this.progressList.forEach((one) => {
        if (!one.isEnd) {
          list.push(one);
        }
      });
      this.progressList = list;
    },
    cleanAll() {
      let list = [];
      this.progressList.forEach((one) => {
        if (one.waitActionIng) {
          list.push(one);
        }
      });
      this.progressList = list;
    },
    workContextmenu(e) {
      e = e || window.event;
      let menus = [];

      menus.push({
        text: "清理已完成",
        onClick: () => {
          this.cleanEnd();
        },
      });

      menus.push({
        text: "清理所有",
        onClick: () => {
          this.cleanAll();
        },
      });

      this.tool.showContextmenu(menus);
    },
  },
  created() {},
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.isDestroyed = true;
    this.unbindFileWorkProgress();
  },
};
</script>

<style>
.toolbox-file-manager-progress-box {
  width: 100%;
  height: 100%;
}
.toolbox-file-manager-progress-box .progress-list-box {
  width: 100%;
  height: 100%;
  user-select: text;
}
.toolbox-file-manager-progress-box .progress-box {
  display: flex;
  line-height: 20px;
  font-size: 12px;
  padding: 0px 5px;
}
.toolbox-file-manager-progress-box .progress-box .progress-icon {
  padding: 0px 5px;
}
.toolbox-file-manager-progress-box .progress-box .progress-text {
  padding: 0px 5px;
  flex: 1;
}
.toolbox-file-manager-progress-box .progress-box .progress-size {
  padding: 0px 5px;
  font-size: 12px;
}
</style>
