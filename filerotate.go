package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type fileRotate struct {
	sync.Mutex
	ageLimit  time.Duration
	numLimit  int
	sizeLimit int64
	size      int64
	file      *os.File
	sTime     time.Time
}

//age may be 0. size may be 0
func NewFileRotate(fileName string, age time.Duration, size int64, num int) (io.WriteCloser, error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	ret := &fileRotate{
		ageLimit:  age,
		numLimit:  num,
		sizeLimit: size,
		file:      f,
		sTime:     time.Now(),
	}

	var fi os.FileInfo
	if fi, err = f.Stat(); err == nil {
		ret.size = fi.Size()

		if age > 0 && time.Since(fi.ModTime()) >= age {
			err = ret.rotate()
		} else if size > 0 && fi.Size() > int64(size) {
			err = ret.rotate()
		}
	}

	if err != nil {
		f.Close()

		return nil, err
	} else {
		return ret, nil
	}
}

func (l *fileRotate) Write(data []byte) (n int, err error) {
	if l.file == nil {
		return 0, os.ErrInvalid
	}

	l.Lock()
	defer l.Unlock()

	if l.ageLimit > 0 && time.Since(l.sTime) >= l.ageLimit {
		l.sTime = time.Now()
		if err = l.rotate(); err != nil {
			return
		}
	} else if l.sizeLimit > 0 && l.size+int64(len(data)) > l.sizeLimit {
		if err = l.rotate(); err != nil {
			return
		}
	}

	n, err = l.file.Write(data)
	l.size += int64(n)

	return
}

func (l *fileRotate) Close() error {
	l.Lock()
	defer l.Unlock()

	return l.close()
}

func (l *fileRotate) close() (err error) {
	if l.file != nil {
		err = l.file.Close()
		l.file = nil
	}

	return
}

func (l *fileRotate) Rotate() error {
	l.Lock()
	defer l.Unlock()

	return l.rotate()
}

func (l *fileRotate) rotate() error {
	name := l.file.Name()
	l.close()

	var fileName string
	fileNameTo := ""

	for i := l.numLimit; i > 0; i-- {
		if i > 1 {
			fileName = fmt.Sprintf("%s.%d.gz", name, i)
		} else {
			fileName = fmt.Sprintf("%s.1", name)
		}

		if _, err := os.Stat(fileName); err == nil {
			if fileNameTo == "" {
				os.Remove(fileName)
			} else if i == 1 {
				if err := compress(fileName, fileNameTo); err != nil {
					return nil
				}
			} else if err := os.Rename(fileName, fileNameTo); err != nil {
				return err
			}
		}
		fileNameTo = fileName
	}

	if err := os.Rename(name, fileNameTo); err != nil {
		return err
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	l.file = f
	l.size = 0

	return nil
}

func compress(name, nameTo string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	arc, err := os.OpenFile(nameTo, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer arc.Close()

	z := gzip.NewWriter(arc)
	defer z.Close()

	_, err = io.Copy(z, f)

	return err
}
