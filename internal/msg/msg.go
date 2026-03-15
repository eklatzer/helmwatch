package msg

func NewRender(diff string, dir *string) Render {
	return Render{
		Diff:      diff,
		Directory: dir,
	}
}

type Render struct {
	Diff      string
	Directory *string
}

type FileChanged struct{}
