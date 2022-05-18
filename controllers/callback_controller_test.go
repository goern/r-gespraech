/**
 * Copyright (C) 2022 Christoph Görn
 *
 * This file is part of rebuldah.
 *
 * rebuldah is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * rebuldah is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with rebuldah.  If not, see <http://www.gnu.org/licenses/>.
 */

package controllers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webhookv1alpha1 "github.com/goern/r-gespraech/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Callback controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		CallbackName      = "test"
		CallbackNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Callback has an unparsable URL", func() {
		It("Should state the problem", func() {
			By("By adding a new Condition to the Callback")
			ctx := context.Background()
			c := &webhookv1alpha1.Callback{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "erinnerung.thoth-station.ninja/v1alpha1",
					Kind:       "Callback",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      CallbackName,
					Namespace: CallbackNamespace,
				},
				Spec: webhookv1alpha1.CallbackSpec{
					URL:     "broken://localhost.local:9191/page.asp",
					Payload: "{}",
				},
			}
			Expect(k8sClient.Create(ctx, c)).Should(Succeed())
		})
	})
})
