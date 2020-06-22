# openkg链改

## 方案描述

0. 统一以`JSON POST`方式交互
1. 为每个用户，调用`/generateOntIdByUserId`接口，生成OntID
2. 发布时，对于每个数据集，调用`/publish`接口，发布
3. 下载时，对应的数据集，调用`/buyAndUse`接口，下载

### 环境

1. 测试环境： 