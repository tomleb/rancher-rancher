package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleLifecycleConverter struct {
	lifecycle RoleLifecycle
}

func (w *roleLifecycleConverter) CreateContext(_ context.Context, obj *v1.Role) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *roleLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Role) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *roleLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Role) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RoleLifecycle interface {
	Create(obj *v1.Role) (runtime.Object, error)
	Remove(obj *v1.Role) (runtime.Object, error)
	Updated(obj *v1.Role) (runtime.Object, error)
}

type RoleLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Role) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Role) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Role) (runtime.Object, error)
}

type roleLifecycleAdapter struct {
	lifecycle RoleLifecycleContext
}

func (w *roleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *roleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *roleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *roleLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *roleLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *roleLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRoleLifecycleAdapter(name string, clusterScoped bool, client RoleInterface, l RoleLifecycle) RoleHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleGroupVersionResource)
	}
	adapter := &roleLifecycleAdapter{lifecycle: &roleLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Role) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRoleLifecycleAdapterContext(name string, clusterScoped bool, client RoleInterface, l RoleLifecycleContext) RoleHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleGroupVersionResource)
	}
	adapter := &roleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Role) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
