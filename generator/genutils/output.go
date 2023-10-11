//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package genutils

import (
	"bytes"
	"go/format"
	"io"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

// WriteFormattedSourceFile creates a Go source file in a given package, dumps to it the content provided by the `writeContents` function
// and formats the result through go/fmt.
// If formatting cannot be applied (due to some syntax error probably), it returns an error.
func WriteFormattedSourceFile(filename string, ctx *genall.GenerationContext, root *loader.Package, writeContents func(*bytes.Buffer)) {
	buf := new(bytes.Buffer)
	buf.WriteString(`
package ` + root.Name + `
`)

	writeContents(buf)

	outContents, err := format.Source(buf.Bytes())
	if err != nil {
		root.AddError(err)
		return
	}

	fullname := "zz_generated." + filename + ".go"
	outputFile, err := ctx.Open(root, fullname)
	if err != nil {
		root.AddError(err)
		return
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outContents)
	if err != nil {
		root.AddError(err)
		return
	}
	if n < len(outContents) {
		root.AddError(io.ErrShortWrite)
		return
	}
}
