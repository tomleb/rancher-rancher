package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type kontainerDriverLifecycleConverter struct {
	lifecycle KontainerDriverLifecycle
}

func (w *kontainerDriverLifecycleConverter) CreateContext(_ context.Context, obj *v3.KontainerDriver) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *kontainerDriverLifecycleConverter) RemoveContext(_ context.Context, obj *v3.KontainerDriver) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *kontainerDriverLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.KontainerDriver) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type KontainerDriverLifecycle interface {
	Create(obj *v3.KontainerDriver) (runtime.Object, error)
	Remove(obj *v3.KontainerDriver) (runtime.Object, error)
	Updated(obj *v3.KontainerDriver) (runtime.Object, error)
}

type KontainerDriverLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.KontainerDriver) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.KontainerDriver) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.KontainerDriver) (runtime.Object, error)
}

type kontainerDriverLifecycleAdapter struct {
	lifecycle KontainerDriverLifecycleContext
}

func (w *kontainerDriverLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *kontainerDriverLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *kontainerDriverLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *kontainerDriverLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *kontainerDriverLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *kontainerDriverLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *kontainerDriverLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *kontainerDriverLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewKontainerDriverLifecycleAdapter(name string, clusterScoped bool, client KontainerDriverInterface, l KontainerDriverLifecycle) KontainerDriverHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(KontainerDriverGroupVersionResource)
	}
	adapter := &kontainerDriverLifecycleAdapter{lifecycle: &kontainerDriverLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.KontainerDriver) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewKontainerDriverLifecycleAdapterContext(name string, clusterScoped bool, client KontainerDriverInterface, l KontainerDriverLifecycleContext) KontainerDriverHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(KontainerDriverGroupVersionResource)
	}
	adapter := &kontainerDriverLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.KontainerDriver) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
