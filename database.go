package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/buntdb"
)

// DBDir 数据文件地址
const DBDir = "./database"

// NewDB 初始化数据库引擎
func NewDB(schema string) (*buntdb.DB, error) {
	if m, ok := AllModel[schema]; ok {
		return m.GetDB(), nil
	}

	p := DBDir + "/" + schema + ".db"
	return buntdb.Open(p)
}

// DBInit 数据库初始化
func DBInit() {
	logger.Info("开始进行加载数据库 ...")

	// 目录不存在, 则创建
	if _, err := os.Stat(DBDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(DBDir, 0700); err != nil {
				logger.Error("创建数据库目录失败: " + err.Error())
				return
			}
		} else {
			logger.Error("检索数据库目录失败: " + err.Error())
			return
		}
	}

	// 遍历目录下所有 .db 文件
	filepath.Walk(DBDir, func(f string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(f)
		if ext == ".db" {
			_, fname := filepath.Split(filepath.ToSlash(f))
			schema := strings.TrimSuffix(fname, ext)

			if _, err := NewModel(schema); err != nil {
				return err
			}
		}

		return nil
	})
}
