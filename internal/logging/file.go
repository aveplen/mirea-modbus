package logging

import (
	"fmt"
	"os"
)

type FileSubscriber struct {
	file    *os.File
	buffers [][]string
	current int
}

// data race
func (f *FileSubscriber) Consume(message string) {
	if cap(f.buffers[f.current]) == 0 {
		go func(buffer *[]string) {
			for _, line := range *buffer {
				if _, err := f.file.WriteString(line); err != nil {
					panic(fmt.Errorf("logging error: %w", err))
				}

				if _, err := f.file.WriteString("\n"); err != nil {
					panic(fmt.Errorf("logging error: %w", err))
				}
			}

			cleaned := (*buffer)[:0]
			buffer = &cleaned
		}(&f.buffers[f.current])

		f.current = (f.current + 1) % len(f.buffers)
	}

	f.buffers[f.current] = append(f.buffers[f.current], message)
}

func NewFileSubscriber(filename string) *FileSubscriber {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("open file %s: %w", filename, err))
	}

	fileSubscriber := &FileSubscriber{
		file: file,
		buffers: [][]string{
			make([]string, 100),
			make([]string, 100),
		},
	}

	return fileSubscriber
}
