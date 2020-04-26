# ont-gateway


## 中台

核心功能：
1. 为某个租户按需部署一套某个addon的运行时环境，并且可以修改运行时配置
    1. `SetConfig(addonID, tenantID, config)`
    2. `Deploy(addonID, tenantID)`
2. 为某个租户部署的某个addon运行时环境，提供流量转发功能，并统一监控运行时状态
    1. 根据请求头部的`addonID`和`tenantID`查询路由并转发
3. 轮询`ontology`上的某个`tx`，并可靠地将结果进行`callback`
    1. 结果包含：是否存在（确定存在、超时）
    2. 对于`ONG`的`tx`，结果中可以额外带上`amount`信息
4. 记账系统
    1. 首先为每个商品创建一个`payment_config`，包括支付周期、周期对应的金额、币种、支付方式等
    2. 记录每个商品的每笔支付信息，并关联上述的`payment_config`
    3. 查询某比支付是否结清
    4. 系统将在需要续费以及欠费时自动进行`callback`
5. 钱包管理
    1. 导入钱包，并通过接口获取
6. 通用`shell`执行环境    
