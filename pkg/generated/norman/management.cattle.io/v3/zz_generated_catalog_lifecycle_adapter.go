package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type catalogLifecycleConverter struct {
	lifecycle CatalogLifecycle
}

func (w *catalogLifecycleConverter) CreateContext(_ context.Context, obj *v3.Catalog) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *catalogLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Catalog) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *catalogLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Catalog) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CatalogLifecycle interface {
	Create(obj *v3.Catalog) (runtime.Object, error)
	Remove(obj *v3.Catalog) (runtime.Object, error)
	Updated(obj *v3.Catalog) (runtime.Object, error)
}

type CatalogLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Catalog) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Catalog) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Catalog) (runtime.Object, error)
}

type catalogLifecycleAdapter struct {
	lifecycle CatalogLifecycleContext
}

func (w *catalogLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *catalogLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *catalogLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *catalogLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Catalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *catalogLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Catalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *catalogLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Catalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCatalogLifecycleAdapter(name string, clusterScoped bool, client CatalogInterface, l CatalogLifecycle) CatalogHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogGroupVersionResource)
	}
	adapter := &catalogLifecycleAdapter{lifecycle: &catalogLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Catalog) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCatalogLifecycleAdapterContext(name string, clusterScoped bool, client CatalogInterface, l CatalogLifecycleContext) CatalogHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogGroupVersionResource)
	}
	adapter := &catalogLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Catalog) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
