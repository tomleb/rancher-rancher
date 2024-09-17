package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type resourceQuotaLifecycleConverter struct {
	lifecycle ResourceQuotaLifecycle
}

func (w *resourceQuotaLifecycleConverter) CreateContext(_ context.Context, obj *v1.ResourceQuota) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *resourceQuotaLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ResourceQuota) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *resourceQuotaLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ResourceQuota) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ResourceQuotaLifecycle interface {
	Create(obj *v1.ResourceQuota) (runtime.Object, error)
	Remove(obj *v1.ResourceQuota) (runtime.Object, error)
	Updated(obj *v1.ResourceQuota) (runtime.Object, error)
}

type ResourceQuotaLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ResourceQuota) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ResourceQuota) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ResourceQuota) (runtime.Object, error)
}

type resourceQuotaLifecycleAdapter struct {
	lifecycle ResourceQuotaLifecycleContext
}

func (w *resourceQuotaLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *resourceQuotaLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *resourceQuotaLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *resourceQuotaLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ResourceQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *resourceQuotaLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *resourceQuotaLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ResourceQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *resourceQuotaLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *resourceQuotaLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ResourceQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewResourceQuotaLifecycleAdapter(name string, clusterScoped bool, client ResourceQuotaInterface, l ResourceQuotaLifecycle) ResourceQuotaHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ResourceQuotaGroupVersionResource)
	}
	adapter := &resourceQuotaLifecycleAdapter{lifecycle: &resourceQuotaLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ResourceQuota) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewResourceQuotaLifecycleAdapterContext(name string, clusterScoped bool, client ResourceQuotaInterface, l ResourceQuotaLifecycleContext) ResourceQuotaHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ResourceQuotaGroupVersionResource)
	}
	adapter := &resourceQuotaLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ResourceQuota) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
