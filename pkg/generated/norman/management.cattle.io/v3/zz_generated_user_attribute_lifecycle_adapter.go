package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type userAttributeLifecycleConverter struct {
	lifecycle UserAttributeLifecycle
}

func (w *userAttributeLifecycleConverter) CreateContext(_ context.Context, obj *v3.UserAttribute) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *userAttributeLifecycleConverter) RemoveContext(_ context.Context, obj *v3.UserAttribute) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *userAttributeLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.UserAttribute) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type UserAttributeLifecycle interface {
	Create(obj *v3.UserAttribute) (runtime.Object, error)
	Remove(obj *v3.UserAttribute) (runtime.Object, error)
	Updated(obj *v3.UserAttribute) (runtime.Object, error)
}

type UserAttributeLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.UserAttribute) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.UserAttribute) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.UserAttribute) (runtime.Object, error)
}

type userAttributeLifecycleAdapter struct {
	lifecycle UserAttributeLifecycleContext
}

func (w *userAttributeLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *userAttributeLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *userAttributeLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *userAttributeLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.UserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *userAttributeLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *userAttributeLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.UserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *userAttributeLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *userAttributeLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.UserAttribute))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewUserAttributeLifecycleAdapter(name string, clusterScoped bool, client UserAttributeInterface, l UserAttributeLifecycle) UserAttributeHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(UserAttributeGroupVersionResource)
	}
	adapter := &userAttributeLifecycleAdapter{lifecycle: &userAttributeLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.UserAttribute) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewUserAttributeLifecycleAdapterContext(name string, clusterScoped bool, client UserAttributeInterface, l UserAttributeLifecycleContext) UserAttributeHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(UserAttributeGroupVersionResource)
	}
	adapter := &userAttributeLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.UserAttribute) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
