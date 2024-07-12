# nfctl

# 通过 nfctl 快速开启后台项目

### 1. Installation

- ① `go install github.com/loveuer/nf/nft/nfctl@latest`
- ② download prebuild binary [release](https://github.com/loveuer/nf/releases)

### 2. Usage

- `nfctl new {project}`
- `nfctl new project -t ultone`
- `nfctl new project -t https://github.com/xxx/yyy.git`
- `nfctl new project --template https://gitcode/loveuer/ultone.git`
- `nfctl new project --template https://{username}:{password/token}@my.gitlab.com/name/project.git`

### 3. nfctl init script

- `为方便模版的初始化, 可以采用 nfctl init script, 当 nfctl new project -t xxx 从模版开始项目时会自动执行`
- `具体的编写规则如下:`
  * [init 脚本规则](https://github.com/loveuer/nf/nft/nfctl/script.md) 或者
  * [国内](https://gitcode.com/loveuer/nf/nft/nfctl/script.md)