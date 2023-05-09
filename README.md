# FlyDB

English | [简体中文](https://github.com/qishenonly/flydb/blob/master/README_CN.md)

FlyDB is a simple and lightweight kv-type database written in Go. It provides an easy-to-use API that allows users to store and retrieve data in their applications.

## **Note**: This project is currently under iterative development and should not be used in production environments!

## Introduction

FlyDB is a fast and easy-to-use kv-type database based on bitcask, designed to be lightweight and simple. With FlyDB, you can easily store and retrieve data in your Go applications. FlyDB is optimized for speed, making it ideal for applications that require fast data access.

## Features

Some features of FlyDB include:

- Easy to use: FlyDB provides a simple and intuitive API that makes storing and retrieving data very easy.

- Lightweight: FlyDB is designed to be lightweight and efficient, making it ideal for use in resource-constrained environments.
- Reliable: FlyDB supports transactional operations, ensuring that data is not lost or corrupted during the storage process.
- Fast: FlyDB uses in-memory data structures, making it fast and responsive, especially for applications that require fast read/write speeds.
- Scalable: FlyDB provides many configurable options that allow users to adjust its performance and features to suit different usage scenarios.

## Installation

You can install FlyDB using the Go command line tool:

```go
go get github.com/qishenonly/flydb
```

Or clone this project from GitHub:

```bash
git clone https://github.com/qishenonly/flydb
```

## Usage

Here is a simple Linux usage example:

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

For more detailed usage, please refer to the flydb/examples directory.

## Contact

If you have any questions or want to contact us, you can reach out to our development team and we will respond to your email promptly:

Team email: jiustudio@qq.com

Or add my wechat account and invite you to join the project community to exchange and learn with Daniu.

> Note when adding wechat : Github

<img src="./assets/vx-1683193364673-1.png" alt="vx" style="width:33%;"  />

## Contributions

If you have any ideas or suggestions for FlyDB, please feel free to submit issues or pull requests on GitHub. We welcome your contributions!

For full contributing guidelines, please refer to [CONTRIBUTING](https://github.com/qishenonly/flydb/blob/master/CONTRIBUTING.md)

## License

FlyDB is released under the Apache license.  please see the LICENSE file.
