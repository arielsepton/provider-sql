/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exec

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpcontroller "github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-sql/apis/mssql/v1alpha1"
	"github.com/crossplane-contrib/provider-sql/pkg/clients/mssql"
	"github.com/crossplane-contrib/provider-sql/pkg/clients/xsql"
)

const (
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errNoSecretRef  = "ProviderConfig does not reference a credentials Secret"
	errGetSecret    = "cannot get credentials Secret"

	errNotQuery               = "managed resource is not a Query custom resource"
	errSelectQuery            = "cannot select Query"
	errCreateQuery            = "cannot create Query"
	errDropQuery              = "error dropping Query %s"
	errDropLogin              = "error dropping login %s"
	errCannotQueryuteQuery    = "cannot get current logins %s"
	errCannotKillLoginSession = "error killing session %d for login %s"

	maxConcurrency = 5
)

// Setup adds a controller that reconciles Query managed resources.
func Setup(mgr ctrl.Manager, o xpcontroller.Options) error {
	name := managed.ControllerName(v1alpha1.QueryGroupKind)

	t := resource.NewProviderConfigUsageTracker(mgr.GetClient(), &v1alpha1.ProviderConfigUsage{})
	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.QueryGroupVersionKind),
		managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), usage: t, newClient: mssql.New}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Query{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: maxConcurrency,
		}).
		Complete(r)
}

type connector struct {
	kube      client.Client
	usage     resource.Tracker
	newClient func(creds map[string][]byte, database string) xsql.DB
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Exec)
	if !ok {
		return nil, errors.New(errNotQuery)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	// ProviderConfigReference could theoretically be nil, but in practice the
	// DefaultProviderConfig initializer will set it before we get here.
	pc := &v1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	// We don't need to check the credentials source because we currently only
	// support one source (MySQLConnectionSecret), which is required and
	// enforced by the ProviderConfig schema.
	ref := pc.Spec.Credentials.ConnectionSecretRef
	if ref == nil {
		return nil, errors.New(errNoSecretRef)
	}

	s := &corev1.Secret{}
	if err := c.kube.Get(ctx, types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}, s); err != nil {
		return nil, errors.Wrap(err, errGetSecret)
	}

	return &external{
		db:   c.newClient(s.Data, pointer.StringPtrDerefOr(cr.Spec.ForProvider.Database, "")),
		kube: c.kube,
	}, nil
}

type external struct {
	db   xsql.DB
	kube client.Client
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Exec)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotQuery)
	}

	return managed.ExternalObservation{
		ResourceExists:   cr.Status.Synced,
		ResourceUpToDate: cr.Status.Synced,
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Exec)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotQuery)
	}

	Query := cr.Spec.ForProvider.Exec
	_, err := c.db.Query(ctx, xsql.Query{String: Query})
	if err != nil {
		cr.Status.AtProvider.Error = err.Error()
		cr.Status.Synced = true
		if err := c.kube.Status().Update(ctx, cr); err != nil {
			return managed.ExternalCreation{}, errors.Wrap(err, "errFailedToSetStatusCode")
		}

		return managed.ExternalCreation{}, nil
	}

	// defer rows.Close() //nolint:errcheck
	// var results []map[string]string

	// columns, err := rows.Columns()
	// if err != nil {
	// 	return managed.ExternalCreation{}, errors.Wrap(err, "failed to get columns")
	// }

	// for rows.Next() {
	// 	values := make([]interface{}, len(columns))
	// 	rawValues := make([]sql.RawBytes, len(columns))
	// 	for i := range values {
	// 		values[i] = &rawValues[i]
	// 	}

	// 	if err := rows.Scan(values...); err != nil {
	// 		return managed.ExternalCreation{}, errors.Wrap(err, "failed to scan row")
	// 	}

	// 	row := make(map[string]string)
	// 	for i, column := range columns {
	// 		row[column] = string(rawValues[i])
	// 	}

	// 	results = append(results, row)
	// }

	cr.Status.AtProvider.Message = "SQL statement Queried successfully"
	cr.Status.Synced = true
	if err := c.kube.Status().Update(ctx, cr); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "errFailedToSetStatusCode")
	}

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	return nil
}
