CAPNP := $(shell command -v capnp 2> /dev/null)
CAPNPPATH := $(PWD)/serialization/capnpmodels

.PHONY: all check_capnp install_capnp compile_capnp

all: check_capnp compile_capnp

check_capnp:
ifndef CAPNP
	@echo "Cap'n Proto not found, installing..."
	@make install_capnp
else
	@echo "Cap'n Proto is installed."
endif

install_capnp:
	@curl -O https://capnproto.org/capnproto-c++-1.0.1.tar.gz
	@tar zxf capnproto-c++-1.0.1.tar.gz
	@cd capnproto-c++-1.0.1 && ./configure && make -j6 check
	@sudo make -C capnproto-c++-1.0.1 install
	@rm -r capnproto-c++-1.0.1
	@rm capnproto-c++-1.0.1.tar.gz
	@echo "Cap'n Proto installed."

compile_capnp:
	@echo "Compiled Cap'n Proto schema to Go $(CAPNPPATH)."
	@capnp compile -I$(CAPNPPATH) -ogo state_object.capnp
	@echo "Compiled Cap'n Proto schema to Go."