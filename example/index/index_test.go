package informer

import (
	"fmt"
	"forrest/codeCollection/k8s/code/controller/initclient"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"testing"
)

func TestIndex(t *testing.T) {
	client := initclient.ClientSet.Client

	// 创建 Pod 的 ListWatcher
	podListWatcher := cache.NewListWatchFromClient(
		client.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceAll, // 监控所有命名空间的Pod
		fields.Everything(),
	)

	// 创建一个 IndexerInformer，并为 `namespace` 和 `label` 设置索引
	indexer, informer := cache.NewIndexerInformer(
		podListWatcher,
		&v1.Pod{},
		0, // resyncPeriod，0 表示不定期重新同步
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				fmt.Printf("New Pod Added: %s/%s\n", pod.Namespace, pod.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				fmt.Printf("Pod Updated: %s/%s\n", pod.Namespace, pod.Name)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				fmt.Printf("Pod Deleted: %s/%s\n", pod.Namespace, pod.Name)
			},
		},
		// 为 indexer 定义两个索引
		cache.Indexers{
			"namespace": cache.MetaNamespaceIndexFunc,
			"label": func(obj interface{}) ([]string, error) {
				pod := obj.(*v1.Pod)
				if val, ok := pod.Labels["app"]; ok {
					return []string{val}, nil
				}
				return []string{}, nil
			},
		},
	)

	// 启动 informer 以开始监听事件
	stopCh := make(chan struct{})
	defer close(stopCh)
	go informer.Run(stopCh)

	// 等待缓存同步完成
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		klog.Fatalf("Timed out waiting for caches to sync")
	}

	// 查询所有在 `default` 命名空间下的 Pod
	namespacePods, err := indexer.ByIndex("namespace", "default")
	if err != nil {
		klog.Fatalf("Error retrieving pods by namespace index: %v", err)
	}
	fmt.Printf("Pods in default namespace:\n")
	for _, obj := range namespacePods {
		pod := obj.(*v1.Pod)
		fmt.Printf("  - %s/%s\n", pod.Namespace, pod.Name)
	}

	// 查询所有带有 `app=example` 标签的 Pod
	labelPods, err := indexer.ByIndex("label", "example")
	if err != nil {
		klog.Fatalf("Error retrieving pods by label index: %v", err)
	}
	fmt.Printf("Pods with label app=example:\n")
	for _, obj := range labelPods {
		pod := obj.(*v1.Pod)
		fmt.Printf("  - %s/%s\n", pod.Namespace, pod.Name)
	}

	select {}
}
