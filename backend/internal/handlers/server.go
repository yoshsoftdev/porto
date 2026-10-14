package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoshsoftdev/porto/backend/internal/auth"
	"github.com/yoshsoftdev/porto/backend/internal/config"
	"gorm.io/gorm"
)

type Server struct {
	DB  *gorm.DB
	Cfg config.Config
}

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	s := &Server{DB: db, Cfg: cfg}
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.MaxMultipartMemory = cfg.MaxUploadMB << 20

	r.Use(allowCORS(cfg.CORSOrigins))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	api := r.Group("/api")
	{
		api.GET("/profile", s.GetProfile)
		api.GET("/projects", s.ListProjects)
		api.GET("/projects/:slug", s.GetProject)
		api.GET("/images/:id", s.GetImage)
		api.GET("/chat/providers", s.ListProviders)
		api.POST("/chat", s.limitChat(), s.Chat)

		admin := api.Group("/admin")
		admin.POST("/login", s.Login)
		authd := admin.Group("")
		authd.Use(auth.Middleware(cfg.JWTSecret))
		{
			authd.GET("/me", s.Me)
			authd.GET("/profile", s.AdminGetProfile)
			authd.PUT("/profile", s.AdminUpdateProfile)
			authd.GET("/skills", s.AdminListSkills)
			authd.PUT("/skills", s.AdminReplaceSkills)
			authd.GET("/sections", s.AdminListSections)
			authd.PUT("/sections", s.AdminReplaceSections)
			authd.GET("/projects", s.AdminListProjects)
			authd.POST("/projects", s.AdminCreateProject)
			authd.PUT("/projects/:id", s.AdminUpdateProject)
			authd.DELETE("/projects/:id", s.AdminDeleteProject)
			authd.POST("/projects/:id/images", s.AdminAttachImage)
			authd.DELETE("/projects/:id/images/:imageRelId", s.AdminDetachImage)
			authd.POST("/images", s.AdminUploadImage)
			authd.DELETE("/images/:id", s.AdminDeleteImage)
			authd.GET("/llm", s.AdminGetLLM)
			authd.PUT("/llm", s.AdminUpdateLLM)
		}
	}
	return r
}

func allowCORS(origins []string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, o := range origins {
		allowed[o] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Max-Age", "43200")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *Server) imageURL(id *uint) string {
	if id == nil || *id == 0 {
		return ""
	}
	return "/api/images/" + strconv.FormatUint(uint64(*id), 10)
}
