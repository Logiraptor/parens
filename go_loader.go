package main

import (
	"text/template"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"plugin"
	"reflect"
	"io"
	"path"
)

const packageTemplate = `
package main

import pkg "{{.Pkg}}"

var Exported = pkg.{{.Name}}

`

func callFuncDynamically(pkg string, name string, args ...interface{}) (results []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic: %v", r)
		}
	}()

	tempDir := path.Join(os.Getenv("GOPATH"), "src", "parens-dynamic")

	tmpl := template.Must(template.New("root").Parse(packageTemplate))
	srcPackage, err := ioutil.TempDir(tempDir, "parens-dynamic")
	if err != nil {
		return nil, fmt.Errorf("tempDir: %s", err.Error())
	}

	os.MkdirAll(srcPackage, 0700)
	basePkgName := path.Base(srcPackage)
	srcFileName := path.Join(srcPackage, "main.go")
	output, err := os.Create(srcFileName)
	if err != nil {
		return nil, fmt.Errorf("open: %s", err.Error())
	}
	err = tmpl.Execute(io.MultiWriter(os.Stdout, output), struct {
		Pkg  string
		Name string
	}{pkg, name})
	if err != nil {
		return nil, err
	}
	err = output.Close()
	if err != nil {
		return nil, err
	}

	pluginfile, err := ioutil.TempFile(tempDir, "parens-dynamic")
	if err != nil {
		return nil, err
	}
	pluginfile.Close()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginfile.Name(), path.Join("parens-dynamic",basePkgName))
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: output:\n %s", err.Error(), buf)
	}

	loadedPlugin, err := plugin.Open(pluginfile.Name())
	if err != nil {
		return nil, err
	}

	exportedValue, err := loadedPlugin.Lookup("Exported")
	if err != nil {
		return nil, err
	}

	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}
	resultValues := reflect.ValueOf(exportedValue).Elem().Call(argValues)
	results = make([]interface{}, len(resultValues))
	for i, resultValue := range resultValues {
		results[i] = resultValue.Interface()
	}

	return results, nil
}
