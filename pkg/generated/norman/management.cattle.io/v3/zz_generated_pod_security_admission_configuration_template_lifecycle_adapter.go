package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type podSecurityAdmissionConfigurationTemplateLifecycleConverter struct {
	lifecycle PodSecurityAdmissionConfigurationTemplateLifecycle
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleConverter) CreateContext(_ context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleConverter) RemoveContext(_ context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type PodSecurityAdmissionConfigurationTemplateLifecycle interface {
	Create(obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
	Remove(obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
	Updated(obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
}

type PodSecurityAdmissionConfigurationTemplateLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error)
}

type podSecurityAdmissionConfigurationTemplateLifecycleAdapter struct {
	lifecycle PodSecurityAdmissionConfigurationTemplateLifecycleContext
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *podSecurityAdmissionConfigurationTemplateLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.PodSecurityAdmissionConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPodSecurityAdmissionConfigurationTemplateLifecycleAdapter(name string, clusterScoped bool, client PodSecurityAdmissionConfigurationTemplateInterface, l PodSecurityAdmissionConfigurationTemplateLifecycle) PodSecurityAdmissionConfigurationTemplateHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(PodSecurityAdmissionConfigurationTemplateGroupVersionResource)
	}
	adapter := &podSecurityAdmissionConfigurationTemplateLifecycleAdapter{lifecycle: &podSecurityAdmissionConfigurationTemplateLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewPodSecurityAdmissionConfigurationTemplateLifecycleAdapterContext(name string, clusterScoped bool, client PodSecurityAdmissionConfigurationTemplateInterface, l PodSecurityAdmissionConfigurationTemplateLifecycleContext) PodSecurityAdmissionConfigurationTemplateHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(PodSecurityAdmissionConfigurationTemplateGroupVersionResource)
	}
	adapter := &podSecurityAdmissionConfigurationTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.PodSecurityAdmissionConfigurationTemplate) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
