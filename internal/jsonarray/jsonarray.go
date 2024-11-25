package jsonarray

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
)

// MarshalArray writes a JSON array of objects to the writer.
// The objects are streamed from the given iterator.
func MarshalArray[T any](w io.Writer, objects iter.Seq[T]) error {
	first := true

	enc := json.NewEncoder(w)
	w.Write([]byte{'['})

	for obj := range objects {
		if !first {
			w.Write([]byte{','})
		}
		first = false

		if err := enc.Encode(obj); err != nil {
			return fmt.Errorf("encoding object: %w", err)
		}
	}

	w.Write([]byte{']'})
	return nil
}

// UnmarshalArray reads a JSON array of objects from the reader.
// The objects are streamed to the returned iterator.
func UnmarshalArray[T any](r io.Reader) iter.Seq2[T, error] {
	dec := json.NewDecoder(r)

	return func(yield func(T, error) bool) {
		open, err := dec.Token()
		if err != nil {
			var v T
			yield(v, fmt.Errorf("cannot read open delim: %w", err))
			return
		}

		if delim, ok := open.(json.Delim); !ok || delim != '[' {
			var v T
			yield(v, fmt.Errorf("expected array start, got %v", open))
			return
		}

		for dec.More() {
			var v T
			err := dec.Decode(&v)
			if !yield(v, err) {
				return
			}
		}

		close, err := dec.Token()
		if err != nil {
			var v T
			yield(v, fmt.Errorf("cannot read close delim: %w", err))
			return
		}

		if delim, ok := close.(json.Delim); !ok || delim != ']' {
			var v T
			yield(v, fmt.Errorf("expected array end, got %v", close))
			return
		}
	}
}
