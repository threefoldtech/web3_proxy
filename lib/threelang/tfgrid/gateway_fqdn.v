module tfgrid

import freeflowuniverse.crystallib.actionsparser {Action}
import threefoldtech.threebot.tfgrid { GatewayFQDN }
import threefoldtech.threebot.tfgrid.solution
import rand
import log

fn (mut t TFGridHandler) gateway_fqdn(action Action) ! {
	mut logger := log.Logger(&log.Log{
		level: .info
	})

	match action.name {
		'create' {
			node_id := action.params.get_int('node_id')!
			name := action.params.get_default('name', rand.string(10).to_lower())!
			tls_passthrough := action.params.get_default_false('tls_passthrough')
			backend := action.params.get('backend')!
			fqdn := action.params.get('fqdn')!

			gw_deploy := t.solution_handler.tfclient.gateways_deploy_fqdn(GatewayFQDN{
				name: name
				node_id: u32(node_id)
				tls_passthrough: tls_passthrough
				backends: [backend]
				fqdn: fqdn
			})!

			logger.info('${gw_deploy}')
		}
		'delete' {
			name := action.params.get('name')!
			t.solution_handler.tfclient.gateways_delete_fqdn(name)!
		}
		'get' {
			name := action.params.get('name')!
			gw_get := t.solution_handler.tfclient.gateways_get_fqdn(name)!

			logger.info('${gw_get}')
		}
		else {
			return error('action ${action.name} is not supported on gateways')
		}
	}
}