package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleTemplateLifecycleConverter struct {
	lifecycle RoleTemplateLifecycle
}

func (w *roleTemplateLifecycleConverter) CreateContext(_ context.Context, obj *v3.RoleTemplate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *roleTemplateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.RoleTemplate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *roleTemplateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.RoleTemplate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RoleTemplateLifecycle interface {
	Create(obj *v3.RoleTemplate) (runtime.Object, error)
	Remove(obj *v3.RoleTemplate) (runtime.Object, error)
	Updated(obj *v3.RoleTemplate) (runtime.Object, error)
}

type RoleTemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.RoleTemplate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.RoleTemplate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.RoleTemplate) (runtime.Object, error)
}

type roleTemplateLifecycleAdapter struct {
	lifecycle RoleTemplateLifecycleContext
}

func (w *roleTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *roleTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *roleTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *roleTemplateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.RoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *roleTemplateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.RoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *roleTemplateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.RoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRoleTemplateLifecycleAdapter(name string, clusterScoped bool, client RoleTemplateInterface, l RoleTemplateLifecycle) RoleTemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleTemplateGroupVersionResource)
	}
	adapter := &roleTemplateLifecycleAdapter{lifecycle: &roleTemplateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RoleTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRoleTemplateLifecycleAdapterContext(name string, clusterScoped bool, client RoleTemplateInterface, l RoleTemplateLifecycleContext) RoleTemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RoleTemplateGroupVersionResource)
	}
	adapter := &roleTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.RoleTemplate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
