package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterTemplateLifecycleConverter struct {
	lifecycle ClusterTemplateLifecycle
}

func (w *clusterTemplateLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterTemplate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterTemplateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterTemplate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterTemplateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterTemplate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterTemplateLifecycle interface {
	Create(obj *v3.ClusterTemplate) (runtime.Object, error)
	Remove(obj *v3.ClusterTemplate) (runtime.Object, error)
	Updated(obj *v3.ClusterTemplate) (runtime.Object, error)
}

type ClusterTemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterTemplate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterTemplate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterTemplate) (runtime.Object, error)
}

type clusterTemplateLifecycleAdapter struct {
	lifecycle ClusterTemplateLifecycleContext
}

func (w *clusterTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterTemplateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterTemplateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterTemplateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterTemplateLifecycleAdapter(name string, clusterScoped bool, client ClusterTemplateInterface, l ClusterTemplateLifecycle) ClusterTemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterTemplateGroupVersionResource)
	}
	adapter := &clusterTemplateLifecycleAdapter{lifecycle: &clusterTemplateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterTemplateLifecycleAdapterContext(name string, clusterScoped bool, client ClusterTemplateInterface, l ClusterTemplateLifecycleContext) ClusterTemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterTemplateGroupVersionResource)
	}
	adapter := &clusterTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterTemplate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
