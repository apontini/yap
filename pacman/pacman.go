package pacman

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/packagefoundation/yap/constants"
	"github.com/packagefoundation/yap/pack"
	"github.com/packagefoundation/yap/utils"
)

type Pacman struct {
	Pack      *pack.Pack
	pacmanDir string
}

func (p *Pacman) getDepends() (err error) {
	if len(p.Pack.MakeDepends) == 0 {
		return
	}

	args := []string{
		"-S",
		"--noconfirm",
	}
	args = append(args, p.Pack.MakeDepends...)

	err = utils.Exec("", "pacman", args...)
	if err != nil {
		return
	}

	return
}

func (p *Pacman) getUpdates() (err error) {
	err = utils.Exec("", "pacman", "-Sy")
	if err != nil {
		return
	}

	return
}

func (p *Pacman) createInstall() (exists bool, err error) {
	data := ""

	if len(p.Pack.PreInst) > 0 {
		data += "pre_install() {\n"
		for _, line := range p.Pack.PreInst {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	if len(p.Pack.PostInst) > 0 {
		data += "post_install() {\n"
		for _, line := range p.Pack.PostInst {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	if len(p.Pack.PreInst) > 0 {
		data += "pre_upgrade() {\n"
		for _, line := range p.Pack.PreInst {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	if len(p.Pack.PostInst) > 0 {
		data += "post_upgrade() {\n"
		for _, line := range p.Pack.PostInst {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	if len(p.Pack.PreRm) > 0 {
		data += "pre_remove() {\n"
		for _, line := range p.Pack.PreRm {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	if len(p.Pack.PostRm) > 0 {
		data += "post_remove() {\n"
		for _, line := range p.Pack.PostRm {
			data += fmt.Sprintf("    %s\n", line)
		}

		data += "}\n"
	}

	exists = len(data) > 0
	if exists {
		path := filepath.Join(p.pacmanDir, p.Pack.PkgName+".install")
		err = utils.CreateWrite(path, data)

		if err != nil {
			return
		}
	}

	return exists, err
}

func (p *Pacman) createMake() (err error) {
	path := filepath.Join(p.pacmanDir, "PKGBUILD")

	installExists, err := p.createInstall()
	if err != nil {
		return
	}

	data := ""
	data += fmt.Sprintf("# Maintainer: %s\n\n", p.Pack.Maintainer)
	data += fmt.Sprintf("pkgname=%s\n", strconv.Quote(p.Pack.PkgName))
	data += fmt.Sprintf("pkgver=%s\n", strconv.Quote(p.Pack.PkgVer))
	data += fmt.Sprintf("pkgrel=%s\n", strconv.Quote(p.Pack.PkgRel))
	data += fmt.Sprintf("pkgdesc=%s\n", strconv.Quote(p.Pack.PkgDesc))
	data += fmt.Sprintf("arch=(%s)\n",
		strconv.Quote(convertPacman(p.Pack.Arch)))

	data += "license=(\n"
	for _, item := range p.Pack.License {
		data += fmt.Sprintf("    %s\n", strconv.Quote(item))
	}

	data += ")\n"
	data += fmt.Sprintf("url=%s\n", strconv.Quote(p.Pack.URL))

	if len(p.Pack.Depends) > 0 {
		data += "depends=(\n"
		for _, item := range p.Pack.Depends {
			data += fmt.Sprintf("    %s\n", strconv.Quote(item))
		}

		data += ")\n"
	}

	if len(p.Pack.OptDepends) > 0 {
		data += "optdepends=(\n"
		for _, item := range p.Pack.OptDepends {
			data += fmt.Sprintf("    %s\n", strconv.Quote(item))
		}

		data += ")\n"
	}

	if len(p.Pack.Provides) > 0 {
		data += "provides=(\n"
		for _, item := range p.Pack.Provides {
			data += fmt.Sprintf("    %s\n", strconv.Quote(item))
		}

		data += ")\n"
	}

	if len(p.Pack.Conflicts) > 0 {
		data += "conflicts=(\n"
		for _, item := range p.Pack.Conflicts {
			data += fmt.Sprintf("    %s\n", strconv.Quote(item))
		}

		data += ")\n"
	}

	if installExists {
		data += fmt.Sprintf("install=%s\n",
			strconv.Quote(p.Pack.PkgName+".install"))
	}

	data += "options=(\"emptydirs\")\n"

	if len(p.Pack.Backup) > 0 {
		data += "backup=(\n"

		for _, item := range p.Pack.Backup {
			item = strings.TrimPrefix(item, "/")
			data += fmt.Sprintf("    %s\n", strconv.Quote(item))
		}

		data += ")\n"
	}

	data += "package() {\n"
	data += fmt.Sprintf("    rsync -a -A %s/ ${pkgdir}/\n",
		p.Pack.PackageDir)
	data += "}\n"

	err = utils.CreateWrite(path, data)
	if err != nil {
		return
	}

	fmt.Println(data)

	return err
}

func (p *Pacman) pacmanBuild() (err error) {
	err = utils.ChownR(p.pacmanDir, "nobody", "nobody")
	if err != nil {
		return
	}

	err = utils.ChownR(p.Pack.PackageDir, "nobody", "nobody")
	if err != nil {
		return
	}

	err = utils.Exec(p.pacmanDir, "sudo", "-u", "nobody", "makepkg")
	if err != nil {
		return
	}

	return
}

func (p *Pacman) Prep() (err error) {
	err = p.getDepends()
	if err != nil {
		return
	}

	return
}

func (p *Pacman) Update() (err error) {
	err = p.getUpdates()
	if err != nil {
		return
	}

	return
}

func (p *Pacman) makeDirs() (err error) {
	p.pacmanDir = filepath.Join(p.Pack.Root, "pacman")

	err = utils.ExistsMakeDir(p.pacmanDir)
	if err != nil {
		return
	}

	return
}

func (p *Pacman) clean() (err error) {
	if !constants.CleanPrevious {
		return
	}

	pkgPaths, err := utils.FindExt(p.Pack.Home, ".pkg.tar.zst")
	if err != nil {
		return
	}

	for _, pkgPath := range pkgPaths {
		_ = utils.Remove(pkgPath)
	}

	return
}

func (p *Pacman) copy(destination string) (err error) {
	pkgs, err := utils.FindExt(p.pacmanDir, ".pkg.tar.zst")
	if err != nil {
		return
	}

	for _, pkg := range pkgs {
		pkgDestination := filepath.Join(p.Pack.Home, destination)

		err = utils.ExistsMakeDir(pkgDestination)
		if err != nil {
			return
		}

		err = utils.CopyFile("", pkg, pkgDestination, false)
		if err != nil {
			return
		}
	}

	return
}

func (p *Pacman) remDirs() {
	os.RemoveAll(p.pacmanDir)
}

func (p *Pacman) Build(outputDir string) ([]string, error) {
	err := p.makeDirs()
	if err != nil {
		return nil, err
	}

	defer p.remDirs()

	err = p.createMake()
	if err != nil {
		return nil, err
	}

	err = p.pacmanBuild()
	if err != nil {
		return nil, err
	}

	err = p.clean()
	if err != nil {
		return nil, err
	}

	err = p.copy(outputDir)
	if err != nil {
		return nil, err
	}

	pkgs, err := utils.FindExt(p.pacmanDir, ".pkg.tar.zst")
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func (p *Pacman) Install() error {
	pkgs, err := utils.FindExt(p.pacmanDir, ".pkg.tar.zst")
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		absPath, err := filepath.Abs(pkg)
		if err != nil {
			return err
		}

		if err := utils.Exec("", "sudo", "-u", "root", "pacman -U --noconfirm", absPath); err != nil {
			return err
		}
	}

	return nil
}
