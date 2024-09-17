package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type preferenceLifecycleConverter struct {
	lifecycle PreferenceLifecycle
}

func (w *preferenceLifecycleConverter) CreateContext(_ context.Context, obj *v3.Preference) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *preferenceLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Preference) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *preferenceLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Preference) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type PreferenceLifecycle interface {
	Create(obj *v3.Preference) (runtime.Object, error)
	Remove(obj *v3.Preference) (runtime.Object, error)
	Updated(obj *v3.Preference) (runtime.Object, error)
}

type PreferenceLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Preference) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Preference) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Preference) (runtime.Object, error)
}

type preferenceLifecycleAdapter struct {
	lifecycle PreferenceLifecycleContext
}

func (w *preferenceLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *preferenceLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *preferenceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *preferenceLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *preferenceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *preferenceLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *preferenceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *preferenceLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Preference))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPreferenceLifecycleAdapter(name string, clusterScoped bool, client PreferenceInterface, l PreferenceLifecycle) PreferenceHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(PreferenceGroupVersionResource)
	}
	adapter := &preferenceLifecycleAdapter{lifecycle: &preferenceLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Preference) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewPreferenceLifecycleAdapterContext(name string, clusterScoped bool, client PreferenceInterface, l PreferenceLifecycleContext) PreferenceHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(PreferenceGroupVersionResource)
	}
	adapter := &preferenceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Preference) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
