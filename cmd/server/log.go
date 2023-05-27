package main

type LogWriter struct {
	append func(value string)
}

func (w *LogWriter) Write(p []byte) (int, error) {
	w.append(string(p))
	return len(p), nil
}
