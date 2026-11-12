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
        image: mysql:8.0.46
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

`contents: read` 给予检出代码所需的只读权限。

检查顺序也有原因：`gofmt -l .` 列出格式不符合 Go 规范的文件，`test -z` 要求输出为空，否则这一步失败；`go vet ./...` 检查 Go 能识别的一类可疑代码；`go test -count=1 ./...` 运行所有包的测试，包括依赖 MySQL 的集成测试；最后 `go build ./...` 确认项目包可以编译。前一步失败，后面的普通步骤就不会继续执行。`go vet` 和编译通过也不等于业务正确，所以订单测试不可省略。

## 如何确认 CI 真的生效

把 `.github/workflows/ci.yml` 和代码推送到 GitHub 后，进入仓库的 **Actions** 页面，打开 `Go CI` 的运行记录。可以依次查看检出代码、配置 Go 环境、格式检查、静态检查、测试和编译的日志；若从 PR 触发，也可以在 PR 的检查状态里看到结果。[查看工作流运行记录](https://github.c-cpp.com/zh/actions/monitoring-and-troubleshooting-workflows/viewing-workflow-run-history)

如果失败，要看**第一个失败的步骤**。例如格式检查失败时，先在本机执行 `gofmt -l .` 找到文件；测试阶段出现 `failed to ping database`，重点检查 MySQL 健康状态、端口和 DSN；出现 Migration 错误，则检查 `migrations/` 文件和测试库是否一致。当前 `ci.yml` 已包含上述配置；是否在 GitHub Runner 上通过，还需要以实际的 Actions 运行记录为准。

到这一步，CI 已验证代码格式、静态问题、测试和编译。下面继续让测试通过的代码自动构建 Docker 镜像。

# 自动构建 Docker 镜像

第一个 Go CI 已经能在干净的 Runner 上启动 MySQL、执行 Migration、运行订单测试并编译代码。接下来会遇到一个问题：**既然代码已经通过检查，能不能顺手把将来要部署的 Docker 镜像也构建出来？**

我们在上一章已经写好两阶段 `Dockerfile`：第一阶段用 Go 编译 `./cmd/server`，第二阶段只把可执行文件放进 Alpine 运行镜像。现在让 GitHub Actions 在测试通过后执行这份 Dockerfile。我们先验证镜像能构建，再把通过检查的镜像发布到 GHCR。

## 利用 `needs` 建立 Quality Gate

如果直接在 `jobs` 下新增 `docker-build`，它与 `test` 默认可以并行。这样可能出现订单集成测试还没结束，镜像却已经构建成功的情况。镜像构建成功只说明 Dockerfile 和源码可以打包，不能证明业务测试通过。

我们给构建 Job 增加依赖：

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    # Go 格式检查、静态检查、MySQL 集成测试与编译

  docker-build:
    needs: test
    runs-on: ubuntu-latest
    # 构建 Docker 镜像
```

`needs: test` 的意思是：只有 `test` Job 成功，`docker-build` 才开始；如果测试失败或被跳过，构建 Job 默认也会被跳过。这就是这里的 **Quality Gate（质量门禁）**。[GitHub Actions 的 Job 依赖说明](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-jobs)

为什么把两个 Job 分开？一方面，Actions 页面会分别显示“代码检查”和“镜像构建”的结果；另一方面，将来增加镜像发布步骤时，依赖关系仍然清楚。代价是两个 Job 使用各自的 Runner，`docker-build` 需要再次检出代码；`test` Job 中启动的 MySQL 容器也不会带到这个 Job。构建镜像不运行集成测试，因此这里不需要再配置 MySQL。

## Docker Buildx 与 BuildKit

我们本机执行过 `docker build`。在 Actions 里，当前项目用 `docker/setup-buildx-action` 准备 **Buildx**，再用 `docker/build-push-action` 执行构建。Buildx 是 Docker 的构建工具接口，背后使用 **BuildKit** 构建引擎。[Docker 的 GitHub Actions 构建指南](https://docs.docker.com/build/ci/github-actions/)

最初，`.github/workflows/ci.yml` 的第二个 Job 先写成：

```yaml
docker-build:
  needs: test
  runs-on: ubuntu-latest

  steps:
    - name: 检出代码
      uses: actions/checkout@v7

    - name: 配置 Docker Buildx
      uses: docker/setup-buildx-action@v4

    - name: 构建 Docker 镜像
      uses: docker/build-push-action@v7
      with:
        context: .
        file: Dockerfile
        push: false
        cache-from: type=gha
        cache-to: type=gha,mode=max
```

`context: .` 指向仓库根目录作为构建上下文，`file: Dockerfile` 选择根目录的 Dockerfile。当前 `.dockerignore` 排除了 `.env`、`.git/`、日志和临时目录，避免把它们送入构建上下文。Dockerfile 先复制 `go.mod` 和 `go.sum` 下载依赖，再复制源码编译；当只修改 Go 业务代码而依赖文件没变时，依赖下载层有机会复用缓存。

此时的 `push: false` 很关键：它只让构建步骤验证镜像可以生成，**不会把镜像推送到仓库**。这个初版配置也没有 `load: true`，因此不能把这个 Job 的结果当成已经安装到 Runner 本地 Docker 中、可供后续 Job 直接运行的镜像。想在别的机器上拉取它，还需要发布到 Registry。[build-push-action 输入说明](https://github.com/docker/build-push-action)

## Docker Build Cache

每次 GitHub Actions 运行通常使用新 Runner。第一次构建下载过的 Go 模块和构建层，下一次运行不能直接从上一台 Runner 的本地磁盘读取。我们在构建步骤中配置了外部缓存：

```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```

`type=gha` 表示 GitHub Actions Cache 后端。`cache-from` 尝试取回以前的 BuildKit 构建缓存；`cache-to` 把本次构建缓存导出，供后续运行复用。`mode=max` 会尽量导出包括中间构建阶段在内的缓存，因此对我们的两阶段 Dockerfile 更有用。[Docker 的 Actions Cache 说明](https://docs.docker.com/build/cache/backends/gha/)

这与 `test` Job 中 `actions/setup-go` 的缓存不是同一份东西：前者服务于 **Dockerfile 内的 BuildKit 构建**，后者服务于 **Runner 上执行的 Go 命令**。缓存也不是“保证每次都跳过构建”：首次运行没有可用缓存；修改 `go.mod` 或 `go.sum` 会影响依赖层；修改源码会让后续 `COPY . .` 和编译层重新计算。我们可以在 Actions 的“构建 Docker 镜像”日志里观察哪些步骤命中了缓存，再判断它有没有带来实际加速。

## Docker Container Registry：把镜像交给 GHCR

前面 `push: false` 的配置只验证 Dockerfile 能构建。Job 一结束，另一台服务器既拿不到这个 Runner 的构建结果，也不能凭源码提交号直接执行 `docker pull`。我们需要一个 **Container Registry（容器镜像仓库）** 保存并分发构建产物。这里选 GitHub Container Registry，简称 **GHCR**；镜像地址形如 `ghcr.io/<仓库所有者>/<仓库名>:<标签>`。[GitHub Container Registry 文档](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)

### 先想清楚镜像标签

假设今天和明天都向 `main` 提交了代码，两次构建都叫 `latest`。明天的镜像发布后，`latest` 就会指向新的内容。想重新部署今天的版本时，单靠 `latest` 无法准确找到它。

尽管GitHub Actions提供了功能强大的**表达式**，可以让我们手动拼接标签[GitHub Actions 表达式](https://docs.github.com/zh/actions/learn-github-actions/expressions)，但手动拼接容易出错。我们当前让 `docker/metadata-action` 根据触发事件生成标签，自动化程度更高，也更容易维护:

| 触发事件 | 生成的标签示例 | 是否推送 GHCR |
| --- | --- | --- |
| 向 `main` 提交 | `main`、`sha-a81f93c`、`latest` | 是 |
| 向 `main` 发起 PR | `pr-12`，以及对应的 `sha-*` | 否，只构建验证 |
| 推送 `v0.1.0` Git Tag | `0.1.0`、`sha-a81f93c` | 是 |

这里 `main` 和 `latest` 都是会移动的名称；`sha-*` 对应构建所用的提交，更方便追溯源码。Git Tag `v0.1.0` 经 `type=semver,pattern={{version}}` 生成镜像标签 `0.1.0`。但**标签本身可以被重新指向**；需要精确记录部署镜像时，还应保存 Registry 返回的 `sha256` Digest。Digest 标识镜像内容，适合作为日后核对或回滚的依据。[metadata-action 标签规则](https://github.com/docker/metadata-action)

### 在 Workflow 中登录、生成标签、发布

现在看项目 `.github/workflows/ci.yml` 的构建 Job。下面是与 GHCR 有关的配置；原有的检出代码、配置 Buildx 和 `needs: test` 仍然保留：

```yaml
on:
  push:
    branches: ["main"]
    tags: ["v*"]
  pull_request:
    branches: ["main"]

jobs:
  docker-build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      # 前面还有 checkout 和 setup-buildx
      - name: Login to GHCR
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v4
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: 生成tag
        id: tag
        uses: docker/metadata-action@v6
        with:
          images: ghcr.io/${{ github.repository }}
          flavor: latest=false
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=sha
            type=semver,pattern={{version}}
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}

      - name: 构建 Docker 镜像
        uses: docker/build-push-action@v7
        with:
          context: .
          file: ./Dockerfile
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.tag.outputs.tags }}
          labels: ${{ steps.tag.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

先看权限。`contents: read` 用于读取仓库，`packages: write` 允许这个 Job 发布 Package；它们只写在 `docker-build` Job 下，没有扩大到整个 Workflow。`GITHUB_TOKEN` 由 GitHub 为运行自动提供，发布与当前仓库关联的 Package 时无需在文件中写入长期令牌。[GitHub 的 Package 发布说明](https://docs.github.com/en/packages/managing-github-packages-using-github-actions-workflows/publishing-and-installing-a-package-with-github-actions)

再看发布条件。PR 也会触发测试和镜像构建，但登录步骤被 `if` 跳过，`push` 表达式也会得到 `false`；向 `main` 或 `v*` Tag 推送时，测试成功后才登录并推送。**生成 `pr-12` 这样的标签不等于把它发布到 GHCR。** `needs: test` 仍然守在发布前面：测试不通过，构建和推送都不会开始。

`metadata-action` 输出的 `tags` 与 `labels` 被下一步直接引用。`type=raw` 的条件参数必须写成 `enable=`；我们也设置 `flavor: latest=false`，关闭 SemVer Tag 默认自动生成 `latest` 的行为。这样 `latest` 只在 `main` 分支条件成立时由那条 `type=raw` 规则生成，推送版本 Tag 不会意外移动它。[metadata-action 的 `latest` 规则](https://github.com/docker/metadata-action#latest-tag)

这里还有一个权限边界：当前 `packages: write` 授给了整个 `docker-build` Job，而这个 Job 也会处理 PR。虽然 `if` 与 `push` 已阻止 PR 发布，权限本身仍比“只构建的 PR”所需范围大。Fork PR 的令牌通常会受到 GitHub 的只读限制；若以后要让同仓库 PR 也完全拿不到发布权限，应把只读构建与仅在可信推送事件运行的发布 Job 拆开。**不能只把 `push: false` 当作权限隔离。**[GitHub 关于 PR 令牌的说明](https://docs.github.com/en/actions/reference/security/securely-using-pull_request_target)

### 如何验证这一阶段

向 `main` 推送后，在 Actions 中确认 `test` 成功、`docker-build` 成功，再到仓库的 Packages 页面查找对应镜像。如果在 GitHub 看到了第一个 Package，就说明首次发布完成；接下来可以核对它的 `main`、`sha-*`、`latest` 标签与 Digest。若 Package 是私有的，其他机器拉取前还需要具备相应的读取权限。[GHCR 的访问说明](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)

目前我们完成的是 **源码 → 测试 → 构建 → GHCR 发布**。镜像已可分发，但还没有把它接入生产环境的 Docker Compose，也没有完成正式 Release 和部署回滚。下一节再把 Git Tag、镜像 Tag 与 Digest 的关系用于版本化交付。

# Release 与部署回滚
## 版本管理：语义化版本
## 部署回滚
