/*
 * Copyright © 2016-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package ladon

import "context"

// Condition either do or do not fulfill an access request.
type Condition interface {
	// GetName returns the condition's name.
	GetName() string

	// Fulfills returns true if the request is fulfilled by the condition.
	Fulfills(context.Context, interface{}, *Request) bool
}

// Conditions is a collection of conditions.
type Conditions map[string]Condition

// AddCondition adds a condition to the collection.
func (cs Conditions) AddCondition(key string, c Condition) {
	cs[key] = c
}



