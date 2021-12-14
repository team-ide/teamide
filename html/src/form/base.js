import validator from '@/form/validator.js';

let pattern = {
    email: /^([a-zA-Z0-9\w-]+)(@)[a-zA-Z0-9]+(\.)([a-zA-Z]+)$/,
    account: /^[a-zA-Z0-9_]+$/
};

let rule = {
    password: [
        {
            required: true,
            message: `密码不能为空!`
        },
        {
            minLength: 6,
            maxLength: 20,
            message: `密码长度必须大于6位小于20位!`
        },
    ],
    email: [
        {
            required: true,
            message: `邮箱不能为空!`
        },
        {
            minLength: 3,
            maxLength: 20,
            message: `邮箱长度必须大于3位小于20位!`
        },
        {
            pattern: pattern.email,
            message: `请输入正确邮箱格式!`
        },

    ],
    account: [
        {
            required: true,
            message: `账号不能为空!`
        },
        {
            minLength: 4,
            maxLength: 20,
            message: `账号长度必须大于4位小于20位!`
        },
        {
            pattern: pattern.account,
            message: `账号只能包含数字、字母、下划线!`
        },

    ],
};
export {
    pattern,
    validator,
    rule,
};