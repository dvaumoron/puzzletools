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

package initlogindb

import (
	"context"
	"fmt"
	"time"

	pb "github.com/dvaumoron/puzzleloginservice"
	saltclient "github.com/dvaumoron/puzzlesaltclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const timeOut = 5 * time.Second

func InitUser(saltServiceAddr string, loginServiceAddr string, login string, password string) error {
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	salted, err := saltclient.Make(saltServiceAddr, dialOptions, timeOut).Salt(login, password)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(loginServiceAddr, dialOptions)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	response, err := pb.NewLoginClient(conn).Register(ctx, &pb.LoginRequest{Login: login, Salted: salted})
	if err == nil && response.Success {
		fmt.Println("User with id", response.Id, "created.")
	}
	return err
}
