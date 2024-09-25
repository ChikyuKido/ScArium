package routes

import (
	"ScArium/common/log"
	"ScArium/internal/backend/middlewares"
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cache      = make(map[string][]byte)
	cacheMutex = &sync.Mutex{}
	debug      = true
)

func getTemplateFiles(directory string) []string {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}

func getCachedContent(path string, filepath string, data any) []byte {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if content, found := cache[path]; found && !debug {
		return content
	}
	var content []byte
	if strings.HasSuffix(filepath, ".html") {
		byteBuffer := bytes.NewBuffer(make([]byte, 0))
		templates := getTemplateFiles("./static/sites/partials")
		t, err := template.ParseFiles(append([]string{filepath}, templates...)...)
		if err != nil {
			log.I.Fatal("Failed to parse template", err)
			return nil
		}
		err = t.Execute(byteBuffer, data)
		if err != nil {
			log.I.Fatal("Failed to execute template", err)
			return nil
		}
		content = byteBuffer.Bytes()
	} else {
		content, _ = os.ReadFile(filepath)
	}
	var compressedContent bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&compressedContent, gzip.BestCompression)
	_, err := writer.Write(content)
	if err != nil {
		return nil
	}
	writer.Close()

	compressedData := compressedContent.Bytes()
	cache[path] = compressedData
	return compressedData
}

func serveDirectory(rootPath string, baseDir string, r *gin.RouterGroup) {
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseDir, path)
			urlPath := rootPath + relativePath
			servePage(urlPath, path, r, nil)
		}
		return nil
	})
}
func servePage(path string, diskPath string, r *gin.RouterGroup, data any) {
	r.GET(path, func(c *gin.Context) {
		content := getCachedContent(path, diskPath, data)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))
		c.Header("Content-Encoding", "gzip")
		c.Data(200, contentType, content)
	})
}

func InitSitesRoutes(r *gin.Engine) {
	normal := r.Group("/")
	adminRegisterRedirect := r.Group("/", middlewares.AdminRegisterRedirectMiddleware())
	authRedirectGroup := adminRegisterRedirect.Group("/", middlewares.AuthMiddleware(true))

	initStatics(normal)
	initAuthPages(adminRegisterRedirect)
	servePage("/accounts", "./static/sites/accounts.html", authRedirectGroup, nil)
	servePage("/", "./static/sites/index.html", authRedirectGroup, nil)
}

func initStatics(r *gin.RouterGroup) {
	serveDirectory("/css/", "./static/css", r)
	serveDirectory("/js/", "./static/js", r)
	serveDirectory("/imgs/", "./static/imgs", r)
}

func initAuthPages(r *gin.RouterGroup) {
	servePage("/auth/login", "./static/sites/auth/login.html", r, nil)
	var register struct{ Admin bool }
	register.Admin = false
	servePage("/auth/register", "./static/sites/auth/register.html", r, register)
	register.Admin = true
	servePage("/auth/adminRegister", "./static/sites/auth/register.html", r, register)
}
