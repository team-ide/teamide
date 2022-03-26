import source from "../source";

let application = {};


application.apps = null;
application.app = null;
application.context = null;

application.modelTypeOpens = [];

application.tabs = [];
application.activeTab = null;

application.getDataTypeOptions = function (context) {
    let options = [];
    source.dataTypes.forEach(one => {
        options.push(one);
    });
    if (context && context.structs) {
        context.structs.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;

};

application.getStructOptions = function (context) {
    let options = [];
    if (context && context.structs) {
        context.structs.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;
};

application.getActionOptions = function (context) {
    let options = [];
    if (context && context.actions) {
        context.actions.forEach((one) => {
            options.push({
                value: one.name,
                text: one.name,
                comment: one.comment,
            });
        });
    }
    return options;

};
function getKeys(obj) {
    if (obj == null) {
        return [];
    }
    return Object.keys(obj);
}
application.trimObj = function (obj) {
    if (obj == null || obj == "") {
        return
    }
    for (let key in obj) {
        if (obj[key] == null || obj[key] == "") {
            delete obj[key]
        } else {
            if (Array.isArray(obj[key])) {
                application.trimArray(obj[key]);
                if (obj[key].length == 0) {
                    delete obj[key]
                }
            } else if (typeof obj[key] == "object") {
                application.trimObj(obj[key]);
                if (getKeys(obj[key]).length == 0) {
                    delete obj[key]
                }
            }
        }
    }
};
application.trimArray = function (array) {
    let needRemoves = [];
    array.forEach(one => {
        if (one == null || one == "") {
            needRemoves.push(one);
        } else {
            if (Array.isArray(one)) {
                application.trimArray(one);
                if (one.length == 0) {
                    needRemoves.push(one);
                }
            } else if (typeof one == "object") {
                application.trimObj(one);
                if (getKeys(one).length == 0) {
                    needRemoves.push(one);
                }
            }
        }
    });
    needRemoves.forEach(one => {
        array.splice(array.indexOf(one), 1)
    });
};
export default application;