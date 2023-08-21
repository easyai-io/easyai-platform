# Easyai 云原生机器学习平台

> 致力于提供一站式机器学习平台解决方案

**文档细节待补充和完善**

## Features/功能特性

1. 管理和调度计算资源 
   1. cpu/gpu异构资源
   2. 丰富的调度策略
   3. 多种机器学习框架的支持
2. 低学习成本, 友好的使用体验
   1. 尽可能的屏蔽docker/k8s等底层技术细节
   2. 简单易用的web portal
   3. 支持jupyter/vscode等web交互式开发工具和ssh-remote远程开发
   4. 完善的文档和指南
   5. 提供命令行CLI、open API接口和SDK-golang
3. 提供job、develop、serving、pipeline的管理
4. 面向云原生，不绑定云厂商，支持多种云厂商，适配、改造、迁移的成本极低；支持多集群，兼容k8s 1.18~1.25
5. 提供PaaS平台通用能力，如监控、日志、告警等

## Architecture/架构

![architecture](static/architecture.png)