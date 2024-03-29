1. pink https://github.com/busgo/pink
  - 通过Web界面管理操作简单方便，支持各种任务
  - 高可用可以部署 n 台调度集群节点，保证没有单点故障。
  - 部署简单、仅仅需要一个执行文件即可运行。
  - 集成方便，统一语言无关的任务抽象，接入不同语言的sdk。
  - 任务故障转移，任务客户端下线自动切至在线客户端任务机器。
  
2. crocodile https://github.com/labulaka521/crocodile
  - 基于Golang开发的分布式任务调度系统，支持http、golang、python、shell、python3、nodejs、bat等调度任务
  - 依赖redis和mysql
  - 个人觉得界面可以借鉴
  
3. gopherCron https://github.com/holdno/gopherCron
  - 依赖 Etcd 服务注册与发现
  - Gin webapi 提供可视化操作
  - Mysql # 任务日志存储
  - cronexpr # github.com/gorhill/cronexpr cron表达式解析器
  - 秒级定时任务
  - 任务日志查看
  - 随时结束任务进程
  - 分布式扩展
  - 健康节点检测 (分项目显示对应的健康节点IP及节点数)
  
4. gocron https://github.com/ouqiang/gocron
  - Web界面管理定时任务
  - crontab时间表达式, 精确到秒
  - 任务执行失败可重试
  - 任务执行超时, 强制结束
  - 任务依赖配置, A任务完成后再执行B任务
  - 账户权限控制
  - 任务类型
  - shell任务 在任务节点上执行shell命令, 支持任务同时在多个节点上运行
  - HTTP任务 访问指定的URL地址, 由调度器直接执行, 不依赖任务节点
  - 查看任务执行结果日志
  - 任务执行结果通知, 支持邮件、Slack、Webhook
  
5. go-cron https://gitee.com/man0sions/go-cron
  - golang分布式定时任务调度器，支持秒级调度，master节点下发指令，worker节点处理任务
  - 依赖etcd，mongodb

6. jiacrontab https://gitee.com/iwannay/jiacrontab
  - 自定义job执行
  - 允许设置job的最大并发数
  - 每个脚本都可在web界面下灵活配置，如测试脚本运行，查看日志，强杀进程，停止定时…
  - 允许添加脚本依赖（支持跨服务器），依赖脚本提供同步和异步的执行模式
  - 支持异常通知
  - 支持守护脚本进程
  - 支持节点分组
  - iacrontab 由 jiacrontab_admin，jiacrontabd 两部分构成，两者完全独立通过 rpc 通信
  - jiacrontab_admin：管理后台向用户提供web操作界面
  - jiacrontabd：负责job数据存储，任务调度
  
7. clock https://github.com/BruceDone/clock
  - dag图形界面配置
  
8. PPGo_Job https://gitee.com/georgehao/PPGo_Job
  PPGo_Job是一款定时任务可视化的、多人多权限的管理系统，采用golang开发，安装方便，资源消耗少，支持大并发，可同时管理多台服务器上的定时任务。支持PHP,Python,Shell,Java,Go等常见编程语言的定时任务管理，也支持各类unix服务器的各种命令等


9. gojob https://github.com/wj596/gojob
  - 极少依赖： 只依赖MySQL 数据库，分布式环境下使用内建的分布式协调机制，不需要安装第三方分布式协调服务，如Zookeeper、Etcd等；更少的依赖意味着后续需要更少的部署和运维成本。

  - 易部署：原生Native程序，无需安装运行时环境，如JDK、.net framework等；支持单机和集群部署两种部署模式。

  - 任务重试：支持自定义任务重试次数、重试时间间隔。当任务执行失败时，会按照固定的间隔时间进行重试。

  - 任务超时：支持自定义任务超时时间，当任务超时，会强制结束执行。

  - 失败转移：当任务在一个执行节点上执行失败，会转移到下一个可用执行节点上继续执行。如任务在节点A上执行失败，会转移到节点B上继续执行，如果失败会转移到节点C上继续执行，直到执行成功。

  - misfire补偿机制：由于调度服务器宕机、资源耗尽等原因致使任务错过激活时间，称之为哑火(misfire)。比如每天23点整生成日结报表，但是恰巧在23点前服务器宕机、此任务就错失了一次调度。如果我们设置了misfireThreshold为30分钟，如果服务器在23点30分之前恢复，调度器会进行一次执行，以补偿在23点整哑火的调度。

  - 负载均衡：如果集群节点为集群部署，调度服务器可以使用轮询、随机、加权轮询、加权随机等路由策略，为任务选择合适的执行节点。既可以保证执行节点高可用、我单点隐患，也可以将压力分散到不同的执行节点。

  - 任务分片：将大任务拆解为多个小任务均匀的散落在多个节点上并行执行，以协作的方式完成任务。比如订单核对业务，我们有天津、上海、重庆、河北、山西、辽宁、吉林、江苏、浙江、安徽十个省市的账单，如果数据量比较大，单机处理这些订单的核对业务显然不现实。

  - gojob可以将任务分为3片：执行节点1负责–>天津、上海、重庆、河北；执行节点2负责–>山西、辽宁、吉林; 执行节点13负责–>江苏、浙江、安徽。这样可以用3台机器来合力完成这个任务。如果你的机器足够，可以将任务分成更多片，用更多的机器来协同处理。

  - 弹性扩缩容：调度器会感知执行节点的增加和删除、上线和下线，并将执行节点的变化情况应用到下一次的负载均衡算法和任务分片算法中。支持动态的执行节点动态横向扩展，弹性伸缩整个系统的处理能力。

  - 调度唯一性：调度节点集群使用Raft算法进行主节点选举，一个集群中只存在一个主节点。任务在一个执行周期内，只会被主节点调用一次，保证调度的一致性。

  - 调度节点高可用：集群内通过Raft共识算法和数据快照将作业元数据实时进行同步，调度节点收到同步的数据后存在自己内建BoltDB存储引擎中；作业元数据具有强一致性和多副本存储的特性；任务可在任意调度节点被调度，调度节点之间可以无缝衔接，任何一个节点宕机另一个节点可以在毫秒计的时间内接替，保证调度节点无单点隐患。

  - 数据库节点高可用：由于作业元数据保存在节点自己的存储引擎中，MySQL数据库只用来保存调度日志。日志数据的特性使其可容忍短时间内不一致甚至丢失(虽然极少发生但理论上可容忍)，因此将日志数据异步写入多库，无需对数据库做集群或者同步设置。极端情况下，数据库节点全部宕机都不会影响调度业务的正常运行，保证数据库节点无单点隐患。

  - 任务依赖：任务可以设置多个子任务，触发时机。如：任务执行结束触发子任务、任务执行成功触发子任务、任务执行失败触发子任。

  - 告警：支持邮件告警。任务调度失败会发送告警邮件到指定的邮箱，每个任务可配置多个告警邮箱。调度节点出现故障、数据库节点出现故障也会发送告警邮箱。

  - 数字签名：支持HMAC( 哈希消息认证码 )数字签名，调度节点和执行节点之间可以通过数字签名来确认身份。

10. https://github.com/v-mars/jobor
  - 只依赖mysql，真正的的分布式

  - 支持server/controller/master(通过raft一致性算法)的高可用，一个Raft集群通常包含2*N+1个服务器，允许系统有N个故障服务器。
    ldap(支持openldap,AD 认证)
    server <-- gRPC --> worker
    task abort
    task timeout
    api/restful [GET, POST, PUT, DELETE] task
    shell task
    python3 task server task
    
11.https://github.com/yohamta/dagu


12.https://github.com/fieldryand/goflow


13.https://github.com/tsundata/flowline

