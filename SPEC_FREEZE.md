# 规格冻结表

batch_id: solidstate-battery-20260824  
business_direction: 固态电池  
project_variant: backend  
foundation_profile: compact_10  
task_count: 10

| 领域 | 冻结规格 |
|---|---|
| 业务边界 | 电芯批次、材料配方、模块、试验运行、测量、质量冻结、发布结论 |
| 持久化 | PostgreSQL 16；禁止以内存 map 代替生产状态 |
| migration | migrations/001_schema.sql、002_workflow.sql、003_auth_audit.sql，按版本表顺序执行且可重复启动 |
| 事务 | 试验创建+样品分配+审计、质量冻结+状态变更、发布+审计均为跨实体事务 |
| 状态机 | lot draft→qualified→released/hold；run planned→running→paused→succeeded/failed；hold open→cleared |
| 并发 | 行锁、版本号乐观锁、唯一幂等键、每批次单活动运行约束 |
| context | request-id、截止时间、取消信号从 HTTP handler 传入 service、repository、worker |
| worker | 轮询过期运行、有限重试、失败落库、优雅停止、重启后恢复 pending 作业 |
| 错误传播 | domain sentinel + errors.Is/As；HTTP 统一 JSON 错误和 request id，不吞错 |
| HTTP | health/ready、登录/登出、批次、运行、测量、质量冻结、审计查询；JSON 校验和分页 |
| 身份权限 | session token 可撤销、过期拒绝；operator 可执行试验，reviewer 可确认质量，admin 可管理用户 |
| Docker | 真实 ./cmd/server 入口；linux/amd64 与 linux/arm64 镜像独立构建并 inspect、启动、health/ready |
| 测试 | domain 状态机、service 事务、HTTP 合约、PostgreSQL 集成、worker 重试/取消/恢复、race |
| 规模 | compact_10：生产 Go ≥2000 行、≥20 文件、≥10 package；测试 Go ≥1500 行 |
| 禁止题材 | 电商、支付、社交、医疗、博彩、加密货币、武器控制等 |
| 后续容量 | 10 个彼此独立运行时边界：认证、权限、幂等、批次状态、并发版本、事务回滚、测量校验、冻结解除、worker 恢复、审计传播 |

本阶段只冻结无已知 Bug 基线，不创建题面、私测、根因材料或任务分支。
