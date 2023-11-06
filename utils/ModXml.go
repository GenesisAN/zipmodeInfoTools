package util

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type ModXml struct {
	GUID        string `xml:"guid" json:"guid,omitempty"`
	Version     string `xml:"version" json:"version,omitempty"`
	Name        string `xml:"name" json:"name,omitempty"`
	Author      string `xml:"author" json:"author,omitempty"`
	Description string `xml:"description" json:"description,omitempty"`
	Website     string `xml:"website" json:"website,omitempty"`
	Game        string `xml:"game" json:"game,omitempty"`
	//BDMD5       string `gorm:"uniqueIndex"`
	Path   string `json:"path,omitempty"`
	Error  string `json:"error,omitempty"`
	Upload bool   `json:"upload,omitempty"`
}

func ReadZip(src string) (ModXml, error) {
	mod := ModXml{}
	//===============BDMD5 Build ===================
	//bdmd5, err := util.GetFileBDMD5(src)
	//if err != nil {
	//	return mod, fail
	//}

	//==========================================
	mod.Path = src
	zr, err := zip.OpenReader(src) //open modzip
	if err != nil {
		mod.Error = fmt.Sprintf("打开zip失败:%s", err.Error())
		return mod, errors.New(fmt.Sprintf("打开zip失败:%s", err.Error()))
	}

	for _, v := range zr.File {
		if v.FileInfo().Name() == "manifest.xml" { //find manifest.xml
			fr, _ := v.Open()
			data, _ := io.ReadAll(fr) //read manifest.xml
			//构造XML的结构体
			err = xml.Unmarshal(data, &mod) //unmarshal xml to struct
			if err != nil {
				mod.Error = fmt.Sprintf("打开zip失败:%s", err.Error())
				return mod, errors.New(fmt.Sprintf("读取zipmod信息失败:%s", err.Error()))
			}
		}
	}
	return mod, err
}
