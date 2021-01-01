
.PHONY: vnc
vnc:
	go run . `ipconfig getifaddr en0` 5902 "password"
