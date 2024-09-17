package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleBindingLifecycleConverter struct {
	lifecycle RoleBindingLifecycle
}

func (w *roleBindingLifecycleConverter) CreateContext(_ context.Context, obj *v1.RoleBinding) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *roleBindingLifecycleConverter) RemoveContext(_ context.Context, obj *v1.RoleBinding) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *roleBindingLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.RoleBinding) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RoleBindingLifecycle interface {
	Create(obj *v1.RoleBinding) (runtime.Object, error)
	Remove(obj *v1.RoleBinding) (runtime.Object, error)
	Updated(obj *v1.RoleBinding) (runtime.Object, error)
}

type RoleBindingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.RoleBinding) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.RoleBinding) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.RoleBinding) (runtime.Object, error)
}

type roleBindingLifecycleAdapter struct {
	lifecycle RoleBindingLifecycleContext
}

func (w *roleBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *roleBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *roleBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *roleBindingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.RoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *roleBindingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.RoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *roleBindingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.RoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRoleBindingLifecycleAdapter(name string, clusterScoped bool, client RoleBindingInterface, l RoleBindingLifecycle) RoleBindingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleBindingGroupVersionResource)
	}
	adapter := &roleBindingLifecycleAdapter{lifecycle: &roleBindingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.RoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRoleBindingLifecycleAdapterContext(name string, clusterScoped bool, client RoleBindingInterface, l RoleBindingLifecycleContext) RoleBindingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleBindingGroupVersionResource)
	}
	adapter := &roleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.RoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
