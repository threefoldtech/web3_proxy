!!tfgrid.core.login
	mnemonic: 'mom picnic deliver again rug night rabbit music motion hole lion where'
	network: dev

!!tfgrid.sshkeys.new
    name: default
	ssh_key: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCQi7Qp0fs4WowSBQJonYHNWNJ5q776XbFO8HnUggyGse1Z4CFZyVpUdWaIzpFkQdivAACSKmqfE6jAunX7HuujTQhLhVgs/iCQ3ALQfQ118Jh1g2dz7S3/klHJs6JqfXLKtwDHzq2DuEDjls5PPoD6SVipjQo+kFO2tvKUYOrXryGW5VNPSUKtXZJX4kxtLzCANqENMSqZIBiJhXj7+JQq8Kp6E117dkLxh4BmPJmGS4stSAfgSdmEWgm0MgAbHkc2X+fLsWrcEBYaXE1b+n70bVXGDVEfeuMNZjBlsgULVR0DXY5zxegciOSNr1b7yF/ZdoALN0gmQg+AywPy92+oeY7EXLabDoDUKcE+42EHscXEkTHlhCieF/W9worCzGqpMwJuBDNvDu5kP1y/vB+ZfPVTlZ1kS77/OuDTr/zssQI/SgSszVXTyVSFIFIbXLGuUDscnmPHVPV4PnmeOa2aeF1cgX0o/JErQ8+iu2wqQKueZT4QAUFyknIgXloSBVs= mariocs@mario-codescalers'

!!tfgrid.machines.create
    name: testvm
    sshkey: default
	network: skynet
	capacity: small
	times: 2
	gateway: yes
	add_wireguard_access: yes