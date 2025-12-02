package xruntime

import "runtime"

var moduleName, packageName = PackageInfo()

func PackageInfo() (string, string) {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(2, pcs[:])
	//var pc = pcs[0]
	fs := runtime.CallersFrames(pcs[0:1])
	f, _ := fs.Next()
	moduleName, packageName, _ := FuncNameDetailed(f.Function)
	return moduleName, packageName
}

func FuncNameDetailed(name string) (moduleName string, pkgName string, funcName string) {
	// function name is in the form of <module>/<package>.<function>

	// find module name before the last /
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			moduleName = name[:i]
			name = name[i+1:]
			break
		}
	}

	// find package name before the first .
	var j = 0
	for ; j < len(name); j++ {
		if name[j] == '.' {
			break
		}
	}

	// find function name after the last .
	// if no . found, func name is the whole name
	if j == len(name) {
		funcName = name
	} else {
		if j > 1 {
			pkgName = name[:j]
		}
		funcName = name[j+1:]
	}

	return
}
