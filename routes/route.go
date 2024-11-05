package routes

import (
	"file/files"
	// "fmt"
	// "io"
	"net/http"
	"path/filepath"
	"strings"

	// "fmt"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "file no found"})
			return
		}

		extension := strings.ToLower(filepath.Ext(file.Filename))

		// Aquí puedes manejar los tipos de archivo según la extensión
		switch extension {
		case ".xls", ".xlsx":
			if err := ctx.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "file dont save"})
				return
			}
			files.GetExcelData(file.Filename)
			ctx.JSON(http.StatusOK, gin.H{"fileType": "Excel file"})
		case ".txt":
			if err := ctx.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "file dont save"})
				return
			}
			files.GetCsvData(file.Filename)
			ctx.JSON(http.StatusOK, gin.H{"fileType": "Text file"})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		}

		

		// if mimetype == "text/plain; charset=utf-8" {
		// 	files.GetCsvData(file.Filename)
		// } else {
		// 	files.GetExcelData(file.Filename)
		// }

		ctx.JSON(http.StatusOK, gin.H{"message": "file upload success!"})
	})

	server.Run(":8080")
}
