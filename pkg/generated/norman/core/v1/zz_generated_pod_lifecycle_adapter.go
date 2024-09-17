package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podLifecycleConverter struct {
	lifecycle PodLifecycle
}

func (w *podLifecycleConverter) CreateContext(_ context.Context, obj *v1.Pod) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *podLifecycleConverter) RemoveContext(_ context.Context, obj *v1.Pod) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *podLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.Pod) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type PodLifecycle interface {
	Create(obj *v1.Pod) (runtime.Object, error)
	Remove(obj *v1.Pod) (runtime.Object, error)
	Updated(obj *v1.Pod) (runtime.Object, error)
}

type PodLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.Pod) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.Pod) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.Pod) (runtime.Object, error)
}

type podLifecycleAdapter struct {
	lifecycle PodLifecycleContext
}

func (w *podLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *podLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *podLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *podLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.Pod))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *podLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *podLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.Pod))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *podLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *podLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.Pod))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPodLifecycleAdapter(name string, clusterScoped bool, client PodInterface, l PodLifecycle) PodHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(PodGroupVersionResource)
	}
	adapter := &podLifecycleAdapter{lifecycle: &podLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.Pod) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewPodLifecycleAdapterContext(name string, clusterScoped bool, client PodInterface, l PodLifecycleContext) PodHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(PodGroupVersionResource)
	}
	adapter := &podLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.Pod) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
