package utils

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
)

// ImageConverter 图像转换器
type ImageConverter struct {
	Quality int
}

var (
	converterInstance *ImageConverter
	converterOnce     sync.Once
	vipsInitialized   bool
	vipsOnce          sync.Once
)

// NewImageConverter 获取单例实例
func NewImageConverter(quality int) *ImageConverter {
	converterOnce.Do(func() {
		converterInstance = &ImageConverter{
			Quality: quality,
		}
	})
	return converterInstance
}

// 全局初始化
func InitializeVips() error {
	var err error
	vipsOnce.Do(func() {
		vips.Startup(&vips.Config{
			ConcurrencyLevel: 1,
			MaxCacheFiles:    0,
			MaxCacheMem:      50 * 1024 * 1024, //最大缓存内存
			MaxCacheSize:     10,               //最大缓存次数
		})
		vipsInitialized = true
	})
	return err
}

// IsVipsInitialized 检查 VIPS 是否已初始化
func IsVipsInitialized() bool {
	return vipsInitialized
}

// ShutdownVips 关闭 VIPS（应该在应用关闭时调用）
func ShutdownVips() {
	if vipsInitialized {
		vips.Shutdown()
	}
}

// ConvertToJPEG 转换为JPEG
func (ic *ImageConverter) ConvertToJPEG(fileHeader *multipart.FileHeader) (*multipart.FileHeader, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
		}
	}(file)

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 创建vips图像
	image, err := vips.NewImageFromBuffer(fileBytes)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	// 自动旋转
	err = image.AutoRotate()
	if err != nil {
		return nil, err
	}

	// 处理透明通道
	if image.HasAlpha() {
		err = image.Flatten(&vips.Color{R: 255, G: 255, B: 255})
		if err != nil {
			return nil, err
		}
	}

	// 导出为JPEG
	jpegParams := vips.NewJpegExportParams()
	jpegParams.Quality = ic.Quality
	jpegBytes, _, err := image.ExportJpeg(jpegParams)
	if err != nil {
		return nil, err
	}

	// 生成新文件名
	baseName := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	newFilename := baseName + ".jpg"

	return createFileHeaderFromBytes(jpegBytes, newFilename)
}

// 从字节数据创建 FileHeader
func createFileHeaderFromBytes(data []byte, filename string) (*multipart.FileHeader, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(128 << 20)
	if err != nil {
		return nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, errors.New("no file found")
	}

	return files[0], nil
}
