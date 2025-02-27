package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/packagefoundation/yap/constants"
)

func ReadFile(path string) (data []byte, err error) {
	data, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to read file '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func ReadDir(path string) (items []os.FileInfo, err error) {
	items, err = ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to read dir '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func MkdirAll(path string) (err error) {
	err = os.MkdirAll(path, 0o755)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to mkdir '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func Chmod(path string, perm os.FileMode) (err error) {
	err = os.Chmod(path, perm)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to chmod '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func ChownR(path string, user, group string) (err error) {
	err = Exec("",
		"chown",
		"-R",
		fmt.Sprintf("%s:%s", user, group),
		path,
	)

	if err != nil {
		return
	}

	return
}

func Remove(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to remove '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func RemoveAll(path string) (err error) {
	err = os.RemoveAll(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to remove '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func ExistsMakeDir(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = MkdirAll(path)
			if err != nil {
				return
			}
		} else {
			fmt.Printf("%s❌ :: %sfailed to stat '%s'%s\n",
				string(constants.ColorBlue),
				string(constants.ColorYellow),
				path,
				string(constants.ColorWhite))

			return
		}

		return
	}

	return
}

func Create(path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to create '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func CreateWrite(path string, data string) (err error) {
	file, err := Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to write to file '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func Open(path string) (file *os.File, err error) {
	file, err = os.Open(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to open file '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func Move(source, dest string) (err error) {
	err = Exec("", "mv", source, dest)

	if err != nil {
		return
	}

	return
}

func Copy(dir, source, dest string, presv bool) (err error) {
	args := []string{"-r", "-T", "-f"}

	if presv {
		args = append(args, "-p")
	}

	args = append(args, source, dest)

	err = Exec(dir, "cp", args...)
	if err != nil {
		return
	}

	return
}

func CopyFile(dir, source, dest string, presv bool) (err error) {
	args := []string{"-f"}

	if presv {
		args = append(args, "-p")
	}

	args = append(args, source, dest)

	err = Exec(dir, "cp", args...)
	if err != nil {
		return
	}

	return
}

func CopyFiles(source, dest string, presv bool) (err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to read dir '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			source,
			string(constants.ColorWhite))

		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err = CopyFile("", filepath.Join(source, file.Name()), dest, presv)
		if err != nil {
			return
		}
	}

	return
}

func FindExt(path, ext string) (matches []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to read dir '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ext) {
			matches = append(matches, filepath.Join(path, file.Name()))
		}
	}

	return
}

func FindMatch(path, match string) (matches []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to read dir '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), match) {
			matches = append(matches, filepath.Join(path, file.Name()))
		}
	}

	return
}

func Filename(path string) string {
	n := strings.LastIndex(path, "/")
	if n == -1 {
		return path
	}

	return path[n+1:]
}

func GetDirSize(path string) (size int, err error) {
	output, err := ExecOutput("", "du", "-c", "-s", path)
	if err != nil {
		return
	}

	split := strings.Fields(output)

	size, err = strconv.Atoi(split[len(split)-2])
	if err != nil {
		fmt.Printf("%s❌ :: %sfailed to get dir size '%s'%s\n",
			string(constants.ColorBlue),
			string(constants.ColorYellow),
			path,
			string(constants.ColorWhite))

		return
	}

	return
}

func Exists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			fmt.Printf("utils: Exists check error for '%s'\n", path)
			log.Fatal(err)

			return
		}
	} else {
		exists = true
	}

	return
}
