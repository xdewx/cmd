/*
 * @Author: leoking
 * @Date: 2022-09-08 11:43:51
 * @LastEditTime: 2022-09-08 19:50:56
 * @LastEditors: leoking
 * @Description:
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	// "dew/orm/types"
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

func main() {
	// flag.NewFlagSet("提取go文件中所有的struct", flag.ExitOnError)
	src := flag.String("i", "", "待处理源文件路径")
	dst := flag.String("o", "", "目标文件路径 缺省时输出到控制台")
	onlyName := flag.Bool("name", false, "只处理name")
	flag.Parse()

	ss, err := MatchStructsInFile(*src)
	if err != nil {
		output(os.Stderr, err.Error())
		return
	}

	fos := os.Stdout
	if len(*dst) > 0 {
		fos, err = os.OpenFile(*dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			output(os.Stderr, err.Error())
			return
		}
	}

	for i, s := range ss {
		var v string = s.Source
		if *onlyName {
			v = s.Name
		}
		output(os.Stdout, fmt.Sprintf("#%d %s\n", i, s.Name))
		output(fos, v, "\n")
	}
}
