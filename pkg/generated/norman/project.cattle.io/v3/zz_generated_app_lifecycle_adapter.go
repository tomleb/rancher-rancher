package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type appLifecycleConverter struct {
	lifecycle AppLifecycle
}

func (w *appLifecycleConverter) CreateContext(_ context.Context, obj *v3.App) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *appLifecycleConverter) RemoveContext(_ context.Context, obj *v3.App) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *appLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.App) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type AppLifecycle interface {
	Create(obj *v3.App) (runtime.Object, error)
	Remove(obj *v3.App) (runtime.Object, error)
	Updated(obj *v3.App) (runtime.Object, error)
}

type AppLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.App) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.App) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.App) (runtime.Object, error)
}

type appLifecycleAdapter struct {
	lifecycle AppLifecycleContext
}

func (w *appLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *appLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *appLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *appLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.App))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *appLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *appLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.App))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *appLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *appLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.App))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewAppLifecycleAdapter(name string, clusterScoped bool, client AppInterface, l AppLifecycle) AppHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(AppGroupVersionResource)
	}
	adapter := &appLifecycleAdapter{lifecycle: &appLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.App) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewAppLifecycleAdapterContext(name string, clusterScoped bool, client AppInterface, l AppLifecycleContext) AppHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(AppGroupVersionResource)
	}
	adapter := &appLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.App) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
