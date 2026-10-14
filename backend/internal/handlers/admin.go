package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yoshsoftdev/porto/backend/internal/auth"
	"github.com/yoshsoftdev/porto/backend/internal/cryptx"
	"github.com/yoshsoftdev/porto/backend/internal/llm"
	"github.com/yoshsoftdev/porto/backend/internal/models"
)

func (s *Server) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
		return
	}
	if !auth.CheckPassword(s.Cfg.AdminPassword, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	token, err := auth.Issue(s.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "role": "admin"})
}

func (s *Server) AdminGetProfile(c *gin.Context) {
	s.GetProfile(c)
}

func (s *Server) AdminUpdateProfile(c *gin.Context) {
	var p models.Profile
	if err := s.DB.First(&p).Error; err != nil {
		p = models.Profile{}
	}
	var req models.Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Name = req.Name
	p.Title = req.Title
	p.Summary = req.Summary
	p.TargetRole = req.TargetRole
	p.Email = req.Email
	p.Github = req.Github
	p.Linkedin = req.Linkedin
	p.Strengths = req.Strengths
	p.PhotoID = req.PhotoID
	if err := s.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) AdminListSkills(c *gin.Context) {
	var items []models.SkillGroup
	s.DB.Order("sort_order asc, id asc").Find(&items)
	c.JSON(http.StatusOK, items)
}

func (s *Server) AdminReplaceSkills(c *gin.Context) {
	var items []models.SkillGroup
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.DB.Where("1 = 1").Delete(&models.SkillGroup{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range items {
		items[i].ID = 0
		if items[i].SortOrder == 0 {
			items[i].SortOrder = i + 1
		}
		if err := s.DB.Create(&items[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, items)
}

func (s *Server) AdminListSections(c *gin.Context) {
	var items []models.Section
	s.DB.Order("sort_order asc, id asc").Find(&items)
	c.JSON(http.StatusOK, items)
}

func (s *Server) AdminReplaceSections(c *gin.Context) {
	var items []models.Section
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.DB.Where("1 = 1").Delete(&models.Section{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range items {
		items[i].ID = 0
		if items[i].SortOrder == 0 {
			items[i].SortOrder = i + 1
		}
		if err := s.DB.Create(&items[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, items)
}

func (s *Server) AdminListProjects(c *gin.Context) {
	var items []models.Project
	s.DB.Preload("Images").Order("sort_order asc, id asc").Find(&items)
	out := make([]publicProject, 0, len(items))
	for _, p := range items {
		out = append(out, decorateProject(p))
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) AdminCreateProject(c *gin.Context) {
	var p models.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = 0
	p.Images = nil
	if p.Slug == "" {
		p.Slug = slugify(p.Name)
	}
	if err := s.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (s *Server) AdminUpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p models.Project
	if err := s.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var req models.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Slug = req.Slug
	p.Name = req.Name
	p.Subtitle = req.Subtitle
	p.Period = req.Period
	p.Client = req.Client
	p.Status = req.Status
	p.Role = req.Role
	p.Team = req.Team
	p.Lead = req.Lead
	p.Extra = req.Extra
	p.Summary = req.Summary
	p.Tech = req.Tech
	p.Content = req.Content
	p.SortOrder = req.SortOrder
	p.Hidden = req.Hidden
	if err := s.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.DB.Preload("Images").First(&p, p.ID)
	c.JSON(http.StatusOK, decorateProject(p))
}

func (s *Server) AdminDeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s.DB.Where("project_id = ?", id).Delete(&models.ProjectImage{})
	if err := s.DB.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (s *Server) AdminUploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	if file.Size > s.Cfg.MaxUploadMB<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	data, err := io.ReadAll(io.LimitReader(src, s.Cfg.MaxUploadMB<<20+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ct := file.Header.Get("Content-Type")
	if ct == "" || ct == "application/octet-stream" {
		ct = http.DetectContentType(data)
	}
	if !strings.HasPrefix(ct, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only image uploads are allowed"})
		return
	}
	img := models.Image{ContentType: ct, Size: len(data), Data: data}
	if err := s.DB.Create(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":           img.ID,
		"content_type": img.ContentType,
		"size":         img.Size,
		"url":          "/api/images/" + strconv.FormatUint(uint64(img.ID), 10),
	})
}

func (s *Server) AdminDeleteImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s.DB.Where("image_id = ?", id).Delete(&models.ProjectImage{})
	s.DB.Model(&models.Profile{}).Where("photo_id = ?", id).Update("photo_id", nil)
	if err := s.DB.Delete(&models.Image{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (s *Server) AdminAttachImage(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		ImageID   uint   `json:"image_id"`
		Caption   string `json:"caption"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ImageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_id required"})
		return
	}
	var img models.Image
	if err := s.DB.First(&img, req.ImageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}
	rel := models.ProjectImage{ProjectID: uint(pid), ImageID: req.ImageID, Caption: req.Caption, SortOrder: req.SortOrder}
	if err := s.DB.Create(&rel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rel)
}

func (s *Server) AdminDetachImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("imageRelId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := s.DB.Delete(&models.ProjectImage{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type llmView struct {
	Active    string            `json:"active"`
	Providers []llmProviderView `json:"providers"`
}

type llmProviderView struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	HasKey  bool   `json:"has_key"`
	Enabled bool   `json:"enabled"`
}

func (s *Server) AdminGetLLM(c *gin.Context) {
	c.JSON(http.StatusOK, s.llmView())
}

func (s *Server) AdminUpdateLLM(c *gin.Context) {
	var req struct {
		Active    string `json:"active"`
		Providers []struct {
			ID    string `json:"id"`
			Model string `json:"model"`
			Key   string `json:"key"`
		} `json:"providers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Active != "" && !llm.Valid(req.Active) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}
	if req.Active != "" {
		s.putSetting("llm_active", req.Active)
	}
	for _, p := range req.Providers {
		if !llm.Valid(p.ID) {
			continue
		}
		if strings.TrimSpace(p.Model) != "" {
			s.putSetting("llm_"+p.ID+"_model", strings.TrimSpace(p.Model))
		}
		if strings.TrimSpace(p.Key) != "" {
			sealed, err := cryptx.Seal(s.Cfg.EncryptionKey, strings.TrimSpace(p.Key))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			s.putSetting("llm_"+p.ID+"_key", sealed)
		}
	}
	c.JSON(http.StatusOK, s.llmView())
}

func (s *Server) llmView() llmView {
	active := s.getSetting("llm_active")
	if active == "" {
		active = string(llm.Gemini)
	}
	var providers []llmProviderView
	for _, p := range llm.All() {
		id := string(p)
		model := s.getSetting("llm_" + id + "_model")
		if model == "" {
			model = llm.DefaultModel(p)
		}
		has := cryptx.Masked(s.getSetting("llm_" + id + "_key"))
		providers = append(providers, llmProviderView{
			ID:      id,
			Model:   model,
			HasKey:  has,
			Enabled: has,
		})
	}
	return llmView{Active: active, Providers: providers}
}

func (s *Server) getSetting(key string) string {
	var row models.Setting
	if err := s.DB.First(&row, "key = ?", key).Error; err != nil {
		return ""
	}
	return row.Value
}

func (s *Server) putSetting(key, value string) {
	s.DB.Save(&models.Setting{Key: key, Value: value})
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prev := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prev = false
			continue
		}
		if !prev && b.Len() > 0 {
			b.WriteByte('-')
			prev = true
		}
	}
	return strings.Trim(b.String(), "-")
}
