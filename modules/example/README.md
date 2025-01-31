<!--
title: "Example module"
description: "Use this example data collection module, which produces example charts with random values, to better understand how to build your own collector in Go."
custom_edit_url: "https://github.com/khulnasoft/go.plugin/edit/master/modules/example/README.md"
sidebar_label: "Example module in Go"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Mock Collectors"
-->

# Example module

An example data collection module. Use it as an example writing a new module.

## Charts

This module produces example charts with random values. Number of charts, dimensions and chart type is configurable.

## Configuration

Edit the `go/example.conf` configuration file using `edit-config` from the
Khulnasoft [config directory](https://github.com/khulnasoft/khulnasoft/blob/master/docs/configure/nodes.md), which is typically at `/etc/khulnasoft`.

```bash
cd /etc/khulnasoft # Replace this path with your Khulnasoft config directory
sudo ./edit-config go/example.conf
```

Disabled by default. Should be explicitly enabled
in [go.conf](https://github.com/khulnasoft/go.plugin/blob/master/config/go.conf).

```yaml
# go.conf
modules:
  example: yes
```

Here is an example configuration with several jobs:

```yaml
jobs:
  - name: example
    charts:
      num: 3
      dimensions: 5

  - name: hidden_example
    hidden_charts:
      num: 3
      dimensions: 5
```

---

For all available options, see the Example
collector's [configuration file](https://github.com/khulnasoft/go.plugin/blob/master/config/go/example.conf).

## Troubleshooting

To troubleshoot issues with the `example` collector, run the `go.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/khulnasoft/plugins.d/`. If that's not the case on
  your system, open `khulnasoft.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/khulnasoft/plugins.d/
  ```

- Switch to the `khulnasoft` user.

  ```bash
  sudo -u khulnasoft -s
  ```

- Run the `go.plugin` to debug the collector:

  ```bash
  ./go.plugin -d -m example
  ```
