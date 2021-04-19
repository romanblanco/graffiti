package extractor

import (
  "archive/tar"
  "fmt"
  "io"
  "os"
)

type Extractor struct {
  tar *tar.Reader
  counter int64
}

func New(reader io.Reader) *Extractor {
  bodyReader := tar.NewReader(reader)

  ex := &Extractor{tar:bodyReader}
  return ex
}

func (te *Extractor) Next() (os.FileInfo, io.Reader, error) {
  for {
    header, err := te.tar.Next()

    if err != nil && err != io.EOF {
      fmt.Printf("error in extractor: %v\n", err)
      return nil, nil, err
    }

    if header == nil || err == io.EOF {
      fmt.Printf("error 2 in extractor: %v\n", err)
      return nil, nil, err
    }

    switch header.Typeflag {
    case tar.TypeDir, tar.TypeSymlink:
    case tar.TypeReg:
      return header.FileInfo(), te.tar, err
    default:
      return nil, nil, fmt.Errorf("unrecognized tar header type: %d", header.Typeflag)
    }
    te.counter += 1
  }
}
