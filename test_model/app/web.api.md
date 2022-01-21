## Web Server Api

### 访问路径
> http://127.0.0.1:20120/demo.api/

### Token 验证

#### 验证路径
> /**

#### 忽略路径
> /token/create , /file/download/** , /file/open/**

#### 请求头（Header）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|access_token||Y|

### Token创建

#### 接口功能

> Token创建

#### 请求地址

> /token/create

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "account":account, // 字符串(String,string)
  "password":password // 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":value
}
```


### file/download

#### 请求地址

> /file/download/{*path}

#### 请求方式

> GET

#### 路径参数（Path）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|path||Y|

#### 下载文件

```
文件流
```


### file/open

#### 请求地址

> /file/open/{*path}

#### 请求方式

> GET

#### 路径参数（Path）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|path||Y|

#### 打开文件

```
文件流
```


### file/upload

#### 请求地址

> /file/upload

#### 请求方式

> POST

#### 表单文件（Form）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|file||Y|

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{
    "name":name, // 字符串(String,string)
    "type":type, // 字符串(String,string)
    "path":path, // 字符串(String,string)
    "dir":dir, // 字符串(String,string)
    "size":size, // 长整型(long,int64)
    "absolutePath":absolutePath // 字符串(String,string)
  }
}
```


### space/insert

#### 请求地址

> /space/insert

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "name":name, // 空间名称, 字符串(String,string)
  "parentId":parentId, // 空间父ID, 长整型(long,int64)
  "createUserId":createUserId, // 创建用户ID, 长整型(long,int64)
  "createUserName":createUserName // 创建用户名称, 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 空间信息
    "spaceId":spaceId, // 空间ID, 长整型(long,int64)
    "name":name, // 空间名称, 字符串(String,string)
    "parentId":parentId, // 空间父ID, 长整型(long,int64)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createUserId":createUserId, // 创建用户ID, 长整型(long,int64)
    "createUserName":createUserName, // 创建用户名称, 字符串(String,string)
    "updateUserId":updateUserId, // 更新用户ID, 长整型(long,int64)
    "updateUserName":updateUserName, // 更新用户名称, 字符串(String,string)
    "deleteUserId":deleteUserId, // 删除用户ID, 长整型(long,int64)
    "deleteUserName":deleteUserName, // 删除用户名称, 字符串(String,string)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


### Token创建

#### 接口功能

> Token创建

#### 请求地址

> /token/create

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "account":account, // 字符串(String,string)
  "password":password // 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":value
}
```


### token/validate

#### 请求地址

> /token/validate

#### 请求方式

> POST

#### 请求头（Header）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|access_token||Y|

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 访问信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "userName":userName, // 用户名称, 字符串(String,string)
    "timestamp":timestamp // 时间戳, 长整型(long,int64)
  }
}
```


### user/batchInsert

#### 请求地址

> /user/batchInsert

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "users":[{ // 用户信息
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email // 邮箱, 字符串(String,string)
  }],
  "password":password // 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":[{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }]
}
```


### user/delete

#### 请求地址

> /user/delete

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "userId":userId // 长整型(long,int64)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":value
}
```


### user/get

#### 请求地址

> /user/get

#### 请求方式

> GET

#### 请求参数（URL）

|字段名称    |字段说明   |必填   |
| ----------|:--------:|:--------:|
|userId||Y|

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


### user/insert

#### 请求地址

> /user/insert

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "name":name, // 用户名称, 字符串(String,string)
  "account":account, // 账号, 字符串(String,string)
  "email":email // 邮箱, 字符串(String,string)
  "password":password // 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


### user/selectCount

#### 请求地址

> /user/selectCount

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "name":name // 用户名称, 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


### user/selectList

#### 请求地址

> /user/selectList

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "name":name // 用户名称, 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


### user/selectOne

#### 请求地址

> /user/selectOne

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "userId":userId // 长整型(long,int64)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":value
}
```


### user/selectPage

#### 请求地址

> /user/selectPage

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "name":name // 用户名称, 字符串(String,string)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{
    "pageNumber":pageNumber, // 长整型(long,int64)
    "pageSize":pageSize, // 长整型(long,int64)
    "totalPage":totalPage, // 长整型(long,int64)
    "totalSize":totalSize, // 长整型(long,int64)
    "list":[{ // 用户信息
      "userId":userId, // 用户ID, 长整型(long,int64)
      "name":name, // 用户名称, 字符串(String,string)
      "account":account, // 账号, 字符串(String,string)
      "email":email, // 邮箱, 字符串(String,string)
      "activedState":activedState, // 激活状态, 字节型(byte,int8)
      "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
      "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
      "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
      "createTime":createTime, // 创建时间, 字节型(Date,date)
      "updateTime":updateTime // 修改时间, 字节型(Date,date)
    }]
  }
}
```


### user/update

#### 请求地址

> /user/update

#### 请求方式

> POST

#### 请求数据（JSON）

```json

{
  "userId":userId, // 用户ID, 长整型(long,int64)
  "name":name, // 用户名称, 字符串(String,string)
  "account":account, // 账号, 字符串(String,string)
  "email":email, // 邮箱, 字符串(String,string)
  "activedState":activedState, // 激活状态, 字节型(byte,int8)
  "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
  "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
  "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
  "updateTime":updateTime // 修改时间, 字节型(Date,date)
}

```

#### 返回数据（JSON）

```json

{
  "code": "0", // 错误码
  "msg": "", // 错误信息
  "value":{ // 用户信息
    "userId":userId, // 用户ID, 长整型(long,int64)
    "name":name, // 用户名称, 字符串(String,string)
    "account":account, // 账号, 字符串(String,string)
    "email":email, // 邮箱, 字符串(String,string)
    "activedState":activedState, // 激活状态, 字节型(byte,int8)
    "lockedState":lockedState, // 锁定状态, 字节型(byte,int8)
    "enabledState":enabledState, // 启用状态, 字节型(byte,int8)
    "deletedState":deletedState, // 删除状态, 字节型(byte,int8)
    "createTime":createTime, // 创建时间, 字节型(Date,date)
    "updateTime":updateTime // 修改时间, 字节型(Date,date)
  }
}
```


