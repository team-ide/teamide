# 2.2.1

## Version 2.2.1 （2023/5/5）

发布功能

1. 添加thrift模块，配置thrift文件目录
2. 解析目录thrift文件，展示所有service和method
3. 双击method，填写服务地址，执行参数，模拟执行该方法
4. 执行thrift方法，展示写入、读取、总执行时间、异常等信息
5. thrift方法执行，json中长整型精度丢失问题修复
6. thrift服务列表添加搜索，执行页面添加服务基本信息
7. 方法调用参数默认初始化属性
8. thrift添加性能测试，测试报告，图表展示TPS、AVG、Min、Max、T90、T99等
9. thrift加载的服务、方法根据名称排序
10. thrift接口并发测试可输出markdown文件，可在线预览
11. thrift添加ProtocolFactory类型、Buffered、Framed、超时时间配置等