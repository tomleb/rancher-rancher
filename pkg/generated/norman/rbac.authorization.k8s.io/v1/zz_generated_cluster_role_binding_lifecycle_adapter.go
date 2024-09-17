package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleBindingLifecycleConverter struct {
	lifecycle ClusterRoleBindingLifecycle
}

func (w *clusterRoleBindingLifecycleConverter) CreateContext(_ context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterRoleBindingLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterRoleBindingLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterRoleBindingLifecycle interface {
	Create(obj *v1.ClusterRoleBinding) (runtime.Object, error)
	Remove(obj *v1.ClusterRoleBinding) (runtime.Object, error)
	Updated(obj *v1.ClusterRoleBinding) (runtime.Object, error)
}

type ClusterRoleBindingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ClusterRoleBinding) (runtime.Object, error)
}

type clusterRoleBindingLifecycleAdapter struct {
	lifecycle ClusterRoleBindingLifecycleContext
}

func (w *clusterRoleBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRoleBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRoleBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterRoleBindingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ClusterRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterRoleBindingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ClusterRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterRoleBindingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ClusterRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRoleBindingLifecycleAdapter(name string, clusterScoped bool, client ClusterRoleBindingInterface, l ClusterRoleBindingLifecycle) ClusterRoleBindingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleBindingGroupVersionResource)
	}
	adapter := &clusterRoleBindingLifecycleAdapter{lifecycle: &clusterRoleBindingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ClusterRoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterRoleBindingLifecycleAdapterContext(name string, clusterScoped bool, client ClusterRoleBindingInterface, l ClusterRoleBindingLifecycleContext) ClusterRoleBindingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleBindingGroupVersionResource)
	}
	adapter := &clusterRoleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ClusterRoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
