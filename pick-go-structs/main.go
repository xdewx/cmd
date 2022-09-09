/*
 * @Author: leoking
 * @Date: 2022-09-08 11:43:51
 * @LastEditTime: 2022-09-09 13:56:31
 * @LastEditors: leoking
 * @Description:
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// "dew/orm/types"
	"github.com/xdewx/dew-go/logger"
	"github.com/xdewx/dew-go/orm/types"
)

/*
*
  - @description:
    使用正则匹配出文件中的struct
  - @param {string} filename
  - @return {*}
*/
func MatchStructsInFile(filename string) ([]*types.Struct, error) {
	txt, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	PTN_STRUCT := "type (.+?) struct \\{[\\s\\S]+?\\}\\s+"
	expr := regexp.MustCompile(PTN_STRUCT)
	arrs := expr.FindAllStringSubmatch(string(txt), -1)
	results := make([]*types.Struct, 0)
	for _, v := range arrs {
		results = append(results, &types.Struct{
			Source: v[0],
			Name:   v[1],
		})
	}
	return results, nil
}

/*
*
  - @description:
    将内容输出至指定writer
  - @param {io.Writer} w
  - @param {...string} contents
  - @return {*}
*/
func output(w io.Writer, contents ...string) (int, error) {
	cnt := 0
	for _, c := range contents {
		x, err := w.Write([]byte(c))
		cnt += x
		if err != nil {
			break
		}
	}
	return cnt, nil
}

/*
*
  - @description:
    1. 如果是目录
    1.1 是否递归
    2. 如果是文件
  - @param {Option} op
  - @return {*}
*/
func dispatch(op Option) error {
	op.Input, _ = filepath.Abs(op.Input)

	info, err := os.Stat(op.Input)
	if err != nil {
		return err
	}
	logger.DevI(op.Input)
	// TODO：判断是否是文件夹的更好方式
	if !info.IsDir() {
		return handle(op)
	}
	// 没有使用os.Walk，来手动控制是否递归
	files, err := os.ReadDir(op.Input)
	if err != nil {
		return err
	}
	for _, de := range files {
		tmp := op
		tmp.Input = filepath.Join(op.Input, de.Name())
		if !de.IsDir() || op.Recursive {
			_ = dispatch(tmp)
		}
	}
	return nil
}

/*
*
  - @description:
    1. 从文件中提取
    2. 写入输出流
  - @param {Option} op
  - @return {*}
*/
func handle(op Option) error {
	if !strings.HasSuffix(op.Input, ".go") {
		return nil
	}
	ss, err := MatchStructsInFile(op.Input)
	if err != nil {
		output(os.Stderr, err.Error())
		return err
	}

	fos := os.Stdout
	if len(op.Output) > 0 {
		fos, err = os.OpenFile(op.Output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			output(os.Stderr, err.Error())
			return err
		}
	}

	for i, s := range ss {
		var v string = s.Source
		if op.OnlyName {
			v = s.Name
		}
		output(os.Stdout, fmt.Sprintf("#%d %s\n", i, s.Name))
		output(fos, v, "\n")
	}
	return nil
}

type Option struct {
	Input     string
	Output    string
	OnlyName  bool
	Recursive bool
}

func main() {
	// flag.NewFlagSet("提取go文件中所有的struct", flag.ExitOnError)
	op := Option{}
	flag.StringVar(&op.Input, "i", "", "待处理源文件路径")
	flag.StringVar(&op.Output, "o", "", "目标文件路径 缺省时输出到控制台")
	flag.BoolVar(&op.OnlyName, "name", false, "只处理name")
	flag.BoolVar(&op.Recursive, "r", false, "是否递归处理输入目录下的所有源文件")
	flag.Parse()
	if err := dispatch(op); err != nil {
		output(os.Stderr, err.Error())
	}
}
