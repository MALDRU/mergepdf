package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	// HOMEPATH Directorio inicio
	home       = "./"
	namePrefix = "C07_900116413_"
	out        = "./Archivos_Unificados/"
)

func main() {
	fmt.Println("------------------- ING DEVELOPERS -------------------------")
	fmt.Println("--------------------- MERGE PDF ----------------------------")
	fmt.Println("----------------- INICIO DE PROCESO ---------------------")

	if _, err := os.Stat(out); os.IsNotExist(err) {
		os.Mkdir(out, os.ModePerm)
	}

	files, err := ioutil.ReadDir(home)
	if err != nil {
		log.Println(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			log.Println("CARPETA:", file.Name())
			err = mergePDFS(file.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}
	fmt.Println("------------------- FIN DE PROCESO ---------------------")
	p := ""
	fmt.Scanf("%s", &p)
}

func mergePDFS(folder string) error {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", home, folder))
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("------------------ INICIA COMBINADO ", folder, " -----------------------")
	nameFac := namePrefix + folder
	pdfs := make([]string, 1)
	c := 0
	lastFile := ""
	for _, file := range files {
		c++
		fmt.Println("Archivo encontrado: ", file.Name())
		lastFile = file.Name()
		if nameWithoutExt(file.Name()) == nameFac {
			pdfs[0] = folder + "/" + lastFile
		} else {
			pdfs = append(pdfs, folder+"/"+lastFile)
		}
	}
	if c > 1 {
		pdfs = append(pdfs, "output")
		pdfs = append(pdfs, out+nameFac+".pdf")
		cmd := exec.Command("corepdf", pdfs...)
		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput
		err = cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Print(string(cmdOutput.Bytes()))
	} else {
		copy(folder + "/" +lastFile, out+lastFile)
	}

	fmt.Println("------------------ FIN COMBINADO -----------------------")
	return nil
}

func nameWithoutExt(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}

func copy(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
