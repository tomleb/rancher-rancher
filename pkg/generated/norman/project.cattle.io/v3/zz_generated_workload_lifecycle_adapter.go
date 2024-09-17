package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type workloadLifecycleConverter struct {
	lifecycle WorkloadLifecycle
}

func (w *workloadLifecycleConverter) CreateContext(_ context.Context, obj *v3.Workload) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *workloadLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Workload) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *workloadLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Workload) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type WorkloadLifecycle interface {
	Create(obj *v3.Workload) (runtime.Object, error)
	Remove(obj *v3.Workload) (runtime.Object, error)
	Updated(obj *v3.Workload) (runtime.Object, error)
}

type WorkloadLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Workload) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Workload) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Workload) (runtime.Object, error)
}

type workloadLifecycleAdapter struct {
	lifecycle WorkloadLifecycleContext
}

func (w *workloadLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *workloadLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *workloadLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *workloadLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Workload))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *workloadLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *workloadLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Workload))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *workloadLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *workloadLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Workload))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewWorkloadLifecycleAdapter(name string, clusterScoped bool, client WorkloadInterface, l WorkloadLifecycle) WorkloadHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(WorkloadGroupVersionResource)
	}
	adapter := &workloadLifecycleAdapter{lifecycle: &workloadLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Workload) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewWorkloadLifecycleAdapterContext(name string, clusterScoped bool, client WorkloadInterface, l WorkloadLifecycleContext) WorkloadHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(WorkloadGroupVersionResource)
	}
	adapter := &workloadLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Workload) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
