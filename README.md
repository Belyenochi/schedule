### 工程结构
1. calculate: 参赛者只能修改此模块代码，程序编译会替代"cmd", "data", "pkg"模块为demo中对应模块
2. cmd : calculate.go为demo程序启动。score.go为数据结果评测启动
3. data: 数据模块
4. pkg: 基础功能模块，参赛者可使用，也可不使用

### 代码提交
1. fork go-demo工程，创建自己的私有工程，并添加用户信息middleware-show为Reporter权限。提交自定义的私有仓库https地址或ssh地址测试。