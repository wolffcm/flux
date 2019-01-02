
package repl

type completer{
	r *REPL
}

func (c *completer) Do(line []rune, pos int) (options [][]rune, length int) {
	names := c.r.interpreter.GlobalScope().Names()
	sort.Strings(names)

	options := make([][]rune, 0, len(names))
	//s := make([]prompt.Suggest, 0, len(names))
	//for _, n := range names {
	//	if !strings.HasPrefix(n, "_") {
	//		s = append(s, prompt.Suggest{Text: n})
	//	}
	//}

	if len(line) == 0 || strings.HasPrefix(d.Text, "@") {
		prefix := line[:pos]
		l := len(prefix)
		root := "./" + strings.TrimPrefix(d.Text, "@")
		fluxFiles, err := getFluxFiles(root)
		if err == nil {
			for _, fName := range fluxFiles {
				name := "@" + fName
				if readline.HasPrefix(name, prefix) {
					options = append(options, name[l:])
				}
			}
		}
		dirs, err := getDirs(root)
		if err == nil {
			for _, fName := range dirs {
				s = append(s, prompt.Suggest{Text: "@" + fName + string(os.PathSeparator)})
			}
		}
	}

	return options, 5
}
