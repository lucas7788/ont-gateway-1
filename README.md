# ont-gateway


## addon中台

核心功能：
1. 为某个租户按需部署一套某个addon的运行时环境，并且可以修改运行时配置
    1. `SetConfig(addonID, tenantID, config)`
    2. `Deploy(addonID, tenantID)`
2. 为某个租户部署的某个addon运行时环境，提供流量转发功能，并统一监控运行时状态
    1. 根据请求头部的`addonID`和`tenantID`查询路由并转发