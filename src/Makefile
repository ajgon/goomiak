all: zexall zexbit zexdoc z80ccf z80docflags z80doc z80flags z80full z80memptr z80tests

zexall:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/zexall.tap

zexbit:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/zexbit.tap

zexdoc:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/zexdoc.tap

z80ccf:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80ccf.tap

z80docflags:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80docflags.tap

z80doc:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80doc.tap

z80flags:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80flags.tap

z80full:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80full.tap

z80memptr:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80memptr.tap

z80tests:
	@go run main.go -autoload -fullspeed -tapfile ../tests/passing/z80tests.tap
