DEVICE ?= fr55
KEY ?= ../developer_key.der
PROJECT = WhereAreYou

# Find the SDK bin directory automatically
SDK_BIN = $(shell find $(HOME)/.Garmin/ConnectIQ/Sdks -name monkeyc -type f -executable | head -n 1 | xargs dirname)
MONKEYC = $(SDK_BIN)/monkeyc
MONKEYDO = $(SDK_BIN)/monkeydo
CONNECTIQ = $(SDK_BIN)/connectiq

build:
	mkdir -p bin
	$(MONKEYC) \
	-f monkey.jungle \
	-o bin/$(PROJECT).prg \
	-y $(KEY) \
	-d $(DEVICE) \
	-w

run: build
	@pgrep -f "simulator" > /dev/null || ($(CONNECTIQ) & sleep 5)
	$(MONKEYDO) bin/$(PROJECT).prg $(DEVICE)

clean:
	rm -rf bin
