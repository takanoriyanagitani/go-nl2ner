package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	n2 "github.com/takanoriyanagitani/go-nl2ner"
	ej "github.com/takanoriyanagitani/go-nl2ner/engine/jk"
)

var opts ej.Options = ej.OptionsDefault
var ner n2.Ner = opts.AsNer()

var wtr io.Writer = os.Stdout
var sink n2.EntSink = n2.Writer{Writer: wtr}.ToSinkJSON()

var ilmt n2.InputLimit = n2.InputLimitDefault
var isrc n2.InputSource = ilmt.StdinToSource()

func imain(ctx context.Context) error {
	return isrc.ToSink(
		ctx,
		ner,
		sink,
	)
}

func main() {
	err := imain(context.Background())
	if nil != err {
		slog.Error("error got", "error", err.Error())
	}
}
