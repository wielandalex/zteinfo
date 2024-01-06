# zteinfo

Command line app to retrieve information from a ZTE MC801A modem.

I built it to integrate information about the internet connection/usage into [Home Assistant](https://home-assistant.io).

Inspired by: https://github.com/ngarafol/ZTE-MC801A-Home-assistant

## Usage

```sh
zteinfo <ip_address> <password>
```

### Home Assistant

Example configuration for Home Assistant:

```yaml
command_line:
  - sensor:
    command: ~/zteinfo 192.168.0.1 yourPasswordHere
    json_attributes:
      - wa_inner_version
      - network_type
      - signalbar
      - Z5g_SINR
      - Z5g_rsrp
      - nr5g_action_band
      - nr5g_action_channel
      - lte_rsrq
      - lte_rsrp
      - cell_id
      - lte_snr
      - wan_active_channel
      - wan_active_band
      - wan_ipaddr
      - lte_multi_ca_scell_info
      - lte_ca_pcell_band
      - lte_ca_pcell_bandwidth
      - monthly_rx_bytes
      - monthly_tx_bytes
      - realtime_time
      - realtime_rx_thrpt
      - realtime_tx_thrpt
    name: ZTE Router Info
    unique_id: sensor.zte_router_info
    scan_interval: 60

template:
  - sensor:
    - name: "Monthly Usage"
      state: '{{ state_attr("sensor.zte_router_info", "monthly_rx_bytes") | int / 1024 / 1024 }}'
      unit_of_measurement: MB
      device_class: data_size
      state_class: total_increasing
```

## Build

```sh
make
```
