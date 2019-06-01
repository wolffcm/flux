package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/llvm"
	gollvm "github.com/llvm-mirror/llvm/bindings/go/llvm"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a web server that accepts flux and produces wasm",
	Long:  "start a web server that accepts flux and produces wasm",
	Args:  cobra.ExactArgs(0),
	RunE:  serveE,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serveE(cmd *cobra.Command, args []string) error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/generated_files/", generatedFilesHandler)
	fmt.Println("Serving on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}

var wasmFiles = map[string][]byte{}

var jsFiles = map[string][]byte{}

var body string = `<!DOCTYPE html>
<html>
<head>
<style>
body {
  background-color: #000;
  color: #e0a5f7;
  font-family: Roboto, Helvetica, sans-serif;
}
header {
  background-color: #595800;
  color: #F0EC0A;
  padding: 30px;
  text-align: center;
  font-size: 35px;
}


* {
  box-sizing: border-box;
}

textarea {
  width: 100%;
  height: 100%;
  border: solid;
  border-color: #7F4DD6;
  white-space: pre;
  font-family: monospace;
  background-color: #261347;
  font-size: 14px;
  color: #CAABFF;
}

/* Create three equal columns that floats next to each other */
.column {
  margin: 10px;
  background-color: #4B248C;
  color: #DBDBDB;
  float: left;
  width: 30.0%;
  padding: 10px;
}

/* Clear floats after the columns */
.row:after {
  content: "";
  display: table;
  clear: both;
}
</style>
</head>
<body>
<header>
<h2>Compile and Execute Flux</h2>
</header>

<div class="row">
<div class="column">
<h3>Input Flux</h3>
{input}
</div>

<div class="column">
<h3>LLVM IR</h3>
{llvm}
</div>

<div class="column">
<h3>Output</h3>
{output}
</div>

</div> 
</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fluxInput string
	if f, ok := r.PostForm["flux"]; ok {
		fluxInput = f[0]
	}

	inputDiv, llvmDiv, outputDiv := processFlux(fluxInput)

	replacer := strings.NewReplacer("{input}", inputDiv, "{llvm}", llvmDiv, "{output}", outputDiv)
	_, _ = replacer.WriteString(w, body)
}

func generatedFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling file request")
	_, file := path.Split(r.URL.String())
	if strings.Contains(file, ".wasm") {
		bytes, ok := wasmFiles[file]
		if !ok {
			http.Error(w, "WASM file not found", http.StatusNotFound)
		}
		_, _ = w.Write(bytes)
		return
	} else if strings.Contains(file, ".js") {
		bytes, ok := jsFiles[file]
		if !ok {
			http.Error(w, "JS file not found", http.StatusNotFound)
		}
		_, _ = w.Write(bytes)
		return
	}

	http.Error(w, "unknown file "+file, http.StatusBadRequest)
}

func processFlux(inputFlux string) (string, string, string) {
	tempDir, err := ioutil.TempDir("", "emcc")
	if err != nil {
		return makeDivs(inputFlux, err.Error(), "")
	}

	defer func() { _ = os.RemoveAll(tempDir) }()

	var sb strings.Builder
	if inputFlux == "" {
		return makeDivs(inputFlux, "", "")
	}

	astPkg, err := flux.Parse(inputFlux)
	if err != nil {
		return makeDivs(inputFlux, err.Error(), "")
	}

	mod, err := llvm.Build(astPkg)
	if err != nil {
		return makeDivs(inputFlux, err.Error(), "")
	}

	addTextArea(&sb, "llvm", mod.String(), true)
	llvmDiv := sb.String()
	sb.Reset()

	filename, err := compileToWASM(tempDir, mod)
	if err != nil {
		return makeDivs(inputFlux, err.Error(), "")
	}

	log.Println("Now serving ", filename+".js", " and ", filename+".wasm")

	addTextArea(&sb, "output", "", true)
	sb.WriteString(`
<script type="text/javascript">
var Module = {
  preRun: [],
  postRun: [],
  print: (function() {
    var element = document.getElementById('output');
    if (element) element.value = ''; // clear browser cache
    return function(text) {
      if (arguments.length > 1) text = Array.prototype.slice.call(arguments).join(' ');
      // These replacements are necessary if you render to raw HTML
      text = text.replace(/&/g, "&amp;");
      text = text.replace(/</g, "&lt;");
      text = text.replace(/>/g, "&gt;");
      text = text.replace('\n', '<br>', 'g');
      console.log(text);
      if (element) {
        element.value += text + "\n";
        element.scrollTop = element.scrollHeight; // focus on bottom
      }
    };
  })(),
  printErr: function(text) {
    if (arguments.length > 1) text = Array.prototype.slice.call(arguments).join(' ');
    console.error(text);
  },
};
</script>
`)

	sb.WriteString(`<script async type="text/javascript" src="/generated_files/` + filename + `.js"></script>`)
	outputDiv := sb.String()
	return getInputDiv(inputFlux), llvmDiv, outputDiv
}

func makeDivs(inputFlux, llvm, output string) (string, string, string) {
	inputDiv := getInputDiv(inputFlux)

	var sb strings.Builder

	addTextArea(&sb, "llvm", llvm, true)
	llvmDiv := sb.String()
	sb.Reset()

	addTextArea(&sb, "output", "", true)
	outputDiv := sb.String()

	return inputDiv, llvmDiv, outputDiv
}

func addTextArea(sb *strings.Builder, id, text string, readonly bool) {
	sb.WriteString(`<textarea name="` + id + `" id = "` + id + `" rows="30" cols="60"`)
	if readonly {
		sb.WriteString(` readonly`)
	}
	sb.WriteString(` spellcheck="false">` + text + `</textarea>
`)
}

func getInputDiv(fluxInput string) string {
	var sb strings.Builder
	sb.WriteString(`
<form action="/" method="post">
`)
	addTextArea(&sb, "flux", fluxInput, false)
	sb.WriteString(`<input type="submit" value="Submit">
</form> 
`)
	return sb.String()
}

var wrapperText string = `#include <stdio.h>
extern void flux_main();
int main(int argc, char ** argv) {
    flux_main();
    return 0;
}
`

var wasmID int

func newID() int {
	v := wasmID
	wasmID++
	return v
}

func compileToWASM(tempDir string, mod gollvm.Module) (string, error) {
	log.Println("Generating WASM, using temp dir " + tempDir)

	bcFilename := path.Join(tempDir, "flux.bc")
	bcFile, err := os.Create(bcFilename)
	if err != nil {
		return "", err
	}

	if err := gollvm.WriteBitcodeToFile(mod, bcFile); err != nil {
		return "", err
	}

	wrapperFilename := path.Join(tempDir, "wrapper.c")
	if err := ioutil.WriteFile(wrapperFilename, []byte(wrapperText), 0644); err != nil {
		return "", err
	}

	// If this fails, see https://webassembly.org/getting-started/developers-guide/
	emccPath, err := exec.LookPath("emcc")
	if err != nil {
		return "", err
	}

	basename := fmt.Sprintf("flux%d", newID())
	wasmFilename := basename + ".wasm"
	jsFilename := basename + ".js"
	fullWASMFilename := path.Join(tempDir, wasmFilename)
	fullJSFilename := path.Join(tempDir, jsFilename)
	cmd := exec.Command(emccPath, wrapperFilename, bcFilename, "-s", "WASM=1", "-o", fullJSFilename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	log.Println(string(output))

	bytes, err := ioutil.ReadFile(fullWASMFilename)
	if err != nil {
		return "", err
	}
	wasmFiles[wasmFilename] = bytes

	bytes, err = ioutil.ReadFile(fullJSFilename)
	if err != nil {
		return "", err
	}
	jsFiles[jsFilename] = bytes

	return basename, nil
}
