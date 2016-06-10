// +build linux darwin freebsd netbsd
// +build !nofuse

package readonly

import (
	core "github.com/RealImage/go-ipfs/core"
	mount "github.com/RealImage/go-ipfs/fuse/mount"
)

// Mount mounts ipfs at a given location, and returns a mount.Mount instance.
func Mount(ipfs *core.IpfsNode, mountpoint string) (mount.Mount, error) {
	cfg, err := ipfs.Repo.Config()
	if err != nil {
		return nil, err
	}
	allow_other := cfg.Mounts.FuseAllowOther
	fsys := NewFileSystem(ipfs)
	return mount.NewMount(ipfs.Process(), fsys, mountpoint, allow_other)
}
