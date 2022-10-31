// Copyright 2022 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/framework"
	"go.etcd.io/etcd/tests/v3/framework/config"
)

type authRole struct {
	role       string
	permission clientv3.PermissionType
	key        string
	keyEnd     string
}

type authUser struct {
	user string
	pass string
	role string
}

func createRoles(c framework.Client, roles []authRole) error {
	for _, r := range roles {
		// add role
		if _, err := c.RoleAdd(context.TODO(), r.role); err != nil {
			return fmt.Errorf("RoleAdd failed: %w", err)
		}

		// grant permission to role
		if _, err := c.RoleGrantPermission(context.TODO(), r.role, r.key, r.keyEnd, r.permission); err != nil {
			return fmt.Errorf("RoleGrantPermission failed: %w", err)
		}
	}

	return nil
}

func createUsers(c framework.Client, users []authUser) error {
	for _, u := range users {
		// add user
		if _, err := c.UserAdd(context.TODO(), u.user, u.pass, config.UserAddOptions{}); err != nil {
			return fmt.Errorf("UserAdd failed: %w", err)
		}

		// grant role to user
		if _, err := c.UserGrantRole(context.TODO(), u.user, u.role); err != nil {
			return fmt.Errorf("UserGrantRole failed: %w", err)
		}
	}

	return nil
}

func setupAuth(c framework.Client, roles []authRole, users []authUser) error {
	// create roles
	if err := createRoles(c, roles); err != nil {
		return err
	}

	if err := createUsers(c, users); err != nil {
		return err
	}

	// enable auth
	if err := c.AuthEnable(context.TODO()); err != nil {
		return err
	}

	return nil
}
