# eapi-exporter

Prometheus exporter for Arista EOS devices via [eAPI](https://www.arista.com/en/um-eos/eos-command-api).

## Metrics

### Network interfaces (`show interfaces`)

| Metric | Type | Description |
|---|---|---|
| `node_network_up` | gauge | Whether the interface is connected |
| `node_network_info` | gauge | Non-numeric interface metadata (operstate, address, description, hardware) |
| `node_network_mtu_bytes` | gauge | MTU of the interface |
| `node_network_speed_bytes` | gauge | Speed of the interface in bytes per second |
| `node_network_receive_bytes_total` | counter | Total bytes received |
| `node_network_transmit_bytes_total` | counter | Total bytes transmitted |
| `node_network_receive_packets_total` | counter | Total packets received |
| `node_network_transmit_packets_total` | counter | Total packets transmitted |
| `node_network_receive_errs_total` | counter | Total receive errors |
| `node_network_transmit_errs_total` | counter | Total transmit errors |
| `node_network_receive_drop_total` | counter | Total received packets dropped |
| `node_network_transmit_drop_total` | counter | Total transmitted packets dropped |
| `node_network_receive_multicast_total` | counter | Total multicast packets received |
| `node_network_receive_broadcast_total` | counter | Total broadcast packets received |
| `node_network_transmit_multicast_total` | counter | Total multicast packets transmitted |
| `node_network_transmit_broadcast_total` | counter | Total broadcast packets transmitted |
| `node_network_carrier_changes_total` | counter | Total carrier link status changes |

### System (`show version`, `show environment power`)

| Metric | Type | Description |
|---|---|---|
| `node_boot_time_seconds` | gauge | Unix timestamp of system boot time |
| `node_power_supply_info` | gauge | Power supply metadata (model) |
| `node_power_supply_status` | gauge | Power supply status, 1 if ok |
| `node_power_supply_capacity_watts` | gauge | Power supply capacity in watts |
| `node_power_supply_input_current_amperes` | gauge | Input current in amperes |
| `node_power_supply_output_current_amperes` | gauge | Output current in amperes |
| `node_power_supply_output_watts` | gauge | Output power in watts |
| `node_power_supply_uptime_seconds` | gauge | Power supply uptime in seconds |
| `node_power_supply_temp_celsius` | gauge | Power supply temperature in celsius |
| `node_power_supply_fan_speed` | gauge | Power supply fan speed |
| `node_power_supply_fan_status` | gauge | Power supply fan status, 1 if ok |

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `EAPI_HOST` | yes | | Hostname or IP of the EOS device |
| `EAPI_USERNAME` | yes | | eAPI username |
| `EAPI_PASSWORD` | yes | | eAPI password |
| `EAPI_PROTOCOL` | no | `https` | `http` or `https` |
| `EAPI_PORT` | no | `443` | eAPI port |

## Usage

```sh
export EAPI_HOST=10.0.0.1
export EAPI_USERNAME=admin
export EAPI_PASSWORD=secret

go run ./cmd/exporter
```

Metrics are served at `:9120/metrics`.

## Endpoints

| Path | Description |
|---|---|
| `/metrics` | Prometheus metrics |
| `/health` | Liveness probe (always 200) |
| `/ready` | Readiness probe (queries the device) |

## Deploying with Kustomize

A base manifest lives in `deploy/kustomize/`. The deployment expects a Secret named `eapi-exporter` containing the environment variables.

### Using an overlay

Create an overlay directory with your environment-specific secret or configmap:

```
deploy/kustomize/overlays/production/
├── kustomization.yaml
└── secret.yaml
```

**secret.yaml**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: eapi-exporter
stringData:
  EAPI_HOST: "10.0.0.1"
  EAPI_USERNAME: "admin"
  EAPI_PASSWORD: "secret"
```

**kustomization.yaml**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: monitoring

resources:
  - ../../           # points to the base
  - secret.yaml
```

You can also generate the secret directly in the overlay instead of checking it in:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: monitoring

resources:
  - ../../

secretGenerator:
  - name: eapi-exporter
    literals:
      - EAPI_HOST=10.0.0.1
      - EAPI_USERNAME=admin
      - EAPI_PASSWORD=secret

generatorOptions:
  disableNameSuffixHash: true
```

Apply with:

```sh
kubectl apply -k deploy/kustomize/overlays/production/
```

## TLS and older EOS devices

Older Arista EOS versions negotiate TLS with RSA key exchange, which Go disables by default. If you see TLS handshake errors, set the `GODEBUG` environment variable before starting the exporter:

```sh
export GODEBUG=tlsrsakex=1
```

## Building

```sh
go build -o exporter ./cmd/exporter
```
