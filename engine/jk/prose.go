package jk

import (
	"context"
	"iter"

	"github.com/jdkato/prose/v3"
	n2 "github.com/takanoriyanagitani/go-nl2ner"
)

type Options []prose.Option

var OptionsDefault Options = nil

func (o Options) NewDocument(
	ctx context.Context,
	txt string,
) (*prose.Document, error) {
	return prose.NewDocumentContext(
		ctx,
		txt,
		o...,
	)
}

type Doc struct{ *prose.Document }

func (d Doc) Entities() []prose.Entity {
	return d.Document.Entities()
}

type Entity struct{ prose.Entity }

func (e Entity) Text() string  { return e.Entity.Text }
func (e Entity) Label() string { return e.Entity.Label }

func (o Options) ToEntities(
	ctx context.Context,
	txt string,
) iter.Seq2[n2.NamedEntity, error] {
	var empty n2.NamedEntity

	doc, err := o.NewDocument(ctx, txt)
	if nil != err {
		return func(yield func(n2.NamedEntity, error) bool) {
			yield(empty, err)
		}
	}

	return func(yield func(n2.NamedEntity, error) bool) {
		var entities []prose.Entity = Doc{Document: doc}.Entities()
		for _, ent := range entities {
			var text string = ent.Text
			var label string = ent.Label
			nent := n2.NamedEntity{
				Text:  text,
				Label: label,
			}

			if !yield(nent, nil) {
				return
			}
		}
	}
}

func (o Options) AsNer() n2.Ner { return o.ToEntities }
