/*
 *
 * Copyright 2023 puzzletools authors.
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
 */
package initrightdb

import (
	"context"
	"fmt"
	"time"

	pb "github.com/dvaumoron/puzzlerightservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const adminGroupId = 1 // groupId corresponding to role administration
const administratorName = "Administrator"

func MakeUserAdmin(rightServiceAddr string, id uint64) error {
	conn, err := grpc.Dial(rightServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewRightClient(conn)
	response, err := client.UpdateRole(ctx, &pb.Role{
		Name: administratorName, ObjectId: adminGroupId, List: []pb.RightAction{
			pb.RightAction_ACCESS, pb.RightAction_CREATE,
			pb.RightAction_UPDATE, pb.RightAction_DELETE,
		},
	})
	if err != nil {
		return err
	}
	if response.Success {
		fmt.Println(administratorName, "role updated.")
	}

	response, err = client.UpdateUser(ctx, &pb.UserRight{
		UserId: id, List: []*pb.RoleRequest{{Name: administratorName, ObjectId: adminGroupId}},
	})
	if err == nil && response.Success {
		fmt.Println("User with the id", id, "has the", administratorName, "role.")
	}
	return err
}
