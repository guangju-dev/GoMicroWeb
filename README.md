# MicroWeb
一个Go Micro的项目，用到的技术栈挺多的
golang + docker + consul + grpc + protobuf + gin + mysql + redis + fastDFS + nginx
# 2.20~3.1 使用Go-Mirco v3 + Go 1.17重置项目(以其中的User微服务为例)

## 一.环境准备

- consul
- redis
- mysql
- micro v3
	使用micro框架还需要以下插件，我还使用了protoc-gen-gofast
	- [protobuf](https://github.com/golang/protobuf)
	- [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
	- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)
- go开启GOMODULE111=on方便导包
## 二.写代码(以User微服务为例)
1. 在项目的services文件夹下使用命令行键入
	```Go
	micro new user
	```
	
	然后进入该文件夹修改一些文件
	- Makefile 
		protoc --proto_path=. --micro_out=. --go_out=:. proto/user.proto改为
		protoc --proto_path=. --micro_out=. --gofast_out=:. proto/user.proto
	- man.go（**注意导入的包和原来不一样了**）
		```Go
		package main
		
		import (
		  "user/handler"
		  "user/model"
		  pb "user/proto"
		  
		  "github.com/asim/go-micro/plugins/registry/consul/v3"
		  "github.com/asim/go-micro/v3"
		  "github.com/asim/go-micro/v3/logger"
		  "github.com/asim/go-micro/v3/registry"
		)
		
		func main() {
		  model.InitRedis()
		  model.InitDb()
		  // Create service
		  reg := consul.NewRegistry(func(options *registry.Options) {
		    options.Addrs = []string{"127.0.0.1:8500"}
		  })
		  srv := micro.NewService(
		    micro.Name("user"),
		    micro.Version("latest"),
		    micro.Registry(reg),
		  )
		  // Register handler
		  pb.RegisterUserHandler(srv.Server(), new(handler.User))
		
		  // Run service
		  if err := srv.Run(); err != nil {
		    logger.Fatal(err)
		  }
		}
		```
		
2. 根据业务逻辑修改user.proto
	```Protocol Buffers
	syntax = "proto3";
	
	package user;
	
	option go_package = "./proto;user";
	
	service User {
	  //获取用户信息
	  rpc MicroGetUser(Request) returns (Response) {};
	  //更改用户姓名
	  rpc UpdateUserName(UpdateReq)returns(UpdateResp){};
	  //上传头像
	  rpc UploadAvatar(UploadReq)returns(UploadResp){};
	  //实名认证
	  rpc AuthUpdate(AuthReq)returns(AuthResp){};
	}
	
	message AuthReq{
	  string id_card = 1;
	  string real_name = 2;
	  string userName = 3;
	}
	
	message AuthResp{
	  string errno = 1;
	  string errmsg = 2;
	}
	
	
	message UploadData{
	  string avatar_url = 1;
	}
	
	message UploadResp{
	  string errno = 1;
	  string errmsg = 2;
	  UploadData data = 3;
	}
	
	message UploadReq{
	  bytes avatar = 1;
	  string userName = 2;
	  string fileExt = 3;
	}
	
	
	message UpdateReq{
	  string newName = 1;
	  string oldName = 2;
	}
	
	message UpdateResp{
	  string errno = 1;
	  string errmsg = 2;
	  NameData data = 3;
	}
	
	message NameData{
	  string name = 1;
	}
	
	message Request {
	  string name = 1;
	}
	
	message Response {
	  string errno = 1;
	  string errmsg = 2;
	  UserInfo data = 3;
	}
	
	message UserInfo{
	  int32 user_id = 1;
	  string name = 2;
	  string mobile = 3;
	  string real_name = 4;
	  string id_card = 5;
	  string avatar_url = 6;
	}
	
	
	
	```
	
3. 在命令行这个微服务目录下输入`make proto`命令
4. 修改生成的user.pb.micro.go的包的依赖
	```Go
	import (
	  context "context"
	  api "github.com/asim/go-micro/v3/api"
	  client "github.com/asim/go-micro/v3/client"
	  server "github.com/asim/go-micro/v3/server"
	)
	//原来的是
	// import (
	//   context "context"
	//   api "github.com/micro/micro/v3/service/api"
	//   client "github.com/micro/micro/v3/service/client"
	//   server "github.com/micro/micro/v3/service/server"
	// )
	
	```
	
5. 编写handler中的user.go逻辑代码和与数据库打交道的model包
	user.go
	```Go
	package handler
	
	import (
	  "context"
	  "fmt"
	
	  "user/model"
	  user "user/proto"
	  "user/utils"
	
	  "github.com/tedcy/fdfs_client"
	)
	
	type User struct{}
	
	func (e *User) MicroGetUser(ctx context.Context, req *user.Request, rsp *user.Response) error {
	  //根据用户名获取用户信息 在mysql数据库中
	  myUser, err := model.GetUserInfo(req.Name)
	  if err != nil {
	    rsp.Errno = utils.RECODE_USERERR
	    rsp.Errmsg = utils.RecodeText(utils.RECODE_USERERR)
	    return err
	  }
	  rsp.Errno = utils.RECODE_OK
	  rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	
	  //获取一个结构体对象
	  var userInfo user.UserInfo
	  userInfo.UserId = int32(myUser.ID)
	  userInfo.Name = myUser.Name
	  userInfo.Mobile = myUser.Mobile
	  userInfo.RealName = myUser.Real_name
	  userInfo.IdCard = myUser.Id_card
	  userInfo.AvatarUrl = "http://192.168.87.128:8888/" + myUser.Avatar_url
	
	  rsp.Data = &userInfo
	
	  return nil
	}
	
	func (e *User) UpdateUserName(ctx context.Context, req *user.UpdateReq, resp *user.UpdateResp) error {
	  //根据传递过来的用户名更新数据中新的用户名
	  err := model.UpdateUserName(req.NewName, req.OldName)
	  if err != nil {
	    fmt.Println("更新失败", err)
	    resp.Errno = utils.RECODE_DATAERR
	    resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	    //micro规定如果有错误,服务端只给客户端返回错误信息,不返回resp,如果没有错误,就返回resp
	    return nil
	  }
	
	  resp.Errno = utils.RECODE_OK
	  resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	  var nameData user.NameData
	  nameData.Name = req.NewName
	
	  resp.Data = &nameData
	
	  return nil
	}
	
	func (e *User) UploadAvatar(ctx context.Context, req *user.UploadReq, resp *user.UploadResp) error {
	  //存入到fastdfs中
	  fClient, _ := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	  //上传文件到fdfs
	  remoteId, err := fClient.UploadByBuffer(req.Avatar, req.FileExt[1:])
	  if err != nil {
	    fmt.Println("上传文件错误", err)
	    resp.Errno = utils.RECODE_DATAERR
	    resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	    return nil
	  }
	
	  //把存储凭证写入数据库
	  err = model.UpdateAvatar(req.UserName, remoteId)
	  if err != nil {
	    fmt.Println("存储用户头像错误", err)
	    resp.Errno = utils.RECODE_DBERR
	    resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	    return nil
	  }
	
	  resp.Errno = utils.RECODE_OK
	  resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	
	  var uploadData user.UploadData
	  uploadData.AvatarUrl = "http://192.168.87.128:8888/" + remoteId
	  resp.Data = &uploadData
	  return nil
	}
	
	func (e *User) AuthUpdate(ctx context.Context, req *user.AuthReq, resp *user.AuthResp) error {
	  //调用借口校验realName和idcard是否匹配
	
	  //存储真实姓名和真是身份证号  数据库
	  err := model.SaveRealName(req.UserName, req.RealName, req.IdCard)
	  if err != nil {
	    resp.Errno = utils.RECODE_DBERR
	    resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	    return nil
	  }
	
	  resp.Errno = utils.RECODE_OK
	  resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	
	  return nil
	}
	
	```
	
	model.go
	```Go
	package model
	
	import (
	  "time"
	
	  "github.com/gomodule/redigo/redis"
	  "github.com/jinzhu/gorm"
	  _ "github.com/jinzhu/gorm/dialects/mysql"
	)
	
	/* 用户 table_name = user */
	type User struct {
	  ID            int           //用户编号
	  Name          string        `gorm:"size:32;unique"`  //用户名
	  Password_hash string        `gorm:"size:128" `       //用户密码加密的
	  Mobile        string        `gorm:"size:11;unique" ` //手机号
	  Real_name     string        `gorm:"size:32" `        //真实姓名  实名认证
	  Id_card       string        `gorm:"size:20" `        //身份证号  实名认证
	  Avatar_url    string        `gorm:"size:256" `       //用户头像路径       通过fastdfs进行图片存储
	  Houses        []*House      //用户发布的房屋信息  一个人多套房
	  Orders        []*OrderHouse //用户下的订单       一个人多次订单
	}
	
	/* 房屋信息 table_name = house */
	type House struct {
	  gorm.Model                    //房屋编号
	  UserId          uint          //房屋主人的用户编号  与用户进行关联
	  AreaId          uint          //归属地的区域编号   和地区表进行关联
	  Title           string        `gorm:"size:64" `                 //房屋标题
	  Address         string        `gorm:"size:512"`                 //地址
	  Room_count      int           `gorm:"default:1" `               //房间数目
	  Acreage         int           `gorm:"default:0" json:"acreage"` //房屋总面积
	  Price           int           `json:"price"`
	  Unit            string        `gorm:"size:32;default:''" json:"unit"`               //房屋单元,如 几室几厅
	  Capacity        int           `gorm:"default:1" json:"capacity"`                    //房屋容纳的总人数
	  Beds            string        `gorm:"size:64;default:''" json:"beds"`               //房屋床铺的配置
	  Deposit         int           `gorm:"default:0" json:"deposit"`                     //押金
	  Min_days        int           `gorm:"default:1" json:"min_days"`                    //最少入住的天数
	  Max_days        int           `gorm:"default:0" json:"max_days"`                    //最多入住的天数 0表示不限制
	  Order_count     int           `gorm:"default:0" json:"order_count"`                 //预定完成的该房屋的订单数
	  Index_image_url string        `gorm:"size:256;default:''" json:"index_image_url"`   //房屋主图片路径
	  Facilities      []*Facility   `gorm:"many2many:house_facilities" json:"facilities"` //房屋设施   与设施表进行关联
	  Images          []*HouseImage `json:"img_urls"`                                     //房屋的图片   除主要图片之外的其他图片地址
	  Orders          []*OrderHouse `json:"orders"`                                       //房屋的订单    与房屋表进行管理
	}
	
	/* 区域信息 table_name = area */ //区域信息是需要我们手动添加到数据库中的
	type Area struct {
	  Id     int      `json:"aid"`                  //区域编号     1    2
	  Name   string   `gorm:"size:32" json:"aname"` //区域名字     昌平 海淀
	  Houses []*House `json:"houses"`               //区域所有的房屋   与房屋表进行关联
	}
	
	/* 设施信息 table_name = "facility"*/ //设施信息 需要我们提前手动添加的
	type Facility struct {
	  Id     int      `json:"fid"`     //设施编号
	  Name   string   `gorm:"size:32"` //设施名字
	  Houses []*House //都有哪些房屋有此设施  与房屋表进行关联的
	}
	
	/* 房屋图片 table_name = "house_image"*/
	type HouseImage struct {
	  Id      int    `json:"house_image_id"`      //图片id
	  Url     string `gorm:"size:256" json:"url"` //图片url     存放我们房屋的图片
	  HouseId uint   `json:"house_id"`            //图片所属房屋编号
	}
	
	/* 订单 table_name = order */
	type OrderHouse struct {
	  gorm.Model            //订单编号
	  UserId      uint      `json:"user_id"`       //下单的用户编号   //与用户表进行关联
	  HouseId     uint      `json:"house_id"`      //预定的房间编号   //与房屋信息进行关联
	  Begin_date  time.Time `gorm:"type:datetime"` //预定的起始时间
	  End_date    time.Time `gorm:"type:datetime"` //预定的结束时间
	  Days        int       //预定总天数
	  House_price int       //房屋的单价
	  Amount      int       //订单总金额
	  Status      string    `gorm:"default:'WAIT_ACCEPT'"` //订单状态
	  Comment     string    `gorm:"size:512"`              //订单评论
	  Credit      bool      //表示个人征信情况 true表示良好
	}
	
	// 创建 数据库链接句柄
	var GlobalConn *gorm.DB
	
	func InitDb() (*gorm.DB, error) {
	  db, err := gorm.Open("mysql",
	    "root:1201@tcp(127.0.0.1:3306)/test?parseTime=True&loc=Local")
	
	  if err == nil {
	    // 初始化 全局连接池句柄
	    GlobalConn = db
	    GlobalConn.DB().SetMaxIdleConns(10)
	    GlobalConn.DB().SetConnMaxLifetime(100)
	
	    db.SingularTable(true)
	    db.AutoMigrate(new(User), new(House), new(Area), new(Facility), new(HouseImage), new(OrderHouse))
	    return db, nil
	  }
	  return nil, err
	}
	
	//直接用创建一个实例，这样所有的函数就都可以用
	var RedisPool redis.Pool
	
	//初始化Redis连接池
	func InitRedis() {
	  RedisPool = redis.Pool{
	    MaxIdle:         20,
	    MaxActive:       50,
	    MaxConnLifetime: 60 * 5,
	    IdleTimeout:     60,
	    Dial: func() (redis.Conn, error) {
	      return redis.Dial("tcp", "172.19.0.1:6379")
	    },
	  }
	}
	
	```
	
	modelfunc.go
	```Go
	package model
	
	func GetUserInfo(userName string) (User, error) {
	  var user User
	  err := GlobalConn.Where("name = ?", userName).First(&user).Error
	
	  return user, err
	}
	func UpdateUserName(newName, oldName string) error {
	  // update user set name = 'itcast' where name = 旧用户名
	  return GlobalConn.Model(new(User)).Where("name = ?", oldName).Update("name", newName).Error
	}
	
	func UpdateAvatar(userName, avatar string) error {
	  // update user set avatar_url = avatar, where name = username
	  return GlobalConn.Model(new(User)).Where("name = ?", userName).
	    Update("avatar_url", avatar).Error
	}
	func SaveRealName(userName, realName, idCard string) error {
	  return GlobalConn.Model(new(User)).Where("name = ?", userName).
	    Updates(map[string]interface{}{"real_name": realName, "id_card": idCard}).Error
	}
	
	```
	
6. 使用go mod tidy命令拉取依赖的包
## 三.运行项目
1. 首先运行服务发现consul
	```Bash
	consul agent -dev
	
	```
	
2. 运行micro服务
	```Bash
	micro serve
	```
	
3. 开启redis服务
	```Bash
	 redis-server /etc/redis/redis.conf
	```
	
4. 开启mysql服务(一般默认安装完了会自动启动)
5. 在web文件夹下使用
	```Bash
	go run main.go
	
	```
	
6. 在想要开启的微服务项目下使用
	```Bash
	micro run .
	```
	
	注意使用micro status查看当前运行的微服务，如果有error的微服务，直接使用micro kill 服务名即可关闭
7. 在浏览器中查看
