export default {
    newItemsWorker(obj) {
        return Object.assign(obj || {}, {
            items: [],
            activeItem: null,
            getItem(item) {
                if (item == null) {
                    return
                }
                let res = null;
                this.items.forEach(one => {
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
                let find = this.getItem(item);
                if (find) {
                    return -1
                }
                return this.items.indexOf(find);
            },
            getActiveItemIndex() {
                return this.items.indexOf(this.activeItem);
            },
            addItem(item, before) {
                let find = this.getItem(item);
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
                    beforeIndex = this.getItemIndex(before)
                }
                if (beforeIndex >= 0) {
                    this.items.splice(beforeIndex + 1, 0, item);
                } else {
                    this.items.push(item)
                }
            },
            toDeleteOtherItem(item) {
                let list = [];
                this.items.forEach((one) => {
                    if (one != item) {
                        list.push(one);
                    }
                });
                this.toRemoveItems(list);
                if (this.activeItem != item) {
                    this.toActiveItem(item);
                }
            },
            toRemoveAll() {
                let list = [];
                this.items.forEach((one) => {
                    list.push(one);
                });
                this.toRemoveItems(list);
            },
            toRemoveItems(list) {
                list.forEach((one) => {
                    let index = this.items.indexOf(one);
                    if (index >= 0) {
                        this.items.splice(index, 1);
                        this.onRemoveItem && this.onRemoveItem(one);
                    }
                });
            },
            toRemoveItem(item) {
                let find = this.getItem(item);
                if (find == null) {
                    return;
                }
                let index = this.items.indexOf(find);
                this.items.splice(index, 1);
                this.onRemoveItem && this.onRemoveItem(find);
                if (find == this.activeItem) {
                    let nextIndex = index - 1;
                    if (nextIndex < 0) {
                        nextIndex = 0;
                    }
                    this.toActiveItem(this.items[nextIndex]);
                }
            },
            activeNextItem(item) {
                let find = this.getItem(item);
                if (find == null) {
                    return;
                }
                let index = this.items.indexOf(find);
                let next = null;
                if (index >= 0) {
                    if (this.items[index + 1]) {
                        next = this.items[index + 1];
                    } else {
                        next = this.items[index - 1];
                    }
                }
                this.toActiveItem(next);
            },
            toActiveItem(item) {
                item = this.getItem(item);
                this.activeItem = item;
                if (item != null) {
                    this.onActiveItem && this.onActiveItem(item);
                }
            },
        });
    },
}