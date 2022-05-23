/*
Copyright (C) 2022 Christoph Görn

This file is part of r-gespraech.

r-gespraech is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

r-gespraech is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with r-gespraech.  If not, see <http://www.gnu.org/licenses/>.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/goern/r-gespraech/api/v1alpha1"
)

var _ = Describe("CallbackUrl controller", func() {
	const (
		CallbackUrlName      = "abc123"
		CallbackPayloadAName = "abc123"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a CallbackUrl not having associated CallbackPayload", func() {
		It("Should have NoAssociatedPayloads Condition", func() {
			By("By creating a new CallbackUrl")
			labels := make(map[string]string)
			labels["adviser.thoth-station.ninja/adviser-id"] = "abc123"

			ctx := context.Background()
			url := &v1alpha1.CallbackUrl{
				TypeMeta:   metav1.TypeMeta{APIVersion: "erinnerung.thoth-station.ninja/v1alpha1", Kind: "CallbackUrl"},
				ObjectMeta: metav1.ObjectMeta{Name: CallbackUrlName, Namespace: "default", Labels: labels},
				Spec: v1alpha1.CallbackUrlSpec{
					URL: "https://localhost.local:8181/webhook/xyz_callback",
					Selector: metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
				Status: v1alpha1.CallbackUrlStatus{},
			}
			Expect(k8sClient.Create(ctx, url)).Should(Succeed())

			lookupKey := types.NamespacedName{Name: CallbackUrlName, Namespace: "default"}
			createdCallbackUrl := &v1alpha1.CallbackUrl{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdCallbackUrl)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(meta.IsStatusConditionPresentAndEqual(url.Status.Conditions, v1alpha1.NoAssociatedPayloads, metav1.ConditionTrue))
		})
	})
	Context("When creating an associated CallbackPayload", func() {
		It("Should have AssociatedPayloads Condition", func() {
			By("By creating a new CallbackPayload")
			labels := make(map[string]string)
			labels["adviser.thoth-station.ninja/adviser-id"] = "abc123"

			ctx := context.Background()
			url := &v1alpha1.CallbackPayload{
				TypeMeta:   metav1.TypeMeta{APIVersion: "erinnerung.thoth-station.ninja/v1alpha1", Kind: "CallbackPayload"},
				ObjectMeta: metav1.ObjectMeta{Name: CallbackPayloadAName, Namespace: "default", Labels: labels},
				Spec: v1alpha1.CallbackPayloadSpec{
					Data: "{'adviser-document-id': 'abc123'}",
					Selector: metav1.LabelSelector{
						MatchLabels: labels,
					},
				},
				Status: v1alpha1.CallbackPayloadStatus{},
			}
			Expect(k8sClient.Create(ctx, url)).Should(Succeed())

		})
	})
})
