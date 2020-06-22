# openkg链改

## 方案描述

0. 为每个用户，调用`generateOntIdByUserId`接口，生成OntID
1. 发布时，对于每个数据集，调用`publish`接口，发布
2. 下载时，对应的数据集，调用`buyAndUse`接口，下载