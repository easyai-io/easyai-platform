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

# 文档

**开发相关**
- [整体架构](./docs/design/easyai_platform.md)
- [设计文档](./docs/design)
- [产品文档](./docs/about/todo.md)
- [部署和运维](./docs/about/todo.md)

**使用相关**
- [快速开始](./docs/quick-start/quick_start.md)
- [开发环境使用](./docs/about/todo.md)
- [提交训练作业](./docs/about/todo.md)
- [client命令行工具](./docs/about/todo.md)

# 开发

## 项目结构

```
.
├── CHANGELOG.md
├── Dockerfile
├── Makefile 
├── cmd
│   ├── client # client命令行入口
│   └── server # 服务端启动入口
├── docs 用户文档和设计文档
├── internal
│   ├── client # client业务代码
│   └── server # 服务端业务代码
├── pkg
│   ├── apischema # openapi schema定义
└── www # 前端静态资源
```

## Build && Run && Deploy

```shell
# build

make lint
make app

# run
make run

# 打包镜像
make image
make image-push

# 在k8s中升级
make upgrade
```

## 项目规划

- Web前端
  - [] 登录界面
  - [] 作业管理
  - [] 在线开发
  - ...
- 服务端
  - [] 作业管理
    - [] 作业定义
    - [] crud接口
    - [] 作业编排(单机+tf分布式)
    - [] 作业状态机(单机+tf分布式)
    - [] 作业调度
    - [] 作业存储支持
  - [] 资源管理
  - [] 用户管理
  - [] 集群管理
  - [] go-sdk
  - [] ...
- client命令行工具
  - [] 身份认证
  - [] 提交作业
  - [] 构建镜像
  - [] ...

## contribute

欢迎大家以各种形式参与到项目中来，包括但不限于：

+ 新的功能需求
+ 新的设计方案
+ bug反馈/修复
+ 文档完善

