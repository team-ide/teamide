package filework

import "io"

type FileWorker interface {
	Upload(path string, reader io.Reader)
	Download(path string, writer io.Writer)
	Write(reader io.Reader, writer io.Writer)
}
