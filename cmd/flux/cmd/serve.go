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
<link href="https://fonts.googleapis.com/css?family=Roboto|Roboto+Mono&display=swap" rel="stylesheet">
<style>

* {
  box-sizing: border-box;
}

html, body {
  background: linear-gradient(180deg,#202028 0,#0f0e15);
  color: #DBDBDB;
  font-family: Roboto,Helvetica,Arial,Tahoma,Verdana,sans-serif;
  height: 100%;
  margin: 0;
}

.page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.row {
  display: flex;
  flex-direction: row;
  margin-left: 142px;
  margin-right: 142px;
}

.spacer {
  height: 15px;
}

.panelspacer {
  width: 15px;
}

.page-header {
  display: flex;
  height: 120px;
  align-items: center;
}

.panel {
  flex: 1 1 auto;
  background-color: #292933;
  border-radius: 15px;
}

.panel-header {
  display: flex;
  flex-direction: row;
  height: 50px;
  align-items: center;
  justify-content: space-between;
  height: 30px;
  align-items: center;
  padding: 10px;
  justify-content: space-between;
}

.panel-header-left {
  align-content: flex-start;
}

.panel-header-right {
  align-content: flex-end;
}

.panel-body {
  padding-left: 5px;
  padding-right: 7px;
  padding-bottom: 5px;
}

::selection {
    background-color: #22adf6;
    color: #fff;
}


input {
  border-radius: 10px;
  background-color: #22adf6;
  border-color: #22adf6;
  color: #f0fcff;
}

textarea {
  width: 100%;
  border: solid;
  border-color: #383846;
  white-space: pre;
  font-family: RobotoMono,monospace;
  background-color: #000;
  font-size: 14px;
  color: #CAABFF;
  border-radius: 10px;
  padding: 5px;
}

</style>
</head>
<body>

<div class="page">

<div class="row">
  <div class="page-header">
    <h3>Compile and Execute Flux</h3>
  </div>
</div>

<div class="row">

  <div class="panel"><form action="/" method="post">
    <div class="panel-header">
      <div class="panel-header-left">
        <b>FLUX SOURCE</b>
      </div>
      <div class="panel-header-right">
        <input type="submit" value="Submit">
      </div>
    </div>
    <div class="panel-body">
      {input}
    </div>
  </form></div>

  <div class="panelspacer"></div>

  <div class="panel">
    <div class="panel-header">
      <b>OUTPUT</b>
    </div>
    <div class="panel-body">
      {output}
    </div>
  </div> 

</div><!-- end row -->

<div class="row spacer"></div>

<div class="row">
  <div class="panel">
    <div class="panel-header">
      <b>LLVM IR</b>
    </div>
  <div class="panel-body">
    {llvm}
  </div>
</div><!-- end row -->

</div><!-- end page -->
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
	sb.WriteString(`<textarea name="` + id + `" id = "` + id + `" rows="20" cols="60"`)
	if readonly {
		sb.WriteString(` readonly`)
	}
	sb.WriteString(` spellcheck="false">` + text + `</textarea>
`)
}

func getInputDiv(fluxInput string) string {
	var sb strings.Builder
	addTextArea(&sb, "flux", fluxInput, false)
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
