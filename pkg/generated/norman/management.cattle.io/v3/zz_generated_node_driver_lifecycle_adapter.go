package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type nodeDriverLifecycleConverter struct {
	lifecycle NodeDriverLifecycle
}

func (w *nodeDriverLifecycleConverter) CreateContext(_ context.Context, obj *v3.NodeDriver) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *nodeDriverLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NodeDriver) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *nodeDriverLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NodeDriver) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NodeDriverLifecycle interface {
	Create(obj *v3.NodeDriver) (runtime.Object, error)
	Remove(obj *v3.NodeDriver) (runtime.Object, error)
	Updated(obj *v3.NodeDriver) (runtime.Object, error)
}

type NodeDriverLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NodeDriver) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NodeDriver) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NodeDriver) (runtime.Object, error)
}

type nodeDriverLifecycleAdapter struct {
	lifecycle NodeDriverLifecycleContext
}

func (w *nodeDriverLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *nodeDriverLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *nodeDriverLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *nodeDriverLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NodeDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeDriverLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *nodeDriverLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NodeDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeDriverLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *nodeDriverLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NodeDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNodeDriverLifecycleAdapter(name string, clusterScoped bool, client NodeDriverInterface, l NodeDriverLifecycle) NodeDriverHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeDriverGroupVersionResource)
	}
	adapter := &nodeDriverLifecycleAdapter{lifecycle: &nodeDriverLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NodeDriver) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNodeDriverLifecycleAdapterContext(name string, clusterScoped bool, client NodeDriverInterface, l NodeDriverLifecycleContext) NodeDriverHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeDriverGroupVersionResource)
	}
	adapter := &nodeDriverLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NodeDriver) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
