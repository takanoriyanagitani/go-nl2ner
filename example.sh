#!/bin/bash

wsm="./nl2ner.wasm"

input1(){
	echo Alice bought fresh apples, bananas, and mangos
	echo at a market in Orange County.
	echo
	echo Later that afternoon, she met with representatives
	echo from BlackBerry and Apple Inc. in California.
	echo
	echo They discussed shipping fresh pineapples from
	echo Hawaii to London next month.
}

input2(){
	echo Marie Curie studied in Paris.
	echo She won the Nobel Prize twice.
}

run_wasi(){
	cat /dev/stdin |
		wasmtime run "${wsm}"
}

input2 | run_wasi
