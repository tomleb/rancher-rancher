package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type globalRoleLifecycleConverter struct {
	lifecycle GlobalRoleLifecycle
}

func (w *globalRoleLifecycleConverter) CreateContext(_ context.Context, obj *v3.GlobalRole) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *globalRoleLifecycleConverter) RemoveContext(_ context.Context, obj *v3.GlobalRole) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *globalRoleLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.GlobalRole) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GlobalRoleLifecycle interface {
	Create(obj *v3.GlobalRole) (runtime.Object, error)
	Remove(obj *v3.GlobalRole) (runtime.Object, error)
	Updated(obj *v3.GlobalRole) (runtime.Object, error)
}

type GlobalRoleLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.GlobalRole) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.GlobalRole) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.GlobalRole) (runtime.Object, error)
}

type globalRoleLifecycleAdapter struct {
	lifecycle GlobalRoleLifecycleContext
}

func (w *globalRoleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalRoleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalRoleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *globalRoleLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.GlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRoleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *globalRoleLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.GlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRoleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *globalRoleLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.GlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalRoleLifecycleAdapter(name string, clusterScoped bool, client GlobalRoleInterface, l GlobalRoleLifecycle) GlobalRoleHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalRoleGroupVersionResource)
	}
	adapter := &globalRoleLifecycleAdapter{lifecycle: &globalRoleLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GlobalRole) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGlobalRoleLifecycleAdapterContext(name string, clusterScoped bool, client GlobalRoleInterface, l GlobalRoleLifecycleContext) GlobalRoleHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalRoleGroupVersionResource)
	}
	adapter := &globalRoleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.GlobalRole) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
