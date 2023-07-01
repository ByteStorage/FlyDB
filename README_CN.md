

<img src="./assets/FlyDB-logo.png" alt="FlyDB-logo" style="width: 45%;" />

![GitHub top language](https://img.shields.io/github/languages/top/ByteStorage/flydb)   [![Go Reference](https://pkg.go.dev/badge/github.com/ByteStorage/flydb)](https://pkg.go.dev/github.com/ByteStorage/flydb)   ![LICENSE](https://img.shields.io/github/license/ByteStorage/flydb)   ![GitHub stars](https://img.shields.io/github/stars/ByteStorage/flydb)   ![GitHub forks](https://img.shields.io/github/forks/ByteStorage/flydb)   [![Go Report Card](https://goreportcard.com/badge/github.com/qishenonly/flydb)](https://goreportcard.com/report/github.com/qishenonly/flydb)![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/ByteStorage/FlyDB)![GitHub contributors](https://img.shields.io/github/contributors/ByteStorage/FlyDB)

## 

## 该项目正在迭代开发中，请勿在生产环境中使用!

**FlyDB**旨在在某些情况下作为内存键值存储(如**Redis**)的替代方案，旨在在性能和存储成本之间取得平衡。它通过优化资源分配和使用经济有效的存储介质来实现这一点。通过智能管理数据，**FlyDB**确保高效操作，同时最大限度地降低存储成本。它为需要在性能和存储成本之间取得平衡的场景提供了可靠的解决方案。

## 👋 什么是 FlyDB ?

**FlyDB**是基于高效bitcask模型的高性能键值(KV)存储引擎。它提供了快速可靠的数据检索和存储功能。通过利用bitcask模型的简单性和有效性，**FlyDB**确保了高效的读写操作，从而提高了整体性能。它提供了一种简化的方法来存储和访问键值对，使其成为需要快速响应数据访问的场景的绝佳选择。**FlyDB**对速度和简单性的关注使其成为在平衡存储成本的同时优先考虑性能的应用程序的有价值的替代方案。

## 🏁  快速入门 : FlyDB

您可以使用Go命令行工具安装FlyDB:

```GO
go get github.com/ByteStorage/FlyDB@v1.0.4
```

或者从github克隆这个项目:

```bash
git clone https://github.com/ByteStorage/FlyDB
```

## 🚀 如何使用 FlyDB ?

下面是一个如何使用Linux版本的简单示例:

> 详情请参阅 flydb/examples。

```go
package main

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/flydb"
	"github.com/ByteStorage/FlyDB/config"
)

func main() {
    	options := config.DefaultOptions
	options.DirPath = "/tmp/flydb"
	db, _ := flydb.NewFlyDB(options)

    	err := db.Put([]byte("name"), []byte("flydb-example"))
    	if err != nil {
        	fmt.Println("Put Error => ", err)
    	}


	val, err := db.Get([]byte("name"))
	if err != nil {
		fmt.Println("Get Error => ", err)
	}
    	fmt.Println("name value => ", string(val))
    
    
    	err := db.Delete([]byte("name"))
    	if err != nil {
        	fmt.Println("Delete Error => ", err)
    	}
}


```

## 🚀 性能测试

我们对V1.0.4版本的FlyDB做了一个简单的性能测试。本次测试主要针对大规模数据的读写，我们随机选取了50万条数据进行测试。

通过测试，我们发现在V1.0.4中，在50万条数据的情况下:

#### BTree 索引

PUT 性能: 572.265968ms

GET 性能: 355.943926ms

<img src="./assets/v1.0.4-btree" alt="v1.0.4-btree" style="width: 33%;"  />

#### ARTree 索引

PUT 性能: 569.610614ms

GET 性能: 297.781977ms

<img src="./assets/v1.0.4--art" alt="v1.0.4-art" style="width: 33%;"  />



如果您有更好的方法来优化读写性能，请提交您的“pr”。

## 🔮 如何联系我们?

如果您有任何疑问并想与我们联系，您可以联系我们的开发团队，我们会回复您的邮件:

团队邮箱:bytestoragecommunity@gmail.com

或者加我微信，邀请大家进入项目社区，和大牛一起交流学习。

> 添加微信请备注Github

<img src="./assets/vx.png" alt="vx" style="width: 33%;"  />

## ✅ TODO List

- [ ] 扩展数据结构支持:包括但不限于字符串、列表、散列、集合等。
- [ ] 兼容Redis协议和命令。
- [ ] 支持http服务。
- [x] 支持tcp服务。
- [x] 集成日志。
- [ ] 数据备份
- [ ] 分布式集群模型。

## 📜 版本更新信息

V1.0.4: 短期支持版本v1.0.4使用mmap来优化索引部分，与v1.0.3版本相比，显著提高了读写性能。50w数据写入速度从1.35秒加快到0.56秒，读取速度从1.06秒加快到0.355秒。

## 👀 贡献者列表

<a href="https://github.com/ByteStorage/FlyDB/graphs/contributors">

<img src="https://contrib.rocks/image?repo=ByteStorage/FlyDB" /> 
</a>

## 📝 如何贡献 ?

如果您对FlyDB有任何想法或建议，请随时在GitHub上提交“问题”或“pr”。我们欢迎您的贡献!

> 完整的规范步骤请参考：[CONTRIBUTEING](https://github.com/ByteStorage/flydb/blob/master/CONTRIBUTING.md)

## 📋 Licence

FlyDB在Apache许可下发布。请参见LICENSE文件。
