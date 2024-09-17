package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleTemplateBindingLifecycleConverter struct {
	lifecycle ClusterRoleTemplateBindingLifecycle
}

func (w *clusterRoleTemplateBindingLifecycleConverter) CreateContext(_ context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterRoleTemplateBindingLifecycleConverter) RemoveContext(_ context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterRoleTemplateBindingLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterRoleTemplateBindingLifecycle interface {
	Create(obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
	Remove(obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
	Updated(obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
}

type ClusterRoleTemplateBindingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error)
}

type clusterRoleTemplateBindingLifecycleAdapter struct {
	lifecycle ClusterRoleTemplateBindingLifecycleContext
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.ClusterRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.ClusterRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterRoleTemplateBindingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.ClusterRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRoleTemplateBindingLifecycleAdapter(name string, clusterScoped bool, client ClusterRoleTemplateBindingInterface, l ClusterRoleTemplateBindingLifecycle) ClusterRoleTemplateBindingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleTemplateBindingGroupVersionResource)
	}
	adapter := &clusterRoleTemplateBindingLifecycleAdapter{lifecycle: &clusterRoleTemplateBindingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterRoleTemplateBindingLifecycleAdapterContext(name string, clusterScoped bool, client ClusterRoleTemplateBindingInterface, l ClusterRoleTemplateBindingLifecycleContext) ClusterRoleTemplateBindingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleTemplateBindingGroupVersionResource)
	}
	adapter := &clusterRoleTemplateBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.ClusterRoleTemplateBinding) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
