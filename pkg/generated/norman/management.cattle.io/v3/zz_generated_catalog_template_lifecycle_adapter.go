package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type catalogTemplateLifecycleConverter struct {
	lifecycle CatalogTemplateLifecycle
}

func (w *catalogTemplateLifecycleConverter) CreateContext(_ context.Context, obj *v3.CatalogTemplate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *catalogTemplateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.CatalogTemplate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *catalogTemplateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.CatalogTemplate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CatalogTemplateLifecycle interface {
	Create(obj *v3.CatalogTemplate) (runtime.Object, error)
	Remove(obj *v3.CatalogTemplate) (runtime.Object, error)
	Updated(obj *v3.CatalogTemplate) (runtime.Object, error)
}

type CatalogTemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.CatalogTemplate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.CatalogTemplate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.CatalogTemplate) (runtime.Object, error)
}

type catalogTemplateLifecycleAdapter struct {
	lifecycle CatalogTemplateLifecycleContext
}

func (w *catalogTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *catalogTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *catalogTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *catalogTemplateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.CatalogTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *catalogTemplateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.CatalogTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *catalogTemplateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.CatalogTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCatalogTemplateLifecycleAdapter(name string, clusterScoped bool, client CatalogTemplateInterface, l CatalogTemplateLifecycle) CatalogTemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogTemplateGroupVersionResource)
	}
	adapter := &catalogTemplateLifecycleAdapter{lifecycle: &catalogTemplateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.CatalogTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCatalogTemplateLifecycleAdapterContext(name string, clusterScoped bool, client CatalogTemplateInterface, l CatalogTemplateLifecycleContext) CatalogTemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogTemplateGroupVersionResource)
	}
	adapter := &catalogTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.CatalogTemplate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
