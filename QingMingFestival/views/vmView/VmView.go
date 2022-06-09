package vmView

import (
	"QingMingFestival/service/vmService"
	"QingMingFestival/tools"
	"QingMingFestival/tools/Variable"
	"QingMingFestival/tools/middleware"
	"QingMingFestival/types/vmTypes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	excelize "github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type VMView struct {
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (vv *VMView) Router(api *gin.RouterGroup) {
	vmApi := api.Group("/vmware")
	vmApi.GET("/clusters", vv.getClusters)     // 获取vc账户信息列表
	vmApi.POST("/clusters/add", vv.addCluster) // 新加vc账户信息
	vmApi.POST("/clusters/del", vv.delCluster) // 删除vc账户信息

	vmClusterApi := vmApi.Group("/cluster")
	vmClusterApi.GET("/:id", vv.getCluster)      // 获取单个vc账户信息
	vmClusterApi.POST("/login", vv.loginCluster) // 登录某个vc集群
	vmClusterApi.GET("/hosts", vv.getHosts)      // 显示所有主机详情
	vmClusterApi.GET("/host", vv.getHost)        // 显示所有主机详情
	vmClusterApi.GET("/vms", vv.getVms)          // 显示虚拟机详情

	vmCloneApi := vmApi.Group("/clone-vm")
	vmCloneApi.GET("/metadata", vv.getCloneVmMetaData)     // 显示虚拟机与模版详情
	vmCloneApi.GET("/jobs", vv.getCloneVmJobs)             // 获取虚拟机任务
	vmCloneApi.POST("/jobs", vv.createCloneVmJob)          // 创建虚拟机任务
	vmCloneApi.POST("/run", vv.cloneVm)                    // 克隆虚拟机
	vmCloneApi.GET("/template", vv.getCloneVmTemplate)     // 下载克隆虚拟机模版
	vmCloneApi.POST("/template", vv.uploadCloneVmTemplate) // 上传克隆虚拟机模版

	vmClusterApi.POST("/migrate-local-vm", vv.migrateLocalVM)   // 不知道
	vmClusterApi.POST("/relocate-local-vm", vv.relocateLocalVM) // 迁移虚拟机
	vmClusterApi.GET("/clone-vm/message", vv.cloneVmMessage)
}

// 获取vc集群列表
func (vv *VMView) getClusters(ctx *gin.Context) {
	vcResult, err := vmService.NewVmService().GetClusters()
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, "获取vc集群账户列表信息失败.")
		return
	}
	if len(vcResult) == 0 {
		middleware.Logf.Error("数据库中无集群账户列表信息.")
		tools.Failed(ctx, "数据库中无集群账户列表信息.")
		return
	}
	middleware.Logf.Info("获取vc集群账户列表成功.")
	tools.Success(ctx, vcResult)
}

func (vv *VMView) addCluster(ctx *gin.Context) {
	var vcRequest vmTypes.RequestMgmtVcClusterByObj
	err := tools.JsonDecode(ctx.Request.Body, &vcRequest)
	if err != nil {
		middleware.Logf.Error("vc账户的参数有问题.")
		tools.Failed(ctx, "vc账户的参数有问题.")
		return
	}
	fmt.Println("vcRequest: ", vcRequest)
	vcResult, err := vmService.NewVmService().AddCluster(vcRequest.Data)
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, "新加vcenter集群账户失败.")
		return
	}
	middleware.Logf.Info("vc账户添加成功了.")
	tools.Success(ctx, vcResult)

}

func (vv *VMView) delCluster(ctx *gin.Context) {
	var vcRequest vmTypes.RequestClusterById
	err := tools.JsonDecode(ctx.Request.Body, &vcRequest)
	if err != nil {
		middleware.Logf.Error("vc账户的参数有问题.")
		tools.Failed(ctx, "vc账户的参数有问题.")
		return
	}
	fmt.Println("vcRequest: ", vcRequest)
	err = vmService.NewVmService().DelCluster(vcRequest.Data.Id)
	if err != nil {
		middleware.Logf.Error("删除vcenter集群账户失败, ", err.Error())
		tools.Failed(ctx, "删除vcenter集群账户失败.")
		return
	}
	middleware.Logf.Info("vc账户删除成功了.")
	tools.Success(ctx, "")
}

func (vv *VMView) getCluster(ctx *gin.Context) {
	id := ctx.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, "id参数有问题")
		return
	}
	vcAccount, err := vmService.NewVmService().GetCluster(id64)
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, fmt.Sprintf("id 为 %s vc账户查询失败", id))
		return
	}
	tools.Success(ctx, vcAccount)
}

func (vv *VMView) loginCluster(ctx *gin.Context) {
	var vcRequest vmTypes.RequestClusterById
	err := tools.JsonDecode(ctx.Request.Body, &vcRequest)
	if err != nil {
		middleware.Logf.Error("vc ID的参数有问题.")
		tools.Failed(ctx, "vc ID的参数有问题.")
		return
	}
	//id64, err := strconv.ParseInt(id, 10, 64)
	//if err != nil {
	//	middleware.Logf.Error(err)
	//	tools.Failed(ctx, "id参数有问题")
	//	return
	//}
	err = vmService.NewVmService().LoginCluster(vcRequest.Data.Id)
	if err != nil {
		middleware.Logf.Errorf("vc ID的参数有问题. err: %s", err)
		tools.Failed(ctx, "登录失败, err: "+err.Error())
		return
	}
	tools.Success(ctx, "")
}

func (vv *VMView) getHosts(ctx *gin.Context) {
	hosts, err := vmService.NewVmService().GetAllHosts()
	if err != nil {
		middleware.Logf.Error(err.Error())
		tools.Failed(ctx, err.Error())
		return
	}
	middleware.Logf.Info("获取所有宿主机成功.")
	tools.Success(ctx, hosts)
}

func (vv *VMView) getHost(ctx *gin.Context) {
	hostQuery := ctx.Query("host")
	fmt.Println("hostQuery", hostQuery)
	if len(hostQuery) == 0 {
		middleware.Logf.Error("未传入正确的host ip")
		tools.Failed(ctx, "未传入正确的host ip")
		return
	}
	host, err := vmService.NewVmService().GetHost(hostQuery)
	if err != nil {
		middleware.Logf.Errorf("获取宿主机 %s 主机信息失败, error: %s", hostQuery, err.Error())
		tools.Failed(ctx, err.Error())
		return
	}
	middleware.Logf.Infof("获取宿主机 %s 成功.", hostQuery)
	tools.Success(ctx, *host)
}

func (vv *VMView) getVms(ctx *gin.Context) {
	allVms, err := vmService.NewVmService().GetAllVms()
	if err != nil {
		middleware.Logf.Error(err.Error())
		tools.Failed(ctx, err.Error())
		return
	}
	middleware.Logf.Info("获取所有虚拟机成功.")
	tools.Success(ctx, allVms)
}

func (vv *VMView) getCloneVmMetaData(ctx *gin.Context) {
	res, err := vmService.NewVmService().GetCloneVmMetaData()
	if err != nil {
		middleware.Logf.Error(err.Error())
		tools.Failed(ctx, err.Error())
		return
	}
	middleware.Logf.Info("获取所有虚拟机与模版成功.")
	tools.Success(ctx, res)
}

func (vv *VMView) getCloneVmJobs(ctx *gin.Context) {
	jobs, err := vmService.NewVmService().GetCloneVmJobs()
	if err != nil {
		middleware.Logf.Errorf("获取虚拟机任务失败, error: %s", err.Error())
		tools.Failed(ctx, err.Error())
		return
	}
	middleware.Logf.Info("获取虚拟机任务成功.")
	tools.Success(ctx, jobs)
}

func (vv *VMView) createCloneVmJob(ctx *gin.Context) {
	var cloneVmDatas []vmTypes.CloneVmData
	err := tools.JsonDecode(ctx.Request.Body, &cloneVmDatas)
	if err != nil {
		middleware.Logf.Error("克隆虚拟机参数有问题.")
		tools.Failed(ctx, "克隆虚拟机参数有问题.")
		return
	}
	jobs, err := vmService.NewVmService().CreateCloneVmJob(cloneVmDatas)
	if err != nil {
		middleware.Logf.Errorf("创建克隆任务有问题, error: %s", err.Error())
		tools.Failed(ctx, "创建克隆任务有问题.")
		return
	}
	tools.Success(ctx, jobs)
}

// 克隆虚拟机
func (vv *VMView) cloneVm(ctx *gin.Context) {
	var jobIds vmTypes.RequestCloneVmJob
	err := tools.JsonDecode(ctx.Request.Body, &jobIds)
	if err != nil {
		middleware.Logf.Error("克隆虚拟机任务id参数有问题.")
		tools.Failed(ctx, "克隆虚拟机任务id参数有问题.")
		return
	}
	err = vmService.NewVmService().CloneVm(jobIds.Ids)
	if err != nil {
		middleware.Logf.Errorf("克隆虚拟机任务有问题, error: %s", err.Error())
		tools.Failed(ctx, "克隆虚拟机任务有问题.")
		return
	}

	tools.Success(ctx, "")
}

func (vv *VMView) getCloneVmTemplate(ctx *gin.Context) {
	fmt.Println("getCloneVmTemplate")
	//projectPath, _ := os.Getwd()
	//filepath := projectPath + "\\file\\cloneVmTemplate.xlsx"
	filepath := "file\\download\\cloneVmTemplate.xlsx"
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename=cloneVmTemplate.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File(filepath)
}

func (vv *VMView) uploadCloneVmTemplate(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		middleware.Logf.Infof("文件上传参数有误, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("文件上传参数有误, error: %v", err.Error()))
		return
	}
	fileExt := strings.Split(file.Filename, ".")
	filename := fileExt[0] + "-" + time.Now().Format("2006-01-02-15-04-05") + "." + fileExt[1]
	filePath := "file\\upload\\" + filename
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		middleware.Logf.Infof("文件上传失败, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("文件上传失败, error: %v", err.Error()))
		return
	}
	openFile, err := file.Open()
	if err != nil {
		middleware.Logf.Infof("文件打开失败, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("文件打开失败, error: %v", err.Error()))
		return
	}
	xlsx, err := excelize.OpenReader(openFile)
	if err != nil {
		middleware.Logf.Infof("表格打开失败, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("表格打开失败, error: %v", err.Error()))
		return
	}
	defer xlsx.Close()
	rows, err := xlsx.GetRows("Sheet1", excelize.Options{RawCellValue: true})
	if err != nil {
		middleware.Logf.Infof("获取表格数据失败, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("获取表格数据失败, error: %v", err.Error()))
		return
	}
	vs := vmService.NewVmService()
	cloneVmDatas, err := vs.ParseCloneVmXlsx(rows)
	if err != nil {
		middleware.Logf.Infof("解析表格数据失败, error: %v", err.Error())
		tools.Failed(ctx, fmt.Sprintf("解析表格数据失败, error: %v", err.Error()))
		return
	}
	fmt.Println("cloneVmDatas", cloneVmDatas)
	jobs, err := vs.CreateCloneVmJob(cloneVmDatas)
	if err != nil {
		middleware.Logf.Errorf("创建克隆任务有问题, error: %s", err.Error())
		tools.Failed(ctx, "创建克隆任务有问题.")
		return
	}

	tools.Success(ctx, jobs)

}

// 不知道干嘛用的
func (vv *VMView) migrateLocalVM(ctx *gin.Context) {
	var migrateVMDatas []vmTypes.MigrateVMData
	err := tools.JsonDecode(ctx.Request.Body, &migrateVMDatas)
	if err != nil {
		middleware.Logf.Error("迁移虚拟机参数有问题.")
		tools.Failed(ctx, "迁移虚拟机参数有问题.")
		return
	}
	service := vmService.NewVmService()
	for _, migrateVMData := range migrateVMDatas {
		err := service.MigrateVM(migrateVMData)
		if err != nil {
			middleware.Logf.Error("迁移失败", err.Error())
			tools.Failed(ctx, err.Error())
			return
		}

	}
	tools.Success(ctx, "")
}

// 迁移虚拟机
func (vv *VMView) relocateLocalVM(ctx *gin.Context) {
	var relocateVMDatas []vmTypes.RelocateVMData
	err := tools.JsonDecode(ctx.Request.Body, &relocateVMDatas)
	if err != nil {
		middleware.Logf.Error("迁移虚拟机参数有问题.")
		tools.Failed(ctx, "迁移虚拟机参数有问题.")
		return
	}
	service := vmService.NewVmService()
	for _, relocateVMData := range relocateVMDatas {
		err = service.RelocateVM(relocateVMData)
		if err != nil {
			middleware.Logf.Error("迁移失败", err.Error())
			tools.Failed(ctx, err.Error())
			return
		}
	}
	tools.Success(ctx, "迁移成功")
}

func (vv *VMView) cloneVmMessage(ctx *gin.Context) {
	// 可以接受socket参数
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for message := range Variable.CloneVmMessageChan {
		fmt.Println("message", message)
		err = ws.WriteJSON(map[string]interface{}{
			"code":    tools.SUCCESS,
			"message": "成功",
			"data":    message,
		})
		if err != nil {
			break
		}
	}
}
