const newItemsWorker = function (workerOption) {
    workerOption = workerOption || {};
    const worker = {
        items: [],
        activeItem: null,
        getItem(item) {
            if (item == null) {
                return
            }
            let res = null;
            worker.items.forEach(one => {
                if (one == item ||
                    one.key == item ||
                    one.key == item.key
                ) {
                    res = one;
                }
            })
            return res;
        },
        getItemIndex(item) {
            let find = worker.getItem(item);
            if (find) {
                return -1
            }
            return worker.items.indexOf(find);
        },
        getActiveItemIndex() {
            return worker.items.indexOf(worker.activeItem);
        },
        addItem(item, before) {
            let find = worker.getItem(item);
            if (find) {
                return
            }
            if (item.show == undefined) {
                item.show = true
            }
            if (item.name == undefined) {
                item.name = ""
            }
            if (item.title == undefined) {
                item.title = ""
            }
            let beforeIndex = -1
            if (before) {
                beforeIndex = worker.getItemIndex(before)
            }
            if (beforeIndex >= 0) {
                worker.items.splice(beforeIndex + 1, 0, item);
            } else {
                worker.items.push(item)
            }
        },
        toDeleteOtherItem(item) {
            let list = [];
            worker.items.forEach((one) => {
                if (one != item) {
                    list.push(one);
                }
            });
            worker.toRemoveItems(list);
            if (worker.activeItem != item) {
                worker.toActiveItem(item);
            }
        },
        toRemoveAll() {
            let list = [];
            worker.items.forEach((one) => {
                list.push(one);
            });
            worker.toRemoveItems(list);
        },
        toRemoveItems(list) {
            list.forEach((one) => {
                let index = worker.items.indexOf(one);
                if (index >= 0) {
                    worker.items.splice(index, 1);
                    workerOption.onRemoveItem && workerOption.onRemoveItem(one);
                }
            });
        },
        toRemoveItem(item) {
            let find = worker.getItem(item);
            if (find == null) {
                return;
            }
            let index = worker.items.indexOf(find);
            worker.items.splice(index, 1);
            workerOption.onRemoveItem && workerOption.onRemoveItem(find);
            if (find == worker.activeItem) {
                let nextIndex = index - 1;
                if (nextIndex < 0) {
                    nextIndex = 0;
                }
                worker.toActiveItem(worker.items[nextIndex]);
            }
        },
        activeNextItem(item) {
            let find = worker.getItem(item);
            if (find == null) {
                return;
            }
            let index = worker.items.indexOf(find);
            let next = null;
            if (index >= 0) {
                if (worker.items[index + 1]) {
                    next = worker.items[index + 1];
                } else {
                    next = worker.items[index - 1];
                }
            }
            worker.toActiveItem(next);
        },
        toActiveItem(item) {
            item = worker.getItem(item);
            worker.activeItem = item;
            if (item != null) {
                workerOption.onActiveItem && workerOption.onActiveItem(item);
            }
        },
    };

    return worker
};
export default {
    newItemsWorker
}