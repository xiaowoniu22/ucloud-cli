// Copyright © 2018 NAME HERE tony.li@ucloud.cn
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

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

//NewCmdGssh ucloud gssh
func NewCmdGssh() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gssh",
		Short: "GlobalSSH management",
		Long:  `GlobalSSH management, such as create,modify,list and delete`,
	}
	cmd.AddCommand(NewCmdGsshList())
	cmd.AddCommand(NewCmdGsshCreate())
	cmd.AddCommand(NewCmdGsshDelete())
	cmd.AddCommand(NewCmdGsshModify())
	return cmd
}

//NewCmdGsshList ucloud gssh list
func NewCmdGsshList() *cobra.Command {
	req := client.NewDescribeGlobalSSHInstanceRequest()
	var cmd = &cobra.Command{
		Use:     "ls",
		Short:   "List all GlobalSSH instances",
		Long:    `List all GlobalSSH instances`,
		Example: "ucloud gssh ls",
		Run: func(cmd *cobra.Command, args []string) {
			bindGlobalParam(req)
			resp, err := client.DescribeGlobalSSHInstance(req)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				if resp.RetCode == 0 {
					for _, ins := range resp.InstanceSet {
						fmt.Printf("InstanceID:%s, AcceleratingDomain:%s, TargetIP:%v, Port:%v, Remark:%s\n", ins.InstanceId, ins.AcceleratingDomain, ins.TargetIP, ins.Port, ins.Remark)
					}
				} else {
					fmt.Printf("Something wrong, RetCode:%d, Message:%s\n", resp.RetCode, resp.Message)
				}
			}
		},
	}
	return cmd
}

//NewCmdGsshCreate ucloud gssh create
func NewCmdGsshCreate() *cobra.Command {
	var gsshCreateReq = client.NewCreateGlobalSSHInstanceRequest()
	var cmd = &cobra.Command{
		Use:     "create",
		Short:   "Create GlobalSSH instance",
		Long:    "Create GlobalSSH instance",
		Example: "ucloud gssh create --area Washington --target-ip 8.8.8.8",
		Run: func(cmd *cobra.Command, args []string) {
			bindGlobalParam(gsshCreateReq)
			var areaMap = map[string]string{
				"LosAngeles": "洛杉矶",
				"Singapore":  "新加坡",
				"HongKong":   "香港",
				"Tokyo":      "东京",
				"Washington": "华盛顿",
				"Frankfurt":  "法兰克福",
			}

			port, err := strconv.Atoi(gsshCreateReq.Port)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if port <= 1 || port >= 65535 || port == 80 || port == 443 {
				fmt.Println("The port number should be between 1 and 65535, and cannot be equal to 80 or 443")
				return
			}

			if area, ok := areaMap[gsshCreateReq.Area]; ok {
				gsshCreateReq.Area = area
			} else {
				fmt.Println("Area should be one of LosAngeles,Singapore,HongKong,Tokyo,Washington,Frankfurt.")
				return
			}
			resp, err := client.CreateGlobalSSHInstance(gsshCreateReq)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				if resp.RetCode == 0 {
					fmt.Println("Succeed, GlobalSSHInstanceId:", resp.InstanceId)
				} else {
					fmt.Printf("Something wrong. RetCode:%d,Message:%s\n", resp.RetCode, resp.Message)
				}
			}
		},
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&gsshCreateReq.Area, "area", "", "Location of the source server.Only supports six cities,LosAngeles,Singapore,HongKong,Tokyo,Washington,Frankfurt. Required")
	cmd.Flags().StringVar(&gsshCreateReq.TargetIP, "target-ip", "", "IP of the source server. Required")
	cmd.Flags().StringVar(&gsshCreateReq.Port, "port", "22", "Port of The SSH service between 1 and 65535. Do not use ports such as 80,443.")
	cmd.Flags().StringVar(&gsshCreateReq.Remark, "remark", "", "Remark of your GlobalSSH.")
	cmd.Flags().StringVar(&gsshCreateReq.CouponId, "coupon-id", "", "Coupon ID, The Coupon can deduct part of the payment")
	cmd.MarkFlagRequired("area")
	cmd.MarkFlagRequired("target-ip")
	return cmd
}

//NewCmdGsshDelete ucloud gssh delete
func NewCmdGsshDelete() *cobra.Command {
	var gsshDeleteReq = client.NewDeleteGlobalSSHInstanceRequest()
	var gsshIds []string
	var cmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete GlobalSSH instance",
		Long:    "Delete GlobalSSH instance",
		Example: "ucloud gssh delete --id uga-xx1  --id uga-xx2",
		Run: func(cmd *cobra.Command, args []string) {
			bindGlobalParam(gsshDeleteReq)
			for _, id := range gsshIds {
				gsshDeleteReq.InstanceId = id

				if global.projectID != "" {
					gsshDeleteReq.ProjectId = global.projectID
				}
				resp, err := client.DeleteGlobalSSHInstance(gsshDeleteReq)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					if resp.RetCode == 0 {
						fmt.Printf("GlobalSSH(%s) was successfully deleted\n", id)
					} else {
						fmt.Printf("Something wrong. RetCode:%d, Message:%s\n", resp.RetCode, resp.Message)
					}
				}
			}
		},
	}
	cmd.Flags().StringArrayVar(&gsshIds, "id", make([]string, 0), "ID of the GlobalSSH instances you want to delete. Multiple values specified by multiple flags. Required")
	cmd.MarkFlagRequired("id")
	return cmd
}

//NewCmdGsshModify ucloud gssh modify
func NewCmdGsshModify() *cobra.Command {
	var gsshModifyPortReq = client.NewModifyGlobalSSHPortRequest()
	var gsshModifyRemarkReq = client.NewModifyGlobalSSHRemarkRequest()
	var cmd = &cobra.Command{
		Use:     "modify",
		Short:   "Modify GlobalSSH instance",
		Long:    "Modify GlobalSSH instance, including port and remark attribute",
		Example: "ucloud gssh modify --id uga-xxx --port 22",
		Run: func(cmd *cobra.Command, args []string) {
			bindGlobalParam(gsshModifyPortReq)
			bindGlobalParam(gsshModifyRemarkReq)
			if gsshModifyPortReq.Port == "" && gsshModifyRemarkReq.Remark == "" {
				fmt.Println("port or remark required")
			}
			if gsshModifyPortReq.Port != "" {
				port, err := strconv.Atoi(gsshModifyPortReq.Port)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				if port <= 1 || port >= 65535 || port == 80 || port == 443 {
					fmt.Println("The port number should be between 1 and 65535, and cannot be equal to 80 or 443")
					return
				}
				gsshModifyPortReq.InstanceId = gsshModifyRemarkReq.InstanceId
				resp, err := client.ModifyGlobalSSHPort(gsshModifyPortReq)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					if resp.RetCode == 0 {
						fmt.Println("Successfully modified")
					} else {
						fmt.Printf("Something wrong. RetCode:%d, Message: %s\n", resp.RetCode, resp.Message)
					}
				}
			}
			if gsshModifyRemarkReq.Remark != "" {
				resp, err := client.ModifyGlobalSSHRemark(gsshModifyRemarkReq)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					if resp.RetCode == 0 {
						fmt.Println("Successfully modified")
					} else {
						fmt.Printf("Something wrong. RetCode:%d, Message: %s\n", resp.RetCode, resp.Message)
					}
				}
			}
		},
	}
	cmd.Flags().StringVar(&gsshModifyPortReq.Port, "port", "", "Port of SSH service.")
	cmd.Flags().StringVar(&gsshModifyRemarkReq.Remark, "remark", "", "Remark of your GlobalSSH.")
	cmd.Flags().StringVar(&gsshModifyRemarkReq.InstanceId, "id", "", "InstanceID of your GlobalSSH. Required")
	cmd.MarkFlagRequired("id")
	return cmd
}
