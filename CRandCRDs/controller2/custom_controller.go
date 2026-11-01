package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	group     = "cad.com"
	version   = "v1"
	resource  = "customconfigmaps"
	namespace = "default"
)

// createConfigMap creates a ConfigMap from the custom resource's spec data.
func createConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string, data map[string]string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
	return err
}

// deleteConfigMap deletes the ConfigMap associated with the custom resource.
func deleteConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
	return clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// specToStringMap converts the unstructured spec map to map[string]string for ConfigMap data.
func specToStringMap(spec map[string]interface{}) map[string]string {
	data := make(map[string]string)
	for k, v := range spec {
		if s, ok := v.(string); ok {
			data[k] = s
		}
	}
	return data
}

func main() {
	// Use in-cluster configuration (equivalent to config.load_incluster_config())
	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("failed to load in-cluster config: %v", err)
	}

	// Typed client for core resources (ConfigMaps)
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("failed to create kubernetes clientset: %v", err)
	}

	// Dynamic client for custom resources (equivalent to client.CustomObjectsApi())
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("failed to create dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	ctx := context.Background()

	log.Printf("starting controller: watching %s.%s/%s in namespace %s", resource, group, version, namespace)

	// Outer loop restarts the watch if the connection drops (mirrors the while True in Python).
	for {
		watcher, err := dynamicClient.Resource(gvr).Namespace(namespace).Watch(ctx, metav1.ListOptions{})
		if err != nil {
			log.Fatalf("failed to start watch on %s: %v", resource, err)
		}

		// Drain events from the watch channel (equivalent to the for event in stream loop).
		for event := range watcher.ResultChan() {
			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				continue
			}

			name := obj.GetName()
			spec, _, _ := unstructured.NestedMap(obj.Object, "spec")
			data := specToStringMap(spec)

			switch event.Type {
			case watch.Added:
				if err := createConfigMap(ctx, clientset, namespace, name, data); err != nil {
					log.Printf("ERROR creating configmap %q: %v", name, err)
				} else {
					fmt.Printf("ADDED  → created configmap %q\n", name)
				}

			case watch.Deleted:
				if err := deleteConfigMap(ctx, clientset, namespace, name); err != nil {
					log.Printf("ERROR deleting configmap %q: %v", name, err)
				} else {
					fmt.Printf("DELETED → removed configmap %q\n", name)
				}
			}
		}

		log.Println("watch channel closed, restarting...")
	}
}
