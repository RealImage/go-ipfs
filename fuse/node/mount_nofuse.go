// +build linux darwin freebsd netbsd
// +build nofuse

package node

import (
	"errors"

	core "github.com/RealImage/go-ipfs/core"
)

func Mount(node *core.IpfsNode, fsdir, nsdir string) error {
	return errors.New("not compiled in")
}
