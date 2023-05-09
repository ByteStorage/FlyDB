# FlyDB

**FlyDB** 是一款简单轻量的 `Go` 语言`kv`型数据库。它提供了一组易于使用的 API，允许用户在应用程序中存储和检索数据。

## 项目正在迭代开发中，请不要用于生产环境！

## 简介

**FlyDB** 是一个快速且易于使用的基于`bitcask`的`kv`型数据库，旨在轻量和简单。使用 **FlyDB**，您可以轻松地在 `Go` 应用程序中存储和检索数据。**FlyDB** 优化了速度，这使得它非常适合需要快速数据访问的应用程序。

## 特点

**FlyDB** 的一些特点包括：

- `易于使用`：`FlyDB` 提供了一个简单直观的 API，使得存储和检索数据非常容易。
- `轻量`：`FlyDB` 设计为轻量级和高效，这使得它非常适合在资源受限的环境中使用。
- `可靠`：`FlyDB` 支持事务操作，确保在存储过程中不会丢失或损坏数据。
- `快速`：`FlyDB` 使用内存数据结构，这使得它快速响应，特别适合需要快速读写速度的应用程序。
- `可扩展`：`FlyDB` 提供了许多可配置选项，允许用户调整其性能和功能以适应不同的使用情况。

## 安装

您可以使用 Go 命令行工具安装 FlyDB：

```go
go get github.com/qishenonly/flydb
```

或者从github上clone本项目：

```bash
git clone https://github.com/qishenonly/flydb
```

## 用法

以下是一个简单的`Linux版`使用示例：

> 详细使用请看flydb/examples.

```go
package main

import (
	"fmt"
	"github.com/qishenonly/flydb"
)

func main() {
    options := flydb.DefaultOptions
	options.DirPath = "/tmp/flydb"
	db, _ := flydb.NewFlyDB(options)

	db.Put([]byte("name"), []byte("flydb-example"))

	val, err := db.Get([]byte("name"))
	if err != nil {
		fmt.Println("name value => ", string(val))
	}
}
```

## 联系

如果您有任何问题想联系我们，可以联系我们的开发者团队,我们会及时回复您的邮件：

团队邮箱：jiustudio@qq.com

或者添加我的微信，邀你进入项目社群，和大牛一起交流学习。

> 添加微信请备注 GIthub

<img src="./assets/vx.png" alt="vx" style="width: 33%;"  />

## TODO List

- [ ] 性能优化：包括但不限于索引更换、merge优化等
- [ ] 数据备份功能
- [ ] 支持HTTP服务
- [ ] 拓展数据结构的支持：包括但不限于string、list、hash、set等
- [ ] 兼容Redis协议及命令
- [ ] 分布式集群模式
- [ ] 其他待办事项...

## 贡献

如果您有任何想法或建议 FlyDB，请随时在 GitHub 上提交`issues`或`pr`。我们欢迎您的贡献！

> 完整的规范步骤请参考：[CONTRIBUTEING](https://github.com/qishenonly/flydb/blob/master/CONTRIBUTING.md)

## 许可证

FlyDB 在 Apache 许可下发布。请参见 LICENSE 文件。
