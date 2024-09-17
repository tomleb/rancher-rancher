package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterLifecycleConverter struct {
	lifecycle ClusterLifecycle
}

func (w *clusterLifecycleConverter) CreateContext(_ context.Context, obj *v3.Cluster) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Cluster) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Cluster) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterLifecycle interface {
	Create(obj *v3.Cluster) (runtime.Object, error)
	Remove(obj *v3.Cluster) (runtime.Object, error)
	Updated(obj *v3.Cluster) (runtime.Object, error)
}

type ClusterLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Cluster) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Cluster) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Cluster) (runtime.Object, error)
}

type clusterLifecycleAdapter struct {
	lifecycle ClusterLifecycleContext
}

func (w *clusterLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Cluster))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Cluster))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Cluster))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterLifecycleAdapter(name string, clusterScoped bool, client ClusterInterface, l ClusterLifecycle) ClusterHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterGroupVersionResource)
	}
	adapter := &clusterLifecycleAdapter{lifecycle: &clusterLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Cluster) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterLifecycleAdapterContext(name string, clusterScoped bool, client ClusterInterface, l ClusterLifecycleContext) ClusterHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterGroupVersionResource)
	}
	adapter := &clusterLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Cluster) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
