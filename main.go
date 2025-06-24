package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	projectDir := flag.String("dir", "", "Path to the Java project directory")
	oldName := flag.String("old", "", "Old project name (e.g., com.example.oldproject)")
	newName := flag.String("new", "", "New project name (e.g., com.newcompany.newproject)")
	flag.Parse()

	if *projectDir == "" || *oldName == "" || *newName == "" {
		flag.Usage()
		log.Fatal("All flags --dir, --old, and --new are required.")
	}

	fmt.Printf("Renaming Java project in %s from '%s' to '%s'\n", *projectDir, *oldName, *newName)

	err := filepath.Walk(*projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git and target directories
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "target") {
			return filepath.SkipDir
		}

		// Process .java files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".java") {
			return processJavaFile(path, *oldName, *newName)
		}

		// Process build files (e.g., pom.xml, build.gradle)
		if !info.IsDir() && (info.Name() == "pom.xml" || info.Name() == "build.gradle") {
			return processBuildFile(path, *oldName, *newName)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}

	fmt.Println("Renaming complete. Please verify the changes and rebuild your Java project.")
}

func processJavaFile(filePath, oldName, newName string) error {
	fmt.Printf("Processing Java file: %s\n", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	originalContent := string(content)
	modifiedContent := originalContent

	// Replace package declarations
	modifiedContent = strings.ReplaceAll(modifiedContent, "package "+oldName, "package "+newName)
	modifiedContent = strings.ReplaceAll(modifiedContent, "import "+oldName, "import "+newName)

	// Determine old and new base package paths for file renaming
	oldPackagePath := strings.ReplaceAll(oldName, ".", string(filepath.Separator))
	newPackagePath := strings.ReplaceAll(newName, ".", string(filepath.Separator))

	// Attempt to derive old and new class names based on common patterns
	// This is a simplified approach and might need refinement for complex cases

	// If the oldName is a package, and the class name is part of it, try to infer new class name
	if strings.HasPrefix(oldName, newName) { // e.g., old: com.foo.bar.MyClass, new: com.foo.bar
		// This case is tricky, might need more sophisticated parsing
	} else {
		// Simple case: oldName is a full package + class name, newName is a full package + class name
		// This part needs to be more robust. For now, focus on package replacement.
	}

	// Replace class name occurrences (this is very basic and might replace too much)
	// A more robust solution would involve parsing the Java code.
	// For now, let's focus on package and file path renaming.
	// modifiedContent = strings.ReplaceAll(modifiedContent, oldClassName, newClassName)

	if modifiedContent != originalContent {
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
		fmt.Printf("Updated content of %s\n", filePath)
	}

	// Rename file path if package path changes
	if strings.Contains(filePath, oldPackagePath) {
		newFilePath := strings.Replace(filePath, oldPackagePath, newPackagePath, 1)
		if newFilePath != filePath {
			// Ensure the new directory exists
			newDir := filepath.Dir(newFilePath)
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				err = os.MkdirAll(newDir, 0755)
				if err != nil {
					return fmt.Errorf("failed to create directory %s: %w", newDir, err)
				}
			}

			err = os.Rename(filePath, newFilePath)
			if err != nil {
				return fmt.Errorf("failed to rename file from %s to %s: %w", filePath, newFilePath, err)
			}
			fmt.Printf("Renamed file from %s to %s\n", filePath, newFilePath)
		}
	}

	return nil
}

func processBuildFile(filePath, oldName, newName string) error {
	fmt.Printf("Processing build file: %s\n", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	originalContent := string(content)
	modifiedContent := originalContent

	// Simple string replacement for build files. This might need more sophisticated XML/Gradle parsing.
	// For Maven pom.xml: groupId, artifactId
	// For Gradle build.gradle: group, artifactId
	modifiedContent = strings.ReplaceAll(modifiedContent, oldName, newName)

	// Attempt to replace common Maven/Gradle artifact/group IDs if oldName is a package
	// This is a heuristic and might not cover all cases.
	oldParts := strings.Split(oldName, ".")
	newParts := strings.Split(newName, ".")

	if len(oldParts) > 0 && len(newParts) > 0 {
		oldArtifact := oldParts[len(oldParts)-1]
		newArtifact := newParts[len(newParts)-1]
		modifiedContent = strings.ReplaceAll(modifiedContent, oldArtifact, newArtifact)

		oldGroup := strings.Join(oldParts[:len(oldParts)-1], ".")
		newGroup := strings.Join(newParts[:len(newParts)-1], ".")
		if oldGroup != "" && newGroup != "" {
			modifiedContent = strings.ReplaceAll(modifiedContent, oldGroup, newGroup)
		}
	}


	if modifiedContent != originalContent {
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
		fmt.Printf("Updated content of %s\n", filePath)
	}

	return nil
}

// Helper function to rename directories
func renameDirectory(oldPath, newPath string) error {
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil // Directory doesn't exist, nothing to rename
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("new directory %s already exists", newPath)
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename directory from %s to %s: %w", oldPath, newPath, err)
	}
	fmt.Printf("Renamed directory from %s to %s\n", oldPath, newPath)
	return nil
}
