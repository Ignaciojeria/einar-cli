package business

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Ignaciojeria/einar/app/domain"
	"github.com/Ignaciojeria/einar/app/domain/ports/in"
	"github.com/Ignaciojeria/einar/app/shared/utils"
)

var EinarInit in.EinarInit = func(ctx context.Context, templateFilePath string, project string) error {
	if err := createInitialFilesFromTemplate(templateFilePath, project); err != nil {
		fmt.Println(err)
		return err
	}
	if err := createInitialDirectoriesFromTemplate(templateFilePath, project); err != nil {
		fmt.Println(err)
		return err
	}
	template, err := utils.ReadEinarTemplateFromBinaryPath(templateFilePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dependencyTree := make([]string, 0)

	for _, installBase := range template.InstallationsBase {
		dependencyTree = append(dependencyTree, installBase.Library)
	}

	if err := initializeGoModule(dependencyTree, project); err != nil {
		return err
	}
	//IMPLEMENT YOUR BUSINESS USECASE HERE
	return nil
}

func createInitialFilesFromTemplate(templateFilePath string, project string) error {

	moduleName, err := utils.ReadTemplateModuleName(templateFilePath)

	if err != nil {
		return fmt.Errorf("error reading template module path")
	}

	fmt.Println(templateFilePath)
	// Construct the path to the JSON file relative to the binary
	jsonFilePath := filepath.Join(templateFilePath, ".einar.template.json")

	// Read the JSON file content
	jsonContentBytes, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v for project %v", err, project)
	}

	// Unmarshal the JSON content into the EinarTemplate struct
	var template domain.EinarTemplate
	err = json.Unmarshal(jsonContentBytes, &template)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON file: %v for project %v", err, project)
	}

	// Extraer el tag del templateFilePath
	normalizedPath := strings.ReplaceAll(templateFilePath, "\\", "/")
	pathParts := strings.Split(normalizedPath, "/")
	if len(pathParts) < 1 {
		return fmt.Errorf("invalid template file path: %s", templateFilePath)
	}
	latestGitTag := pathParts[len(pathParts)-1]

	// Iterate over the Files slice
	for _, file := range template.BaseTemplate.Files {
		// Construct the source and destination paths
		sourcePath := filepath.Join(templateFilePath, file.SourceFile)
		destinationPath := file.DestinationFile

		// Copy the file
		err = utils.CopyFile(sourcePath, destinationPath, []string{`"` + moduleName, "${project}", "${latest-git-tag}"}, []string{`"` + project, project, latestGitTag})
		if err != nil {
			return fmt.Errorf("error copying file from %s to %s: %v for project %v", sourcePath, destinationPath, err, project)
		}

		fmt.Printf("File copied successfully from %s to %s.\n", sourcePath, destinationPath)
	}

	return nil
}

func createInitialDirectoriesFromTemplate(templateFilePath string, project string) error {
	// Construct the path to the JSON file relative to the binary
	jsonFilePath := filepath.Join(templateFilePath, ".einar.template.json")

	// Read the JSON file content
	jsonContentBytes, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v for project %v", err, project)
	}

	// Unmarshal the JSON content into the EinarTemplate struct
	var template domain.EinarTemplate
	err = json.Unmarshal(jsonContentBytes, &template)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON file: %v for project %v", err, project)
	}

	moduleName, err := utils.ReadTemplateModuleName(templateFilePath)

	if err != nil {
		return fmt.Errorf("error reading template module path")
	}

	// Iterate over the Folders slice
	for _, folder := range template.BaseTemplate.Folders {
		// Construct the source and destination paths
		sourceDir := filepath.Join(templateFilePath, folder.SourceDir)
		destinationDir := folder.DestinationDir

		// Copy the directory
		err = utils.CopyDirectory(
			sourceDir, destinationDir,
			[]string{`"` + moduleName, "${project}"},
			[]string{`"` + project, project})

		if err != nil {
			return fmt.Errorf("error copying directory from %s to %s: %v for project %v", sourceDir, destinationDir, err, project)
		}

		fmt.Printf("Directory copied successfully from %s to %s.\n", sourceDir, destinationDir)
	}

	return nil
}

func initializeGoModule(dependencies []string, project string) error {
	// Initialize a new Go module
	goModCmd := exec.Command("go", "mod", "init", project)
	goModCmd.Dir = ""
	err := goModCmd.Run()
	if err != nil {
		err := fmt.Errorf("error initializing go module %s", err)
		fmt.Println(err)
		return err
	}

	goGetCmd := exec.Command("go", "get")
	goGetCmd.Dir = ""
	err = goGetCmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Print success message
	fmt.Printf("Go module '%s' generated successfully with dependencies:", dependencies)
	return nil
}
