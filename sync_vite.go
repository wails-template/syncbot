package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	cp "github.com/otiai10/copy"
)

const (
	VITE_REPO_URL       = "https://github.com/vitejs/vite"
	VITE_TEMPLATES_PATH = "packages/create-vite"
)

func SyncVite() error {
	viteDir, err := os.MkdirTemp("", "vite")
	if err != nil {
		return err
	}
	defer os.RemoveAll(viteDir)

	viteRepo, err := NewRepository(viteDir, VITE_REPO_URL)
	if err != nil {
		return err
	}
	err = viteRepo.Clone()
	if err != nil {
		return err
	}

	viteTemplateDir := filepath.Join(viteDir, VITE_TEMPLATES_PATH)
	viteTemplateNames, _ := doublestar.Glob(os.DirFS(viteTemplateDir), "template-*")
	if len(viteTemplateNames) == 0 {
		return ErrNoTemplatesFound
	}

	// TODO: use goroutines
	for _, templateName := range viteTemplateNames {
		templateDir, err := os.MkdirTemp("", templateName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer os.RemoveAll(templateDir)

		templateRepoUrl := TARGET_ORG_URL + "/" + strings.ReplaceAll(templateName, "template", "vite")

		fmt.Println("Current repo:", templateRepoUrl)
		templateRepo, err := NewRepository(templateDir, templateRepoUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// clone
		err = templateRepo.Clone()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// delete all files but .git
		Clean(templateDir, []string{".git"})

		// copy ./template to templateDir
		err = cp.Copy("template", templateDir, cp.Options{
			Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
				return strings.Contains(srcinfo.Name(), "gitkeep"), nil
			},
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		// copy upstream template files to frontend/
		err = cp.Copy(filepath.Join(viteTemplateDir, templateName), filepath.Join(templateDir, "frontend"), cp.Options{
			Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
				skip := false
				skipfiles := []string{
					// NOTE: skip folder not working, throw file does not exist
					".vscode/extensions.json",
					"_gitignore",
					"README.md",
				}

				for _, file := range skipfiles {
					if strings.Contains(src, file) {
						skip = true
						break
					}
				}

				return skip, nil
			},
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		// check repo status is clean
		status, err := templateRepo.Status()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if status.IsClean() {
			fmt.Println("Skip, nothing changed\n")
			continue
		}

		// add all files to git
		err = templateRepo.AddAll()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// commit changes
		err = templateRepo.Commit(fmt.Sprintf("chore: sync vite %s", templateName))
		if err != nil {
			fmt.Println(err)
			continue
		}

		// push changes
		err = templateRepo.Push()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Done!")
		fmt.Println("")
	}

	return nil
}
