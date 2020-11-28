package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/markbates/pkger"
)

// addFilesToZip adds files in a directory to a zip file.
func addFilesToZip(writer *zip.Writer, baseZipPath, baseInZip string) error {
	// Read the files
	files, err := ioutil.ReadDir(baseZipPath)
	if err != nil {
		return err
	}

	// Loop over the files and add them to the zip writer.
	for _, file := range files {
		fileName := path.Join(baseZipPath, file.Name())
		// If the file is a directory we need to add the contents
		// of the directory to the zip.
		if file.IsDir() {
			newCurrentPath := path.Join(baseInZip, file.Name())
			addFilesToZip(writer, fileName, newCurrentPath)
		} else {
			// Add the file to the zip.
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return err
			}

			// Create the file in the zip path.
			zipFilePath := path.Join(baseInZip, file.Name())
			zipFile, err := writer.Create(zipFilePath)
			if err != nil {
				return err
			}

			// Write the file
			_, err = zipFile.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ZipDirectory zips a directory.
func ZipDirectory(pathToDir string, pathToZip string) error {
	// Create the new zip file.
	file, err := os.Create(path.Join(pathToZip))
	if err != nil {
		return fmt.Errorf("Error creating new zip file: %v", err)
	}
	defer file.Close()

	// Create the zip writer.
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Write the files.
	err = addFilesToZip(zipWriter, pathToDir, "")
	if err != nil {
		return fmt.Errorf("Error writing files to zip: %v", err)
	}

	// Close the writer.
	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("Error closing zip writer: %v", err)
	}

	return nil
}

// ZipFile zips a file.
func ZipFile(fileToZipPath string, zipFilePath string, zipFileName string) error {
	// Create the new zip file.
	file, err := os.Create(path.Join(zipFilePath, zipFileName))
	if err != nil {
		return fmt.Errorf("Error creating new zip file: %v", err)
	}
	defer file.Close()

	// Create the zip writer.
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Open file to zip
	fileToZip, err := os.Open(fileToZipPath)
	if err != nil {
		return fmt.Errorf("Error opening file to zip: %v", err)
	}
	defer fileToZip.Close()

	// Get the file info
	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("Error getting file to zip stats: %v", err)
	}

	// Get file header info.
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("Error getting file to zip header info: %v", err)
	}

	// Create the zip file.
	header.Name = zipFileName
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write out the file.
	_, err = io.Copy(writer, fileToZip)
	return err
}

// DoesFileExist checks if a file exists
func DoesFileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsCurrentDirectoryEmpty checks to see if the current working
// directory is empty.
func IsCurrentDirectoryEmpty() (bool, error) {
	// Get the current working directory.
	path, err := os.Getwd()
	if err != nil {
		return false, err
	}

	// Open the current working directory.
	directory, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer directory.Close()

	// Check if the directory is empty.
	_, err = directory.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	// The directory is either not empty or
	// has returned an error.
	return false, err
}

// WriteNewFile writes out a new file with the provided
// content.
func WriteNewFile(pathName string, fileName string, content string) error {
	// Create the new file.
	filePath := path.Join(pathName, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file.
	byteContent := []byte(content)
	_, err = file.Write(byteContent)
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

// CopyPackagedDirectory copies a directory packaged in the binary and its contents to a new location.
func CopyPackagedDirectory(oldDirPath, newDirPath string, exclustionList []string) error {
	// Read the directory.
	dir, err := pkger.Open(oldDirPath)
	if err != nil {
		return fmt.Errorf("Error reading directory: %v", err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return fmt.Errorf("Error reading directory contents: %v", err)
	}

	// Loop over the files and copy them to their new location.
	for _, file := range files {
		// If the file name is in the exlusion list we will skip it.
		fileName := file.Name()
		excludeFile := false
		for _, exlusionName := range exclustionList {
			if fileName == exlusionName {
				excludeFile = true
			}
		}
		if excludeFile {
			continue
		}

		// If the file is directory we need to read and copy
		// the contents.
		if file.IsDir() {
			// Create the new directory.
			oldBase := path.Join(oldDirPath, fileName)
			newDirPath := path.Join(newDirPath, fileName)
			err = CreateNewDirectory(newDirPath)
			if err != nil {
				return fmt.Errorf("Error creating new directory: %v", err)
			}

			err = CopyPackagedDirectory(oldBase, newDirPath, exclustionList)
			if err != nil {
				return fmt.Errorf("Error copy sub directory: %v", err)
			}
		} else {
			oldFilePath := path.Join(oldDirPath, fileName)
			newFilePath := path.Join(newDirPath, fileName)

			// Read the old file.
			oldFile, err := pkger.Open(oldFilePath)
			if err != nil {
				return fmt.Errorf("Error reading file to copy: %v", err)
			}

			contents, err := ioutil.ReadAll(oldFile)
			if err != nil {
				return fmt.Errorf("Error reading packaged file: %v", err)
			}

			// Create the new file.
			newFile, err := os.Create(newFilePath)
			if err != nil {
				return fmt.Errorf("Error creating new file for copying: %v", err)
			}

			// Write the file contents to the new file.
			_, err = newFile.Write(contents)
			if err != nil {
				return fmt.Errorf("Error writing contents to new file: %v", err)
			}
		}
	}

	return nil
}

// CopyDirectory copies a directory and its contents to a new location.
func CopyDirectory(oldDirPath, newDirPath string, exclustionList []string) error {
	// Read the directory.
	files, err := ioutil.ReadDir(oldDirPath)
	if err != nil {
		return fmt.Errorf("Error reading directory: %v", err)
	}

	// Loop over the files and copy them to their new location.
	for _, file := range files {
		// If the file name is in the exlusion list we will skip it.
		fileName := file.Name()
		excludeFile := false
		for _, exlusionName := range exclustionList {
			if fileName == exlusionName {
				excludeFile = true
			}
		}
		if excludeFile {
			continue
		}

		// If the file is directory we need to read and copy
		// the contents.
		if file.IsDir() {
			// Create the new directory.
			oldBase := path.Join(oldDirPath, fileName)
			newDirPath := path.Join(newDirPath, fileName)
			err = CreateNewDirectory(newDirPath)
			if err != nil {
				return fmt.Errorf("Error creating new directory: %v", err)
			}

			err = CopyDirectory(oldBase, newDirPath, exclustionList)
			if err != nil {
				return fmt.Errorf("Error copy sub directory: %v", err)
			}
		} else {
			oldFilePath := path.Join(oldDirPath, fileName)
			newFilePath := path.Join(newDirPath, fileName)

			// Read the old file.
			contents, err := ioutil.ReadFile(oldFilePath)
			if err != nil {
				return fmt.Errorf("Error reading file to copy: %v", err)
			}

			// Create the new file.
			newFile, err := os.Create(newFilePath)
			if err != nil {
				return fmt.Errorf("Error creating new file for copying: %v", err)
			}

			// Write the file contents to the new file.
			_, err = newFile.Write(contents)
			if err != nil {
				return fmt.Errorf("Error writing contents to new file: %v", err)
			}
		}
	}

	return nil
}

// CopyFile copies a file to a different file.
func CopyFile(oldFilePath, newFilePath string) error {
	// Open the old file.
	source, err := os.Open(oldFilePath)
	if err != nil {
		return fmt.Errorf("Error opening source file for copying: %v", err)
	}
	defer source.Close()

	// Create the new file.
	destination, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("Error creating new file for copying: %v", err)
	}
	defer destination.Close()

	// Copy the contents of the old file to the new file.
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("Error copy file: %v", err)
	}

	return nil
}

// TemporaryDirectory is a temporary directory that will
// be cleaned up.
type TemporaryDirectory struct {
	Name string
}

func (tmp *TemporaryDirectory) Create() error {
	err := CreateNewDirectory(tmp.Name)
	if err != nil {
		return err
	}

	return nil
}

func (tmp *TemporaryDirectory) Clean() {
	err := os.RemoveAll(tmp.Name)
	if err != nil {
		fmt.Printf("Error remove %s directory: %v", tmp.Name, err)
	}
}

// CreateNewDirectory creates a new directory in the
// current working directory.
func CreateNewDirectory(name string) error {
	// Get the current working directory.
	pathName, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create the directory.
	dirPath := path.Join(pathName, name)
	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
