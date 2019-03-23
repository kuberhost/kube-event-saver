package main

import (
	"flag"
	_ "fmt"

	_ "github.com/kr/pretty"
	"github.com/paxa/kube-event-saver/pkg"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func watchEvents(eventQueue chan pkg.IncomingEvent, client kubernetes.Interface, stopCh chan struct{}) {
	eventListWatcher := cache.NewListWatchFromClient(
		client.CoreV1().RESTClient(), "events", v1.NamespaceAll, fields.Everything())

	_, informer := cache.NewIndexerInformer(eventListWatcher, &v1.Event{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Printf("AddFunc %# v\n", pretty.Formatter(obj))
			event, ok := obj.(*v1.Event)
			if ok {
				eventQueue <- pkg.IncomingEvent{Action: "add", Event: event}
			} else {
				klog.Error("AddFunc: not an event %# v\n", obj)
			}
		},
		UpdateFunc: func(old interface{}, obj interface{}) {
			//fmt.Printf("UpdateFunc %# v\n", pretty.Formatter(obj))
			event, ok := obj.(*v1.Event)
			if ok {
				eventQueue <- pkg.IncomingEvent{Action: "update", Event: event}
			} else {
				klog.Error("UpdateFunc: not an event %# v\n", obj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			//fmt.Printf("DeleteFunc %# v\n", pretty.Formatter(obj))
		},
	}, cache.Indexers{})

	informer.Run(stopCh)
}

func saveEvents(eventQueue chan pkg.IncomingEvent) {
	for {
		incEvent, ok := <-eventQueue
		if !ok {
			break
		}
		pkg.SaveEvent(incEvent)
	}
}

func main() {
	klog.InitFlags(nil)

	flag.Set("logtostderr", "true")
	//flag.Set("v", "10")
	flag.Parse()

	client, err := pkg.GetClient()
	if err != nil {
		panic(err.Error())
	}

	eventQueue := make(chan pkg.IncomingEvent, 100)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go watchEvents(eventQueue, client, stopCh)
	go saveEvents(eventQueue)

	select {}
}
