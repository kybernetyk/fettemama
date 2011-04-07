# Copyright 2009 The Go Authors. All rights reserved. 
# Use of this source code is governed by a BSD-style 
# license that can be found in the LICENSE file. 
# prerequisite: GOROOT and GOARCH must be defined 
# defines $(GC) (compiler), $(LD) (linker) and $(O) (architecture) 
include $(GOROOT)/src/Make.inc
.SUFFIXES : .$O .go 
# how to get an object from a .go file 
.go.$O: 
	$(GC) -o $@ $< 
# build target executable 
TARG = main 

# list of objects to be used to build the executable 
# note: $O is the target architecture - see above 
# note: list order of dependence with least dependent first 
OBJS = blog.$O dumbdb.$O main.$O
$(TARG): $(OBJS) 
	$(LD) -o $@ $@.$O 
	@echo "Executable is: $@" 
clean: 
	rm -rf *.[$(OS)o] 
nuke: 
	rm -rf *.[$(OS)o] *~ $(TARG)
	