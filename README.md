# Helper Tools to test a cifs-mount-control snappy interface

https://forum.snapcraft.io/t/interface-mount-umount-cifs-share-permission/5689

## Prepare samba server

It exports a share named "sambashare" without any access restrictions.

    cd samba-host
    ./build_container.sh
    ./run_container.sh

Adjust `run_container.sh` as needed. You can also just pick up the smb.conf and
put it on some vm.

## Prepare smb-copier snap

A simple daemon which spawns a simple REST API to mount a cifs directory on a local directory.

Adjust the `snap/snapcraft.yaml` so the smbShare commandline parameter points
to your samba container/vm/host.

e.g.:

```
apps:

  smb-copier:
    command: smb-copier -smbShare "//192.168.123.1/sambashare"
```

Then build the snap:

    ./build_snap.sh

### Ensure basic functionality in devmode

Copy the snap over to an ubuntu core box and install it. First install it as devmode to
ensure the basic setup works.

    sudo snap install --dangerous --devmode smb-copier_0.1.0_amd64.snap

Now execute the REST calls to check basic functionality.

* /mount : executes the smb mount
* /umount: executes the unmount
* / (or any other path): show all cifs mounts

e.g.:

    curl http://192.168.123.67:9090
    curl http://192.168.123.67:9090/mount
    curl http://192.168.123.67:9090
    curl http://192.168.123.67:9090/umount
    curl http://192.168.123.67:9090

### Reinstall in strict mode and check interfaces

Now reinstall the snap in strict mode and connect the necessary interfaces:

    sudo snap install --dangerous smb-copier_0.1.0_amd64.snap
	sudo snap connect smb-copier:mount-observe
    sudo snap connect smb-copier:cifs-mount-control
    sudo snap connect smb-copier:network-bind

Now try again to execute the curl commands on the minimal rest api.
