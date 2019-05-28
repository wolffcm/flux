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
	Long: "start a web server that accepts flux and produces wasm",
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

var wasmFiles = map[string][]byte{
}

var jsFiles = map[string][]byte{
}

var body string = `<!DOCTYPE html>
<html>
<body>

<h2>Compile and Execute Flux</h2>

<form action="/" method="post">
  Enter Flux text:<br>
  <textarea name="flux" cols="80" rows="8">%s</textarea>
  <br><br>
  <input type="submit" value="Submit">
</form> 

%s

</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling request...")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var existingFlux string

	var responseHTML string
	if flx, ok := r.PostForm["flux"]; ok {
		existingFlux = flx[0]
		log.Println("got some flux")
		responseHTML = processFlux(flx[0])
	}
	_, _ = fmt.Fprintf(w, body, existingFlux, responseHTML)
}

func generatedFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Files handler invoked")
	_, file := path.Split(r.URL.String())
	if strings.Contains(file, ".wasm") {
		bytes, ok := wasmFiles[file]
		if ! ok {
			http.Error(w, "WASM file not found", http.StatusNotFound)
		}
		_, _ = w.Write(bytes)
		return
	} else if strings.Contains(file, ".js") {
		bytes, ok := jsFiles[file]
		if ! ok {
			http.Error(w, "JS file not found", http.StatusNotFound)
		}
		_, _ = w.Write(bytes)
		return
	}

	http.Error(w, "unknown file " + file, http.StatusNotFound)
}

var fluxResponse string = `
<h2>LLVM IR</h2>

<p>%s</p>
`

func processFlux(flx string) string {
	tempDir, err := ioutil.TempDir("", "emcc")
	if err != nil {
		return fmt.Sprintf(fluxResponse, err)
	}

	//defer func() {_ = os.RemoveAll(tempDir)}()

	astPkg, err := flux.Parse(flx)
	if err != nil {
		return fmt.Sprintf(fluxResponse, err.Error())
	}

	mod, err := llvm.Build(astPkg)
	if err != nil {
		return fmt.Sprintf(fluxResponse, err.Error())
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(fluxResponse, "<pre>" + mod.String() + "</pre>"))

	filename, err := compileToWASM(tempDir, mod)
	if err != nil {
		return fmt.Sprintf(fluxResponse, err.Error())
	}

	log.Println("Now serving ", filename + ".js", " and ", filename + ".wasm")

	sb.WriteString("<h2>Output</h2>")
	sb.WriteString(`<textarea id="output" rows="8" cols="80"></textarea>`)
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
