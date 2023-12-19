tinygo-build:
	tinygo build -no-debug -panic=trap -x -target=spec.json && strip --strip-all -R .comment -R .note -R .eh_frame -R .shstrtab test
# tinygo build -panic="trap" -size=full -x -target=spec.json main.go 

# 

# & strip -s -R .eh_frame main
# strip -s -R .comment -R .gnu.version main

go-build:
	go build -ldflags="-w -s" && strip -s -R .comment -R .gnu.version -- strip-unneded main
