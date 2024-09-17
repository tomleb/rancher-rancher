package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type daemonSetLifecycleConverter struct {
	lifecycle DaemonSetLifecycle
}

func (w *daemonSetLifecycleConverter) CreateContext(_ context.Context, obj *v1.DaemonSet) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *daemonSetLifecycleConverter) RemoveContext(_ context.Context, obj *v1.DaemonSet) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *daemonSetLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.DaemonSet) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type DaemonSetLifecycle interface {
	Create(obj *v1.DaemonSet) (runtime.Object, error)
	Remove(obj *v1.DaemonSet) (runtime.Object, error)
	Updated(obj *v1.DaemonSet) (runtime.Object, error)
}

type DaemonSetLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.DaemonSet) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.DaemonSet) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.DaemonSet) (runtime.Object, error)
}

type daemonSetLifecycleAdapter struct {
	lifecycle DaemonSetLifecycleContext
}

func (w *daemonSetLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *daemonSetLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *daemonSetLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *daemonSetLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.DaemonSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *daemonSetLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *daemonSetLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.DaemonSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *daemonSetLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *daemonSetLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.DaemonSet))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDaemonSetLifecycleAdapter(name string, clusterScoped bool, client DaemonSetInterface, l DaemonSetLifecycle) DaemonSetHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(DaemonSetGroupVersionResource)
	}
	adapter := &daemonSetLifecycleAdapter{lifecycle: &daemonSetLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.DaemonSet) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewDaemonSetLifecycleAdapterContext(name string, clusterScoped bool, client DaemonSetInterface, l DaemonSetLifecycleContext) DaemonSetHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(DaemonSetGroupVersionResource)
	}
	adapter := &daemonSetLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.DaemonSet) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
