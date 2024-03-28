package entity

import (
	"time"
)

// Directory

type Content struct {
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	LastModified *time.Time  `json:"last_modified"`
	Created      *time.Time  `json:"created"`
	Content      interface{} `json:"content"` // It can be a string, struct, etc., based on the actual content
	Format       string      `json:"format"`
	MimeType     string      `json:"mimetype"`
	Size         *int        `json:"size"`
	Writable     bool        `json:"writable"`
	Type         string      `json:"type"`
}

type Directory struct {
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	LastModified *time.Time `json:"last_modified"`
	Created      *time.Time `json:"created"`
	Content      []Content  `json:"content"`
	Format       string     `json:"format"`
	MimeType     string     `json:"mimetype"`
	Size         *int       `json:"size"`
	Writable     bool       `json:"writable"`
	Type         string     `json:"type"`
}

// Notebook

type Metadata struct {
	Kernelspec struct {
		DisplayName string `json:"display_name"`
		Language    string `json:"language"`
		Name        string `json:"name"`
	} `json:"kernelspec"`
	LanguageInfo struct {
		CodemirrorMode struct {
			Name    string `json:"name"`
			Version int    `json:"version"`
		} `json:"codemirror_mode"`
		FileExtension     string `json:"file_extension"`
		Mimetype          string `json:"mimetype"`
		Name              string `json:"name"`
		NbconvertExporter string `json:"nbconvert_exporter"`
		PygmentsLexer     string `json:"pygments_lexer"`
		Version           string `json:"version"`
	} `json:"language_info"`
}

type CodeCell struct {
	CellType       string                 `json:"cell_type"`
	ExecutionCount *int                   `json:"execution_count"`
	ID             string                 `json:"id"`
	Metadata       map[string]interface{} `json:"metadata"`
	Outputs        []struct {
		Name       string `json:"name"`
		OutputType string `json:"output_type"`
		Text       string `json:"text"`
	} `json:"outputs"`
	Source string `json:"source"`
}

type Notebook struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	LastModified time.Time `json:"last_modified"`
	Created      time.Time `json:"created"`
	Content      struct {
		Cells         []CodeCell `json:"cells"`
		Metadata      Metadata   `json:"metadata"`
		Nbformat      int        `json:"nbformat"`
		NbformatMinor int        `json:"nbformat_minor"`
	} `json:"content"`
	Format   string `json:"format"`
	Mimetype string `json:"mimetype"`
	Size     *int   `json:"size"`
	Writable bool   `json:"writable"`
	Type     string `json:"type"`
}
