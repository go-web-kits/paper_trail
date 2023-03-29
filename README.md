# PaperTrail

Track Changes to Your Models  
审计：自动记录基于模型的数据操作日志

Maintainer: @will.huang

## Features

1. 全自动无感知（基于 GORM callbacks）
2. 支持指定需要追踪的字段，追踪字段未发生变化将跳过
3. 90% 测试覆盖率，支持 `Jsonb` 字段
4. 提供手动创建审计的接口
5. 提供获取审计记录的接口
6. 记录值变化集
7. 提供版本回退功能（TODO）

## Setup

1. 注册回调
    ```go
    _ = paper_trail.Register()
    ```
2. 模型声明（打上 `trail:"true"` tag 即表明为追踪字段，无则被忽略）
    ```go
    type User struct {
        Name string `trail:"true"`
    
        paper_trail.EnableTrail
    }
    ```

That's all
