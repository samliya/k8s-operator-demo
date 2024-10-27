

## introduction
这是一个简易的controller。
当用户expose pod或 deployment时，即创建了svc，
且这个svc带有`ingress/http` 这样的annotation字段时，自动给svc创建ingress。
其中ingress-controller可以参考[ngix-controller-deploy](https://kubernetes.github.io/ingress-nginx/deploy/#local-development-clusters)

## TODO
将controller部署到k8s集群中
1. role 
2. roleBinding 
3. serviceAccount 
4. docker image 
5. deploy the controller