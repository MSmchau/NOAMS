\# 网络运维自动化管理系统实施文档

  
\## 一、项目背景与现状分析

\### 1.1 现有网络环境概况

当前网络环境由中心机房、宿舍楼两个层级构成：

\*\*中心机房设备：\*\*  
\- 核心交换机：H3C S10504（1台）  
\- 无线控制器AC：H3C EWPXM1MAC0F（1台）

\*\*宿舍楼设备（总计15栋，每栋配置相同）：\*\*  
\- 楼栋本体交换机：H3C S5560X-30F-EI（1台/栋，共15台）  
\- PoE接入交换机：H3C S5130S-52P-EI（6台/栋，共90台）  
\- 有线接入交换机：H3C S5130S-52P-EI（6台/栋，共90台）  
\- 无线AP：H3C WTU420（120台/栋，总计1,800台）

\*\*网络设备总量统计：\*\*

| 设备类型 | 数量 | 管理方式 |  
|---|---|---|  
| 核心交换机 S10504 | 1 | SSH / Web |  
| 无线控制器 AC | 1 | Web |  
| 本体交换机 S5560X | 15 | SSH / Web |  
| PoE交换机 S5130S | 90 | SSH / Web |  
| 有线交换机 S5130S | 90 | SSH / Web |  
| 无线AP WTU420 | 1,800 | 仅AC Web |

\> \*\*合计：1,997 台设备\*\*

\### 1.2 现有运维痛点分析

\*\*痛点一：AP管理高度依赖AC Web界面\*\*——1,800台AP的状态监测、故障定位、上线/下线/重启操作均只能通过AC的Web管理页面进行，无法通过SSH直接访问AP设备，操作路径单一且效率低下。

\*\*痛点二：缺乏统一的设备管理入口\*\*——核心、AC、本体交换机、接入交换机分散管理，运维人员需要记忆多台设备的IP地址和登录凭据，分别登录查看状态。传统手动巡检需逐台登录设备，耗时耗力且准确性难以保障。

\*\*痛点三：自动化能力缺失\*\*——缺乏定时巡检、配置自动备份、告警统一汇聚等自动化能力。在拥有近2,000台设备的网络中，如果每台设备手动巡检需2分钟，全量巡检将耗时超过66小时。

\*\*痛点四：知识孤岛与人员风险\*\*——运维经验高度依赖个人，设备配置变更无记录追溯，存在人员流动带来的运维断层风险。

\*\*痛点五：故障响应被动\*\*——当前缺乏集中告警机制，故障感知依赖人工巡检或用户报修，响应速度慢、被动挨打。

  
\## 二、系统建设目标与原则

\### 2.1 建设目标

构建一套统一的网络设备运维管理平台，实现以下核心目标：

1\. \*\*统一纳管\*\*：将核心、AC、本体交换机、接入交换机、AP全部纳入统一平台管理，消除信息孤岛。  
2\. \*\*自动化巡检\*\*：通过SSH自动采集设备CPU、内存、接口状态、硬件信息、运行时长，替代人工逐台登录。  
3\. \*\*配置集中管理\*\*：一键备份设备运行配置，支持历史版本追溯与一键回滚。  
4\. \*\*实时状态监控\*\*：Dashboard大屏展示全网设备在线状态与资源使用率，故障一目了然。  
5\. \*\*告警统一汇聚\*\*：设备异常（离线、CPU/内存超阈值、端口Down等）自动触发告警，支持分级通知。  
6\. \*\*定时任务调度\*\*：内置调度器，支持定时全量巡检和配置备份，真正做到无人值守。

\### 2.2 设计原则

\- \*\*多厂商兼容\*\*：系统需支持H3C、华为、思科、锐捷等主流厂商设备，统一抽象管理接口。  
\- \*\*前后端分离\*\*：前端与后端解耦，接口标准化，便于后续扩展和团队协作。  
\- \*\*容器化部署\*\*：基于Docker Compose一键部署，降低部署门槛与运维成本。  
\- \*\*安全性优先\*\*：设备凭据加密存储、API认证鉴权、操作审计日志齐全。  
\- \*\*渐进式交付\*\*：分期实施，一期聚焦核心功能快速上线，二期迭代增强。

  
\## 三、技术架构设计

\### 3.1 技术栈选型

根据“当前最主流技术栈”要求，结合网络运维场景特点，选定如下技术方案：

| 层级 | 技术选型 | 选型理由 |  
|---|---|---|  
| \*\*前端框架\*\* | Vue 3 + TypeScript + Vite | Vue 3 Composition API生态成熟，TypeScript保障代码健壮性，Vite构建速度极快 |  
| \*\*前端UI库\*\* | Element Plus | 企业级中后台UI组件库，与Vue 3深度适配，组件丰富，适合运维平台 |  
| \*\*状态管理\*\* | Pinia | Vue官方推荐的状态管理库，API简洁，TypeScript友好 |  
| \*\*图表可视化\*\* | ECharts 5 | 国产图表库，支持仪表盘、实时曲线等监控场景 |  
| \*\*后端语言\*\* | Go 1.22+ | 高并发、低内存占用，编译为单一二进制文件，部署极为便捷 |  
| \*\*后端框架\*\* | Gin | 高性能HTTP框架，社区活跃，中间件生态丰富 |  
| \*\*ORM\*\* | GORM 2.0 | Go生态最成熟的ORM，支持MySQL、自动迁移、关联查询 |  
| \*\*数据库\*\* | MySQL 8.0 | 结构化数据存储（设备信息、配置历史、告警记录、用户数据） |  
| \*\*缓存\*\* | Redis 7 | 会话缓存、设备在线状态缓存、任务队列、分布式锁 |  
| \*\*SSH自动化\*\* | Netmiko（Python） | 多厂商SSH自动化事实标准，原生支持H3C Comware、Huawei VRP、Cisco IOS等 |  
| \*\*配置备份引擎\*\* | Oxidized + Git | 轻量级网络配置备份工具，支持130+种OS，原生Git版本管理 |  
| \*\*定时任务\*\* | robfig/cron（Go） | Go语言成熟的Cron调度库，支持秒级精度 |  
| \*\*容器化\*\* | Docker + Docker Compose | 一键编排部署，环境一致性保障 |  
| \*\*反向代理\*\* | Nginx | 前端静态资源服务 + API反向代理 |

\> \*\*关于后端语言的选择说明\*\*：后端选用Go语言（而非Python），主要基于以下考量：  
\> - Go编译为单一二进制文件，部署时无需Python运行环境，Docker镜像更轻量（约20MB vs Python的200MB+）  
\> - Go原生高并发（goroutine），巡检调度、WebSocket实时推送等场景性能表现优异  
\> - 自动巡检的底层SSH交互使用Python Netmiko微服务（Docker容器），兼顾开发效率与运行时性能

\### 3.2 系统总体架构

\`\`\`  
┌──────────────────────────────────────────────────────────────────┐  
│                          前端层 (Vue 3 + Element Plus)            │  
│  设备管理 │ 自动巡检 │ 配置备份 │ 状态监控 │ 告警管理 │ 定时任务  │  
└──────────────────────────┬───────────────────────────────────────┘  
                           │ HTTP / WebSocket  
┌──────────────────────────▼───────────────────────────────────────┐  
│                       API网关层 (Nginx)                           │  
│                  反向代理 │ 静态资源 │ 负载均衡                     │  
└──────────────────────────┬───────────────────────────────────────┘  
                           │  
┌──────────────────────────▼───────────────────────────────────────┐  
│                     后端服务层 (Go + Gin)                         │  
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐   │  
│  │ 设备管理  │ │ 巡检调度  │ │ 配置管理  │ │ 告警引擎         │   │  
│  │ Service  │ │ Service  │ │ Service  │ │ Service          │   │  
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘   │  
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐   │  
│  │ 用户认证  │ │ 权限管理  │ │ 任务调度  │ │ WebSocket推送    │   │  
│  │ Service  │ │ (RBAC)   │ │ (Cron)   │ │ Service          │   │  
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘   │  
└──────────────────────────┬───────────────────────────────────────┘  
                           │  
┌──────────────────────────▼───────────────────────────────────────┐  
│                      自动化执行层                                  │  
│  ┌────────────────────┐    ┌────────────────────┐               │  
│  │  Netmiko 微服务     │    │  Oxidized 配置备份  │               │  
│  │  (Python Docker)   │    │  (Ruby Docker)     │               │  
│  │  SSH → 设备CLI     │    │  Git → 版本管理    │               │  
│  └────────────────────┘    └────────────────────┘               │  
└──────────────────────────┬───────────────────────────────────────┘  
                           │ SSH / SNMP / NETCONF  
┌──────────────────────────▼───────────────────────────────────────┐  
│                    被管网络设备层                                  │  
│  H3C S10504 │ H3C AC │ S5560X │ S5130S │ WTU420 │ 华为/思科/锐捷 │  
└──────────────────────────────────────────────────────────────────┘  
\`\`\`

\### 3.3 多厂商SSH适配方案

系统通过Netmiko库实现多厂商设备的统一SSH交互。Netmiko在Paramiko基础上深度封装，原生支持Cisco IOS/NX-OS、Juniper Junos、Arista EOS、Huawei VRP、H3C Comware等数十种主流厂商设备类型，具备自动识别设备型号、智能等待提示符、命令执行超时控制、异常捕获与日志记录等高级特性。

\*\*核心设计思路\*\*：

1\. 在设备管理中为每台设备标记\`device\_type\`（如\`hp\_comware\`、\`huawei\`、\`cisco\_ios\`、\`ruijie\_os\`）  
2\. 自动巡检调度器根据\`device\_type\`调用对应驱动，发送厂商特定命令，解析返回结果  
3\. 配置备份时同样依据\`device\_type\`选择正确的命令（如H3C的\`display current-configuration\`、Cisco的\`show running-config\`）  
4\. 对于支持NETCONF的设备（如H3C S10504、AC），优先使用NETCONF协议进行结构化数据获取与配置管理，Netmiko的SSH方式作为补充和兜底方案

  
\## 四、功能模块设计

\### 4.1 设备管理模块

\*\*功能描述\*\*：提供网络设备的全生命周期管理，包括添加、编辑、删除、分组、搜索、导入导出等操作。

\*\*核心字段设计\*\*：

| 字段 | 说明 | 示例 |  
|---|---|---|  
| 设备名称 | 自定义标识 | \`宿舍楼1-本体交换机\` |  
| 设备类型 | 厂商型号 | \`h3c\_comware\` / \`huawei\` / \`cisco\_ios\` / \`ruijie\` |  
| 管理IP | SSH连接IP | \`192.168.1.1\` |  
| SSH端口 | 默认22 | \`22\` |  
| 所属分组 | 逻辑分组 | \`中心机房\` / \`宿舍楼1\` / \`宿舍楼2\` |  
| 设备角色 | 网络角色 | \`核心\` / \`AC\` / \`本体\` / \`接入\` |  
| 设备型号 | 具体型号 | \`S5560X-30F-EI\` |  
| SNMP团体字 | 可选 | \`public\` |  
| 凭据关联 | 关联SSH凭据ID | — |

\*\*设备类型映射表（Netmiko device\_type）\*\* ：

| 厂商 | Netmiko device\_type | 适用设备 |  
|---|---|---|  
| H3C Comware | \`hp\_comware\` | S10504, S5560X, S5130S, AC |  
| 华为 VRP | \`huawei\` | CE/S系列交换机 |  
| 思科 IOS | \`cisco\_ios\` | Catalyst系列 |  
| 锐捷 | \`ruijie\_os\` | RG系列 |

\*\*凭据管理\*\*：SSH用户名/密码采用AES-256加密存储，支持多套凭据，设备添加时通过ID关联，避免凭据重复录入和明文暴露。

\### 4.2 自动巡检模块

\*\*功能描述\*\*：通过SSH自动登录设备执行巡检命令，采集CPU使用率、内存使用率、接口状态、硬件信息、系统运行时长等关键指标，生成巡检报告。

\*\*巡检项定义（按设备类型区分）\*\* ：

| 巡检项 | H3C Comware命令 | 华为VRP命令 | Cisco IOS命令 |  
|---|---|---|---|  
| CPU使用率 | \`display cpu-usage\` | \`display cpu-usage\` | \`show processes cpu\` |  
| 内存使用率 | \`display memory\` | \`display memory-usage\` | \`show memory statistics\` |  
| 接口状态 | \`display interface brief\` | \`display interface brief\` | \`show ip interface brief\` |  
| 硬件信息 | \`display device\` | \`display device\` | \`show inventory\` |  
| 系统运行时长 | \`display version\` | \`display version\` | \`show version\` |  
| 设备温度 | \`display environment\` | \`display environment\` | \`show environment\` |

\*\*巡检执行流程\*\*：

\`\`\`  
1\. Cron触发 / 手动触发巡检任务  
2\. 查询待巡检设备列表，按设备分组并发执行（goroutine池，最大并发数可配置）  
3\. 对每台设备：  
  a. 从连接池获取SSH连接（或新建连接）  
  b. 依次执行各巡检命令  
  c. 使用TextFSM模板解析命令输出为结构化数据  
  d. 写入MySQL巡检结果表  
  e. 检测指标是否超过告警阈值，触发告警  
4\. 巡检结束后生成巡检报告（正常/异常统计）  
\`\`\`

\*\*并发能力\*\*：Go语言goroutine原生支持高并发，200台设备全量巡检可在5分钟内完成（10并发下），相比人工逐台登录效率提升超100倍。

\### 4.3 配置备份模块

\*\*功能描述\*\*：基于Oxidized引擎实现设备运行配置的一键备份与定时全量备份，通过Git进行版本管理，支持历史配置追溯与一键回滚。

\*\*技术方案\*\*：Oxidized是一款轻量级网络配置备份工具，支持超过130种网络设备操作系统，原生Git集成使每次备份都是一次Git提交，天然具备完整的版本历史、差异比较和回滚能力。在500台设备的实测中，全量配置备份耗时仅约30分钟，远优于传统方案的数小时甚至数天。

\*\*核心功能\*\*：

| 功能 | 说明 |  
|---|---|  
| 一键备份 | 选择单台或多台设备，立即触发配置备份 |  
| 定时备份 | 通过内置Cron调度器，每日/每周定时全量备份 |  
| 版本历史 | 每次备份自动生成Git commit，记录时间、操作人、变更内容 |  
| 配置对比 | 任意两个版本之间的diff对比，高亮显示变更行 |  
| 配置回滚 | 选择历史版本，通过SSH将配置回写至设备 |  
| 备份下载 | 支持下载指定版本的配置文件 |

\*\*回滚安全机制\*\*：  
\- 回滚前自动备份当前运行配置（二次确认保护）  
\- 回滚操作需管理员二次审批  
\- 全量操作审计日志记录（谁、何时、对哪台设备、回滚到哪个版本）

\### 4.4 状态监控模块

\*\*功能描述\*\*：Dashboard大屏实时展示全网设备在线状态、关键资源使用率，支持设备拓扑分组视图。

\*\*监控维度\*\*：

| 维度 | 展示方式 | 数据来源 |  
|---|---|---|  
| 设备在线率 | 饼图 / 统计卡片 | 最近一次巡检结果 |  
| CPU使用率排行 | 柱状图 Top N | 巡检数据 |  
| 内存使用率排行 | 柱状图 Top N | 巡检数据 |  
| 接口状态统计 | 统计卡片（Up/Down数量） | 巡检数据 |  
| 设备在线状态 | 表格（绿色在线/红色离线） | 实时心跳检测 |  
| 历史趋势 | 折线图（7天/30天CPU/内存趋势） | 历史巡检数据 |

\*\*实时推送\*\*：后端通过WebSocket向前端推送设备在线状态变更、告警触发等实时消息，无需手动刷新页面。

\### 4.5 告警管理模块

\*\*功能描述\*\*：统一汇聚设备告警信息，支持告警规则配置、分级通知（钉钉/邮件/Webhook）、告警确认与处理。

\*\*告警类型与阈值（可自定义）\*\* ：

| 告警类型 | 默认阈值 | 严重级别 | 通知方式 |  
|---|---|---|---|  
| 设备离线 | 连续3次Ping不通 | 🔴 严重 | 钉钉 + 邮件 |  
| CPU使用率过高 | > 90% | 🟠 警告 | 钉钉 |  
| 内存使用率过高 | > 90% | 🟠 警告 | 钉钉 |  
| 接口Down | 任意核心端口Down | 🔴 严重 | 钉钉 + 邮件 |  
| 设备温度过高 | > 65℃ | 🟠 警告 | 钉钉 |  
| 配置变更 | 配置内容Hash变化 | 🟡 提示 | 系统内通知 |

\*\*告警生命周期\*\*：  
\`\`\`  
触发 → 通知发送 → 待确认 → \[已确认 / 自动恢复\] → \[处理中\] → 已关闭  
\`\`\`

\### 4.6 定时任务模块

\*\*功能描述\*\*：内置Cron调度器，支持创建、编辑、启用/禁用定时任务，涵盖定时巡检、定时配置备份、自定义脚本执行等场景。

\*\*预置任务模板\*\*：

| 任务名称 | 默认Cron | 执行内容 |  
|---|---|---|  
| 每日设备巡检 | \`0 6 \* \* \*\`（每日06:00） | 全量设备执行巡检 |  
| 每周全量配置备份 | \`0 2 \* \* 0\`（每周日02:00） | 全量设备配置备份 |  
| AP状态定时采集 | \`\*/30 \* \* \* \*\`（每30分钟） | 通过AC采集AP在线状态 |

\*\*任务执行机制\*\*：  
\- 基于robfig/cron库实现秒级精度的定时调度  
\- 使用Redis分布式锁确保单实例执行（避免多副本重复执行）  
\- 任务执行日志完整记录（开始时间、结束时间、执行结果、错误信息）  
\- 失败任务支持自动重试（可配置重试次数和间隔）

  
\## 五、数据库设计（核心表）

\### 5.1 设备表（devices）

| 字段 | 类型 | 说明 |  
|---|---|---|  
| id | BIGINT AUTO\_INCREMENT | 主键 |  
| name | VARCHAR(128) | 设备名称 |  
| device\_type | VARCHAR(64) | Netmiko类型标识 |  
| vendor | VARCHAR(32) | 厂商：h3c/huawei/cisco/ruijie |  
| model | VARCHAR(64) | 设备型号 |  
| role | VARCHAR(32) | 设备角色：core/ac/aggregation/access/ap |  
| management\_ip | VARCHAR(45) | 管理IP |  
| ssh\_port | INT | SSH端口，默认22 |  
| credential\_id | BIGINT | 关联凭据表外键 |  
| group\_id | BIGINT | 关联设备分组表外键 |  
| building | VARCHAR(64) | 所属楼栋 |  
| floor | INT | 所属楼层 |  
| status | TINYINT | 在线状态：0-离线 1-在线 |  
| last\_seen | DATETIME | 最后在线时间 |  
| created\_at | DATETIME | 创建时间 |  
| updated\_at | DATETIME | 更新时间 |

\### 5.2 巡检结果表（inspection\_results）

| 字段 | 类型 | 说明 |  
|---|---|---|  
| id | BIGINT AUTO\_INCREMENT | 主键 |  
| device\_id | BIGINT | 关联设备ID |  
| cpu\_usage | DECIMAL(5,2) | CPU使用率(%) |  
| memory\_usage | DECIMAL(5,2) | 内存使用率(%) |  
| temperature | DECIMAL(5,2) | 设备温度(℃) |  
| uptime | VARCHAR(128) | 运行时长 |  
| interface\_status | JSON | 接口状态（结构化数据） |  
| raw\_output | TEXT | 原始命令输出 |  
| is\_anomaly | TINYINT | 是否异常 |  
| inspected\_at | DATETIME | 巡检时间 |

\### 5.3 配置备份表（config\_backups）

| 字段 | 类型 | 说明 |  
|---|---|---|  
| id | BIGINT AUTO\_INCREMENT | 主键 |  
| device\_id | BIGINT | 关联设备ID |  
| config\_hash | VARCHAR(64) | 配置内容SHA256 |  
| git\_commit\_id | VARCHAR(40) | Git commit hash |  
| file\_path | VARCHAR(512) | 配置文件存储路径 |  
| triggered\_by | VARCHAR(64) | 触发方式：manual/scheduled |  
| operator | VARCHAR(64) | 操作人 |  
| created\_at | DATETIME | 备份时间 |

\### 5.4 告警表（alerts）

| 字段 | 类型 | 说明 |  
|---|---|---|  
| id | BIGINT AUTO\_INCREMENT | 主键 |  
| device\_id | BIGINT | 关联设备ID |  
| alert\_type | VARCHAR(64) | 告警类型 |  
| severity | VARCHAR(16) | 严重级别：critical/warning/info |  
| message | VARCHAR(512) | 告警消息 |  
| status | VARCHAR(16) | 状态：triggered/confirmed/resolved |  
| triggered\_at | DATETIME | 触发时间 |  
| resolved\_at | DATETIME | 恢复时间 |

  
\## 六、Docker Compose 部署方案

\### 6.1 服务清单

| 服务 | 镜像 | 说明 |  
|---|---|---|  
| nginx | nginx:1.25-alpine | 反向代理 + 前端静态资源 |  
| backend | 自构建（Go） | 后端API服务 |  
| netmiko-worker | 自构建（Python） | Netmiko SSH巡检微服务 |  
| oxidized | oxidized/oxidized:latest | 配置备份引擎 |  
| mysql | mysql:8.0 | 关系型数据库 |  
| redis | redis:7-alpine | 缓存与任务队列 |

\### 6.2 docker-compose.yml 参考

\`\`\`yaml  
version: "3.8"

services:  
 nginx:  
   image: nginx:1.25-alpine  
   container\_name: netops-nginx  
   ports:  
     - "80:80"  
     - "443:443"  
   volumes:  
     - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro  
     - ./frontend/dist:/usr/share/nginx/html:ro  
   depends\_on:  
     - backend  
   restart: always

 backend:  
   build: ./backend  
   container\_name: netops-backend  
   ports:  
     - "8080:8080"  
   environment:  
     - DB\_HOST=mysql  
     - DB\_PORT=3306  
     - DB\_USER=netops  
     - DB\_PASSWORD=${DB\_PASSWORD}  
     - DB\_NAME=netops  
     - REDIS\_ADDR=redis:6379  
     - NETMIKO\_WORKER\_URL=http://netmiko-worker:5000  
     - OXIDIZED\_API\_URL=http://oxidized:8888  
     - JWT\_SECRET=${JWT\_SECRET}  
   depends\_on:  
     mysql:  
       condition: service\_healthy  
     redis:  
       condition: service\_started  
   restart: always

 netmiko-worker:  
   build: ./netmiko-worker  
   container\_name: netops-netmiko  
   ports:  
     - "5000:5000"  
   environment:  
     - REDIS\_ADDR=redis:6379  
   depends\_on:  
     - redis  
   restart: always

 oxidized:  
   image: oxidized/oxidized:latest  
   container\_name: netops-oxidized  
   ports:  
     - "8888:8888"  
   volumes:  
     - ./oxidized/config:/root/.config/oxidized  
     - ./oxidized/configs:/root/.config/oxidized/configs  
   restart: always

 mysql:  
   image: mysql:8.0  
   container\_name: netops-mysql  
   ports:  
     - "3306:3306"  
   environment:  
     - MYSQL\_ROOT\_PASSWORD=${MYSQL\_ROOT\_PASSWORD}  
     - MYSQL\_DATABASE=netops  
     - MYSQL\_USER=netops  
     - MYSQL\_PASSWORD=${DB\_PASSWORD}  
   volumes:  
     - mysql\_data:/var/lib/mysql  
   healthcheck:  
     test: \["CMD", "mysqladmin", "ping", "-h", "localhost"\]  
     interval: 10s  
     timeout: 5s  
     retries: 3  
   restart: always

 redis:  
   image: redis:7-alpine  
   container\_name: netops-redis  
   ports:  
     - "6379:6379"  
   volumes:  
     - redis\_data:/data  
   restart: always

volumes:  
 mysql\_data:  
 redis\_data:  
\`\`\`

\### 6.3 一键部署步骤

\`\`\`bash  
\# 1. 克隆项目  
git clone <repository-url> && cd netops

\# 2. 配置环境变量  
cp .env.example .env  
\# 编辑 .env，设置数据库密码、JWT密钥等

\# 3. 构建并启动所有服务  
docker compose up -d

\# 4. 初始化数据库  
docker compose exec backend ./netops migrate

\# 5. 验证服务状态  
docker compose ps  
\`\`\`

部署完成后，访问 \`http://<服务器IP>\` 即可进入运维管理平台。

  
\## 七、AP管理专项方案

鉴于当前1,800台AP的唯一管理入口是AC Web页面，这是用户最核心的痛点，本节专门设计AP自动化管理方案。

\### 7.1 当前AP管理限制分析

H3C WTU420 AP作为瘦AP（Fit AP），自身不提供独立的SSH/Telnet管理接口，所有管理操作必须通过无线控制器AC（H3C EWPXM1MAC0F）执行。这决定了\*\*AP管理自动化的核心思路是通过AC来实现间接管理\*\*。

\### 7.2 方案设计：通过AC实现AP自动化管理

\*\*AC支持的自动化接口\*\*：H3C Comware设备支持RESTful API、NETCONF over SSH等方式进行可编程管理。

\*\*实现路径\*\*：

| 管理操作 | 实现方式 | 说明 |  
|---|---|---|  
| AP在线状态采集 | SSH登录AC，执行\`display wlan ap all\`，解析输出 | 通过Netmiko自动执行，结构化解析AP列表及状态 |  
| AP详细诊断 | SSH登录AC，执行\`display wlan ap name xxx verbose\` | 获取单台AP的详细信息（客户端数、信道、信号强度等） |  
| AP重启 | SSH登录AC，执行\`reset wlan ap name xxx\` | 远程重启指定AP |  
| AP配置下发 | SSH登录AC，执行配置命令 | 批量修改AP模板、信道、功率等 |

\*\*AP数据采集频率设计\*\*：

\- 基础状态（在线/离线）：每5分钟采集一次，通过Redis缓存  
\- 详细诊断数据：每日巡检时全量采集  
\- 离线告警：连续2次采集均离线则触发告警

\### 7.3 AP在系统中的展示

在运维平台的设备管理模块中，AP作为特殊设备类型存在：  
\- 设备类型标记为\`ap\`  
\- 管理IP字段填写AC的管理IP  
\- 增加\`ap\_name\`字段存储AP在AC上的名称  
\- 自动巡检时，系统识别设备类型为\`ap\`，自动走AC查询路径而非直连AP

  
\## 八、实施计划

\### 8.1 分期实施规划

\*\*一期：基础平台搭建（6周）\*\*

| 阶段 | 工作内容 | 产出物 |  
|---|---|---|  
| 第1-2周 | 环境搭建、数据库设计、项目初始化 | 可运行的骨架项目 |  
| 第3-4周 | 设备管理、凭据管理、用户认证模块开发 | 设备增删改查功能可用 |  
| 第5周 | Netmiko微服务开发、设备SSH连通性测试 | 自动巡检能力就绪 |  
| 第6周 | 前端Dashboard、设备列表页面开发 | 一期功能联调上线 |

\*\*二期：核心自动化能力（4周）\*\*

| 阶段 | 工作内容 | 产出物 |  
|---|---|---|  
| 第7-8周 | 自动巡检模块完善（解析引擎、巡检报告） | 全量设备自动巡检可用 |  
| 第9周 | 配置备份（Oxidized集成、Git版本管理） | 一键/定时配置备份可用 |  
| 第10周 | 告警引擎（规则配置、钉钉/邮件通知） | 告警统一汇聚与推送可用 |

\*\*三期：增强与优化（3周）\*\*

| 阶段 | 工作内容 | 产出物 |  
|---|---|---|  
| 第11周 | AP管理专项（AC SSH自动化、AP状态面板） | 1,800台AP纳入统一管理 |  
| 第12周 | 定时任务调度、配置回滚、权限细化 | 全功能就绪 |  
| 第13周 | 性能测试、安全加固、文档完善、培训 | 正式上线 |

\*\*总工期：约 13 周（3个月）\*\*

\### 8.2 里程碑节点

\`\`\`  
M1（第2周末）：项目骨架就绪，前后端联调通过  
M2（第6周末）：一期功能上线，设备统一纳管可用  
M3（第10周末）：二期功能上线，自动化巡检+备份+告警可用  
M4（第13周末）：三期全功能上线，系统正式交付  
\`\`\`

  
\## 九、可行性分析

\### 9.1 技术可行性 ✅ 可行

\- \*\*成熟的技术栈\*\*：Go + Gin + Vue 3 + Element Plus 均为业界主流框架，社区活跃，文档丰富，学习成本低。  
\- \*\*可靠的开源组件\*\*：Netmiko已被数千家企业用于网络自动化，支持H3C、华为、思科等主流厂商，稳定性经过充分验证。Oxidized同样广泛应用于网络配置备份场景，支持130+种OS。  
\- \*\*多厂商兼容性\*\*：H3C Comware设备原生支持SSH、NETCONF、RESTful API，自动化接口完备。  
\- \*\*容器化部署\*\*：基于Docker Compose编排，6个服务可一键部署，无复杂依赖。  
\- \*\*已有成功案例\*\*：国网白银供电公司基于类似方案建设的运维智能体，将每日巡检耗时从2小时压缩至30分钟，全量配置备份从48小时缩短至30分钟。

\### 9.2 经济可行性 ✅ 可行

| 成本项 | 估算 | 说明 |  
|---|---|---|  
| 人力成本 | 1-2名全栈工程师 × 13周 | 核心开发投入 |  
| 服务器资源 | 1台虚拟机（4C8G 200G SSD） | 现有资源可复用，无需采购新硬件 |  
| 软件许可 | ¥0 | 全部采用开源组件，无商业许可费用 |  
| 外部依赖 | ¥0 | 无第三方付费API或服务 |

\> \*\*投资回报分析\*\*：系统上线后，每日巡检耗时预计从2小时降至5分钟（效率提升96%），每月节省约50工时；全量配置备份从48小时降至30分钟（效率提升99%）。按运维人员平均人力成本计算，\*\*预计3-4个月即可收回全部开发投入\*\*。

\### 9.3 人员可行性 ✅ 可行

\- \*\*开发人员\*\*：需要1-2名具备Go/Vue开发经验的全栈工程师，该技能组合在行业内较为普遍，招聘或内部调配均可。  
\- \*\*运维人员\*\*：系统上线后，运维人员的使用门槛低——通过Web界面即可完成日常巡检、备份、告警查看，无需掌握编程技能。  
\- \*\*培训成本\*\*：运维人员半天即可熟练掌握系统操作；开发交接文档完备，新人亦可快速接手维护。  
\- \*\*降低人员依赖\*\*：系统将运维知识和操作流程固化为自动化任务，有效消除“知识孤岛”风险，降低人员流动带来的运维断层。

\### 9.4 风险及应对措施

| 风险类别 | 风险描述 | 风险等级 | 应对措施 |  
|---|---|---|---|  
| \*\*网络安全\*\* | SSH凭据泄露可能导致设备被非法访问 | 🔴 高 | 凭据AES-256加密存储；API JWT认证；操作全量审计日志；SSH密钥+密码双因素认证 |  
| \*\*误操作\*\* | 批量配置回滚误操作导致大面积断网 | 🔴 高 | 回滚前强制备份当前配置；二次审批机制；先单台测试再批量执行；变更窗口限制 |  
| \*\*设备兼容\*\* | 部分老旧设备SSH版本不兼容 | 🟡 中 | 预留Telnet回退通道；设备适配表逐步完善；上线前逐型号验证 |  
| \*\*性能瓶颈\*\* | 大规模并发SSH导致设备CPU过高 | 🟡 中 | 并发数可配置（建议≤10）；巡检错峰执行；设备侧设置SSH限速 |  
| \*\*系统单点\*\* | 运维平台自身宕机 | 🟡 中 | Docker自动重启策略；数据库定期备份；故障时仍可手动SSH登录设备（不依赖平台） |  
| \*\*厂商API变更\*\* | H3C固件升级导致CLI输出格式变化 | 🟢 低 | TextFSM模板化解析，变更时仅需更新模板；建立设备版本兼容性矩阵 |

\> \*\*自动化不是替代运维人员，而是将人从重复劳动中解放出来，将精力投入到更高价值的架构设计、安全优化等工作中。关键操作仍需“人工兜底”。\*\*

  
\## 十、附录

\### 附录A：TextFSM模板示例（H3C display cpu-usage 解析）

\`\`\`textfsm  
Value CPU\_USAGE (\\d+)

Start  
 ^.\*\\s+${CPU\_USAGE}%\\s+in\\s+last\\s+5\\s+seconds.\* -> Record  
\`\`\`

\### 附录B：Netmiko微服务API接口示例

\`\`\`  
POST /api/v1/execute  
{  
   "device\_type": "hp\_comware",  
   "ip": "192.168.1.1",  
   "username": "admin",  
   "password": "\*\*\*",  
   "commands": \["display cpu-usage", "display memory"\]  
}

Response:  
{  
   "success": true,  
   "results": \[  
       {"command": "display cpu-usage", "output": "..."},  
       {"command": "display memory", "output": "..."}  
   \]  
}  
\`\`\`

\### 附录C：RESTful API接口清单（后端核心接口）

| 接口路径 | 方法 | 说明 |  
|---|---|---|  
| \`/api/v1/auth/login\` | POST | 用户登录，返回JWT |  
| \`/api/v1/devices\` | GET/POST | 设备列表查询 / 添加设备 |  
| \`/api/v1/devices/:id\` | GET/PUT/DELETE | 设备详情 / 编辑 / 删除 |  
| \`/api/v1/devices/:id/inspect\` | POST | 手动触发单设备巡检 |  
| \`/api/v1/inspections/batch\` | POST | 批量巡检 |  
| \`/api/v1/inspections/report\` | GET | 获取巡检报告 |  
| \`/api/v1/configs/backup\` | POST | 手动触发配置备份 |  
| \`/api/v1/configs/history/:deviceId\` | GET | 配置备份历史 |  
| \`/api/v1/configs/rollback\` | POST | 配置回滚 |  
| \`/api/v1/alerts\` | GET | 告警列表 |  
| \`/api/v1/alerts/:id/resolve\` | PUT | 确认/恢复告警 |  
| \`/api/v1/tasks\` | GET/POST | 定时任务管理 |  
| \`/api/v1/monitor/dashboard\` | GET | Dashboard汇总数据 |

\---

以上为网络运维自动化管理系统的完整实施文档，涵盖了项目背景、架构设计、功能规划、数据库设计、部署方案、实施计划和可行性分析。系统建成后，将彻底改变当前“逐台登录Web页面查看AP状态”的被动运维模式，实现近2,000台网络设备的统一纳管、自动化巡检、配置集中管理和告警智能推送，显著提升运维效率、降低运维风险和人力成本。