package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterTemplateRevisionLifecycleConverter struct {
	lifecycle ClusterTemplateRevisionLifecycle
}

func (w *clusterTemplateRevisionLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterTemplateRevisionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterTemplateRevisionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterTemplateRevisionLifecycle interface {
	Create(obj *v3.ClusterTemplateRevision) (runtime.Object, error)
	Remove(obj *v3.ClusterTemplateRevision) (runtime.Object, error)
	Updated(obj *v3.ClusterTemplateRevision) (runtime.Object, error)
}

type ClusterTemplateRevisionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterTemplateRevision) (runtime.Object, error)
}

type clusterTemplateRevisionLifecycleAdapter struct {
	lifecycle ClusterTemplateRevisionLifecycleContext
}

func (w *clusterTemplateRevisionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterTemplateRevisionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterTemplateRevisionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterTemplateRevisionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterTemplateRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterTemplateRevisionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterTemplateRevisionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterTemplateRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterTemplateRevisionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterTemplateRevisionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterTemplateRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterTemplateRevisionLifecycleAdapter(name string, clusterScoped bool, client ClusterTemplateRevisionInterface, l ClusterTemplateRevisionLifecycle) ClusterTemplateRevisionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterTemplateRevisionGroupVersionResource)
	}
	adapter := &clusterTemplateRevisionLifecycleAdapter{lifecycle: &clusterTemplateRevisionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterTemplateRevision) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterTemplateRevisionLifecycleAdapterContext(name string, clusterScoped bool, client ClusterTemplateRevisionInterface, l ClusterTemplateRevisionLifecycleContext) ClusterTemplateRevisionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterTemplateRevisionGroupVersionResource)
	}
	adapter := &clusterTemplateRevisionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterTemplateRevision) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
