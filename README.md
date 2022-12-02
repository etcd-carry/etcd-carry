# etcd-carry

English | [简体中文](./README-zh_CN.md)

## Introduction

etcd-carry provides the ability to synchronize resources in the K8s cluster that meet custom rules to the standby k8s cluster in real time.

## Status

This is still very experimental and should be used with caution in production environments. At present, it does not support the correct synchronization of Service resources, and does not support the synchronization of data in PVs of stateful services.It is recommended to back up the standby etcd before use.

## Installation

Supports binary deployment and deployment to K8s with helm.

## Options

- debug -- enable client-side debug logging
- mirror-rule -- Specify the custom rules to start mirroring
- encryption-provider-config --  The file containing configuration for encryption providers to be used for storing K8s secrets in etcd
- kube-prefix -- the prefix to all kubernetes resources passed to etcd (default "/registry")
- max-txn-ops -- Maximum number of operations permitted in a transaction during syncing updates (default 128)
- rev -- Specify the kv revision to start to mirror (default 0)
- dial-timeout -- dial timeout for client connections
- keepalive-time -- keepalive time for client connections
- keepalive-timeout -- keepalive timeout for client connections
- master-endpoints -- List of source etcd servers to connect with (scheme://ip:port), comma separated
- master-cacert -- verify certificates of TLS-enabled secure servers using this CA bundle
- master-cert -- identify secure client using this TLS certificate file
- master-key -- identify secure client using this TLS key file
- master-insecure-skip-tls-verify -- skip server certificate verification
- master-insecure-transport -- disable transport security for client connections
- slave-endpoints -- List of sink etcd servers to connect with (scheme://ip:port), comma separated
- slave-cacert -- verify certificates of TLS-enabled secure servers using this CA bundle
- slave-cert -- identify secure client using this TLS certificate file
- slave-key -- identify secure client using this TLS key file
- slave-insecure-skip-tls-verify -- skip server certificate verification
- slave-insecure-transport -- disable transport security for client connections
- db-path -- the path where kv-db stores data
- bind-address -- the address the metric endpoint and ready/healthz binds to
- bind-port -- the port on which to serve restful

## Custom rules

At present, it mainly supports the synchronization of K8s resources, and more data sources will be supported in the future.
If no synchronization rule is specified, running etcd-carry will not synchronize any data on the source cluster to the sink cluster.
Synchronization rules are in the form of yaml, the format is as follows:

```yaml
filters:
  sequential: []
  secondary: []
```
`sequential` is used to configure resources that require priority order synchronization; `secondary` is used to configure resources that do not require priority, and these will be synchronized after `sequential` finished.

### Synchronize high-priority resources sequentially

Example 1:
```yaml
filters:
  sequential:
    - group: apiextensions.k8s.io
      resources:
        - group: apiextensions.k8s.io
          version: v1beta1
          kind: CustomResourceDefinition
    - group: ""
      resources:
        - version: v1
          kind: Namespace
  secondary: []
```
According to the rules, the CRDs will be synchronized first, and then the Namespaces will be synchronized.

Example 2:
```yaml
filters:
  sequential:
    - group: ""
      resources:
        - version: v1
          kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-name
              operator: In
              values:
                - test1
                - test2
  secondary: []
```
According to the rules, Namespaces with the label `test.io/namespace-name:test1` or `test.io/namespace-name:test2` will be synchronized first.

### Synchronize low-priority resources

Example 1:
```yaml
filters:
  sequential:
    - group: ""
      resources:
        - version: v1
          kind: Namespace
  secondary:
    - group: "monitoring.coreos.com"
      resources:
        - group: monitoring.coreos.com
          version: v1alpha1
          kind: AlertmanagerConfig
        - group: monitoring.coreos.com
          version: v1
          kind: PrometheusRule
        - group: monitoring.coreos.com
          version: v1
          kind: ServiceMonitor
        - group: monitoring.coreos.com
          version: v1
          kind: PodMonitor
      namespaceSelectors:
        - matchExpressions:
            - key: test.io/namespace-name
              operator: In
              values:
                - test1
                - test2
```
According to the rules, all Namespaces will be synchronized first, and then the `AlertmanagerConfig`, `PrometheusRule`, `ServiceMonitor`, and `PodMonitor` resources belonging to group `monitoring.coreos.com` in the specific namespaces will be synchronized. The relevant namespaces must be labeled `test.io/namespace-name:test1` or `test.io/namespace-name:test2`.

Example 2:
```yaml
filters:
  sequential:
    - group: ""
      resources:
        - version: v1
          kind: Namespace
  secondary:
    - group: ""
      resources:
        - version: v1
          Kind: Secret
      namespaceSelectors:
        - matchExpressions:
            - key: test.io/namespace-name
              operator: In
              values:
                - test1
                - test2
      fieldSelectors:
        - matchExpressions:
            - key: type
              operator: NotIn
              values:
                - kubernetes.io/service-account-token
```
According to the rules, all Namespaces will be synchronized first, and then the Secrets whose `type` is not `kubernetes.io/service-account-token` in the specific namespaces will be synchronized. The relevant namespaces must be labeled `test .io/namespace-name:test1` and `test.io/namespace-name:test2`.

Example 3:
```yaml
filters:
  sequential:
    - group: ""
      resources:
        - version: v1
          kind: Namespace
  secondary:
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
        - version: v1
          Kind: Secret
      namespace: test1
      labelSelectors:
        - matchExpressions:
            - key: test.io/configuration
              operator: Exists
        - matchExpressions:
            - key: test.io/credential
              operator: Exists
        - matchExpressions:
            - key: repo.test.io/docker-registry
              operator: Exists
      excludes:
        - resource:
            version: v1
            kind: Secret
          name: test-password-secret
          namespace: test1
```
According to the rules, all Naemspaces will be synchronized first, then all ConfigMaps and Secrets that were labeled with any of `test.io/configuration`, `test.io/credential` and `repo.test.io/docker-registry` in namespace test1 will be synchronized, except the Secret test-password-secret will not be synchronized.

Example4:
```yaml
filters:
  sequential:
    - group: ""
      resources:
        - version: v1
          kind: Namespace
  secondary:
    - group: "*.test.io"
      namespace: test1
      excludes:
        - resource:
            group: "app.test.io"
            version: v1beta1
            kind: TestCell
          labelSelectors:
            - matchExpressions:
                - key: app.test.io/required-on-controller
                  operator: In
                  values:
                    - "true"
```
According to the rules, all Naemspaces will be synchronized first, then synchronize the resources whose group name matches the `*.test.io` regular rule in namespace test1, except the `app.test.io/v1beta1 TestCell` that are labeled with `app.test.io/required-on-controller:true`.

## quick start

### Checking out source code

```shell
$ git clone https://github.com/etcd-carry/etcd-carry.git
```

### Build/Package

```shell
$ cd etcd-carry
$ make
```
### Usage

```shell
$ ./bin/etcd-carry --help

A simple command line for etcd mirroring

Usage:
  etcd-carry [flags]

Generic flags:

      --debug                enable client-side debug logging
      --mirror-rule string   Specify the rules to start mirroring (default "/etc/mirror/rules.yaml")
      --mode string          running mode, standalone or active-standby (default "standalone")

Etcd flags:

      --encryption-provider-config string   The file containing configuration for encryption providers to be used for storing secrets in etcd (default "/etc/mirror/secrets-encryption.yaml")
      --kube-prefix string                  the prefix to all kubernetes resources passed to etcd (default "/registry")
      --max-txn-ops uint                    Maximum number of operations permitted in a transaction during syncing updates (default 128)
      --rev int                             Specify the kv revision to start to mirror

Transport flags:

      --dial-timeout duration             dial timeout for client connections (default 2s)
      --keepalive-time duration           keepalive time for client connections (default 2s)
      --keepalive-timeout duration        keepalive timeout for client connections (default 6s)
      --master-cacert string              verify certificates of TLS-enabled secure servers using this CA bundle (default "/etc/kubernetes/master/etcd/ca.crt")
      --master-cert string                identify secure client using this TLS certificate file (default "/etc/kubernetes/master/etcd/server.crt")
      --master-endpoints strings          List of etcd servers to connect with (scheme://ip:port), comma separated
      --master-insecure-skip-tls-verify   skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)
      --master-insecure-transport         disable transport security for client connections (default true)
      --master-key string                 identify secure client using this TLS key file (default "/etc/kubernetes/master/etcd/server.key")
      --slave-cacert string               Verify certificates of TLS enabled secure servers using this CA bundle for the destination cluster (default "/etc/kubernetes/slave/etcd/ca.crt")
      --slave-cert string                 Identify secure client using this TLS certificate file for the destination cluster (default "/etc/kubernetes/slave/etcd/server.crt")
      --slave-endpoints strings           List of etcd servers to connect with (scheme://ip:port) for the destination cluster, comma separated
      --slave-insecure-skip-tls-verify    skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)
      --slave-insecure-transport          Disable transport security for client connections for the destination cluster (default true)
      --slave-key string                  Identify secure client using this TLS key file for the destination cluster (default "/etc/kubernetes/slave/etcd/server.key")

KeyValue flags:

      --db-path string   the path where kv-db stores data (default "/var/lib/mirror/db")

Daemon flags:

      --bind-address ip   the address the metric endpoint and ready/healthz binds to (default 0.0.0.0)
      --bind-port int     the port on which to serve restful (default 10520)
```
### Test the etcd-carry

Prepare the cert and key required to access the source etcd and sink etcd. For example, the path of the cert and key file is as follows:
```shell
/etc/etcd-carry/
├── master
│   ├── ca.crt
│   ├── server.crt
│   └── server.key
└── slave
    ├── ca.crt
    ├── server.crt
    └── server.key
```

Apply the yamls in `./deploy/examples/kube` direcotry to source K8s cluster.
```shell
# copy yamls in ./deploy/examples/kube directory to the master node of source K8s cluster
scp -r ./deploy/examples/kube {user}@{source_K8s_cluster_ip}:/opt/
# apply to source K8s cluster
kubectl apply -f /opt/kube/
```

Execute the following command to start synchronizing test resources that match the rules.
```shell
./bin/etcd-carry --master-cacert=/etc/etcd-carry/master/ca.crt --master-cert=/etc/etcd-carry/master/server.crt --master-key=/etc/etcd-carry/master/server.key --master-endpoints=10.20.144.29:2379 --slave-cacert=/etc/etcd-carry/slave/ca.crt --slave-cert=/etc/etcd-carry/slave/server.crt --slave-key=/etc/etcd-carry/slave/server.key --slave-endpoints=192.168.48.220:2379 --encryption-provider-config=./deploy/examples/secrets-encryption.yaml --mirror-rule=./deploy/examples/rules.yaml

```

According to the rules in ./deploy/examples/rules.yaml, only the following resources should be synchronized to the sink cluster:
```shell
namespace unique
namespace unit-test
secret named influxdb in namespace unique
configmap named influxdb1 in namespace unit-test
```
