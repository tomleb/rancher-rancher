package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type nodeLifecycleConverter struct {
	lifecycle NodeLifecycle
}

func (w *nodeLifecycleConverter) CreateContext(_ context.Context, obj *v1.Node) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *nodeLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Node) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *nodeLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Node) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NodeLifecycle interface {
	Create(obj *v1.Node) (runtime.Object, error)
	Remove(obj *v1.Node) (runtime.Object, error)
	Updated(obj *v1.Node) (runtime.Object, error)
}

type NodeLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Node) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Node) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Node) (runtime.Object, error)
}

type nodeLifecycleAdapter struct {
	lifecycle NodeLifecycleContext
}

func (w *nodeLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *nodeLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *nodeLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *nodeLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Node))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *nodeLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Node))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *nodeLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Node))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNodeLifecycleAdapter(name string, clusterScoped bool, client NodeInterface, l NodeLifecycle) NodeHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeGroupVersionResource)
	}
	adapter := &nodeLifecycleAdapter{lifecycle: &nodeLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Node) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNodeLifecycleAdapterContext(name string, clusterScoped bool, client NodeInterface, l NodeLifecycleContext) NodeHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeGroupVersionResource)
	}
	adapter := &nodeLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Node) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
