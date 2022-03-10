# 阿里云函数计算基于 Custom Golang Event 函数的钉钉机器人

本目录提供一个钉钉机器人示例，该机器人基于golang Custom Runtime（非函数计算官方内置runtime）

如何使用：

0. 创建您的机器人，获得webhook地址，请您参考[钉钉官方文档](https://open.dingtalk.com/document/robots/custom-robot-access)
1. 修改robot文件，填入机器人的webhook地址
2. 执行`s deploy` 部署
3. 执行`s invoke` 向您的钉钉群发送一个机器人消息
4. 您可以基于本例子，修改函数的触发器和代码，定制您自己的机器人
