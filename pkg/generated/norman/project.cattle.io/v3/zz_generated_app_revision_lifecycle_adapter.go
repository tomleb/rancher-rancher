package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type appRevisionLifecycleConverter struct {
	lifecycle AppRevisionLifecycle
}

func (w *appRevisionLifecycleConverter) CreateContext(_ context.Context, obj *v3.AppRevision) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *appRevisionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.AppRevision) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *appRevisionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.AppRevision) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type AppRevisionLifecycle interface {
	Create(obj *v3.AppRevision) (runtime.Object, error)
	Remove(obj *v3.AppRevision) (runtime.Object, error)
	Updated(obj *v3.AppRevision) (runtime.Object, error)
}

type AppRevisionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.AppRevision) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.AppRevision) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.AppRevision) (runtime.Object, error)
}

type appRevisionLifecycleAdapter struct {
	lifecycle AppRevisionLifecycleContext
}

func (w *appRevisionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *appRevisionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *appRevisionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *appRevisionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.AppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *appRevisionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *appRevisionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.AppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *appRevisionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *appRevisionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.AppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewAppRevisionLifecycleAdapter(name string, clusterScoped bool, client AppRevisionInterface, l AppRevisionLifecycle) AppRevisionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(AppRevisionGroupVersionResource)
	}
	adapter := &appRevisionLifecycleAdapter{lifecycle: &appRevisionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.AppRevision) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewAppRevisionLifecycleAdapterContext(name string, clusterScoped bool, client AppRevisionInterface, l AppRevisionLifecycleContext) AppRevisionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(AppRevisionGroupVersionResource)
	}
	adapter := &appRevisionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.AppRevision) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
