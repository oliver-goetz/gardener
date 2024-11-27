// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	nodeagentv1alpha1 "github.com/gardener/gardener/pkg/nodeagent/apis/config/v1alpha1"
)

const (
	FileCreated  FileOperation = "created"
	FileModified FileOperation = "modified"

	oscFilesPath = nodeagentv1alpha1.BaseDir + "/osc-files.yaml"
)

// NodeAgentFileSystem can apply
type NodeAgentFileSystem interface {
	afero.Fs
	GetFileOperation(path string) FileOperation
}

// FileOperation represents the file operation performed by NodeAgentFileSystem.
type FileOperation string

// NewNodeAgentFileSystem creates a new NodeAgentFileSystem.
func NewNodeAgentFileSystem(fs afero.Afero) (NodeAgentFileSystem, error) {
	nodeAgentFilesRaw, err := fs.ReadFile(oscFilesPath)
	if errors.Is(err, afero.ErrFileNotFound) {
		return &nodeAgentFileSystem{
			fs: fs,
			nodeAgentFiles: nodeAgentFiles{
				FileOperation: map[string]FileOperation{},
			},
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to read %q file: %w", oscFilesPath, err)
	}

	nodeAgentFiles := nodeAgentFiles{}
	err = yaml.Unmarshal(nodeAgentFilesRaw, &nodeAgentFiles)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %q file: %w", oscFilesPath, err)
	}

	return &nodeAgentFileSystem{
		fs:             fs,
		nodeAgentFiles: nodeAgentFiles,
	}, nil
}

type nodeAgentFiles struct {
	FileOperation map[string]FileOperation `json:"FileOperation"`
}

type nodeAgentFileSystem struct {
	fs             afero.Afero
	mutex          sync.Mutex
	nodeAgentFiles nodeAgentFiles
}

// GetFileOperation returns the file operation performed on the file.
func (n *nodeAgentFileSystem) GetFileOperation(path string) FileOperation {
	return n.nodeAgentFiles.FileOperation[path]
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func (n *nodeAgentFileSystem) Create(name string) (afero.File, error) {
	operation, err := n.beforeWrite(name)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare file handler %q for writing: %w", name, err)
	}
	file, err := n.fs.Create(name)
	if err != nil {
		return file, err
	}

	if err := n.saveState(name, operation); err != nil {
		return nil, fmt.Errorf("unable to save state %q: %w", name, err)
	}

	return file, nil
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (n *nodeAgentFileSystem) Mkdir(name string, perm os.FileMode) error {
	return n.fs.Mkdir(name, perm)
}

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func (n *nodeAgentFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return n.fs.MkdirAll(path, perm)
}

// Open opens a file, returning it or an error, if any happens.
func (n *nodeAgentFileSystem) Open(name string) (afero.File, error) {
	return n.fs.Open(name)
}

// OpenFile opens a file using the given flags and the given mode.
func (n *nodeAgentFileSystem) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if flag == os.O_RDONLY {
		return n.fs.OpenFile(name, flag, perm)
	}

	operation, err := n.beforeWrite(name)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare file handler %q for writing: %w", name, err)
	}

	file, err := n.fs.OpenFile(name, flag, perm)
	if err != nil {
		return file, err
	}

	if err := n.saveState(name, operation); err != nil {
		return nil, fmt.Errorf("unable to save state %q: %w", name, err)
	}

	return file, nil
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (n *nodeAgentFileSystem) Remove(name string) error {
	if operation := n.nodeAgentFiles.FileOperation[name]; operation != FileCreated {
		return nil
	}

	if err := n.fs.Remove(name); err != nil {
		return err
	}

	if err := n.deleteState(name); err != nil {
		return fmt.Errorf("unable to save state %q: %w", name, err)
	}

	return nil
}

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func (n *nodeAgentFileSystem) RemoveAll(path string) error {
	isDir, err := afero.IsDir(n.fs, path)
	if err != nil {
		return fmt.Errorf("unable to check if %q is a directory: %w", path, err)
	}

	if !isDir {
		if err := n.Remove(path); err != nil && !errors.Is(err, afero.ErrFileNotFound) {
			return fmt.Errorf("unable to remove file %q: %w", path, err)
		}

		return nil
	}

	files, err := afero.ReadDir(n.fs, path)
	if err != nil {
		return fmt.Errorf("unable to read directory %q: %w", path, err)
	}

	for _, file := range files {
		if err := n.Remove(filepath.Join(path, file.Name())); err != nil {
			return fmt.Errorf("unable to remove file %q: %w", filepath.Join(path, file.Name()), err)
		}
	}

	if empty, err := afero.IsEmpty(n.fs, path); err != nil {
		return fmt.Errorf("unable to check if directory %q is empty: %w", path, err)
	} else if empty {
		return n.fs.RemoveAll(path)
	}

	return nil
}

// Rename renames a file.
func (n *nodeAgentFileSystem) Rename(oldname, newname string) error {
	operation, err := n.beforeWrite(newname)
	if err != nil {
		return fmt.Errorf("unable to prepare file handler %q for writing: %w", newname, err)
	}

	if err := n.fs.Rename(oldname, newname); err != nil {
		return err
	}

	if err := n.saveState(newname, operation); err != nil {
		return fmt.Errorf("unable to save state %q: %w", newname, err)
	}

	if err := n.deleteState(oldname); err != nil {
		return fmt.Errorf("unable to save state %q: %w", oldname, err)
	}

	return nil
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func (n *nodeAgentFileSystem) Stat(name string) (os.FileInfo, error) {
	return n.fs.Stat(name)
}

// The name of this FileSystem
func (n *nodeAgentFileSystem) Name() string {
	return n.fs.Name()
}

// Chmod changes the mode of the named file to mode.
func (n *nodeAgentFileSystem) Chmod(name string, mode os.FileMode) error {
	return n.fs.Chmod(name, mode)
}

// Chown changes the uid and gid of the named file.
func (n *nodeAgentFileSystem) Chown(name string, uid, gid int) error {
	return n.fs.Chown(name, uid, gid)
}

// Chtimes changes the access and modification times of the named file
func (n *nodeAgentFileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return n.fs.Chtimes(name, atime, mtime)
}

func (n *nodeAgentFileSystem) beforeWrite(path string) (*FileOperation, error) {
	if _, ok := n.nodeAgentFiles.FileOperation[path]; ok {
		return nil, nil
	}

	exists, err := n.fs.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("unable to check if file %q exists: %w", path, err)
	}

	operation := FileCreated
	if exists {
		operation = FileModified
	}

	return &operation, nil
}

func (n *nodeAgentFileSystem) deleteState(path string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	delete(n.nodeAgentFiles.FileOperation, path)

	return n.marshallStateAndSave()
}

func (n *nodeAgentFileSystem) saveState(path string, operation *FileOperation) error {
	if operation == nil {
		return nil
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.nodeAgentFiles.FileOperation[path] = *operation

	return n.marshallStateAndSave()
}

func (n *nodeAgentFileSystem) marshallStateAndSave() error {
	nodeAgentFilesRaw, err := yaml.Marshal(n.nodeAgentFiles)
	if err != nil {
		return fmt.Errorf("unable to marshal %q file: %w", oscFilesPath, err)
	}

	if err = n.fs.WriteFile(oscFilesPath, nodeAgentFilesRaw, 0600); err != nil {
		return fmt.Errorf("unable to write %q file: %w", oscFilesPath, err)
	}

	return nil
}
