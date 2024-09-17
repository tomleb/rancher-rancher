package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type dynamicSchemaLifecycleConverter struct {
	lifecycle DynamicSchemaLifecycle
}

func (w *dynamicSchemaLifecycleConverter) CreateContext(_ context.Context, obj *v3.DynamicSchema) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *dynamicSchemaLifecycleConverter) RemoveContext(_ context.Context, obj *v3.DynamicSchema) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *dynamicSchemaLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.DynamicSchema) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type DynamicSchemaLifecycle interface {
	Create(obj *v3.DynamicSchema) (runtime.Object, error)
	Remove(obj *v3.DynamicSchema) (runtime.Object, error)
	Updated(obj *v3.DynamicSchema) (runtime.Object, error)
}

type DynamicSchemaLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.DynamicSchema) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.DynamicSchema) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.DynamicSchema) (runtime.Object, error)
}

type dynamicSchemaLifecycleAdapter struct {
	lifecycle DynamicSchemaLifecycleContext
}

func (w *dynamicSchemaLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *dynamicSchemaLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *dynamicSchemaLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *dynamicSchemaLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dynamicSchemaLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *dynamicSchemaLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dynamicSchemaLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *dynamicSchemaLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDynamicSchemaLifecycleAdapter(name string, clusterScoped bool, client DynamicSchemaInterface, l DynamicSchemaLifecycle) DynamicSchemaHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(DynamicSchemaGroupVersionResource)
	}
	adapter := &dynamicSchemaLifecycleAdapter{lifecycle: &dynamicSchemaLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.DynamicSchema) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewDynamicSchemaLifecycleAdapterContext(name string, clusterScoped bool, client DynamicSchemaInterface, l DynamicSchemaLifecycleContext) DynamicSchemaHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(DynamicSchemaGroupVersionResource)
	}
	adapter := &dynamicSchemaLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.DynamicSchema) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
