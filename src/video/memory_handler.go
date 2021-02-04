package video

const ramStartAddress uint16 = 0x4000
const ramEndAddress uint16 = 0x5aff

type VideoMemoryHandler struct {
	isDirty        bool
	dirtyAddresses [6912]bool
}

func (vmh *VideoMemoryHandler) Name() string {
	return "video"
}

func (vmh *VideoMemoryHandler) MarkRangeAsDirty(start, end uint16) {
	if start < ramStartAddress {
		start = ramStartAddress
	}

	if end > ramEndAddress {
		end = ramEndAddress
	}

	if start > end {
		return
	}

	vmh.isDirty = true

	for i := start; i <= end; i++ {
		vmh.dirtyAddresses[i-ramStartAddress] = true
	}
}

func (vmh *VideoMemoryHandler) MarkAsDirty(address uint16) {
	if address < ramStartAddress || address > ramEndAddress {
		return
	}

	vmh.isDirty = true
	vmh.dirtyAddresses[address-ramStartAddress] = true
}

func (vmh *VideoMemoryHandler) IsMemoryDirty() bool {
	return vmh.isDirty
}

func (vmh *VideoMemoryHandler) CheckAddressDirtiness(address uint16) bool {
	if address < ramStartAddress || address > ramEndAddress {
		panic("Invalid dirty address to check")
	}

	return vmh.dirtyAddresses[address-ramStartAddress]
}

func (vmh *VideoMemoryHandler) MarkAsFresh() {
	vmh.dirtyAddresses = [6912]bool{}
	vmh.isDirty = false
}

func VideoMemoryHandlerNew() *VideoMemoryHandler {
	return new(VideoMemoryHandler)
}
