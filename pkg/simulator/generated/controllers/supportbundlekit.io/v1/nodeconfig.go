/*
Copyright 2022 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/support-bundle-kit/pkg/simulator/apis/supportbundlekit.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type NodeConfigHandler func(string, *v1.NodeConfig) (*v1.NodeConfig, error)

type NodeConfigController interface {
	generic.ControllerMeta
	NodeConfigClient

	OnChange(ctx context.Context, name string, sync NodeConfigHandler)
	OnRemove(ctx context.Context, name string, sync NodeConfigHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() NodeConfigCache
}

type NodeConfigClient interface {
	Create(*v1.NodeConfig) (*v1.NodeConfig, error)
	Update(*v1.NodeConfig) (*v1.NodeConfig, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.NodeConfig, error)
	List(namespace string, opts metav1.ListOptions) (*v1.NodeConfigList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.NodeConfig, err error)
}

type NodeConfigCache interface {
	Get(namespace, name string) (*v1.NodeConfig, error)
	List(namespace string, selector labels.Selector) ([]*v1.NodeConfig, error)

	AddIndexer(indexName string, indexer NodeConfigIndexer)
	GetByIndex(indexName, key string) ([]*v1.NodeConfig, error)
}

type NodeConfigIndexer func(obj *v1.NodeConfig) ([]string, error)

type nodeConfigController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewNodeConfigController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) NodeConfigController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &nodeConfigController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromNodeConfigHandlerToHandler(sync NodeConfigHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.NodeConfig
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.NodeConfig))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *nodeConfigController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.NodeConfig))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateNodeConfigDeepCopyOnChange(client NodeConfigClient, obj *v1.NodeConfig, handler func(obj *v1.NodeConfig) (*v1.NodeConfig, error)) (*v1.NodeConfig, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *nodeConfigController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *nodeConfigController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *nodeConfigController) OnChange(ctx context.Context, name string, sync NodeConfigHandler) {
	c.AddGenericHandler(ctx, name, FromNodeConfigHandlerToHandler(sync))
}

func (c *nodeConfigController) OnRemove(ctx context.Context, name string, sync NodeConfigHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromNodeConfigHandlerToHandler(sync)))
}

func (c *nodeConfigController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *nodeConfigController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *nodeConfigController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *nodeConfigController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *nodeConfigController) Cache() NodeConfigCache {
	return &nodeConfigCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *nodeConfigController) Create(obj *v1.NodeConfig) (*v1.NodeConfig, error) {
	result := &v1.NodeConfig{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *nodeConfigController) Update(obj *v1.NodeConfig) (*v1.NodeConfig, error) {
	result := &v1.NodeConfig{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *nodeConfigController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *nodeConfigController) Get(namespace, name string, options metav1.GetOptions) (*v1.NodeConfig, error) {
	result := &v1.NodeConfig{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *nodeConfigController) List(namespace string, opts metav1.ListOptions) (*v1.NodeConfigList, error) {
	result := &v1.NodeConfigList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *nodeConfigController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *nodeConfigController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.NodeConfig, error) {
	result := &v1.NodeConfig{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type nodeConfigCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *nodeConfigCache) Get(namespace, name string) (*v1.NodeConfig, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.NodeConfig), nil
}

func (c *nodeConfigCache) List(namespace string, selector labels.Selector) (ret []*v1.NodeConfig, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NodeConfig))
	})

	return ret, err
}

func (c *nodeConfigCache) AddIndexer(indexName string, indexer NodeConfigIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.NodeConfig))
		},
	}))
}

func (c *nodeConfigCache) GetByIndex(indexName, key string) (result []*v1.NodeConfig, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.NodeConfig, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.NodeConfig))
	}
	return result, nil
}
