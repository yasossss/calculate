# project01

##练习gRPC-go

用go实现基于grpc的任务批处理服务，client端可以向server端一次提交一组数，server端同时计算这组数的最大值，最小值和平均值。
###要求：
①利用etcd实现服务的注册发现。client端可以动态从注册中心（etcd）获取服务端的信息。
②可以运行多个server实例，client可以利用随机、轮询等算法进行负载均衡。
③服务端要求能控制goroutine的数量实现协程池的功能，基于channel进行任务的分发。
