# go-retry

[English](./README.md) | [中文](./README-ZH.md)

## 简介

一个简单的重试方法。

现实工作中，经常有需求针对 `IO操作` 、 `API交互` 、 `远程调用` 等可能会失败的过程进行重试，于是就有当前这个小项目。

可以简单配置充实次数、间隔时间、停止充实的条件等，如果经过充实最终结果不符合预期，将抛出 `goretry.ErrMaxApt` 错误。

## 如何获取

下载模块后导入使用。

```shell
go get -u github.com/weston-shih/go-retry
```

或者

```go
import "github.com/weston-shih/go-retry"
```

## 使用

首先进行重试选项配置。

```go
package main

import (
   "log"

   "github.com/weston-shih/go-retry"
)

func main() {
   // First of all, create a retry option, 
   // then set max attempt number and interval.
   op := NewRetryOption()
   
   // Set retry attempts
   _, err := op.SetAttempt(3)
   if err != nil {
      log.Fatal("Set attempt number ended with failure: ", err)
   }

   // Set retry attempt intervals
   _, err := op.SetBackoff(1)
   if err != nil {
      log.Fatal("Set retry interval ended with failure: ", err)
   }

   // Or you can force the configuration,
   // and invalid input will result in a panic
   op = NewRetryOption().MustSetAttempt(3).MustSetBackoff(1)
}
```

可以自定义判断条件，即当 `Judgment` 方法返回结果为 `true` 时候，认为已达到预期，重试流程将会终止。

```go
// e.g. the retry will continue if err equals ErrTest.
var ErrTest = errors.New("Just a test.")
op.SetJudgment(func(err ...interface{}) bool { return err[0] == ErrTest })
```

- 仅有错误返回值的函数使用重试

  ```go
   var (
      ErrOdd  = errors.New("Odd number")
      ErrEven = errors.New("Even number")
   )
   got := op.ReDo(
         func() error {
            if mod := test.seed % 2; mod == 0 {
               return ErrEven
            }
            return ErrOdd
         })
  ```

- 含有数据与错误返回值的函数使用重试

   ```go
   import "time"

   // Set a judgment condition to get a even number.
   op.SetJudgment(func(arg ...interface{}) bool { return arg[0] != 1 })
   got, err := op.ReTry(
      func() (interface{}, error){
         mod := time.Now().UTC().Second() % 2
         return mod, nil
      }
   )
   // Need to check if it fails after max retry attempts
   if err != nil {
      log.Fatal("Failed after max retry: ", err)
   }
   ```

## 反馈

使用中遇到任何问题欢迎反馈，将尽快跟进。

同时也欢迎一起优化代码：）
