# 2.6.33

## Version 2.6.33 （2025/1/14）

发布功能

* 数据库 `date、datetime、timestamp` 类型 数据插入、数据更新 转为 时间 插入
  * 需要满足 `date、datetime、timestamp` 数据类型
  * 传入的值是 数值类型 且是毫秒或秒的时间戳
  * `date` 类型 直接插入 日期 字符串 "2006-01-02" 格式 
  * 插入、更新 转换的 SQL 语句 `datetime、timestamp` 转换为 日期 字符串 "2006-01-02 15:04:05" 格式 
