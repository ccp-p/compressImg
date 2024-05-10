package auto

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func DeleteNodeModules(root string) error {
	fmt.Println("Deleting node_modules...")
	var wg sync.WaitGroup
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "node_modules" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("Deleting: ", path)
				_, err := os.Stat(path)
				if err == nil {
					err = os.RemoveAll(path)
					if err != nil {
						fmt.Println("Error deleting: ", path, err)
					}
				}
			}()
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error walking the path: %s, %w", root, err)
	}

	wg.Wait()
	return nil
}
