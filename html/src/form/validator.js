let cloneForm = function (form) {
    form = form || [];
    let res = {};
    Object.assign(res, form);
    res.fields = [];
    if (form.fields) {
        form.fields.forEach(field => {
            let f = {};
            res.fields.push(f);
            Object.assign(f, field);
            f.rules = [];
            if (field.rules) {
                field.rules.forEach(rule => {
                    let r = {};
                    f.rules.push(r);
                    Object.assign(r, rule);
                });
            }
        });
    }

    return res;
};
let buildFormValidator = function (form) {
    let validatorForm = cloneForm(form);


    validatorForm.fields.forEach(field => {
        if (isEmpty(field.name)) {
            return;
        }
        field.valid = undefined;
        field.required = false;
        field.validMessage = undefined;
        let rules = field.rules || [];
        rules.forEach(rule => {
            if (rule.required) {
                field.required = true;
            }
        });

        field.validate = function (data) {
            return validateField(data, field);
        };
    });

    validatorForm.validate = function (data) {
        return validateFields(data, validatorForm.fields)
    };
    validatorForm.validateAll = function (data) {
        return validateFields(data, validatorForm.fields, true)
    };
    validatorForm.validateReset = function (data) {
        return validateReset(data, validatorForm.fields)
    };
    validatorForm.newDefaultData = function () {
        let data = {};
        validatorForm.fields.forEach(field => {
            if (isEmpty(field.name)) {
                return;
            }
            data[field.name] = field.defaultValue || null;
        });
        return data;
    };

    return validatorForm;
};

let validateFields = function (data, fields, all) {
    validateReset(fields);
    return new Promise((resolve, reject) => {
        if (fields == null || fields.length == 0) {
            resolve({
                valid: true,
            })
            return
        }
        let errors = [];
        let process = function (index) {
            if (index >= fields.length) {
                resolve({
                    valid: errors.length == 0,
                    errors: errors,
                })
                return;
            }
            if (!all && errors.length > 0) {
                resolve({
                    valid: errors.length == 0,
                    errors: errors,
                })
                return;
            }
            let field = fields[index];
            validateField(data, field).then(valid => {
                if (!valid) {
                    errors.push(field);
                }
                process(index + 1);
            }).catch(err => {
                reject(err);
            });
        }
        process(0);
    })
};

// new Promise((resolve, reject)=>{

// })
let validateField = function (data, field) {
    return new Promise((resolve, reject) => {
        if (field.type == 'jsonView') {
            resolve(true)
            return;
        }
        let value = data[field.name];
        if (value != null) {
            if (field.isNumber) {
                data[field.name] = Number(value);
            } else if (field.type == 'json') {
                let jsonValue = null;
                if (field.jsonStringValue != "") {
                    try {
                        jsonValue = JSON.parse(field.jsonStringValue);
                    } catch (error) {
                        try {
                            jsonValue = eval("(" + field.jsonStringValue + ")");
                        } catch (error2) {
                            field.valid = false;//????????? ????????????
                            field.validMessage = error;
                            resolve(false)
                            return;
                        }
                    }
                }
                data[field.name] = jsonValue;
            }
        }
        let rules = field.rules || [];
        let valid = true;
        let process = function (index) {
            if (index >= rules.length) {
                resolve(valid)
                return;
            }
            if (!valid) {
                resolve(valid)
                return;
            }
            let rule = rules[index];
            validateRule(data, field, rule).then(res => {
                valid = res;
                process(index + 1);
            }).catch(err => {
                reject(err);
            });
        }
        process(0);
    })

};

let validateRule = function (data, field, rule) {

    return new Promise((resolve, reject) => {
        let value = data[field.name];
        let valid = true;
        // required ??????
        // pattern ??????
        // range ??????
        // length ??????
        // enum ????????????
        if (valid && rule.required && isEmpty(value)) {
            valid = false;
        }
        if (valid && rule.length && ('' + value).length > rule.length) {
            valid = false;
        }
        if (valid && rule.minLength && ('' + value).length < rule.minLength) {
            valid = false;
        }
        if (valid && rule.maxLength && ('' + value).length > rule.maxLength) {
            valid = false;
        }
        if (valid && rule.min && value < rule.min) {
            valid = false;
        }
        if (valid && rule.max && value > rule.max) {
            valid = false;
        }
        if (valid && rule.pattern && !rule.pattern.test(value)) {
            valid = false;
        }
        let msg = null;
        let process = function () {
            if (valid) {
                field.valid = true;//????????? ????????????
                field.validMessage = null;
            } else {
                msg = msg || rule.message || field.message;
                field.valid = false;//????????? ????????????
                field.validMessage = msg;
            }
            resolve(valid);
        }
        if (valid && rule.validate) {
            rule.validate().then((resValid, resMsg) => {
                valid = resValid;
                msg = resMsg;
                process();
            }).then(err => {
                reject(err);
            });
        } else {
            process();
        }
    })

};
let validateReset = function (fields) {
    fields.forEach(field => {
        field.valid = undefined;
        field.validMessage = undefined;
    });
};

let isEmpty = function (arg) {
    if (arg == null || arg == "") {
        return true;
    }
    return false;
};
let isTrimEmpty = function (arg) {
    if (arg == null) {
        return true;
    }
    return isEmpty(('' + arg).trim());
};

let isNotTrimEmpty = function (arg) {
    return !isTrimEmpty(arg);
};
export default {
    buildFormValidator
};