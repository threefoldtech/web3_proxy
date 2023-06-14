module tfgrid

[params]
pub struct VM {
pub mut:
	name                 string // This is the vm's name. If multiple vms are to be deployed, index is appended to the vm's name. If not provided, a random name is generated.
	network              string // Identifier for the network that these VMs will be a part of
	farm_id              u32    // farm id to deploy on, if 0, a random eligible node on a random farm will be selected
	capacity             string // capacity for the vms, could be 'small', 'medium', 'large', or 'extra-large'
	times                u32 = 1 // indicates how many vms will be deployed with the configuration defined by this object
	disk_size            u32    // disk size to mount on vms in GB
	ssh_key              string // this is the public key that will allow you to ssh into the VM at a later stage
	gateway              bool   // if true, a gateway will deployed for each vm. vms should listen for traffic coming from the gateway at port 9000
	add_wireguard_access bool   // if true, a wireguard access point will be added to the network
	add_public_ipv4      bool   // if true, a public ipv4 will be added to each vm
	add_public_ipv6      bool   // if true, a public ipv6 will be added to each vm
}

[params]
pub struct RemoveVMArgs {
pub:
	network string
	vm_name string
}

// Deploys a vm with the posibility to add a gateway. if the there is already a network with the same name, the the vms are added to this network
pub fn (mut t TFGridClient) deploy_vm(vm VM) !VMResult {
	return t.client.send_json_rpc[[]VM, VMResult]('tfgrid.DeployVM', [vm], default_timeout)!
}

// Removes a vm from a network
pub fn (mut t TFGridClient) remove_vm(args RemoveVMArgs) !VMResult {
	return t.client.send_json_rpc[[]RemoveVMArgs, VMResult]('tfgrid.RemoveVM', [
		args,
	], default_timeout)!
}

// Gets a deployed network of vms
pub fn (mut t TFGridClient) get_vm(network string) !VMResult {
	return t.client.send_json_rpc[[]string, VMResult]('tfgrid.GetVM', [
		network,
	], default_timeout)!
}

// Deletes a deployed network of vms
pub fn (mut t TFGridClient) delete_vm(network string) ! {
	_ := t.client.send_json_rpc[[]string, string]('tfgrid.DeleteVM', [
		network,
	], default_timeout)!
}
