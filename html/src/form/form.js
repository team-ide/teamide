import { pattern, rule } from '@/form/base.js';

let login = {
    fields: [
        {
            label: "账号",
            name: "account",
            rules: rule.account,
        },
        {
            label: "密码",
            name: "password",
            type: "password",
            rules: rule.password,
        }
    ],
};

let register = {
    fields: [
        {
            label: "名称",
            name: "name",
            rules: [
                {
                    required: true,
                    message: `名称不能为空!`
                },
            ],
        },
        {
            label: "账号",
            name: "account",
            rules: rule.account,
        },
        {
            label: "邮箱",
            name: "email",
            type: "email",
            rules: rule.email,
        },
        {
            label: "密码",
            name: "password",
            type: "password",
            rules: rule.password,
        }
    ],
};

export default {
    login,
    register,
};