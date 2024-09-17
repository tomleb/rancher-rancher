package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type catalogTemplateVersionLifecycleConverter struct {
	lifecycle CatalogTemplateVersionLifecycle
}

func (w *catalogTemplateVersionLifecycleConverter) CreateContext(_ context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *catalogTemplateVersionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *catalogTemplateVersionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type CatalogTemplateVersionLifecycle interface {
	Create(obj *v3.CatalogTemplateVersion) (runtime.Object, error)
	Remove(obj *v3.CatalogTemplateVersion) (runtime.Object, error)
	Updated(obj *v3.CatalogTemplateVersion) (runtime.Object, error)
}

type CatalogTemplateVersionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.CatalogTemplateVersion) (runtime.Object, error)
}

type catalogTemplateVersionLifecycleAdapter struct {
	lifecycle CatalogTemplateVersionLifecycleContext
}

func (w *catalogTemplateVersionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *catalogTemplateVersionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *catalogTemplateVersionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *catalogTemplateVersionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.CatalogTemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogTemplateVersionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *catalogTemplateVersionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.CatalogTemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *catalogTemplateVersionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *catalogTemplateVersionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.CatalogTemplateVersion))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewCatalogTemplateVersionLifecycleAdapter(name string, clusterScoped bool, client CatalogTemplateVersionInterface, l CatalogTemplateVersionLifecycle) CatalogTemplateVersionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogTemplateVersionGroupVersionResource)
	}
	adapter := &catalogTemplateVersionLifecycleAdapter{lifecycle: &catalogTemplateVersionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.CatalogTemplateVersion) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewCatalogTemplateVersionLifecycleAdapterContext(name string, clusterScoped bool, client CatalogTemplateVersionInterface, l CatalogTemplateVersionLifecycleContext) CatalogTemplateVersionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(CatalogTemplateVersionGroupVersionResource)
	}
	adapter := &catalogTemplateVersionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.CatalogTemplateVersion) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
