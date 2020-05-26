# ddxf framework

## flow chart

1. 上架流程 
    1. 卖家在ddxf公共客户端通过MP registry选择MP，
    2. TODO 查看MP设置（audit规则，手续费，挑战期）
    3. 从MP获取jsonlod格式的商品表单(mp item meta schema)，
        1. 如果卖家有data，没有data meta和dataid，跳到4
        2. 如果卖家没有token meta，跳到5
        3. 卖家通过ddxf公共客户端，通过jsonld对齐算法自动填部分字段（data meta(value)+token meta(value) + mp item meta schema(context) => partial mp item meta），未填部分手动填，填完后生成mp jsonld context+实填数据（mp item meta）。
        4. 卖家找地方存mp item meta，全量存
        5. 卖家通过ddxf公共客户端提交mp item meta和ddxf publish tx（并签名）给MP
        6. MP审核
            1. 自动审核
                1. tx的descHash与meta是否一致
                2. tx反解，得到卖家的token的Endpoint，根据tokenHash反查token meta，校验tokenHash和token meta
                3. tx反解，如果包含dataid，根据dataid反查data meta，校验hash和data meta
                4. 每个meta都可以做jsonld的完整性校验
            2. 人工审核
                1. Not include
        7. MP签名、上架
        8. 跳到5
    4. 卖家通过ddxf公共客户端自定义data meta，
        1. 找地方存（选项是卖家自己搭建服务或者onchain提供的SP），
        2. 生成dataid（构造一笔dataid tx，卖家签名后上链），如果是静态，内容需要有dataHash和metaHash+Endpoint，
        3. 返回
    5. 卖家通过ddxf公共客户端自定义token meta，
        1. 包含token操作endpoint和token meta反查endpoint
        2. 找地方存（选项是卖家自己搭建服务或者onchain提供的SP），
        3. 返回
    6. 结束
2. 浏览流程
    1. 买家通过ddxf公共客户端，选mp
    2. 从mp获得商品列表和列表项jsonld schema（context， mp item meta schema子集）
    3. 根据context渲染商品列表
    4. 单选商品，进入详情（根据mp item meta渲染）
    5. 详情页有查看源头按钮（点击后渲染data meta和token meta，复用审核逻辑）
3. 下单流程
    1. 买家通过ddxf公共客户端构造交易，签名，上链，拿到token（链上）和endpoint，找地方存（选项是卖家自己搭建服务或者onchain提供的SP X）
4. 使用流程
    1. 买家通过公共客户端连接X，选择token，找到endpoint，去Endpoint使用token
        1. 如果token带dataid，则根据dataid的endpoint获取data meta
        2. 根据data meta+token meta，获得数据类型+操作对应的客户端渲染器
        3. 通过客户端构造ddxf tx签名后发给卖家的token操作endpoint
    2. 卖家链上核销，返回结果（mvp为同步）
    3. 买家客户端渲染器渲染卖家返回结果
    4. 如果数据绑定dataid，且是静态数据源，比较datahash和返回结果的hash
        1. 如果不匹配，发起挑战（TBD）
