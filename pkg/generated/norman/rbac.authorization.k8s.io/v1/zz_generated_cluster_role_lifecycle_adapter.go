package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleLifecycleConverter struct {
	lifecycle ClusterRoleLifecycle
}

func (w *clusterRoleLifecycleConverter) CreateContext(_ context.Context, obj *v1.ClusterRole) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *clusterRoleLifecycleConverter) RemoveContext(_ context.Context, obj *v1.ClusterRole) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *clusterRoleLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.ClusterRole) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type ClusterRoleLifecycle interface {
	Create(obj *v1.ClusterRole) (runtime.Object, error)
	Remove(obj *v1.ClusterRole) (runtime.Object, error)
	Updated(obj *v1.ClusterRole) (runtime.Object, error)
}

type ClusterRoleLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.ClusterRole) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.ClusterRole) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.ClusterRole) (runtime.Object, error)
}

type clusterRoleLifecycleAdapter struct {
	lifecycle ClusterRoleLifecycleContext
}

func (w *clusterRoleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRoleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRoleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *clusterRoleLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.ClusterRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *clusterRoleLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.ClusterRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRoleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *clusterRoleLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.ClusterRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRoleLifecycleAdapter(name string, clusterScoped bool, client ClusterRoleInterface, l ClusterRoleLifecycle) ClusterRoleHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleGroupVersionResource)
	}
	adapter := &clusterRoleLifecycleAdapter{lifecycle: &clusterRoleLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.ClusterRole) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewClusterRoleLifecycleAdapterContext(name string, clusterScoped bool, client ClusterRoleInterface, l ClusterRoleLifecycleContext) ClusterRoleHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(ClusterRoleGroupVersionResource)
	}
	adapter := &clusterRoleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.ClusterRole) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
