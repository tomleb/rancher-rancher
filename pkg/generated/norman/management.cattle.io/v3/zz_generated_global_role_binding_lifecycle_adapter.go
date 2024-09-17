package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type globalRoleBindingLifecycleConverter struct {
	lifecycle GlobalRoleBindingLifecycle
}

func (w *globalRoleBindingLifecycleConverter) CreateContext(_ context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *globalRoleBindingLifecycleConverter) RemoveContext(_ context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *globalRoleBindingLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type GlobalRoleBindingLifecycle interface {
	Create(obj *v3.GlobalRoleBinding) (runtime.Object, error)
	Remove(obj *v3.GlobalRoleBinding) (runtime.Object, error)
	Updated(obj *v3.GlobalRoleBinding) (runtime.Object, error)
}

type GlobalRoleBindingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.GlobalRoleBinding) (runtime.Object, error)
}

type globalRoleBindingLifecycleAdapter struct {
	lifecycle GlobalRoleBindingLifecycleContext
}

func (w *globalRoleBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalRoleBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalRoleBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *globalRoleBindingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.GlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRoleBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *globalRoleBindingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.GlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRoleBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *globalRoleBindingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.GlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalRoleBindingLifecycleAdapter(name string, clusterScoped bool, client GlobalRoleBindingInterface, l GlobalRoleBindingLifecycle) GlobalRoleBindingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalRoleBindingGroupVersionResource)
	}
	adapter := &globalRoleBindingLifecycleAdapter{lifecycle: &globalRoleBindingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GlobalRoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewGlobalRoleBindingLifecycleAdapterContext(name string, clusterScoped bool, client GlobalRoleBindingInterface, l GlobalRoleBindingLifecycleContext) GlobalRoleBindingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalRoleBindingGroupVersionResource)
	}
	adapter := &globalRoleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.GlobalRoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
