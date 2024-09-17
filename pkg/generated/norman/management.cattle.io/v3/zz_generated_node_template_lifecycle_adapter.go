package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type nodeTemplateLifecycleConverter struct {
	lifecycle NodeTemplateLifecycle
}

func (w *nodeTemplateLifecycleConverter) CreateContext(_ context.Context, obj *v3.NodeTemplate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *nodeTemplateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NodeTemplate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *nodeTemplateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NodeTemplate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NodeTemplateLifecycle interface {
	Create(obj *v3.NodeTemplate) (runtime.Object, error)
	Remove(obj *v3.NodeTemplate) (runtime.Object, error)
	Updated(obj *v3.NodeTemplate) (runtime.Object, error)
}

type NodeTemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NodeTemplate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NodeTemplate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NodeTemplate) (runtime.Object, error)
}

type nodeTemplateLifecycleAdapter struct {
	lifecycle NodeTemplateLifecycleContext
}

func (w *nodeTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *nodeTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *nodeTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *nodeTemplateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NodeTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *nodeTemplateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NodeTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *nodeTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *nodeTemplateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NodeTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNodeTemplateLifecycleAdapter(name string, clusterScoped bool, client NodeTemplateInterface, l NodeTemplateLifecycle) NodeTemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeTemplateGroupVersionResource)
	}
	adapter := &nodeTemplateLifecycleAdapter{lifecycle: &nodeTemplateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NodeTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNodeTemplateLifecycleAdapterContext(name string, clusterScoped bool, client NodeTemplateInterface, l NodeTemplateLifecycleContext) NodeTemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NodeTemplateGroupVersionResource)
	}
	adapter := &nodeTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NodeTemplate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
