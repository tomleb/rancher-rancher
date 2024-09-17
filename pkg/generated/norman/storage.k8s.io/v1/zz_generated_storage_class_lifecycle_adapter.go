package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type storageClassLifecycleConverter struct {
	lifecycle StorageClassLifecycle
}

func (w *storageClassLifecycleConverter) CreateContext(_ context.Context, obj *v1.StorageClass) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *storageClassLifecycleConverter) RemoveContext(_ context.Context, obj *v1.StorageClass) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *storageClassLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.StorageClass) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type StorageClassLifecycle interface {
	Create(obj *v1.StorageClass) (runtime.Object, error)
	Remove(obj *v1.StorageClass) (runtime.Object, error)
	Updated(obj *v1.StorageClass) (runtime.Object, error)
}

type StorageClassLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.StorageClass) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.StorageClass) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.StorageClass) (runtime.Object, error)
}

type storageClassLifecycleAdapter struct {
	lifecycle StorageClassLifecycleContext
}

func (w *storageClassLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *storageClassLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *storageClassLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *storageClassLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.StorageClass))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *storageClassLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *storageClassLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.StorageClass))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *storageClassLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *storageClassLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.StorageClass))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewStorageClassLifecycleAdapter(name string, clusterScoped bool, client StorageClassInterface, l StorageClassLifecycle) StorageClassHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(StorageClassGroupVersionResource)
	}
	adapter := &storageClassLifecycleAdapter{lifecycle: &storageClassLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.StorageClass) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewStorageClassLifecycleAdapterContext(name string, clusterScoped bool, client StorageClassInterface, l StorageClassLifecycleContext) StorageClassHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(StorageClassGroupVersionResource)
	}
	adapter := &storageClassLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.StorageClass) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
