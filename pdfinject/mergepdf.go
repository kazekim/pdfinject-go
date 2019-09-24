package pdfinject

func MergePDF(dir string, out string, in ...string) error{

	args := make([]string, len(in) + 3)

	for i,_ := range in {
		args[i] = in[i]
	}

	args[len(in)] = "cat"
	args[len(in) + 1] = "output"
	args[len(in) + 2] = out

	cmd := NewShellCommand(pdfFormPkgName)
	err := cmd.RunInPath(dir,args...)

	return err

}