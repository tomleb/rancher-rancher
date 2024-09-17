package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type rancherUserNotificationLifecycleConverter struct {
	lifecycle RancherUserNotificationLifecycle
}

func (w *rancherUserNotificationLifecycleConverter) CreateContext(_ context.Context, obj *v3.RancherUserNotification) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *rancherUserNotificationLifecycleConverter) RemoveContext(_ context.Context, obj *v3.RancherUserNotification) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *rancherUserNotificationLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.RancherUserNotification) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type RancherUserNotificationLifecycle interface {
	Create(obj *v3.RancherUserNotification) (runtime.Object, error)
	Remove(obj *v3.RancherUserNotification) (runtime.Object, error)
	Updated(obj *v3.RancherUserNotification) (runtime.Object, error)
}

type RancherUserNotificationLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.RancherUserNotification) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.RancherUserNotification) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.RancherUserNotification) (runtime.Object, error)
}

type rancherUserNotificationLifecycleAdapter struct {
	lifecycle RancherUserNotificationLifecycleContext
}

func (w *rancherUserNotificationLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *rancherUserNotificationLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *rancherUserNotificationLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *rancherUserNotificationLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.RancherUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rancherUserNotificationLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *rancherUserNotificationLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.RancherUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *rancherUserNotificationLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *rancherUserNotificationLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.RancherUserNotification))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRancherUserNotificationLifecycleAdapter(name string, clusterScoped bool, client RancherUserNotificationInterface, l RancherUserNotificationLifecycle) RancherUserNotificationHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(RancherUserNotificationGroupVersionResource)
	}
	adapter := &rancherUserNotificationLifecycleAdapter{lifecycle: &rancherUserNotificationLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.RancherUserNotification) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewRancherUserNotificationLifecycleAdapterContext(name string, clusterScoped bool, client RancherUserNotificationInterface, l RancherUserNotificationLifecycleContext) RancherUserNotificationHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(RancherUserNotificationGroupVersionResource)
	}
	adapter := &rancherUserNotificationLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.RancherUserNotification) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
