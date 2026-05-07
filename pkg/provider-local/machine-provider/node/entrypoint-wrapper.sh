#!/bin/bash
# If /var/lib/containerd is backed by an overlay filesystem, containerd cannot
# create nested overlay mounts (overlay-on-overlay). This happens when the machine
# pod runs inside another machine pod (e.g. shoot-in-managedseed), where the
# containerd PVC is provisioned by local-path-provisioner on top of an overlay fs.
# Switch to the 'native' snapshotter which uses bind mounts and works on any filesystem.
if [ "$(stat -f -c %T /var/lib/containerd 2>/dev/null)" = "overlayfs" ]; then
  export KIND_EXPERIMENTAL_CONTAINERD_SNAPSHOTTER=native
fi
exec /usr/local/bin/entrypoint "$@"
