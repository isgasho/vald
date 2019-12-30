//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package compress provides compress functions
package compress

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4"
)

type lz4Compressor struct {
	gobc Compressor
}

func NewLZ4() Compressor {
	return &lz4Compressor{
		gobc: NewGob(),
	}
}

func (l *lz4Compressor) CompressVector(vector []float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := lz4.NewWriter(buf)

	gob, err := l.gobc.CompressVector(vector)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(gob)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (l *lz4Compressor) DecompressVector(bs []byte) ([]float64, error) {
	buf := new(bytes.Buffer)
	zr := lz4.NewReader(bytes.NewReader(bs))
	_, err := io.Copy(buf, zr)
	if err != nil {
		return nil, err
	}

	bufbytes := buf.Bytes()

	vec, err := l.gobc.DecompressVector(bufbytes)
	if err != nil {
		return nil, err
	}

	return vec, nil
}