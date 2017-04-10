package main

import (
	"strings"
	"os"
	"path"
	"strconv"
	"io/ioutil"

	"github.com/qiniu/xlog.v1"
	"github.com/qiniu/log.v1"
	"github.com/qiniu/errors"

)
func main ()  {
	xl := xlog.NewDummy()
	pathName := setUp()
	names, err := GetFiles(xl,pathName)
	if err != nil {
		xl.Error("GetFiles", err)
	}
	xl.Info("len of names", len(names))
	xl.Info("names", names)
}

func setUp() string {
	execDirAbsPath, _ := os.Getwd()
	pathName := path.Join(execDirAbsPath, "run")
	logInfo := "123456789"
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			uid := strconv.Itoa(i)
			tbl := strconv.Itoa(j)
			dirName := path.Join(pathName, "2017-02-01-15-20", "130000"+uid, "bucket"+tbl)
			err := os.MkdirAll(dirName, 0744)
			if err != nil {
				log.Fatal(" os.MkdirAll ", err)
			}
			fname := path.Join(dirName, "io.log")
			err = ioutil.WriteFile(fname, []byte(logInfo), 0744)
			if err != nil {
				log.Fatal("ioutil.WriteFile", err)
			}
		}
	}
	return pathName
}

//遍历uploadPath下所有文件夹,找到io.log文件,`<uploadPath>/YYYY-MM-DD-HH-mm/<uid>/<bucket>/io.log`
func GetFiles(xl *xlog.Logger, uploadPath string) ([]string, error) {
	if !strings.HasSuffix(uploadPath, "/") {
		uploadPath = uploadPath + "/"
	}
	xl.Info("uploadPath", uploadPath)
	var fileNames []string
	fi, err := os.Stat(uploadPath)
	if err != nil {
		err = errors.Info(err, "os.Stat", uploadPath).Detail(err)
		return nil, err
	}
	if !fi.IsDir() {
		err = errors.New("not dir")
		err = errors.Info(err, uploadPath).Detail(err)
		xl.Error("fi.IsDir:", uploadPath, err)
		return nil, err
	}
	dir, err := os.Open(uploadPath)
	if err != nil {
		err = errors.Info(err, "os.Open", uploadPath).Detail(err)
		xl.Error("os.Open:", err)
		return nil, err
	}
	dateFolders, err := dir.Readdir(-1)
	if err != nil {
		err = errors.Info(err, "dir.Readdir", dir.Name()).Detail(err)
		xl.Error("dir.Readdir:", err)
		return nil, err
	}
	dir.Close()
	for _, dateFolder := range dateFolders {
		if !dateFolder.IsDir() {
			err = errors.New("not dir")
			err = errors.Info(err, dateFolder).Detail(err)
			xl.Error("dateFolder.IsDir:", dateFolder, err)
			return nil, err
		}
		dateName := path.Join(uploadPath, dateFolder.Name())
		dirDate, err := os.Open(dateName)
		if err != nil {
			err = errors.Info(err, "os.Open", uploadPath).Detail(err)
			xl.Error("os.Open:", err)
			return nil, err
		}

		uidFolders, err := dirDate.Readdir(-1)
		if err != nil {
			err = errors.Info(err, "dir.Readdir", dirDate.Name()).Detail(err)
			xl.Error("dirDate.Readdir:", err)
			return nil, err
		}
		dirDate.Close()
		for _, uidFolder := range uidFolders {
			if !uidFolder.IsDir() {
				err = errors.New("not dir")
				err = errors.Info(err, uidFolder.Name()).Detail(err)
				xl.Error("uidFolder.IsDir:", err)
				return nil, err
			}
			uidName := path.Join(dateName, uidFolder.Name())
			dirUid, err := os.Open(uidName)
			if err != nil {
				err = errors.Info(err, "os.Open", uidName).Detail(err)
				xl.Error("os.Open:", err)
				return nil, err
			}

			tblFolders, err := dirUid.Readdir(-1)
			if err != nil {
				err = errors.Info(err, "dir.Readdir").Detail(err)
				xl.Error("dirUid.Readdir", err)
				return nil, err
			}
			dirUid.Close()
			for _, tblFolder := range tblFolders {
				if !tblFolder.IsDir() {
					err = errors.New("not dir")
					err = errors.Info(err, tblFolder).Detail(err)
					xl.Error("tblFolder.IsDir", err)
					return nil, err
				}
				logName := path.Join(uidName, tblFolder.Name(), "io.log")
				_, err := os.Stat(logName)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					} else {
						xl.Error("os.Stat", err)
						return nil, err
					}
				}

				fileNames = append(fileNames, logName)
			}
		}
	}
	return fileNames, nil
}
