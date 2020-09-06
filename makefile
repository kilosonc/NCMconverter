NCMconverter:
	go build -o $@
clean:
	-rm ./*.ncm
	-rm ./*.mp3
	-rm -r ./temp
.PHONY: clean NCMconverter
