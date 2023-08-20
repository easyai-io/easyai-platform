easyai-platform
=====

基于k8s和云原生生态的机器学习平台， 目标是提供一套简单易用、功能完整的机器学习平台，包括算力管理、作业管理、模型服务发布等功能。

# 项目背景

## 设计目标

1. 面向云原生：支持多集群和不同云厂商，支持k8s 1.19+，不引入非云原生的依赖
2. 简单易用：对用户屏蔽docker/k8s的底层细节，降低使用门槛
3. 易扩展：提供web ui、client cli、open api、sdk等多种接入方式
4. 平台化：用户和作业管理，提供监控、日志、告警等解决方案

## 技术栈

1. 后端：golang + gin + ent-go + wire + swagger-go + docsify
2. client命令行工具：cobra + resty + docker/podman-sdk
3. web前端：vue/reactjs

## 其他依赖

1. MySQL + Redis， for http server
2. nfs server, minIO, for file storage 
3. k8s，kubernetes云原生集群管理和容器编排系统，paas核心
4. kubeflow training operator/volcano-job, AI工作负载管理
5. volcano(scheduler+controller)，k8s批处理调度器，支持gang-scheduler等调度策略
6. nginx-ingress-controller, 提供服务访问路由

# 设计方案

请参考[总体设计文档](./docs/design/README.md)

# 开发

## 项目结构

```
todo 待补充
```

## Build && Run && Deploy

```shell
todo 待补充
```

## contribute

欢迎大家以各种形式参与到项目中来，包括但不限于：

+ 新的功能需求
+ 新的设计方案
+ bug反馈/修复
+ 文档完善

