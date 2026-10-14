package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yoshsoftdev/porto/backend/internal/models"
)

type publicProfile struct {
	models.Profile
	PhotoURL string              `json:"photo_url"`
	Skills   []models.SkillGroup `json:"skills"`
	Sections []models.Section    `json:"sections"`
}

type publicProject struct {
	models.Project
	ImageURLs []string `json:"image_urls"`
}

func (s *Server) GetProfile(c *gin.Context) {
	var p models.Profile
	if err := s.DB.First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not seeded"})
		return
	}
	var skills []models.SkillGroup
	s.DB.Order("sort_order asc, id asc").Find(&skills)
	var sections []models.Section
	s.DB.Where("hidden = ?", false).Order("sort_order asc, id asc").Find(&sections)
	c.JSON(http.StatusOK, publicProfile{
		Profile:  p,
		PhotoURL: s.imageURL(p.PhotoID),
		Skills:   skills,
		Sections: sections,
	})
}

func (s *Server) ListProjects(c *gin.Context) {
	var items []models.Project
	q := s.DB.Preload("Images").Where("hidden = ?", false).Order("sort_order asc, id asc")
	if err := q.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]publicProject, 0, len(items))
	for _, p := range items {
		out = append(out, decorateProject(p))
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) GetProject(c *gin.Context) {
	var p models.Project
	if err := s.DB.Preload("Images").Where("slug = ? AND hidden = ?", c.Param("slug"), false).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, decorateProject(p))
}

func decorateProject(p models.Project) publicProject {
	urls := make([]string, 0, len(p.Images))
	for _, img := range p.Images {
		urls = append(urls, "/api/images/"+strconv.FormatUint(uint64(img.ImageID), 10))
	}
	return publicProject{Project: p, ImageURLs: urls}
}

func (s *Server) GetImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var img models.Image
	if err := s.DB.First(&img, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, img.ContentType, img.Data)
}
