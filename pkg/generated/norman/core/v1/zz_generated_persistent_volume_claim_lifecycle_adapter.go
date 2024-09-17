package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type persistentVolumeClaimLifecycleConverter struct {
	lifecycle PersistentVolumeClaimLifecycle
}

func (w *persistentVolumeClaimLifecycleConverter) CreateContext(_ context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *persistentVolumeClaimLifecycleConverter) RemoveContext(_ context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *persistentVolumeClaimLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type PersistentVolumeClaimLifecycle interface {
	Create(obj *v1.PersistentVolumeClaim) (runtime.Object, error)
	Remove(obj *v1.PersistentVolumeClaim) (runtime.Object, error)
	Updated(obj *v1.PersistentVolumeClaim) (runtime.Object, error)
}

type PersistentVolumeClaimLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.PersistentVolumeClaim) (runtime.Object, error)
}

type persistentVolumeClaimLifecycleAdapter struct {
	lifecycle PersistentVolumeClaimLifecycleContext
}

func (w *persistentVolumeClaimLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *persistentVolumeClaimLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *persistentVolumeClaimLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *persistentVolumeClaimLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.PersistentVolumeClaim))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *persistentVolumeClaimLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *persistentVolumeClaimLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.PersistentVolumeClaim))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *persistentVolumeClaimLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *persistentVolumeClaimLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.PersistentVolumeClaim))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPersistentVolumeClaimLifecycleAdapter(name string, clusterScoped bool, client PersistentVolumeClaimInterface, l PersistentVolumeClaimLifecycle) PersistentVolumeClaimHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(PersistentVolumeClaimGroupVersionResource)
	}
	adapter := &persistentVolumeClaimLifecycleAdapter{lifecycle: &persistentVolumeClaimLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.PersistentVolumeClaim) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewPersistentVolumeClaimLifecycleAdapterContext(name string, clusterScoped bool, client PersistentVolumeClaimInterface, l PersistentVolumeClaimLifecycleContext) PersistentVolumeClaimHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(PersistentVolumeClaimGroupVersionResource)
	}
	adapter := &persistentVolumeClaimLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.PersistentVolumeClaim) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
