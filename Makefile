build:
	@go build \
		-v \
		-o bin/t-walker

install: build
	 @cp ./bin/t-walker /usr/local/bin
	 @cp ./t-wrapper.sh /usr/local/bin # We need t-wrapper.sh so that the 'c' change directory on the shell command works
	 @chmod +x /usr/local/bin/t-walker
	 @chmod +x /usr/local/bin/t-wrapper.sh
	 @grep -qxF 'alias t="source /usr/local/bin/t-wrapper.sh"' ~/.zshrc || echo '\nalias t="source /usr/local/bin/t-wrapper.sh"' >> ~/.zshrc
