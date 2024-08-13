package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Ignaciojeria/einar/app/domain"
	"github.com/Ignaciojeria/einar/app/domain/ports/in"
	"github.com/Ignaciojeria/einar/app/shared/utils"
)

var EinarInstall in.EinarInstall = func(ctx context.Context, project, commandName string) error {

	// read einar.cli.json
	cliPath := filepath.Join(".einar.cli.json")
	cliBytes, err := ioutil.ReadFile(cliPath)
	if err != nil {
		return fmt.Errorf("failed to read .einar.cli.json: %v", err)
	}

	var cli domain.EinarCli
	err = json.Unmarshal(cliBytes, &cli)
	if err != nil {
		return fmt.Errorf("failed to unmarshal .einar.cli.json: %v", err)
	}

	tagFolder := ""
	if cli.Template.Tag != "" {
		tagFolder = "/" + cli.Template.Tag
	}

	templateFolderPath, err := utils.GetTemplateFolderPath(cli.Template.URL + tagFolder)
	if err != nil {
		return err
	}

	jsonFilePath := filepath.Join(templateFolderPath, ".einar.template.json")
	jsonContentBytes, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		utils.GitCloneTemplateInBinaryPath(cli.Template.URL, "no-auth", cli.Template.Tag)
		jsonFilePath := filepath.Join(templateFolderPath, ".einar.template.json")
		jsonContentBytes, err = ioutil.ReadFile(jsonFilePath)
	}

	if err != nil {
		return fmt.Errorf("error reading JSON file: %v for project %v", err, project)
	}

	var template domain.EinarTemplate
	err = json.Unmarshal(jsonContentBytes, &template)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON file: %v for project %v", err, project)
	}

	var installCommand domain.InstallationCommand
	for _, command := range template.InstallationCommands {
		if command.Name == commandName {
			installCommand = command
			break
		}
	}

	if installCommand.Name == "" {
		return fmt.Errorf("%s command not found in .einar.template.json", commandName)
	}

	if installCommand.SourceDir != "" && installCommand.DestinationDir != "" {
		installCommand.Folders = append(installCommand.Folders,
			domain.InstallationFolder{
				SourceDir:      installCommand.SourceDir,
				DestinationDir: installCommand.DestinationDir,
				IocDiscovery:   true,
			})
	}

	installationsMap := make(map[string]bool)
	for _, installation := range cli.Installations {
		installationsMap[installation.Name] = true
	}

	// Verificar si las dependencias están presentes
	var dependsOn []string
	for _, dependency := range installCommand.DependsOn {
		if !installationsMap[dependency] {
			dependsOn = append(dependsOn, dependency)
			fmt.Println("einar install " + dependency)
		}
	}

	// Devolver error si hay dependencias faltantes
	if len(dependsOn) > 0 {
		return errors.New("dependencies are not present")
	}

	placeHolders := []string{`"archetype`, "${project}"}
	placeHoldersReplace := []string{`"` + project, project}
	for _, folder := range installCommand.Folders {
		sourceDir := filepath.Join(templateFolderPath, folder.SourceDir)
		destDir := filepath.Join( /*project*/ "", folder.DestinationDir)

		err = utils.CopyDirectory(sourceDir, destDir, placeHolders, placeHoldersReplace)
		if err != nil {
			return fmt.Errorf("error cloning %s directory: %v", commandName, err)
		}

		fmt.Printf("%s directory cloned successfully to %s.\n", commandName, destDir)
		/*
			for _, lib := range installCommand.Libraries {
				cmd := exec.Command("go", "get", lib)
				cmd.Dir = ""
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("error installing %s library %s: %v", commandName, lib, err)
				}
			}
		*/

		if !folder.IocDiscovery {
			continue
		}

		setupFilePath := filepath.Join( /*project*/ "", "main.go")

		err = utils.AddImportStatement(setupFilePath, fmt.Sprintf(project+"/"+folder.SourceDir))
		if err != nil {
			return fmt.Errorf("failed to add import statement to setup.go: %v", err)
		}

		firstLevelDirs, err := utils.ListFirstLevelDirs(sourceDir)
		if err != nil {
			return fmt.Errorf("failed to list first level directories: %v", err)
		}

		for _, v := range firstLevelDirs {
			err = utils.AddImportStatement(setupFilePath, fmt.Sprintf(project+"/"+folder.SourceDir+"/"+v))
			if err != nil {
				return fmt.Errorf("failed to add import statement to setup.go: %v", err)
			}
		}
	}

	for _, file := range installCommand.Files {

		sourceDir := filepath.Join(templateFolderPath, file.SourceFile)
		destDir := filepath.Join( /*project*/ "", file.DestinationDir+"/"+filepath.Base(file.SourceFile))

		err = utils.CopyFile(sourceDir, destDir, placeHolders, placeHoldersReplace)
		if err != nil {
			return fmt.Errorf("error cloning %s directory: %v", commandName, err)
		}

		fmt.Printf("%s directory cloned successfully to %s.\n", commandName, destDir)

		/*
			for _, lib := range installCommand.Libraries {
				cmd := exec.Command("go", "get", lib)
				cmd.Dir = ""
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("error installing %s library %s: %v", commandName, lib, err)
				}
			}
		*/

		if !file.IocDiscovery {
			continue
		}

		setupFilePath := filepath.Join( /*project*/ "", "main.go")

		err = utils.AddImportStatement(setupFilePath, fmt.Sprintf(project+"/"+file.DestinationDir))
		if err != nil {
			return fmt.Errorf("failed to add import statement to setup.go: %v", err)
		}
	}

	if err := addInstallationInsideCli( /*"project"*/ "", commandName); err != nil {
		return fmt.Errorf("failed to update .einar.template.json: %v", err)
	}

	cmd := exec.Command("go", "get")
	cmd.Dir = ""
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error execugint go get for %s : %v", commandName, err)
	}

	return nil
}

func addInstallationInsideCli(project, commandName string) error {

	// read einar.cli.json
	cliPath := filepath.Join(".einar.cli.json")
	cliBytes, err := ioutil.ReadFile(cliPath)
	if err != nil {
		return fmt.Errorf("failed to read .einar.cli.json: %v", err)
	}

	var cli domain.EinarCli
	err = json.Unmarshal(cliBytes, &cli)
	if err != nil {
		return fmt.Errorf("failed to unmarshal .einar.cli.json: %v", err)
	}

	tagFolder := ""
	if cli.Template.Tag != "" {
		tagFolder = "/" + cli.Template.Tag
	}

	templateFolderPath, err := utils.GetTemplateFolderPath(cli.Template.URL + tagFolder)
	if err != nil {
		return err
	}

	// read einar.template.json
	templatePath := filepath.Join(templateFolderPath, ".einar.template.json")

	fmt.Println(templatePath)
	templateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read .einar.template.json: %v", err)
	}

	var template domain.EinarTemplate
	err = json.Unmarshal(templateBytes, &template)
	if err != nil {
		return fmt.Errorf("failed to unmarshal .einar.template.json: %v", err)
	}

	// find the command
	var command domain.InstallationCommand
	for _, cmd := range template.InstallationCommands {
		if cmd.Name == commandName {
			command = cmd
			break
		}
	}

	if command.Name == "" {
		return fmt.Errorf("command %s not found in .einar.template.json", commandName)
	}

	// add the command to the CLI
	cli.Installations = append(cli.Installations, domain.Installation{
		Name:      command.Name,
		Libraries: command.Libraries,
	})

	// write back the updated einar.cli.json
	cliBytes, err = json.MarshalIndent(cli, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal .einar.cli.json: %v", err)
	}

	err = ioutil.WriteFile(cliPath, cliBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write .einar.cli.json: %v", err)
	}

	return nil
}
