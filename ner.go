package ner

import (
	"bufio"
	"context"
	j2 "encoding/json/v2"
	"io"
	"iter"
	"os"
	"strings"
)

type NamedEntity struct {
	Text  string `json:"text"`
	Label string `json:"label"`
}

type Ner func(context.Context, string) iter.Seq2[NamedEntity, error]

type Entities iter.Seq2[NamedEntity, error]

type EntSink func(context.Context, Entities) error

type Writer struct{ io.Writer }

func (w Writer) ToSinkJSON() EntSink {
	return func(ctx context.Context, ient Entities) error {
		var bwtr *bufio.Writer = bufio.NewWriter(w.Writer)
		for ent, err := range ient {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if nil != err {
				return err
			}

			werr := j2.MarshalWrite(bwtr, ent)
			if nil != werr {
				return werr
			}

			_, nerr := bwtr.WriteString("\n")
			if nil != nerr {
				return nerr
			}
		}

		return bwtr.Flush()
	}
}

type InputSource func(context.Context) (string, error)

func (i InputSource) ToEntities(ctx context.Context, ner Ner) Entities {
	var empty NamedEntity

	src, err := i(ctx)

	if nil != err {
		return func(yield func(NamedEntity, error) bool) {
			yield(empty, err)
		}
	}

	return Entities(ner(ctx, src))
}

func (i InputSource) ToSink(
	ctx context.Context,
	ner Ner,
	sink EntSink,
) error {
	var ent Entities = i.ToEntities(ctx, ner)
	return sink(ctx, ent)
}

type InputLimit int64

const InputLimitDefault InputLimit = 1048576

func (l InputLimit) ReaderToSource(rdr io.Reader) InputSource {
	return func(ctx context.Context) (string, error) {
		var bldr strings.Builder
		lmtd := &io.LimitedReader{
			R: rdr,
			N: int64(l),
		}
		_, err := io.Copy(&bldr, lmtd)
		return bldr.String(), err
	}
}

func (l InputLimit) StdinToSource() InputSource {
	return l.ReaderToSource(os.Stdin)
}
