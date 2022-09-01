# transaction-matching-engine
一套基于内存的数字货币交易所撮合引擎
1. 使用grpc作为服务端框架，传输更加高效。
2. 数据结构方面，实现了自定义跳表，数据操作也更加高效。
3. 系统启动和终止都将与持久化数据交互，存储或者读取。

- 挂单：采用了异步通知，引擎只返回当前挂单成功或者失败。其实也可以将FOK,IOC类型订单做成类似深度一样的同步通知，但是效率会降低，调用侧平均等待时间会加长。【使用者可根据自己需要自行修改,示例如盘口深度查询】
- 盘口深度：查询采用了同步通知，同步开销比较大，虽然加了盘口深度快照，但涉及到快照更新，由于单跳表只能单线程操作，所以会影响到挂单效率。
- 成交推送：为了降低服务耦合，提高撮合引擎性能，使用者须将成交接入消息中间件。以订阅成交消息推送。【使用者也可根据自己需要自定义增加 grpc stream流接口，实时推送成交至业务端，与websocket类似】


### 挂单

订单类型|描述|限价单|市价单(以对手价成交)
-|-|-|-
IOC|无法立即成交的部分就撤销,订单在失效前会尽量多的成交|支持|支持
FOK|无法全部立即成交就撤销,如果无法全部成交，订单会失效|支持|支持
GTC|订单会一直有效,直到被成交或者取消|支持|支持

### 盘口

档位数据： [][3]string{用户id，价格，数量}

### 使用

- start grpc 后跟交易对参数
    ```
    go run main.go start grpc BTC-USDT,ETH-USDT
    ```

- dump 文件内配置交易对，和执行程序同目录下，dump文件内，pairs.json文件内配置交易对`[BTC-USDT,ETH-USDT]`,然后直接启动

    ```
    go run main.go start grpc
    ```

### 项目结构

- cmd 启动命令相关
- common 系统运行时组件
- engine 撮合引擎
- grpc rpc对外服务
- models 数据模型
- pool 撮合池

### orders_test.json 订单类型测试数据

- 相同的 id 表示为并行测试数据 如 `[1,2,3(1),3(2)]` 表示测试顺序有2条，分别为 `1,2,3(1)` 和  `1,2,3(2)` 

