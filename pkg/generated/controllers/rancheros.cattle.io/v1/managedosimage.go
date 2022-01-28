/*
Copyright 2021 Rancher Labs, Inc.

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

	v1 "github.com/rancher-sandbox/os2/pkg/apis/rancheros.cattle.io/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
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

type ManagedOSImageHandler func(string, *v1.ManagedOSImage) (*v1.ManagedOSImage, error)

type ManagedOSImageController interface {
	generic.ControllerMeta
	ManagedOSImageClient

	OnChange(ctx context.Context, name string, sync ManagedOSImageHandler)
	OnRemove(ctx context.Context, name string, sync ManagedOSImageHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ManagedOSImageCache
}

type ManagedOSImageClient interface {
	Create(*v1.ManagedOSImage) (*v1.ManagedOSImage, error)
	Update(*v1.ManagedOSImage) (*v1.ManagedOSImage, error)
	UpdateStatus(*v1.ManagedOSImage) (*v1.ManagedOSImage, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.ManagedOSImage, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ManagedOSImageList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ManagedOSImage, err error)
}

type ManagedOSImageCache interface {
	Get(namespace, name string) (*v1.ManagedOSImage, error)
	List(namespace string, selector labels.Selector) ([]*v1.ManagedOSImage, error)

	AddIndexer(indexName string, indexer ManagedOSImageIndexer)
	GetByIndex(indexName, key string) ([]*v1.ManagedOSImage, error)
}

type ManagedOSImageIndexer func(obj *v1.ManagedOSImage) ([]string, error)

type managedOSImageController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewManagedOSImageController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ManagedOSImageController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &managedOSImageController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromManagedOSImageHandlerToHandler(sync ManagedOSImageHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ManagedOSImage
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ManagedOSImage))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *managedOSImageController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ManagedOSImage))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateManagedOSImageDeepCopyOnChange(client ManagedOSImageClient, obj *v1.ManagedOSImage, handler func(obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error)) (*v1.ManagedOSImage, error) {
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

func (c *managedOSImageController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *managedOSImageController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *managedOSImageController) OnChange(ctx context.Context, name string, sync ManagedOSImageHandler) {
	c.AddGenericHandler(ctx, name, FromManagedOSImageHandlerToHandler(sync))
}

func (c *managedOSImageController) OnRemove(ctx context.Context, name string, sync ManagedOSImageHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromManagedOSImageHandlerToHandler(sync)))
}

func (c *managedOSImageController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *managedOSImageController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *managedOSImageController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *managedOSImageController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *managedOSImageController) Cache() ManagedOSImageCache {
	return &managedOSImageCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *managedOSImageController) Create(obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error) {
	result := &v1.ManagedOSImage{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *managedOSImageController) Update(obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error) {
	result := &v1.ManagedOSImage{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *managedOSImageController) UpdateStatus(obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error) {
	result := &v1.ManagedOSImage{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *managedOSImageController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *managedOSImageController) Get(namespace, name string, options metav1.GetOptions) (*v1.ManagedOSImage, error) {
	result := &v1.ManagedOSImage{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *managedOSImageController) List(namespace string, opts metav1.ListOptions) (*v1.ManagedOSImageList, error) {
	result := &v1.ManagedOSImageList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *managedOSImageController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *managedOSImageController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.ManagedOSImage, error) {
	result := &v1.ManagedOSImage{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type managedOSImageCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *managedOSImageCache) Get(namespace, name string) (*v1.ManagedOSImage, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.ManagedOSImage), nil
}

func (c *managedOSImageCache) List(namespace string, selector labels.Selector) (ret []*v1.ManagedOSImage, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ManagedOSImage))
	})

	return ret, err
}

func (c *managedOSImageCache) AddIndexer(indexName string, indexer ManagedOSImageIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ManagedOSImage))
		},
	}))
}

func (c *managedOSImageCache) GetByIndex(indexName, key string) (result []*v1.ManagedOSImage, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.ManagedOSImage, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.ManagedOSImage))
	}
	return result, nil
}

type ManagedOSImageStatusHandler func(obj *v1.ManagedOSImage, status v1.ManagedOSImageStatus) (v1.ManagedOSImageStatus, error)

type ManagedOSImageGeneratingHandler func(obj *v1.ManagedOSImage, status v1.ManagedOSImageStatus) ([]runtime.Object, v1.ManagedOSImageStatus, error)

func RegisterManagedOSImageStatusHandler(ctx context.Context, controller ManagedOSImageController, condition condition.Cond, name string, handler ManagedOSImageStatusHandler) {
	statusHandler := &managedOSImageStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromManagedOSImageHandlerToHandler(statusHandler.sync))
}

func RegisterManagedOSImageGeneratingHandler(ctx context.Context, controller ManagedOSImageController, apply apply.Apply,
	condition condition.Cond, name string, handler ManagedOSImageGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &managedOSImageGeneratingHandler{
		ManagedOSImageGeneratingHandler: handler,
		apply:                           apply,
		name:                            name,
		gvk:                             controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterManagedOSImageStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type managedOSImageStatusHandler struct {
	client    ManagedOSImageClient
	condition condition.Cond
	handler   ManagedOSImageStatusHandler
}

func (a *managedOSImageStatusHandler) sync(key string, obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type managedOSImageGeneratingHandler struct {
	ManagedOSImageGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *managedOSImageGeneratingHandler) Remove(key string, obj *v1.ManagedOSImage) (*v1.ManagedOSImage, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.ManagedOSImage{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *managedOSImageGeneratingHandler) Handle(obj *v1.ManagedOSImage, status v1.ManagedOSImageStatus) (v1.ManagedOSImageStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ManagedOSImageGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
