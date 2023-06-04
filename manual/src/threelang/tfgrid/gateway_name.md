# Gateway Name Namespace

- To deploy a gateway name workload, use the Gateway Name namespace

## Create Operation

- action name: !!tfgrid.gateway_name.create
- parameters:
  - name [required]
  - farm_id [optional]
  - backend [required]
    - the URL that the gateway will pass traffic to.
- arguments:
  - tls_passthrough

## Get Operation

- action name: !!tfgrid.gateway_name.get
- parameters:
  - name [required]

## Update Operations

- Update operations are not allowed on gateway names.

## Delete Operation

- action_name: !!tfgrid.gateway_name.delete
- parameters:
  - name [required]