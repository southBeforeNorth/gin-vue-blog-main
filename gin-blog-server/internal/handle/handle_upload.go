package handle

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/upload"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

// @Summary 上传文件
// @Description 上传文件
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 0 {object} Response[string]
// @Router /upload/file [post]
func (*Upload) UploadFile(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}
	slog.Debug(fmt.Sprintf("UploadFile start:%s", fileHeader.Filename))
	if fileHeader.Size > 128*1024*1024 { // 128MB 限制
		ReturnError(c, g.ErrFileReceive, errors.New("文件过大，最大支持128MB"))
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	needChangeMap := map[string]bool{
		".dng":  true,
		".heic": true,
		".heif": true,
		".svg":  true,
		".jp2k": true,
		".tiff": true,
		".webp": true,
	}
	if needChangeMap[ext] {
		// 检查 VIPS 是否已初始化
		if !utils.IsVipsInitialized() {
			slog.Error("VIPS not initialized")
			ReturnError(c, g.ErrFileReceive, errors.New("图片处理服务未初始化"))
			return
		}
		// 获取转换器实例（不需要初始化和关闭）
		converter := utils.NewImageConverter(100)
		fileHeader, err = converter.ConvertToJPEG(fileHeader)
		if err != nil {
			slog.Error(fmt.Sprintf("UploadFile ConvertToJPEG err:%s", err.Error()))
			ReturnError(c, g.ErrFileReceive, err)
			return
		}
	}
	oss := upload.NewOSS()
	filePath, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}

	ReturnSuccess(c, filePath)
}
