package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterCatalogLifecycleConverter struct {
	lifecycle ClusterCatalogLifecycle
}

func (w *clusterCatalogLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterCatalog) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterCatalogLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterCatalog) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterCatalogLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterCatalog) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterCatalogLifecycle interface {
	Create(obj *v3.ClusterCatalog) (runtime.Object, error)
	Remove(obj *v3.ClusterCatalog) (runtime.Object, error)
	Updated(obj *v3.ClusterCatalog) (runtime.Object, error)
}

type ClusterCatalogLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterCatalog) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterCatalog) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterCatalog) (runtime.Object, error)
}

type clusterCatalogLifecycleAdapter struct {
	lifecycle ClusterCatalogLifecycleContext
}

func (w *clusterCatalogLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterCatalogLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterCatalogLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterCatalogLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterCatalogLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterCatalogLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterCatalogLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterCatalogLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterCatalogLifecycleAdapter(name string, clusterScoped bool, client ClusterCatalogInterface, l ClusterCatalogLifecycle) ClusterCatalogHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterCatalogGroupVersionResource)
	}
	adapter := &clusterCatalogLifecycleAdapter{lifecycle: &clusterCatalogLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterCatalog) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterCatalogLifecycleAdapterContext(name string, clusterScoped bool, client ClusterCatalogInterface, l ClusterCatalogLifecycleContext) ClusterCatalogHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterCatalogGroupVersionResource)
	}
	adapter := &clusterCatalogLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterCatalog) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
