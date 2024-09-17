package v2

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/runtime"
)

type horizontalPodAutoscalerLifecycleConverter struct {
	lifecycle HorizontalPodAutoscalerLifecycle
}

func (w *horizontalPodAutoscalerLifecycleConverter) CreateContext(_ context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *horizontalPodAutoscalerLifecycleConverter) RemoveContext(_ context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *horizontalPodAutoscalerLifecycleConverter) UpdatedContext(_ context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type HorizontalPodAutoscalerLifecycle interface {
	Create(obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
	Remove(obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
	Updated(obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
}

type HorizontalPodAutoscalerLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error)
}

type horizontalPodAutoscalerLifecycleAdapter struct {
	lifecycle HorizontalPodAutoscalerLifecycleContext
}

func (w *horizontalPodAutoscalerLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *horizontalPodAutoscalerLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *horizontalPodAutoscalerLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *horizontalPodAutoscalerLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v2.HorizontalPodAutoscaler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *horizontalPodAutoscalerLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *horizontalPodAutoscalerLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v2.HorizontalPodAutoscaler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *horizontalPodAutoscalerLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *horizontalPodAutoscalerLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v2.HorizontalPodAutoscaler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewHorizontalPodAutoscalerLifecycleAdapter(name string, clusterScoped bool, client HorizontalPodAutoscalerInterface, l HorizontalPodAutoscalerLifecycle) HorizontalPodAutoscalerHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(HorizontalPodAutoscalerGroupVersionResource)
	}
	adapter := &horizontalPodAutoscalerLifecycleAdapter{lifecycle: &horizontalPodAutoscalerLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewHorizontalPodAutoscalerLifecycleAdapterContext(name string, clusterScoped bool, client HorizontalPodAutoscalerInterface, l HorizontalPodAutoscalerLifecycleContext) HorizontalPodAutoscalerHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(HorizontalPodAutoscalerGroupVersionResource)
	}
	adapter := &horizontalPodAutoscalerLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v2.HorizontalPodAutoscaler) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
