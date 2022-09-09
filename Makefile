
CurMakefile:=$(abspath $(firstword $(MAKEFILE_LIST)))
DirPath:=$(shell dirname $(CurMakefile))
DirName:=$(shell basename $(DirPath))

TempDir:=$(shell mktemp -d)
TempBinPath:=$(TempDir)/$(DirName)
BinPath:=$(DirPath)/$(DirName)

default:
	go build .

# 默认似乎只支持linux，需要探索其他选项
upx:
	go build -o $(TempBinPath) . && rm -f $(BinPath) && upx $(TempBinPath) -o $(BinPath)
gzexe:
	go build -o $(TempBinPath) . && gzexe $(TempBinPath) && mv -f $(TempBinPath) $(BinPath)