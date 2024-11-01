## introduction
这是一个简易的controller。
当用户expose pod或 deployment时，即创建了svc，
且这个svc带有`ingress/http` 这样的annotation字段时，自动给svc创建ingress。
其中ingress-controller可以参考[ngix-controller-deploy](https://kubernetes.github.io/ingress-nginx/deploy/#local-development-clusters)



## 部署step by step
### 1、RBAC设置
```bash
#RBAC
k create ns ing-test --dry-run=client -oyaml >./manifest/namespace.yaml
k create serviceaccount ing-sa --dry-run=client -oyaml > ./manifest/sa.yaml  
k create clusterrole ing-clusterrole --verb get,watch,list,create,update,delete --resource ingress --verb watch,list --resource service --dry-run=client  -oyaml > ./manifest/clusterrole.yaml
k create clusterrolebinding ing-cluster-rolebinding --clusterrole ing-clusterrole --serviceaccount default:ing-sa --dry-run=client -oyaml > ./manifest/clusterrolebinding.yaml
```

### 2、controller构建镜像，书写controller-deploy文件
```bash
#将controller打包为docker image
docker build -t controller-demo:v1 .
#ing-controller-controller.yaml
 k create deploy ing-controller -ning-test --image controller-demo:v1 --replicas 1  --dry-run=client -oyaml >./manifest/ing-controller-deploy.yaml
#将本地docker image load进 kind集群
➜  controller-demo git:(main) ✗ kind load docker-image controller-demo:v1           
Image: "controller-demo:v1" with ID "sha256:2026051fbe3b3127c547526602ce1bec1f1ff61e6c27056a9086a8b4333c8d5a" not yet present on node "kind-control-plane", loading...
Image: "controller-demo:v1" with ID "sha256:2026051fbe3b3127c547526602ce1bec1f1ff61e6c27056a9086a8b4333c8d5a" not yet present on node "kind-worker", loading...
Image: "controller-demo:v1" with ID "sha256:2026051fbe3b3127c547526602ce1bec1f1ff61e6c27056a9086a8b4333c8d5a" not yet present on node "kind-worker3", loading...
Image: "controller-demo:v1" with ID "sha256:2026051fbe3b3127c547526602ce1bec1f1ff61e6c27056a9086a8b4333c8d5a" not yet present on node "kind-worker2", loading...
```

### 3、部署
```bash
#部署controller & RBAC
➜  controller-demo git:(main) ✗ k apply -f ./manifest 
clusterrole.rbac.authorization.k8s.io/ing-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/ing-cluster-rolebinding created
namespace/ing-test created
serviceaccount/ing-sa created
Error from server (NotFound): error when creating "manifest/ing-controller-deploy.yaml": namespaces "ing-test" not found
➜  controller-demo git:(main) ✗ k apply -f ./manifest
clusterrole.rbac.authorization.k8s.io/ing-clusterrole unchanged
clusterrolebinding.rbac.authorization.k8s.io/ing-cluster-rolebinding unchanged
deployment.apps/ing-controller created
namespace/ing-test unchanged
serviceaccount/ing-sa unchanged
```

### 4、测试
```bash
#1. 创建一个pod，再expose他，观察ing

#2. annotate这个svc，ingress/http: "true"

#3. 再观察ing是否被创建

#4. 删除这个svc的ingress label

#5.观察ing是否成功被删除
```

## 踩坑记录
1. config配置文件部署到容器中需要用这种方式读取，读取容器内的`/var/run/secrets/kubernetes.io/serviceaccount`，挂载的流程由k8s完成，将必要的凭证和API地址注入到pod中
```bash
var cfg *rest.Config
var err error
cfg, err = rest.InClusterConfig() 
2. deployment没有设置serviceAccount，默认使用的default sa，权限认证过不去
```