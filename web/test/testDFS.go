package test

import (
	"fmt"

	"github.com/tedcy/fdfs_client"
)

func testDFS() {
	client, err := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("客户端初始化失败")
	}
	fdfsresponse, err := client.UploadByFilename("./test/111.JPG")
	if err != nil {
		fmt.Println("上传文件失败", err)
	}
	fmt.Println(fdfsresponse)
}
