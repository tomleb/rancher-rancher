package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type projectCatalogLifecycleConverter struct {
	lifecycle ProjectCatalogLifecycle
}

func (w *projectCatalogLifecycleConverter) CreateContext(_ context.Context, obj *v3.ProjectCatalog) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *projectCatalogLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ProjectCatalog) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *projectCatalogLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ProjectCatalog) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ProjectCatalogLifecycle interface {
	Create(obj *v3.ProjectCatalog) (runtime.Object, error)
	Remove(obj *v3.ProjectCatalog) (runtime.Object, error)
	Updated(obj *v3.ProjectCatalog) (runtime.Object, error)
}

type ProjectCatalogLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ProjectCatalog) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ProjectCatalog) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ProjectCatalog) (runtime.Object, error)
}

type projectCatalogLifecycleAdapter struct {
	lifecycle ProjectCatalogLifecycleContext
}

func (w *projectCatalogLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *projectCatalogLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *projectCatalogLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *projectCatalogLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ProjectCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectCatalogLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *projectCatalogLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ProjectCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectCatalogLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *projectCatalogLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ProjectCatalog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewProjectCatalogLifecycleAdapter(name string, clusterScoped bool, client ProjectCatalogInterface, l ProjectCatalogLifecycle) ProjectCatalogHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectCatalogGroupVersionResource)
	}
	adapter := &projectCatalogLifecycleAdapter{lifecycle: &projectCatalogLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ProjectCatalog) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewProjectCatalogLifecycleAdapterContext(name string, clusterScoped bool, client ProjectCatalogInterface, l ProjectCatalogLifecycleContext) ProjectCatalogHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ProjectCatalogGroupVersionResource)
	}
	adapter := &projectCatalogLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ProjectCatalog) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
