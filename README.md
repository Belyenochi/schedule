### 赛题更新
1. 赛题样例、数据、测评流程: <https://code.aliyun.com/middleware-contest-2020/django>

### 工程结构
1. calculate: 参赛者实现静态布局和动态迁移功能模块。
2. pkg: 基础功能模块，参赛者可使用，也可不使用。
3. data: 数据模块。
4. cmd : 程序启动运行模块。

### 工程注意
1. 参赛者只能修改calculate模块代码。程序在测评编译阶段会替代"cmd", "data", "pkg"模块为demo中对应模块

### 程序运行
1. calculate.go 为demo程序启动。