# 任务

1. 编写完代码
2. 添加注释并整理
3. 阅读groupcache的源码
4. 输出笔记

# TODO

1. 《分布式缓存原理 Go语言实现》要不要买？？？
2. LRU目前仅支持超过容量再淘汰，后续支持根据过期时间惰性删除key
3. raft共识算法保证服务的一致性、网络分区可容性（CP）
4. 缓存持久化，AOF、RDB？

# 功能

1. 基于LRU做缓存淘汰策略；
2. 支持集群部署，节点之间采用Protobuf序列化通信；
3. 集群的负载均衡基于哈希一致性的实现；
4. 利用Go的sync包限流，防止缓存击穿，缓存雪崩现象。

# 源码

1. build.sh用于启动gocache服务
2. cache包是gocache的视图层
3. cachepb包是protobuf的文件
4. consistenthash包是哈希一致性的实现逻辑
5. lru包是LRU算法的实现逻辑
6. singalflight包用于瞬时的流量控制