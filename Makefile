buildDate := $(shell date +%y%m%d)
all:
	make build

build:
	go build -C cmd/gemrunner/ -o ../../gemrunner
	# go build -C cmd/editor/ -o ../../editor
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -C cmd/gemrunner/ -o ../../gemrunner.exe

package:
	zip -r linux-gemrunner-$(buildDate) assets README.md test_level.jpg level_editor.jpg gemrunner
	zip -r windows-gemrunner-$(buildDate) assets README.md test_level.jpg level_editor.jpg gemrunner.exe
	mv linux-gemrunner-$(buildDate).zip $(HOME)/Dropbox/builds/gemrunner/linux/linux-gemrunner-$(buildDate).zip
	mv windows-gemrunner-$(buildDate).zip $(HOME)/Dropbox/builds/gemrunner/windows/windows-gemrunner-$(buildDate).zip