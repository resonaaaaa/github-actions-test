# Docker
现在我们的GO项目依赖于本地的MySQL和环境变量，当换一台PC或者在服务器上部署时，我们需要重新安装MySQL和配置环境变量，这样会很麻烦。Docker就是为了解决这个问题而诞生的，它可以把我们的应用程序和它的依赖打包在一个容器里，这样我们就可以在任何地方运行它，而不需要担心环境配置的问题。

传统部署：
- 安装Go环境
- 安装MySQL
- 下载依赖
- 配置环境变量
- 编译项目
- 启动服务

Docker部署：
- 获取Docker Image
- 运行Docker容器
- 启动服务

我们先明确几个概念：
|概念|解释|
|---|---|
|Docker|一个开源的应用容器引擎，它可以让开发者打包他们的应用以及依赖包到一个轻量级、可移植的容器中，然后发布到任何流行的Linux机器上，也可以实现虚拟化。|
|Dockerfile|一个文本文件，它包含了构建Docker镜像所需的所有指令。|
|image(镜像)|一个只读的运行包，它包含了运行容器所需的所有信息，包括操作系统、应用程序和依赖项。|
|container(容器)|一个运行中的镜像实例，它是镜像的一个可写层。|
|volume(卷)|一个可供容器使用的存储空间，它可以在容器之间共享数据，使得容器重建后数据不会丢失。|
|network(网络)|Docker提供的网络功能，可以让容器之间相互通信。|
|Compose|一个工具，它描述并一起启动项目与其依赖的服务。它使用YAML文件来配置应用程序的服务、网络和卷。|

接下来，我们将使用Docker来部署我们的GO项目和MySQL数据库。我们将创建一个Dockerfile来构建我们的GO应用程序镜像，并使用docker-compose.yml文件来定义和运行我们的多容器应用程序。

# Dockerfile
Dockerfile是一个文本文件，它包含了构建Docker镜像所需的所有指令。
## Dockerignore
.dockerignore文件是一个文本文件，它告诉Docker在构建镜像时应该忽略哪些文件和目录。它的作用类似于.gitignore文件，可以帮助我们减少镜像的大小，提高构建速度。
### Build Context
在构建Docker镜像时，Docker会将指定的目录作为构建上下文(context)，简单来说，就是Docker在构建镜像时可以访问这个目录下的所有文件和子目录。

如果我们现在的项目目录结构如下：
```
go-backend/
│
├── .env
├── .env.example
├── go.mod
├── go.sum
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── application/
│   ├── config/
│   ├── repository/
│   └── service/
│
└── migrations/
```
然后我们执行`build`命令：
```bash
docker build -t go-backend .
```
> `-t`选项用于指定镜像的名称，末尾的`.`表示当前目录，也就是`go-backend`目录作为构建上下文。

在执行上面的`build`命令时，Docker会将整个`go-backend`目录作为构建上下文，这样Docker就会访问目录中的所有文件。这带来了几个问题：
- **增加构建开销**：有的文件可能不需要被包含在镜像中，比如`.env`文件、测试文件、日志文件等，这些文件会增加构建的开销。
- **降低缓存效率**：Docker在构建镜像时会使用缓存机制，如果构建上下文中包含了不必要的文件，可能会导致缓存失效，从而增加构建时间。
- **增加安全风险**：有些文件可能包含敏感信息，比如`.env`文件中的数据库密码，如果这些文件被包含在镜像中，可能会导致安全风险。

### Dockerignore
为了避免这些问题，我们可以使用`.dockerignore`文件来告诉Docker在构建镜像时应该忽略哪些文件和目录。

我们一般排除哪些文件和目录呢？通常我们会排除以下几类文件：
- **环境变量文件**：比如`.env`文件，这些文件通常包含敏感信息，比如数据库密码、API密钥等，不应该被包含在镜像中。
- **编译和测试产物**：比如`bin/`、`obj/`、`*.exe`、`*.dll`等，这些文件通常是编译和测试的产物，不属于源代码的一部分，不应该被包含在镜像中。
- **版本控制文件**：比如`.git/`、`.svn/`、`.hg/`等，这些文件是版本控制系统的元数据，不应该被包含在镜像中。
- **日志文件**：比如`*.log`、`logs/`等，这些文件通常是运行时产生的日志文件，不应该被包含在镜像中。
- **IDE配置文件**：比如`.idea/`、`.vscode/`、`*.iml`等，这些文件是本地IDE的配置文件，不应该被包含在镜像中。

但要注意的是：不要误排除了`migrations/`目录，因为这个目录下的数据库迁移文件是我们应用程序的一部分，应该被包含在镜像中。

还有一点可能的误区：`.dockerignore`文件的作用是告诉Docker在构建镜像时应该忽略哪些文件和目录，而不是阻止Docker在运行容器时访问相应文件和目录。比如我们在`.dockerignore`文件中排除了`.env`文件，这并不意味着我们在运行容器时就不能访问`.env`文件了。在Compose中，我们可以通过`volumes`选项将宿主机的`.env`文件挂载到容器中，这样容器就可以访问`.env`文件了。

## Dockerfile的基本指令
先来认识几个Dockerfile的基本指令：
|指令|作用|
|---|---|
|FROM|指定基础镜像|
|WORKDIR|设置工作目录|
|COPY|复制文件到镜像中|
|ADD|添加文件到镜像中|
|RUN|执行构建命令|
|CMD|指定容器启动时执行的命令|


----
`From`指令用于指定基础镜像，它是Dockerfile的第一条指令。基础镜像可以是官方提供的镜像，也可以是我们自己构建的镜像。比如：
```dockerfile
FROM golang:1.24 AS builder
```
> 这条指令指定了一个名为`golang:1.24`的官方镜像作为基础镜像，并且给这个构建阶段起了一个名字`builder`。我们可以在后续的Dockerfile中使用这个名字来引用这个构建阶段。

----
`WORKDIR`指令用于设置工作目录，它会在镜像中创建一个目录，并将其设置为当前工作目录。比如：
```dockerfile
WORKDIR /app
```
之后就会转入容器内的`/app`目录，后续的指令比如`COPY`、`RUN`等都会在这个目录下执行。

----
`COPY`指令用于将文件从宿主机复制到镜像中。它有两个参数：源路径和目标路径。比如：
```dockerfile
COPY go.mod go.sum ./
```
> 这条指令将宿主机的`go.mod`和`go.sum`文件复制到镜像的当前工作目录下。

注意：`COPY`指令不会复制`.dockerignore`文件中排除的文件和目录，也就是只能复制Docker Build Context中存在的文件和目录。

----
`ADD`指令和`COPY`指令类似，它也用于将文件从宿主机复制到镜像中。不同的是，`ADD`指令可以自动解压tar文件，并且可以从远程URL下载文件。比如：
```dockerfile
ADD https://example.com/file.tar.gz /app/
```
> 这条指令会从远程URL下载一个tar文件，并将其解压到镜像的`/app/`目录下。

----
`RUN`指令用于构建镜像时执行命令。

```dockerfile
RUN go mod download
```
> 这条指令会在构建镜像时执行`go mod download`命令，用于下载项目的依赖包。

我们还可以使用`RUN`指令执行多个命令，比如：
```dockerfile
RUN go get -u github.com/gin-gonic/gin && \
    go get -u github.com/joho/godotenv && \
    go mod download
```
----

`CMD`指令用于指定容器启动时执行的命令。它有多种形式，最常用的是exec形式和shell形式。

Exec形式：
```dockerfile
CMD ["./server"]
```

shell形式：
```dockerfile
CMD ./server
```


这两条指令的含义都是在容器启动时执行相应的命令。但相比于shell形式，exec形式的优点是不会启动一个额外的shell进程，因而能让应用程序成为容器的**主进程**，方便信号的传递和处理。


----
一个比较完整的Dockerfile示例：
```dockerfile
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go build -o server ./cmd/server/main.go
CMD ["./server"]
```

## 多阶段构建
多阶段构建是Dockerfile的一种优化技术，它可以让我们在一个Dockerfile中使用多个`FROM`指令，每个`FROM`指令都会创建一个新的构建阶段。这样我们就可以在一个阶段中构建应用程序，在另一个阶段中创建最终的镜像，从而减少镜像的大小。

我们可以以Go的基础镜像作为构建的开始：
```dockerfile
FROM golang:1.24
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["./server"]
```
如果我们直接将其作为最终镜像，那么这个镜像会包含Go的编译器和所有的依赖包，这样会导致镜像的体积很大。在编译完成后，实际上Go的编译器与依赖包已经不再需要了，我们只需要将编译好的二进制文件放到最终的镜像中即可。

这里有一个细节：先复制`go.mod`和`go.sum`文件，然后执行`go mod download`命令，这样可以利用Docker的缓存机制。如果我们修改了源代码，但没有修改依赖包，那么Docker会使用缓存的依赖包，而不需要重新下载，从而提高构建速度。

----
在最终的镜像中，我们可以使用一个更小的基础镜像，比如`alpine`，它是一个轻量级的Linux发行版。我们可以在最终的镜像中只包含必要的运行时依赖，而将编译好的二进制文件从构建阶段复制到最终的镜像中：
```dockerfile
FROM alpine
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
```
这里最重要的代码是`COPY --from=builder /app/server .`，它的意思是从名为`builder`的构建阶段中复制`/app/server`文件到最终镜像的当前工作目录下。这是我们多阶段构建的关键，它可以让我们在最终的镜像中只包含必要的文件，从而减少镜像的大小。

----
之后就可以写出一个完整的多阶段构建的Dockerfile：
```dockerfile

# Stage 1: Builder
FROM golang:1.27.3-bookworm AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build Linux executable
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.24

WORKDIR /app

# Copy executable from builder
COPY --from=builder /app/server .

# Configure runtime user
USER nobody

# Start application
CMD ["./server"]
```
其中`ENV CGO_ENABLED=0 GOOS=linux`是为了在构建Linux可执行文件时禁用CGO，并指定目标操作系统为Linux，这样可以确保生成的二进制文件可以在Alpine镜像中运行。

# 构建Docker镜像
构建镜像前需要确认Docker 已经安装并启动，可以通过`docker version`命令查看Docker的版本信息。
```bash
docker version
```
当能看到`Client`和`Server`的版本信息时，说明Docker已经安装并启动成功。

确认Dockerfile和.dockerignore文件已经创建好后，就可以使用`docker build`命令来构建Docker镜像了。我们可以在项目根目录下执行以下命令：
```bash
docker build -t go-backend .
```
> `-t`选项用于指定镜像的名称，末尾的`.`表示当前目录，也就是`go-backend`目录作为构建上下文。

可以使用`--progress=plain`选项来查看构建过程的详细信息：
```bash
docker build --progress=plain -t go-backend .
```

构建完成后可以查看镜像的信息：
```bash
docker images
docker image inspect go-backend:v1 --format '{{.Os}}/{{.Architecture}}'
```
> 最终输出应该类似Linux/amd64

为了验证多阶段构建是否成功，我们可以使用下面的命令查看还有没有Go的编译器和依赖包：
```bash
docker run --rm --entrypoint sh go-backend:v1 -c "id; ls -lah /app; command -v go || echo 'Go compiler not found'"
```
> 预期结果中应该包含`Go compiler not found`，说明最终镜像中没有Go的编译器。

# Docker Compose
完成Docker镜像的构建后，我们可以尝试使用`docker run`命令来运行Docker容器。在项目根目录下执行以下命令：
```bash
docker run --rm --name go-order-api go-backend:v1
```
> `--rm`选项表示容器在停止后会自动删除，`--name`选项用于指定容器的名称，`go-backend:v1`是我们刚刚构建的Docker镜像的名称。

但结果会告诉我们：*缺少必须的环境变量*。这是因为我们之前在`.dockerignore`文件中排除了`.env`文件，而我们的应用程序需要读取`.env`文件中的环境变量。Docker提供了`--env`与`--env-file`选项来设置环境变量，前者是手动设置单个环境变量，后者应用指定的环境变量文件：
```bash
docker run --rm --name go-order-api --env APP_ENV=production go-backend:v1 #设置单个环境变量APP_ENV为production
docker run --rm --name go-order-api --env-file .env go-backend:v1   #应用环境变量文件
```
----
有了配置文件，还有一个小问题：容器中的Localhost和宿主机的Localhost是不同的，容器中的Localhost指的是**容器本身**，而不是宿主机。虽然我们可以使用`host.docker.internal`来访问宿主机的服务，但这并不是一个很好的解决方案。

在接下来，我们要学习Docker Compose，它可以让我们把MySQL服务一并部署到Docker容器中，这样就不需要依赖宿主机的MySQL服务了。

## 什么是Docker Compose
Docker Compose是一个用于定义和运行多容器Docker应用程序的工具。通过Compose，我们可以使用一个YAML文件来配置应用程序的服务、网络和卷，然后使用一个命令就可以启动所有的服务。

首先回答一个问题：为什么不把MySQL也打包到我们的GO应用程序的Docker镜像中呢？原因有两个：
1. **职责分离**：一个容器应该只运行一个服务，这样可以更好地管理和维护容器。如果我们把MySQL和GO应用程序打包在同一个容器中，那么当我们需要升级MySQL或者GO应用程序时，就需要重新构建整个镜像，这样会增加维护的复杂度。
2. **数据持久化**：MySQL是一个数据库服务，它需要持久化存储数据。如果我们把MySQL和GO应用程序打包在同一个容器中，那么当容器停止或者删除时，MySQL的数据也会丢失。通过使用Docker Compose，我们可以将MySQL的数据存储在宿主机的卷中，从而实现数据的持久化。

因此需要把MySQL单独部署到一个容器中，这就涉及到了多容器的部署问题。Docker Compose可以帮助我们解决这个问题，它可以让我们在一个YAML文件中定义多个服务，然后使用一个命令就可以启动所有的服务，并且还提供了**内部网络**，让容器之间可以相互通信。

## compose.yaml

`compose.yaml` 是 Compose 项目的运行说明：它声明要启动哪些服务、为容器传入什么配置，以及数据存在哪里。当前先让 MySQL 独立运行；Go API、健康检查和 Migration 任务会在验证数据库后逐步加入。
- **服务**：服务是Docker Compose应用程序的核心概念，它表示一个运行中的容器。我们可以在compose.yml文件中定义多个服务，每个服务都可以有自己的镜像、环境变量、端口映射、卷挂载等配置。
- **网络**：网络是Docker Compose应用程序的另一个核心概念，它表示容器之间的通信通道。Docker Compose会为每个应用程序创建一个默认的网络，所有的服务都会连接到这个网络上，从而实现容器之间的通信。
- **卷**：卷是Docker Compose应用程序的另一个核心概念，它表示容器之间共享的数据存储空间。我们可以在compose.yml文件中定义卷，然后将卷挂载到服务的容器中，从而实现数据的持久化和共享。


### 当前的 MySQL 服务

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: go_backend
      MYSQL_USER: ${MYSQL_USER}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
    ports:
      - "127.0.0.1:3307:3306"
    volumes:
      - mysql-data:/var/lib/mysql

volumes:
  mysql-data:
```

`services.mysql` 定义名为 `mysql` 的服务，`image` 指定 MySQL 8.0 镜像。服务名由我们指定；以后 Go API 加入同一个 Compose 项目后，可以在默认网络中使用 `mysql:3306` 访问它，而不必依赖容器 IP。当前只有一个服务，也不需要手写 `networks`。[Compose 网络文档](https://docs.docker.com/compose/how-tos/networking/)

`environment` 把变量传入 **MySQL 容器**。其中 `MYSQL_DATABASE=go_backend` 用于首次初始化时创建数据库；它不会执行项目 `migrations/` 中的 SQL，也不会自动创建订单业务表。`MYSQL_ROOT_PASSWORD` 是管理员密码，`MYSQL_USER` 和 `MYSQL_PASSWORD` 是应用账号及其密码。MySQL 的初始化变量只对空数据目录生效；以后修改 `.env` 中的密码，不会重设已有数据卷里的账号。[MySQL 官方镜像说明](https://hub.docker.com/_/mysql)

### 端口与数据生命周期

`"127.0.0.1:3307:3306"` 按“Windows 监听地址 : Windows 端口 : 容器端口”理解：DBeaver 和 Windows 上运行的 `golang-migrate` 连接 `127.0.0.1:3307`，MySQL 容器内部监听 `3306`。宿主机 MySQL 仍使用 `127.0.0.1:3306`；而Go API 容器可以直接通过服务名 `mysql:3306` 连接，不走 Windows 的 `3307` 端口。

`mysql-data:/var/lib/mysql` 把 MySQL 数据目录放入命名卷。重建 MySQL 容器时，数据卷可以继续使用；`docker compose down -v` 则会删除 Compose 声明的命名卷，数据库数据也会随之丢失。

### `.env` 与变量插值

`${MYSQL_ROOT_PASSWORD}` 等写法是 **Compose 解析 YAML 时的变量插值**。

我们需要在项目根目录的本地 `.env` 中提供这三个值。若变量缺失，普通 `${VAR}` 会产生警告并替换为空值；可以改用 `${VAR:?VAR is required}`，让配置检查直接失败。

PowerShell 中的同名环境变量优先于 `.env`。[Compose 变量插值文档](https://docs.docker.com/compose/how-tos/environment-variables/variable-interpolation/)

此时本机运行 Go 程序的 `MYSQL_DSN` 仍指向本机数据库。将来配置 Go API 容器时，需要单独向该容器传入指向 `mysql:3306` 的 DSN；Compose 不会根据 MySQL 的初始化变量自动生成它。

### 验证配置

在项目根目录的 PowerShell 中执行：

```powershell
docker compose --env-file .env config --quiet
docker compose up -d mysql
docker compose ps
```

`config --quiet` 检查解析后的配置而不打印密码值；`up -d mysql` 只启动数据库服务；`ps` 查看容器和端口状态。[Compose config 命令](https://docs.docker.com/reference/cli/docker/compose/config/)

## 初始化容器中的 MySQL

我们已经能通过 Compose 启动 MySQL 容器，但容器显示为 `Running`，是否就意味着订单系统可以连接数据库了？**不一定。** MySQL 第一次启动时还要初始化数据目录、创建账号和数据库；在这段时间里，容器进程已经运行，数据库却可能尚未接受连接。[MySQL 官方镜像说明](https://hub.docker.com/_/mysql)

目前 `compose.yaml` 将 `mysql-data` 挂载到 `/var/lib/mysql`。当这个数据目录为空时，MySQL 镜像会读取 `MYSQL_ROOT_PASSWORD`、`MYSQL_DATABASE`、`MYSQL_USER` 和 `MYSQL_PASSWORD`，创建 `go_backend` 数据库及应用账号。**这些变量是首次初始化配置，不是每次启动都会重置数据库。** 如果数据卷里已有数据库，仅修改 `.env` 中的密码不会修改已有账号的密码。

我们可以在 DBeaver 中新建一个连接，检查容器数据库是否已经可用：

| 数据库 | DBeaver 连接地址 | 用途 |
|---|---|---|
| Windows 本机 MySQL | `127.0.0.1:3306` | 原有开发数据库 |
| Compose 中的 MySQL | `127.0.0.1:3307` | 当前容器数据库 |

连接时选择数据库 `go_backend`，使用 `.env` 中设置的 `MYSQL_USER` 和 `MYSQL_PASSWORD`。连接成功后执行：

```sql
SELECT DATABASE(), VERSION();
SHOW TABLES;
```

第一条语句确认我们连接到了哪个数据库以及 MySQL 的版本；第二条语句查看业务表。即使 `go_backend` 已创建，`SHOW TABLES` 仍可能没有 `users`、`products`、`orders`：**创建数据库与执行 Migration 是两件事。** 后面还需要用现有的 `golang-migrate` 和 `migrations/` 文件初始化订单系统的表结构。

这里也要记住数据卷的作用：普通的 `docker compose down` 会移除容器，但不会移除声明的命名卷；`docker compose down -v` 会连命名卷一起删除。我们不应为了让初始化变量重新生效而随意删除已有数据。[Compose down 命令](https://docs.docker.com/reference/cli/docker/compose/down/)

## Healthcheck

假如我们接下来把 Go API 加入 Compose，API 启动时会调用 `db.PingContext()`。如果 MySQL 容器刚进入 `Running` 状态、仍在初始化，API 的连接检查就可能失败。我们需要让 Docker 定期检查**数据库服务是否已响应**，这就是 Healthcheck（健康检查）。

Healthcheck 必须写在对应服务的配置里，也就是 `services.mysql` 之下，结构类似：

```yaml
services:
  mysql:
    # image、environment、ports、volumes
    healthcheck:
      test: ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -u root -p\"$${MYSQL_ROOT_PASSWORD}\" --silent"]
      interval: 5s
      timeout: 3s
      retries: 10
      start_period: 30s
```

- `CMD-SHELL` 表示在容器中通过 shell 执行检查命令；
- `-h 127.0.0.1` 让检查通过容器内的 TCP 地址访问 MySQL。
- `$$`是 Compose 的转义写法，`$ [variable]`会被Compose在解析YAML时替换为变量值，而写成`$$[variable]`则避免被Compose解析，留给容器内的 shell 解析。
- `--silent` 让 `mysqladmin ping` 只返回退出码，不输出任何文本。

各参数的作用如下：[Compose 服务配置文档](https://docs.docker.com/reference/compose-file/services/)

| 参数 | 含义 |
|---|---|
| `test` | 执行什么检查；命令退出码为 `0` 表示成功 |
| `interval` | 两次检查之间的时间 |
| `timeout` | 单次检查的超时时间 |
| `retries` | 连续失败多少次后标记为 `unhealthy` |
| `start_period` | 首次启动的宽限时间，期间的失败不计入连续失败次数 |

这里的 `mysqladmin ping` 只检查 MySQL Server 是否响应。即使账号认证失败，它也可能返回成功的退出码；因此 **`healthy` 不能证明应用账号有权限，更不能证明 Migration 已执行完毕**。我们仍要用 DBeaver 的实际连接和 SQL 查询验证数据库，再单独执行 Migration。[MySQL `mysqladmin` 手册](https://dev.mysql.com/doc/refman/8.0/en/mysqladmin.html)

在项目根目录的 PowerShell 中验证：

```powershell
docker compose --env-file .env config --quiet
docker compose up -d mysql
docker compose ps
```

第一条命令检查 Compose 文件，第二条应用 MySQL 服务的配置，第三条查看容器是否从 `starting` 进入 `healthy`。后面加入 Go API 时，可以用 `depends_on` 的 `service_healthy` 条件等待数据库通过检查；但业务表是否就绪还要由 Migration 步骤保证。[Compose 启动顺序文档](https://docs.docker.com/compose/how-tos/startup-order/)

----
完成了某项服务的Healthcheck配置后，我们还需要**让Compose监测服务的就绪状态**！

之后配置 Go API或者其他服务时，`depends_on` 可以写成：

```yaml
services:
  api:
    image: go-backend:v1
    depends_on:
      mysql:
        condition: service_healthy
```
这与普通的：
```yaml
services:
  depends_on:
    - mysql
```
不同，后者只保证 MySQL 容器进入 `Running` 状态，而不保证数据库服务已经响应。而 `condition: service_healthy` 会等待 MySQL 容器通过健康检查后才启动 API 服务。

## 独立的 Migration 任务

我们已经让 MySQL 显示为 `healthy`，但这里还有一个问题：数据库能接受连接，是否意味着订单系统可以处理请求？

**不一定。** `go_backend` 数据库可能已经创建，而 `users`、`products`、`orders` 等业务表还没有创建。Healthcheck 检查数据库服务，Migration 则负责把 Schema 更新到应用需要的版本。

我们需要让这一步也成为 Compose 中明确的一项任务，而不是依赖每个人手动在 Windows 上运行命令。

### 为什么让 Migration 独立运行？

假如把 `migrate up` 写进 Go API 的 `main.go`，以后同时启动多个 API 实例时，每个实例都会尝试迁移数据库。这样数据库结构变更与 API 的启动、重启绑在一起，失败时也更难判断是哪一步出了问题。

我们希望它们有不同的生命周期：

```text
MySQL 通过 Healthcheck
          ↓
Migration 执行并成功退出
          ↓
Go API 开始处理请求
```

Migration 是**一次性任务**：运行已有的 SQL 文件，完成后退出；Go API 才是长期运行的服务。我们继续使用项目原有的 `migrations/` 和 `golang-migrate`，不需要在 `main.go` 中再实现一套迁移逻辑。

### `migrate` 服务的配置

当前 `compose.yaml` 中的 `mysql` 与 `migrate` 都是 `services` 下的服务。新增的配置如下：

```yaml
  migrate:
    image: migrate/migrate:v4.18.3
    depends_on:
      mysql:
        condition: service_healthy
    volumes:
      - ./migrations:/migrations:ro
    command:
      - "-path=/migrations"
      - "-database=${DOCKER_MIGRATE_DATABASE_URL}"
      - "up"
    restart: "no"
```

这里有几个细节：

- `depends_on` 等待 MySQL 的 Healthcheck 通过后才运行迁移，避免数据库仍在初始化时就尝试连接。
- `./migrations:/migrations:ro` 把项目里的 SQL 文件只读挂载到容器中。左侧是 Windows 项目目录，右侧是 Migration 容器中的路径；`:ro` 表示容器只能读取这些文件。
- `command` 相当于对容器执行 `migrate -path=/migrations -database=... up`，按版本顺序应用尚未执行的迁移。
- `restart: "no"` 表示任务完成或失败后不自动重启。迁移失败应当明确暴露出来，而不是反复重试。

`service_healthy` 只说明 MySQL 已响应；它不会检查 Schema 版本。因此 Migration 必须有自己的成功或失败结果。[Compose 服务依赖文档](https://docs.docker.com/reference/compose-file/services/)

### Migration 容器连接哪个地址？

Migration 容器和 MySQL 容器在同一个 Compose 网络中，所以数据库主机名是服务名 `mysql`，端口是容器内部的 `3306`。Windows 上的 DBeaver 或 `migrate` CLI 才使用映射地址 `127.0.0.1:3307`。

在本地 `.env` 中为容器任务单独准备 `DOCKER_MIGRATE_DATABASE_URL`，格式例如：

```dotenv
DOCKER_MIGRATE_DATABASE_URL=mysql://app:<password>@tcp(mysql:3306)/go_backend
```

`app` 应与 `MYSQL_USER` 一致，密码应与 `MYSQL_PASSWORD` 一致。示例中的 `<password>` 只是占位符；如果真实密码含 URL 特殊字符，需要按连接 URL 的规则编码。这个地址也不能直接作为 Go API 的 `MYSQL_DSN`：Go MySQL 驱动使用的 DSN 没有开头的 `mysql://`。本机开发数据库的 DSN 保持独立。

### 验证一次性任务

在项目根目录的 PowerShell 中执行：

```powershell
docker compose --env-file .env config --quiet
docker compose --env-file .env run --rm migrate
```

第一条命令检查 Compose 配置；第二条创建并运行一次 Migration 容器，任务结束后由 `--rm` 删除这个临时容器。`config --quiet` 只检查配置，不会验证数据库连接是否真的可用；也不要公开完整的 `docker compose config` 展开结果，因为其中可能包含密码。

等我们把 Go API 加入 Compose 时，还要让 API 依赖 `migrate` 的 `service_completed_successfully` 状态。这个条件要求一次性任务成功退出，而不只是开始运行；只有这样，API 才能在所需表结构准备好后启动。[Compose 启动顺序文档](https://docs.docker.com/compose/how-tos/startup-order/)

## 加入 Go API 服务

我们现在已经有三个组件：MySQL 容器能通过 Healthcheck，Migration 容器能执行 `up` 并成功退出，Go API 的镜像也已经构建。接下来的问题是：**怎样保证 API 启动时，数据库不仅能连接，而且已经有订单系统需要的表？**

我们希望 Compose 按下面的关系启动服务：

```text
MySQL healthy
      ↓
Migration 成功完成（退出码为 0）
      ↓
Go API 启动
```

当前 `compose.yaml` 已经加入 `api` 服务。它的关键配置可以概括为：

```yaml
  api:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      mysql:
        condition: service_healthy
      migrate:
        condition: service_completed_successfully
    environment:
      APP_ENV: development
      HTTP_ADDR: :8080
      MYSQL_DSN: ${DOCKER_MYSQL_DSN:?DOCKER_MYSQL_DSN is required}
      # 其余连接池和 HTTP 超时配置沿用现有 compose.yaml
    ports:
      - "127.0.0.1:8080:8080"
    restart: unless-stopped
```

这里展示的是学习用的关键片段；实际文件中还明确设置了 MySQL 连接池和 HTTP 超时参数。

### 构建与配置注入

`build.context: .` 表示以项目根目录为构建上下文，`dockerfile: Dockerfile` 指向我们已经写好的多阶段构建文件。因此执行 `docker compose up --build` 时，Compose 可以先构建 Go API 镜像，再创建容器。镜像中只有运行所需的程序,本地 `.env` 没有被打包进去。

Go 程序启动时由 `config.Load()` 读取 `MYSQL_DSN`、`HTTP_ADDR` 等环境变量，`.env` 文件不存在也不会阻止启动。

Compose 中的：
```yaml
MYSQL_DSN: ${DOCKER_MYSQL_DSN:?DOCKER_MYSQL_DSN is required}
```

意思是先从 PowerShell 环境或指定的 `.env` 中读取 `DOCKER_MYSQL_DSN`，再把它作为容器内的 `MYSQL_DSN`。`${VAR:?错误信息}` 要求变量存在且非空；缺失时，Compose 在创建容器前就会报错。[Compose 变量插值文档](https://docs.docker.com/compose/how-tos/environment-variables/variable-interpolation/)

为什么需要单独的 `DOCKER_MYSQL_DSN`？因为同一个 `127.0.0.1` 在 Windows 和 API 容器中指向不同的地方。Windows 本机运行 Go 程序时可以连接本机 MySQL 的 `127.0.0.1:3306`；API 容器则通过 Compose 服务名连接 `mysql:3306`。例如：

```dotenv
DOCKER_MYSQL_DSN=app:<password>@tcp(mysql:3306)/go_backend?parseTime=true
```

这里使用 Go MySQL 驱动的 DSN 格式，与migration 工具使用的格式不同。**不要把带有 `mysql://` 前缀的 `DOCKER_MIGRATE_DATABASE_URL` 直接交给 Go API**。我们只向 API 注入它需要的变量，而不是把本机开发、测试和迁移用的所有配置一并注入。

### 端口与启动依赖

API 的 `HTTP_ADDR: :8080` 让服务监听容器的 8080 端口。`"127.0.0.1:8080:8080"` 再把 Windows 本机的 8080 映射到容器的 8080。这样 Windows 上的客户端可以访问 `http://127.0.0.1:8080/orders`，而 API 访问 MySQL 仍使用 Compose 内部地址 `mysql:3306`，不经过 Windows 映射的 `3307`。[Compose 网络文档](https://docs.docker.com/compose/how-tos/networking/)

`depends_on` 中的两个条件解决不同的问题：

| 依赖 | 等待的条件 | 保证什么 |
|---|---|---|
| `mysql` | `service_healthy` | MySQL 已通过健康检查 |
| `migrate` | `service_completed_successfully` | 一次性迁移任务成功退出 |

Migration 容器只是 `started` 还不够，因为 SQL 可能尚未执行完毕或执行失败。API 要等待它**成功完成**，才能开始处理订单。[Compose 服务依赖文档](https://docs.docker.com/reference/compose-file/services/)

`restart: unless-stopped` 用于长期运行的 API：容器退出后 Docker 会尝试重启，直到我们明确停止它；Migration 使用 `restart: "no"`，成功或失败后都退出。重启策略不能修复错误的 DSN、密码或数据库结构。如果 API 反复退出，应该查看日志并处理原因。

### 验证完整启动

 `api` 配置已写入并成功启动。下面的命令可以在项目根目录的 PowerShell 中复现和检查完整启动过程：

```powershell
docker compose --env-file .env config --quiet
docker compose --env-file .env up -d --build
docker compose --env-file .env ps -a
```
启动后，我们希望看到 MySQL 为 `healthy`、Migration 为 `Exited (0)`、API 为 `Up`。Migration 再次运行时输出 `no change` 是正常的，表示当前数据库已是最新版本。

如果 API 没有保持运行，先查看它的日志：

```powershell
docker compose --env-file .env logs --tail=50 api
```

# 在Docker 中运行测试

我们的三个服务现在都能启动：MySQL 为 `healthy`，Migration 成功退出，Go API 处于 `Up`。接下来我们就可以进行业务测试：在 Windows 发起请求，经过 API 容器访问 MySQL 容器，再检查订单、库存和数据卷中的真实结果。

这里做的是容器化环境的端到端实验：请求从 PowerShell 发出，数据库结果在 DBeaver 中核对。它与项目原有的 `go test` 集成测试相互补充。现有集成测试通过 `MIGRATE_TEST_DATABASE_URL` 管理独立测试数据库，不能直接对保存实验订单的数据库执行数据清理。

## 验证订单创建服务

先确认服务状态：

```powershell
docker compose --env-file .env ps -a
docker compose --env-file .env logs --tail=50 api
```

预期是 MySQL `healthy`、Migration `Exited (0)`、API `Up`，API 日志中没有数据库连接或配置错误。随后在 DBeaver 连接 Docker MySQL 的 `127.0.0.1:3307`，找到一个有效用户和库存足够的商品：

```sql
SELECT id, name, email FROM users;
SELECT id, name, price_cents, stock FROM products;
```

记下商品原有库存。Migration 负责建表，不负责保证这些业务测试数据存在；如果查询不到用户或商品，需要先在容器数据库中准备数据，不能直接沿用本机 MySQL 中的 ID。

我们的接口接收 `POST /orders`。请求体与请求头如下：
```powershell
$body = @{
    user_id = 1       # 改成 Docker MySQL 中实际存在的用户 ID
    product_id = 10   # 改成实际存在且库存足够的商品 ID
    quantity = 2
} | ConvertTo-Json

$headers = @{ "Idempotency-Key" = "docker-note-$(Get-Date -Format 'yyyyMMddHHmmss')" }

$order = Invoke-RestMethod 
    -Uri "http://127.0.0.1:8080/orders" 
    -Method POST 
    -Headers $headers 
    -ContentType "application/json" 
    -Body $body

$order
```

第一次请求成功时，Handler 返回 `201 Created` 和订单的 `id`、`user_id`、`product_id`、`quantity`、`unit_price`、`created_at`。

接下来到DBeaver验证：

```sql
SELECT id, user_id, product_id, quantity, unit_price, request_id
FROM orders
ORDER BY id DESC
LIMIT 5;

SELECT id, name, stock
FROM products
WHERE id = 10;
```

我们要确认订单确实写入 Docker MySQL，且库存恰好减少了购买数量。这里只看到 HTTP 成功还不够，因为我们也需要验证事务带来的数据库变化。

接着复用上面已经创建的 `$body` 和 `$headers`，用**完全相同的请求体与 `Idempotency-Key`** 再发送一次：

```powershell
$repeat = Invoke-RestMethod `
    -Uri "http://127.0.0.1:8080/orders" `
    -Method POST `
    -Headers $headers `
    -ContentType "application/json" `
    -Body $body

$repeat
```

检查返回的订单 ID、`orders` 中对应 `request_id` 的记录数量和商品库存。第二次请求没有创建新订单，也没有再次扣减库存。这验证了我们业务的幂等性。

----

## 验证数据持久化

订单已经创建之后，我们还需要确认：如果 MySQL 容器被删除并重建，订单是否仍然存在？

`mysql-data:/var/lib/mysql` 将数据库文件放在命名卷中，它的生命周期独立于 MySQL 容器。我们希望**只重建容器，而保留数据卷。**[Docker Volume 文档](https://docs.docker.com/engine/storage/volumes/)


在项目根目录的 PowerShell 中依次执行：

```powershell
$beforeId = docker compose --env-file .env ps -q mysql
docker compose --env-file .env stop api
docker compose --env-file .env stop mysql
docker compose --env-file .env rm -f mysql
```
先停 API，是为了避免数据库重建期间出现新的订单请求。`rm -f mysql` 删除的是已停止的 MySQL 容器。

接下来重新打开 MySQL 容器：
```powershell
docker compose --env-file .env up -d mysql
docker compose --env-file .env ps
$afterId = docker compose --env-file .env ps -q mysql
$beforeId -ne $afterId
```

`up -d mysql` 会创建新容器并重新挂载原来的命名卷。最后一行应输出 `True`，表明容器 ID 已变化。等待 `ps` 显示 MySQL `healthy`，再运行：

```powershell
docker compose --env-file .env run --rm migrate
```

如果数据库结构也保留了，Migration 应返回 `no change`。最后在 DBeaver 重复之前的查询，比较同一订单 ID、`request_id` 和商品库存。验收条件是：

1. 原订单仍在，ID 与幂等键没有变化。
2. 商品库存保持扣减后的数值，没有恢复到初始值。
3. Migration 没有重新创建表，版本仍是已应用的 Version 2。
4. MySQL 已由新容器承载，而数据仍来自原来的 `mysql-data` 卷。

实验完成后，可以用 `docker compose --env-file .env up -d api` 恢复 API。这里不能使用 `docker compose down -v`：`-v` 会删除 Compose 声明的命名卷，正好把我们要验证的数据一起删除。[Compose down 命令](https://docs.docker.com/reference/cli/docker/compose/down/)

下面是一些命令的总结：
| 命令 | Container | Volume |
|------|-----------|--------|
| `docker compose up -d` | 创建并启动新的容器 | 使用已存在的卷 |
| `docker compose down` | 停止并删除容器 | 保留卷 |
| `docker compose down -v` | 停止并删除容器 | 删除卷 |
| `docker compose rm ` | 删除容器 | 保留卷 |
| `docker compose stop` | 停止容器，但不删除 | 保留卷 |

