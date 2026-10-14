package seed

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yoshsoftdev/porto/backend/internal/models"
	"github.com/yoshsoftdev/porto/backend/internal/portfolio"
	"gorm.io/gorm"
)

const knowledgeKey = "porto_markdown"

func Run(db *gorm.DB, overridePath string) error {
	md := loadMarkdown(overridePath)
	parsed := parse(md)

	return db.Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&models.Profile{}).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			if err := tx.Create(&parsed.Profile).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&models.SkillGroup{}).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			for i := range parsed.Skills {
				if err := tx.Create(&parsed.Skills[i]).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Model(&models.Project{}).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			for i := range parsed.Projects {
				if err := tx.Create(&parsed.Projects[i]).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Model(&models.Section{}).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			for i := range parsed.Sections {
				if err := tx.Create(&parsed.Sections[i]).Error; err != nil {
					return err
				}
			}
		}

		return tx.Save(&models.Setting{Key: knowledgeKey, Value: md}).Error
	})
}

func loadMarkdown(overridePath string) string {
	candidates := []string{overridePath}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "porto.md"),
			filepath.Join(wd, "..", "porto.md"),
			filepath.Join(wd, "..", "..", "porto.md"),
		)
	}
	for _, p := range candidates {
		if strings.TrimSpace(p) == "" {
			continue
		}
		b, err := os.ReadFile(p)
		if err == nil && len(b) > 0 {
			return string(b)
		}
	}
	return portfolio.Markdown
}

type parsed struct {
	Profile  models.Profile
	Skills   []models.SkillGroup
	Projects []models.Project
	Sections []models.Section
}

var (
	metaRe   = regexp.MustCompile(`(?m)^\*\*([^*]+):\*\*\s*(.+?)\s*$`)
	bulletRe = regexp.MustCompile(`(?m)^- (.+)$`)
	numHead  = regexp.MustCompile(`^(\d+)\.\s+(.+)$`)
)

func parse(md string) parsed {
	out := parsed{
		Profile: models.Profile{
			Title:      "Backend & Integration Engineer",
			TargetRole: "Backend & Integration Engineer",
		},
	}

	sections := splitLevel(md, 1)
	for _, sec := range sections {
		title := strings.TrimSpace(sec.title)
		switch {
		case title == "Portfolio Contact":
			fillContact(&out.Profile, sec.body)
		case strings.HasPrefix(title, "Portfolio"):
			fillProfile(&out.Profile, sec.body)
		case title == "Core Strengths":
			out.Profile.Strengths = bullets(sec.body)
		case title == "Technical Skills":
			out.Skills = parseSkills(sec.body)
		case title == "Projects":
			out.Projects = append(out.Projects, parseProjects(sec.body)...)
		case title == "Selected Projects":
			continue
		case numHead.MatchString(title):
			out.Projects = append(out.Projects, parseOneProject(title, sec.body))
		default:
			out.Sections = append(out.Sections, models.Section{
				Slug:      slugify(title),
				Title:     title,
				Content:   strings.TrimSpace(sec.body),
				SortOrder: len(out.Sections) + 1,
			})
		}
	}
	return out
}

func fillProfile(p *models.Profile, body string) {
	for _, sub := range splitLevel(body, 2) {
		if strings.TrimSpace(sub.title) != "Profile" {
			continue
		}
		text := strings.TrimSpace(strings.ReplaceAll(sub.body, "---", ""))
		p.Summary = strings.TrimSpace(text)
		if line := firstNonEmpty(strings.Split(text, "\n")); strings.HasPrefix(line, "**") {
			if title, _, ok := strings.Cut(strings.TrimPrefix(line, "**"), "**"); ok {
				p.Title = strings.TrimSpace(title)
			}
		}
	}
}

func fillContact(p *models.Profile, body string) {
	for _, m := range metaRe.FindAllStringSubmatch(body, -1) {
		key := strings.ToLower(strings.TrimSpace(m[1]))
		val := strings.TrimSpace(m[2])
		switch key {
		case "name":
			p.Name = val
		case "email":
			p.Email = val
		case "github":
			p.Github = val
		case "linkedin":
			if !strings.HasPrefix(val, "http") {
				val = "https://" + val
			}
			p.Linkedin = val
		}
	}
}

func parseSkills(body string) []models.SkillGroup {
	var groups []models.SkillGroup
	for _, sub := range splitLevel(body, 2) {
		items := bullets(sub.body)
		if len(items) == 0 {
			continue
		}
		groups = append(groups, models.SkillGroup{
			Category:  strings.TrimSpace(sub.title),
			Items:     items,
			SortOrder: len(groups) + 1,
		})
	}
	return groups
}

func parseProjects(body string) []models.Project {
	var projects []models.Project
	for _, sub := range splitLevel(body, 2) {
		if numHead.MatchString(strings.TrimSpace(sub.title)) {
			projects = append(projects, parseOneProject(sub.title, sub.body))
		}
	}
	return projects
}

func parseOneProject(title, body string) models.Project {
	m := numHead.FindStringSubmatch(strings.TrimSpace(title))
	order, _ := strconv.Atoi(m[1])
	name, subtitle := splitName(m[2])
	p := models.Project{
		Slug:      slugify(name),
		Name:      name,
		Subtitle:  subtitle,
		SortOrder: order,
		Content:   strings.TrimSpace(body),
	}
	applyProjectMeta(&p, body)
	p.Summary = overviewSummary(body)
	p.Tech = extractTech(body)
	return p
}

func applyProjectMeta(p *models.Project, body string) {
	known := map[string]*string{
		"period": &p.Period,
		"client": &p.Client,
		"status": &p.Status,
		"role":   &p.Role,
		"team":   &p.Team,
		"lead":   &p.Lead,
	}
	for _, m := range metaRe.FindAllStringSubmatch(body, -1) {
		key := strings.ToLower(strings.TrimSpace(m[1]))
		val := strings.TrimSpace(m[2])
		if dest, ok := known[key]; ok {
			*dest = val
			continue
		}
		p.Extra = append(p.Extra, models.MetaItem{Key: strings.TrimSpace(m[1]), Value: val})
	}
}

func overviewSummary(body string) string {
	for _, sub := range splitLevel(body, 3) {
		if strings.EqualFold(strings.TrimSpace(sub.title), "Overview") {
			for _, para := range strings.Split(sub.body, "\n\n") {
				para = strings.TrimSpace(para)
				if para != "" && !strings.HasPrefix(para, "```") {
					return para
				}
			}
		}
	}
	return ""
}

func extractTech(body string) []string {
	for _, sub := range splitLevel(body, 3) {
		if strings.EqualFold(strings.TrimSpace(sub.title), "Technology") {
			return bullets(sub.body)
		}
	}
	return nil
}

type mdBlock struct {
	title string
	body  string
}

func splitLevel(md string, level int) []mdBlock {
	prefix := strings.Repeat("#", level) + " "
	var blocks []mdBlock
	var cur *mdBlock
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, prefix) && !strings.HasPrefix(line, strings.Repeat("#", level+1)) {
			if cur != nil {
				blocks = append(blocks, *cur)
			}
			cur = &mdBlock{title: strings.TrimSpace(strings.TrimPrefix(line, prefix))}
			continue
		}
		if cur != nil {
			cur.body += line + "\n"
		}
	}
	if cur != nil {
		blocks = append(blocks, *cur)
	}
	return blocks
}

func bullets(s string) []string {
	var out []string
	for _, m := range bulletRe.FindAllStringSubmatch(s, -1) {
		out = append(out, strings.TrimSpace(m[1]))
	}
	return out
}

func splitName(title string) (string, string) {
	if name, sub, ok := strings.Cut(title, " — "); ok {
		return strings.TrimSpace(name), strings.TrimSpace(sub)
	}
	return strings.TrimSpace(title), ""
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func firstNonEmpty(lines []string) string {
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

func Knowledge(db *gorm.DB) string {
	var s models.Setting
	if err := db.First(&s, "key = ?", knowledgeKey).Error; err != nil {
		return portfolio.Markdown
	}
	if strings.TrimSpace(s.Value) == "" {
		return portfolio.Markdown
	}
	return s.Value
}
