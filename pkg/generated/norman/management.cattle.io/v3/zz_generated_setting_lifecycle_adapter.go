package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type settingLifecycleConverter struct {
	lifecycle SettingLifecycle
}

func (w *settingLifecycleConverter) CreateContext(_ context.Context, obj *v3.Setting) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *settingLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Setting) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *settingLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Setting) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type SettingLifecycle interface {
	Create(obj *v3.Setting) (runtime.Object, error)
	Remove(obj *v3.Setting) (runtime.Object, error)
	Updated(obj *v3.Setting) (runtime.Object, error)
}

type SettingLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Setting) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Setting) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Setting) (runtime.Object, error)
}

type settingLifecycleAdapter struct {
	lifecycle SettingLifecycleContext
}

func (w *settingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *settingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *settingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *settingLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Setting))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *settingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *settingLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Setting))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *settingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *settingLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Setting))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSettingLifecycleAdapter(name string, clusterScoped bool, client SettingInterface, l SettingLifecycle) SettingHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(SettingGroupVersionResource)
	}
	adapter := &settingLifecycleAdapter{lifecycle: &settingLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Setting) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewSettingLifecycleAdapterContext(name string, clusterScoped bool, client SettingInterface, l SettingLifecycleContext) SettingHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(SettingGroupVersionResource)
	}
	adapter := &settingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Setting) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
