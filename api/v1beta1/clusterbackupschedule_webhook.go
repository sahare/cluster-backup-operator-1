/*
Copyright 2021.

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

package v1beta1

import (
	"github.com/robfig/cron/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var backupschedulelog = logf.Log.WithName("backupschedule-resource")

// SetupWebhookWithManager sets up the webhook with the Manager.
func (r *BackupSchedule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-cluster-open-cluster-management-io-v1beta1-backupschedule,mutating=false,failurePolicy=fail,sideEffects=None,groups=cluster.open-cluster-management.io,resources=backupschedules,verbs=create;update,versions=v1beta1,name=vbackupschedule.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &BackupSchedule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupSchedule) ValidateCreate() error {
	backupschedulelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateBackupSchedule()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupSchedule) ValidateUpdate(old runtime.Object) error {
	backupschedulelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateBackupSchedule()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *BackupSchedule) ValidateDelete() error {
	backupschedulelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *BackupSchedule) validateBackupSchedule() error {
	var allErrs field.ErrorList
	if err := r.validateScheduleSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: "cluster.open-cluster-management.io",
			Kind:  "BackupSchedule",
		},
		r.Name,
		allErrs,
	)
}

func (r *BackupSchedule) validateScheduleSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateScheduleFormat(
		r.Spec.VeleroSchedule,
		field.NewPath("spec").Child("veleroSchedule"))
}

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}
