## 利用kind搭建k8s一主三从节点 step by step

### step1 前置条件
- kind安装
- kubectl 安装
- docker 安装
```bash
brew install kind
brew install kubectl
#docker可以去官网下载
```

### step2 集群搭建
一个control plane，三个worker,配置文件可以参考 [kind 配置文件](https://kind.sigs.k8s.io/docs/user/configuration/)
```bash
➜  kind git:(main) ✗ kind create cluster --name k8s-cluster-with3worker --config ./config.yaml
Creating cluster "k8s-cluster-with3worker" ...
 ✓ Ensuring node image (kindest/node:v1.31.0) 🖼
 ✓ Preparing nodes 📦 📦 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-k8s-cluster-with3worker"
You can now use your cluster with:

kubectl cluster-info --context kind-k8s-cluster-with3worker

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

