function upf(name, image, cpu, memory, disk) {
    return {
        name: name,
        image: image,
        cpu: { cores: cpu, "passthru": true },
        memory: { capacity: GB(memory) },
        //disks: [ { 'size': disk+"G", 'dev': 'vdb', 'bus': 'virtio' } ],
        //mounts: [{ source: cwd+'../..', point: "/tmp/upf" }]
        mounts: [{ source: cwd+'/../../..', point: "/tmp/upf" }]
    };
}

function node(name, image, cpu, memory) {
    return {
        name: name,
        image: image,
        memory: { capacity: GB(memory) },
    };
}

topo = {
	name: "upf"+Math.random().toString().substr(-6),
	nodes: [
		upf("upf","ubuntu-2204",2,8,64),
	],
    switches: [],
	links: []
}
