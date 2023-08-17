// Copyright 2023 Upbound Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route53

import (
	srcv1alpha1 "github.com/crossplane-contrib/provider-aws/apis/route53/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	targetv1beta1 "github.com/upbound/provider-aws/apis/route53/v1beta1"
	"github.com/upbound/upjet/pkg/migration"
)

func ResourceRecordSetResource(mg resource.Managed) ([]resource.Managed, error) {
	source := mg.(*srcv1alpha1.ResourceRecordSet)
	target := &targetv1beta1.Record{}
	if _, err := migration.CopyInto(source, target, targetv1beta1.Record_GroupVersionKind); err != nil {
		return nil, errors.Wrap(err, "failed to copy source into target")
	}
	target.Spec.ForProvider.Records = []*string{}
	for _, t := range source.Spec.ForProvider.ResourceRecords {
		v := t.Value
		target.Spec.ForProvider.Records = append(target.Spec.ForProvider.Records, &v)
	}

	externalName := source.Annotations["crossplane.io/external-name"]
	if len(externalName) > 0 {
		target.Spec.ForProvider.Name = &externalName
	}

	return []resource.Managed{
		target,
	}, nil
}