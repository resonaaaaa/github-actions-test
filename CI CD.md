# CI/CD
CI/CD（**持续集成 / 持续交付与部署**）是现代软件工程中实现自动化交付的核心实践。简单来说，它是一套自动化的工作流，将代码从开发者的电脑上，安全、高效地一路护送到线上生产环境。

## CI/CD 的三个概念
这三个名称分别代表了自动化交付的不同深度：

1. **持续集成 (CI - Continuous Integration)**
* **是什么**：多个开发者将代码频繁地（每天多次）合并到共享的主干仓库中。每一次代码提交都会自动触发代码编译（Build）和自动化单元测试。核心是**频繁集成代码，并通过自动化检查尽早发现问题**。
* **解决什么问题**：消除“在我的机器上明明能跑！”的集成地狱现象。它通过机器检查确保新加入的代码不会破坏原有功能。


2. **持续交付 (CD - Continuous Delivery)**
* **是什么**：在 CI 测试通过的基础上，自动将代码打包并部署到类生产环境（如测试环境、预发布环境）。此时代码已经具备了随时可以发布到生产环境的条件，但**最终的上线操作通常需要人工点击确认**。
* **解决什么问题**：消除繁琐的手动打包和环境配置过程，让业务方或测试人员可以随时拿最新版本进行验收。


3. **持续部署 (CD - Continuous Deployment)**
* **是什么**：持续交付的终极进阶版。在代码通过所有的自动化测试后，**无需任何人工干预**，直接自动部署到生产环境面对真实用户。
* **解决什么问题**：极致缩短上线时间（Time-to-Market），实现真正的“改完代码几分钟后，用户就能用上”。



---

## 典型的 CI/CD 流水线 (Pipeline)

在实际工程中，CI/CD 是通过流水线工具（如 GitHub Actions、GitLab CI、Jenkins）串联起来的步骤集合。

| 阶段 | 执行动作 | 常用工具举例 |
| --- | --- | --- |
| **1. 提交 (Commit)** | 开发者提交代码到代码仓库，触发流水线运行。 | Git, GitHub, GitLab |
| **2. 构建 (Build)** | 将代码编译成可执行文件，下载依赖，或打包成 Docker 镜像。 | Maven, Webpack, Docker |
| **3. 测试 (Test)** | 运行单元测试、集成测试、代码质量及安全漏洞扫描。 | JUnit, Jest, SonarQube |
| **4. 预发布 (Staging)** | 部署到与线上环境几乎一样的测试环境，进行验收测试。 | Kubernetes, Ansible |
| **5. 部署 (Deploy)** | 自动（或手动确认后）将新版本发布到生产集群，对外提供服务。 | Terraform, ArgoCD |

---

## 在工程中的重要意义

在没有 CI/CD 之前，软件发布往往是按“月”甚至“半年”为周期的，发布前夜工程师们经常需要通宵手动解决冲突和配置问题。CI/CD 的引入彻底改变了研发效能：

* **快速反馈与低成本修复**：代码一旦提交，几分钟内就能知道是否跑通测试。Bug 在开发初期就被拦截，而不是拖到上线前夕才爆发，极大地降低了修复成本。
* **消灭重复劳动与人为失误**：手动打包、上传服务器、修改配置文件极其容易出错。机器自动执行流水线，保证了无论发布多少次，每次部署的步骤都 100% 严格一致。
* **小步快跑与降低风险**：因为发布变得像点击按钮一样容易，团队可以一次只上线两三个小功能，而不是积攒几百个功能一次性发布。单次变更越小，出现故障时排查问题和回滚版本就越容易。
* **业务敏捷性**：一个需求从写完代码到用户可用，时间从几周缩短到几小时甚至几分钟，让产品能够快速试错，灵活响应市场竞争。

CI/CD 本质上不是某一个具体的软件工具，而是一种通过**自动化**建立系统信任的工程文化。它让开发团队不再害怕“发布日”，从而把精力真正聚焦在创造业务价值上。

# Github Actions
**GitHub Actions** 是 GitHub 官方提供的一款强大的**自动化工作流与 CI/CD 平台**。

[GitHub Actions 官方文档](https://github.c-cpp.com/zh/actions/quickstart)

如果说 GitHub 是存放代码的“仓库”，那么 GitHub Actions 就是仓库自带的“智能机器人管家”。当你对代码仓库做出某些动作（比如提交代码、新建 PR，甚至只是有人提了一个 Issue）时，它就能自动帮你执行一系列预先设定好的任务（如测试代码、打包部署、自动回复等）。

理解 GitHub Actions，最核心的是掌握它的“套娃式”结构。从大到小，它由以下几个核心概念组成：

## 核心结构

1. **Event (触发事件) —— “导火索”**
决定工作流什么时候开始。比如：
* `push`：有人把代码推送到主分支。
* `pull_request`：有人发起了合并请求。
* `schedule`：定时任务（比如每天半夜 12 点执行一次）。


2. **Workflow (工作流) —— “剧本”**
这是一份写在 `.github/workflows/` 目录下的 YAML 配置文件。它定义了被事件触发后，整个自动化流水线要经历哪些阶段。一个仓库可以有多个 Workflow（比如一个专门负责测试，一个专门负责部署部署），比如下面的结构：
```plaintext
.github/
└── workflows/
    ├── ci.yml
    ├── docker.yml
    └── release.yml
```

3. **Job (作业) —— “工作小组”**
一个 Workflow 包含一个或多个 Job。默认情况下，多个 Job 是**并行运行**的，但你也可以设置它们串行（比如“部署” Job 必须等“测试” Job 成功后才执行）。
4. **Runner (运行器)**：每个 Job 都需要一台服务器来运行，GitHub 免费提供了云端服务器作为 Runner,也可以自己搭建私有 Runner。Runner的生命周期通常**与 Job 一致**，Job 执行完毕后，Runner 会被销毁。


4. **Step (步骤) —— “具体任务”**
每个 Job 由一系列 Step 组成。Step 必须按顺序**串行**执行。在一个 Step 中，你可以运行命令行脚本（比如 `go test ./...`），也可以直接调用现成的 Action，两者分别对应`run` 和 `uses` 两种操作。
5. **Action (动作) —— “工具包”**
这是 GitHub Actions 最精髓的地方。Action 是**可复用的代码块**。比如“拉取最新代码”、“安装 Node.js 环境”这些极度常用的操作，你不需要自己写脚本，直接去 **GitHub Marketplace（插件市场）** 里拿别人写好的 Action 像搭积木一样拼进来就行。

---
一个简单的例子：

假设我们有一个前端项目，我们希望：**每次有人提交代码，机器就自动帮我们安装依赖并跑一遍测试。**

你只需要在代码目录下建一个 `.github/workflows/test.yml` 文件，写下这段代码：

```yaml
name: 自动测试流水线               # Workflow 的名字

on: 
    push:
        branches: [ main ]            # Event：当有人 push 到 main 分支时触发
    pull_request:
        branches: [ main ]            # Event：当有人发起 PR 到 main 分支时触发

jobs:
  test-code:                     # Job 1：起名叫 test-code
    runs-on: ubuntu-latest       # Runner：向 GitHub 申请一台最新的 Ubuntu 服务器来跑
    
    steps:                       # 开始按顺序执行具体步骤
      - name: 拉取代码库
        uses: actions/checkout@v4  # Action：直接使用官方写好的拉取代码插件
        
      - name: 安装 Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'       # 告诉插件我们要用 Node 20 版本
          
      - name: 安装依赖包
        run: npm install           # 运行自己写的普通命令行
        
      - name: 运行单元测试
        run: npm run test          # 运行测试命令

```

只要把这个文件推送到 GitHub，以后每次提交代码，GitHub 的页面上就会显示流水线的运行状态。你可以点开每个 Job，看到每个 Step 的执行日志，方便排查问题。

# 设计第一个 GO CI

前面的例子让我们知道了 Workflow、Job 和 Step 分别是什么，但把 Node.js 示例原样搬进订单系统并不能完成检查。我们的项目使用 Go，还包含真正连接 MySQL 的订单集成测试。**第一个 CI 要回答的问题是：代码离开本机后，能否在一台干净的机器上完成检查？**

## 从只检查 Go 代码开始

先在项目根目录创建 Workflow 文件：

```plaintext
.github/
└── workflows/
    └── ci.yml
```

最初可以只在推送到 `main` 或向 `main` 提交 Pull Request 时，检出代码、安装 Go、执行 `go vet ./...`，再运行不需要数据库的测试。这样能先确认触发条件与 Go 环境是否正确。不过只检查不依赖 MySQL 的包，**并不能证明订单创建、库存扣减和数据库迁移可用**。

还有一个容易忽略的细节：我们需要让 CI 执行 `go test -count=1 ./...`，把应用层的集成测试也纳入进来。`-count=1` 强制让这次测试实际执行1次，而不是复用 Go 的测试结果缓存。

## 给测试提供 MySQL 服务容器

GitHub 托管的 Runner 不会带着我们的本机 MySQL 和 `.env` 文件启动。如何让订单集成测试获得数据库？在 `test` Job 下配置 **Service container（服务容器）**。GitHub Actions 会为这个 Job 启动临时 MySQL 容器，Job 结束后再清理它。[服务容器说明](https://github.c-cpp.com/zh/actions/using-containerized-services/about-service-containers)

下面摘出当前 `ci.yml` 的关键配置。这里的账号和密码只用于一次性测试数据库；真实环境的凭据不应写进 Workflow。

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    env:
      MYSQL_GO_DSN: "app:app@tcp(localhost:3306)/go_backend_test?parseTime=true"
      MIGRATE_TEST_DATABASE_URL: "mysql://app:app@tcp(localhost:3306)/go_backend_test?parseTime=true"

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: go_backend_test
          MYSQL_USER: app
          MYSQL_PASSWORD: app
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping -h 127.0.0.1 -uroot -proot --silent"
          --health-interval=5s
          --health-timeout=5s
          --health-retries=10
```

这里 `runs-on: ubuntu-latest` 让 Job 直接运行在 Ubuntu Runner 上；MySQL 才运行在旁边的容器里。因此 Go 测试连接的是 Runner 的 `localhost:3306`，`ports: 3306:3306` 将容器的 3306 端口映射到 Runner 的 3306 端口。

它与前面 Docker Compose 中 **API 容器通过服务名 `mysql:3306` 连接数据库** 的场景不同。如果整个 Job 也运行在容器里，连接方式才会改为使用服务名。[服务容器网络说明](https://docs.github.com/en/actions/tutorials/use-containerized-services/use-docker-service-containers)

`MYSQL_DATABASE` 只负责创建空的 `go_backend_test` 数据库，并不会自动创建 `users`、`products`、`orders` 等表。`--health-cmd` 则让 Runner 等待 MySQL 对 `ping` 有响应后再开始步骤；它也不能证明应用账号、权限或表结构都正确。真正的连接和表结构仍要交给测试验证。

为什么要传入两种连接串？`MYSQL_GO_DSN` 使用 Go MySQL 驱动的 DSN 格式，供集成测试执行 SQL；`MIGRATE_TEST_DATABASE_URL` 带有 `mysql://` 前缀，供 `golang-migrate` 使用。两者都指向 `go_backend_test`。当前测试代码会核对数据库名，并在 `TestMain` 中先执行 `migrations/` 的 SQL，再运行订单测试。这样每次新建的空库才能拥有订单系统所需的表。`.env` 文件没有提交到 Runner 也没关系，因为这个 Job 已经直接提供了测试需要的环境变量。

## 完整的检查顺序

在 MySQL 准备好后，同一个 Job 的 Step 按顺序执行。当前 Workflow 使用 `actions/checkout@v7` 拉取仓库，使用 `actions/setup-go@v7` 根据根目录 `go.mod` 中的 `go 1.27.0` 安装 Go。[setup-go 的版本文件说明](https://github.com/actions/setup-go/blob/main/docs/advanced-usage.md)

```yaml
permissions:
  contents: read

jobs:
  test:
    steps:
      - name: 检出代码
        uses: actions/checkout@v7

      - name: 配置 Go 环境
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod

      - name: 检查Go格式
        run: test -z "$(gofmt -l .)"

      - name: 执行静态检查
        run: go vet ./...

      - name: 执行测试
        run: go test -count=1 ./...

      - name: 进行编译
        run: go build ./...
```

上面只摘录了权限和步骤；`runs-on`、数据库环境变量及 `services` 仍使用前一节的配置。`contents: read` 给予检出代码所需的只读权限。

检查顺序也有原因：`gofmt -l .` 列出格式不符合 Go 规范的文件，`test -z` 要求输出为空，否则这一步失败；`go vet ./...` 检查 Go 能识别的一类可疑代码；`go test -count=1 ./...` 运行所有包的测试，包括依赖 MySQL 的集成测试；最后 `go build ./...` 确认项目包可以编译。前一步失败，后面的普通步骤就不会继续执行。`go vet` 和编译通过也不等于业务正确，所以订单测试不可省略。

## 如何确认 CI 真的生效

把 `.github/workflows/ci.yml` 和代码推送到 GitHub 后，进入仓库的 **Actions** 页面，打开 `Go CI` 的运行记录。可以依次查看检出代码、配置 Go 环境、格式检查、静态检查、测试和编译的日志；若从 PR 触发，也可以在 PR 的检查状态里看到结果。[查看工作流运行记录](https://github.c-cpp.com/zh/actions/monitoring-and-troubleshooting-workflows/viewing-workflow-run-history)

如果失败，要看**第一个失败的步骤**。例如格式检查失败时，先在本机执行 `gofmt -l .` 找到文件；测试阶段出现 `failed to ping database`，重点检查 MySQL 健康状态、端口和 DSN；出现 Migration 错误，则检查 `migrations/` 文件和测试库是否一致。当前 `ci.yml` 已包含上述配置；是否在 GitHub Runner 上通过，还需要以实际的 Actions 运行记录为准。

这个 CI 的边界也很明确：它验证代码格式、静态问题、测试和编译，**还没有打包镜像或部署生产环境**。等这些检查稳定后，再讨论 Docker 构建与后续交付流程。
