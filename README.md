<h1 align="center">online chat</h1>
<div align="center">
    <strong>online chat use go and redis</strong>
</div>


## 简介

使用 go 开发的基于 redis 发布订阅功能的即时在线聊天软件。
出于 redis 发布订阅的特性，所以只能接收在程序运行期间发送的消息。（不会有非运行期间的历史消息

处于安全考虑，对公网开放的 redis 不要使用 root 启动，并且配置访问密码。

## 开发计划

* 多频道切换
* 多 redis 切换
* 保存读取历史消息记录

欢迎提 issue 和 pr

## 鸣谢

[A powerful little TUI framework 🏗](https://github.com/charmbracelet/bubbletea)
