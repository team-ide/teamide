# 1.9.0

## Version 1.9.0 （2022/12/14）

发布功能

1. ES添加数据导入
   * 配置数据策略导入
   * 展示数据准备 总、成功、失败等数据
   * 展示数据插入 总、成功、失败等数据
2. ES添加索引状态信息查看
3. 通过非https打开页面，无法直接使用复制、粘贴功能的兼容处理
4. 添加Linux系统，server包，可以运行于服务端
5. 服务包发布Docker
   * docker run -itd --name teamide-21080 -p 21080:21080 -v /data/teamide/data:/opt/teamide/data teamide/teamide-server:1.9.0