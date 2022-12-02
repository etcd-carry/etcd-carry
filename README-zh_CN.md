# etcd-carry

[English](./README.md) | 简体中文

## 介绍

etcd-carry提供一种能力将K8s集群中符合自定义规则的资源实时同步到备用k8s集群。

## 状态

该项目仍然是非常实验性的，生产环境中应谨慎使用，目前还不支持Service资源的正确同步，不支持同步有状态服务在PV中的数据，建议使用前对备用K8s集群的etcd做好备份。

## 安装

支持二进制部署，也支持K8s Deployment一键部署

## 运行参数

- debug -- 开启etcd客户端侧debug日志输出
- mirror-rule -- 自定义同步规则配置文件路径
- encryption-provider-config -- 为k8s secret资源存入etcd时提供加密服务的配置信息
- kube-prefix -- k8s资源存入etcd时key使用的前缀，默认为"/registry"
- max-txn-ops -- etcd-carry每次事务提交时包含的最大操作数，不能超过etcd服务端对应的设置，默认为128
- rev -- 开始同步时指定etcd的key reversion，默认为0
- dial-timeout -- etcd-carry和服务端建立连接的超时时间
- keepalive-time -- etcd-carry持续连接(长连接)保活时间
- keepalive-timeout -- etcd-carry持续连接(长连接)保活超时时间
- master-endpoints -- -master集群的etcd节点列表地址信息，多个节点用逗号隔开
- master-cacert -- master集群的CA证书信息
- master-cert -- etcd-carry和master集群建立安全连接的证书信息
- master-key -- etcd-carry和master集群建立安全连接的密钥信息
- master-insecure-skip-tls-verify -- 是否跳过对master证书合法性校验
- master-insecure-transport -- 未指定证书密钥时，是否和master建立非安全连接
- slave-endpoints -- slave集群的etcd节点列表地址信息，多个节点用逗号隔开
- slave-cacert -- slave集群的CA证书信息
- slave-cert -- etcd-carry和slave集群建立安全连接的证书信息
- slave-key -- etcd-carry和slave集群建立安全连接的密钥信息
- slave-insecure-skip-tls-verify -- 是否跳过对slave证书合法性校验
- slave-insecure-transport -- 未指定证书密钥时，是否和slave建立非安全连接
- db-path -- rocksdb数据目录
- bind-address -- metric/ready/healthz绑定的地址
- bind-port -- 绑定的端口

## 自定义同步规则

目前主要支持对K8s资源的同步，后续会进行扩展以支持更多不同数据源。
未指定同步规则时，运行etcd-carry是不会将master集群上的任何数据同步到slave集群的。
同步规则配置为yaml形式，格式如下：

```yaml
filters:
  sequential: []
  secondary: []
```
`sequential`用于配置需要优先顺序同步的资源；`secondary`用于配置对优先级没有要求的资源，再同步完`sequential`中指定的资源后才开始同步该部分的资源。

### 按顺序同步高优先级资源

示例1：
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
根据上面的规则，将会优先同步CustomResourceDefinition资源，再同步Namespace资源。

示例2：
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
根据上面的规则，将优先同步带有`test.io/namespace-name:test1`标签或`test.io/namespace-name:test2`标签的Namespace资源

### 低优先级资源同步

示例1：
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
根据上面的规则，将优先同步所有Namespace资源，之后才开始同步特定namespace下属于`monitoring.coreos.com`资源组的`AlertmanagerConfig`,`PrometheusRule`,`ServiceMonitor`,`PodMonitor`资源，符合条件的namespace必须带有标签`test.io/namespace-name:test1`或`test.io/namespace-name:test2`。

示例2：
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
根据上面的规则，将优先同步所有Namespace资源，之后才开始同步特定namespace下`type`字段不为`kubernetes.io/service-account-token`的Secret资源，符合条件的namespace必须带有标签`test.io/namespace-name:test1`和`test.io/namespace-name:test2`。

示例3：
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
根据上面的规则，将优先同步所有Naemspace资源，之后才开始同步test1命名空间下带有`test.io/configuration:`、`test.io/credential:`和`repo.test.io/docker-registry`其中任意标签的所有Confimap和Secret资源，除了test-password-secret这个Secret资源不会被同步。

示例4：
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
根据上面的规则，将优先同步所有Namespace资源，之后才开始同步test1命名空间下所属资源组名称符合`*.test.io`正则规则的任意资源，除了带有`app.test.io/required-on-controller:true`标签的所有`app.test.io/v1beta1 TestCell`资源不会被同步。

## 快速开始

克隆etcd-carry：
```shell
$ git clone https://github.com/etcd-carry/etcd-carry.git
```

二进制编译：
```shell
$ cd etcd-carry
$ make
```

下面的命令会输出帮助信息：
```shell
$ ./bin/etcd-carry --help
```
下面是使用详情的示例：
```shell
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

准备访问K8s master集群和slave集群etcd组件所需的证书密钥信息，例如证书密钥文件的路径如下：
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

通过./deploy/examples目录下的示例来演示同步过程，在K8s master集群创建./deploy/examples/kube目录下的K8s资源：
```shell
# 拷贝./deploy/examples/kube目录到master集群的master节点上
scp -r ./deploy/examples/kube root@{K8s master集群的master节点IP}:/opt/
# 创建测试资源
kubectl apply -f /opt/kube/
```

执行以下命令开始同步符合自定义规则的测试资源：
```shell
./bin/etcd-carry --master-cacert=/etc/etcd-carry/master/ca.crt --master-cert=/etc/etcd-carry/master/server.crt --master-key=/etc/etcd-carry/master/server.key --master-endpoints=10.20.144.29:2379 --slave-cacert=/etc/etcd-carry/slave/ca.crt --slave-cert=/etc/etcd-carry/slave/server.crt --slave-key=/etc/etcd-carry/slave/server.key --slave-endpoints=192.168.48.220:2379 --encryption-provider-config=./deploy/examples/secrets-encryption.yaml --mirror-rule=./deploy/examples/rules.yaml
```
按照./deploy/examples/rules.yaml中的自定义规则，最终应该只有以下资源会被同步到slave集群：
```shell
命名空间unique
命名空间unit-test
命名空间unique下名为influxdb的secret资源
命名空间unit-test下名为influxdb1的configmap资源
```

