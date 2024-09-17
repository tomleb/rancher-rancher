package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type etcdBackupLifecycleConverter struct {
	lifecycle EtcdBackupLifecycle
}

func (w *etcdBackupLifecycleConverter) CreateContext(_ context.Context, obj *v3.EtcdBackup) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *etcdBackupLifecycleConverter) RemoveContext(_ context.Context, obj *v3.EtcdBackup) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *etcdBackupLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.EtcdBackup) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type EtcdBackupLifecycle interface {
	Create(obj *v3.EtcdBackup) (runtime.Object, error)
	Remove(obj *v3.EtcdBackup) (runtime.Object, error)
	Updated(obj *v3.EtcdBackup) (runtime.Object, error)
}

type EtcdBackupLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.EtcdBackup) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.EtcdBackup) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.EtcdBackup) (runtime.Object, error)
}

type etcdBackupLifecycleAdapter struct {
	lifecycle EtcdBackupLifecycleContext
}

func (w *etcdBackupLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *etcdBackupLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *etcdBackupLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *etcdBackupLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *etcdBackupLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *etcdBackupLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *etcdBackupLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *etcdBackupLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.EtcdBackup))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewEtcdBackupLifecycleAdapter(name string, clusterScoped bool, client EtcdBackupInterface, l EtcdBackupLifecycle) EtcdBackupHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(EtcdBackupGroupVersionResource)
	}
	adapter := &etcdBackupLifecycleAdapter{lifecycle: &etcdBackupLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.EtcdBackup) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewEtcdBackupLifecycleAdapterContext(name string, clusterScoped bool, client EtcdBackupInterface, l EtcdBackupLifecycleContext) EtcdBackupHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(EtcdBackupGroupVersionResource)
	}
	adapter := &etcdBackupLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.EtcdBackup) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
