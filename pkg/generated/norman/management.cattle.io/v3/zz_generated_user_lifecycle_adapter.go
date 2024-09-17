package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type userLifecycleConverter struct {
	lifecycle UserLifecycle
}

func (w *userLifecycleConverter) CreateContext(_ context.Context, obj *v3.User) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *userLifecycleConverter) RemoveContext(_ context.Context, obj *v3.User) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *userLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.User) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type UserLifecycle interface {
	Create(obj *v3.User) (runtime.Object, error)
	Remove(obj *v3.User) (runtime.Object, error)
	Updated(obj *v3.User) (runtime.Object, error)
}

type UserLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.User) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.User) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.User) (runtime.Object, error)
}

type userLifecycleAdapter struct {
	lifecycle UserLifecycleContext
}

func (w *userLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *userLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *userLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *userLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.User))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *userLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *userLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.User))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *userLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *userLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.User))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewUserLifecycleAdapter(name string, clusterScoped bool, client UserInterface, l UserLifecycle) UserHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(UserGroupVersionResource)
	}
	adapter := &userLifecycleAdapter{lifecycle: &userLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.User) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewUserLifecycleAdapterContext(name string, clusterScoped bool, client UserInterface, l UserLifecycleContext) UserHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(UserGroupVersionResource)
	}
	adapter := &userLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.User) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
