<img src="./assets/FlyDB-logo.png" alt="FlyDB-logo" style="display: block; margin: 0 auto; width: 45%;" />

![GitHub top language](https://img.shields.io/github/languages/top/ByteStorage/flydb)   [![Go Reference](https://pkg.go.dev/badge/github.com/ByteStorage/flydb)](https://pkg.go.dev/github.com/ByteStorage/flydb)   ![LICENSE](https://img.shields.io/github/license/ByteStorage/flydb)   ![GitHub stars](https://img.shields.io/github/stars/ByteStorage/flydb)   ![GitHub forks](https://img.shields.io/github/forks/ByteStorage/flydb)   [![Go Report Card](https://goreportcard.com/badge/github.com/qishenonly/flydb)](https://goreportcard.com/report/github.com/qishenonly/flydb)
## The project is under iterative development, please do not use it in production environment!

English | [中文](https://github.com/ByteStorage/flydb/blob/master/README_CN.md)

**FlyDB** aims to serve as an alternative to in-memory key-value storage (such as **Redis**) in some cases, aiming to strike a balance between performance and storage cost. It does this by optimizing resource allocation and using cost-effective storage media. By intelligently managing data, **FlyDB** ensures efficient operations while minimizing storage costs. It provides a reliable solution for scenarios that require a balance between performance and storage costs.

## 👋 What is FlyDB ?

**FlyDB** is a high-performance key-value (KV) storage engine based on the efficient bitcask model. It offers fast and reliable data retrieval and storage capabilities. By leveraging the simplicity and effectiveness of the bitcask model, **FlyDB** ensures efficient read and write operations, resulting in improved overall performance. It provides a streamlined approach to storing and accessing key-value pairs, making it an excellent choice for scenarios that require fast and responsive data access. **FlyDB's** focus on speed and simplicity makes it a valuable alternative for applications that prioritize performance while balancing storage costs.

## 🏁  Fast Start: FlyDB 

You can install FlyDB using the Go command line tool:

```GO
go get github.com/ByteStorage/flydb
```

Or clone this project from github:

```bash
git clone https://github.com/ByteStorage/flydb.git
```

## 🚀 How to use FlyDB ?

Here is a simple example of how to use the Linux version:

> See flydb/examples for details.

```go
package main

import (
	"fmt"
	"github.com/ByteStorage/flydb"
)

func main() {
    	options := flydb.DefaultOptions
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

## 🔮 How to contact us ?

If you have any questions and want to contact us, you can contact our developer team, we will reply to your email:

Team Email: bytestorage@163.com

Or add my wechat, invite you to enter the project community, and code masters together to exchange learning.

> Add wechat please comment GIthub

<img src="./assets/vx.png" alt="vx" style="width: 33%;"  />

## 📜 TODO List

- [ ] Extended data structure support: including but not limited to string, list, hash, set, etc.
- [ ] Compatible with Redis protocols and commands.
- [ ] Support http services.
- [ ] Support tcp services.
- [x] Log aggregation
- [ ] Data backup
- [ ] Distributed cluster model.

## How to contribute ?

If you have any ideas or suggestions for FlyDB, please feel free to submit 'issues' or' pr 'on GitHub. We welcome your contributions!

> Please refer to the complete specification procedure：[CONTRIBUTEING](https://github.com/ByteStorage/flydb/blob/master/CONTRIBUTING.md)

## Licence

FlyDB is released under the Apache license. For details, see LICENSE file.
