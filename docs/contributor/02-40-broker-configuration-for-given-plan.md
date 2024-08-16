# Kyma Environment Broker Configuration for a Given Kyma Version and Plan

> [!NOTE] 
> This document has been deprecated. Kyma version configuration is not possible.

Some processes in Kyma Environment Broker (KEB) can be configured to deliver different results. KEB needs a ConfigMap with a configuration for the given Kyma plan to process the requests. 
A default configuration must be defined. This configuration must be recognized by KEB as applicable for all the supported plans. You can also set a separate configuration for each plan.
  
While processing requests, KEB reads a configuration from a ConfigMap which holds data for a given plan.
> [!NOTE] 
> Create all configurations in the `kcp-system` namespace.

> [!NOTE] 
> Currently, only the Kyma custom resource template and the additional components list can be configured.

## ConfigMap  

The appropriate ConfigMap is selected by filtering the resources using labels. KEB recognizes the ConfigMaps with configurations when they contain the label:

```yaml
keb-config: "true"
```

> [!NOTE] 
> Each ConfigMap that defines a configuration must have both labels assigned.

The actual configuration is stored in ConfigMap's `data` object. Add `default` key under `data` object:

```yaml
data:
  default: |-
    kyma-template: |-
      apiVersion: operator.kyma-project.io/v1beta2
      kind: Kyma
      metadata:
        labels:
          "operator.kyma-project.io/managed-by": "lifecycle-manager"
          "operator.kyma-project.io/beta": "true"
        name: tbd
        namespace: kcp-system
      spec:
        channel: regular
        modules:
        - name: module1
        - name: module2
    additional-components:
      - name: "additional-component1"
        namespace: "kyma-system"
```

You must define a default configuration that is selected when the supported plan key is missing. This means that, for example, if there are no other plan keys under the `data` object, the default configuration applies to all the plans. You do not have to change `tbd` value of `kyma-template.metadata.name` field as KEB generates the name for Kyma CR during provisioning operation.

> [!NOTE] 
> `kyma-template` and `additional-components` configurations are required.

See the example of a ConfigMap with a default configuration for Kyma and specific configurations for `plan1`, `plan2`, and `trial` plans:

```yaml
# keb-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: keb-config
  labels:
    keb-config: "true"
data:
  default: |-
    kyma-template: |-
      apiVersion: operator.kyma-project.io/v1beta2
      kind: Kyma
      metadata:
        labels:
          "operator.kyma-project.io/managed-by": "lifecycle-manager"
          "operator.kyma-project.io/beta": "true"
        name: tbd
        namespace: kcp-system
      spec:
        channel: regular
        modules:
        - name: api-gateway
        - name: istio
  plan1: |-
    kyma-template: |-
      apiVersion: operator.kyma-project.io/v1beta2
      kind: Kyma
      metadata:
        labels:
          "operator.kyma-project.io/managed-by": "lifecycle-manager"
          "operator.kyma-project.io/beta": "true"
        name: tbd
        namespace: kcp-system
      spec:
        channel: fast
        modules:
        - name: api-gateway
        - name: istio
        - name: btp-operator
  plan2: |-
    kyma-template: |-
      apiVersion: operator.kyma-project.io/v1beta2
      kind: Kyma
      metadata:
        labels:
          "operator.kyma-project.io/managed-by": "lifecycle-manager"
          "operator.kyma-project.io/beta": "true"
        name: tbd
        namespace: kcp-system
      spec:
        channel: fast
        modules:
        - name: api-gateway
        - name: istio
        - name: btp-operator
        - name: keda
        - name: serverless
    additional-components:
      - name: "additional-component2"
        namespace: "kyma-system"
      - name: "additional-component3"
        namespace: "kyma-system"
        source:
          url: "https://example.source.url.local/artifacts/additional-component3-0.0.1.tgz"
  trial: |-
    kyma-template: |-
      apiVersion: operator.kyma-project.io/v1beta2
      kind: Kyma
      metadata:
        labels:
          "operator.kyma-project.io/managed-by": "lifecycle-manager"
          "operator.kyma-project.io/beta": "true"
        name: tbd
        namespace: kcp-system
      spec:
        channel: regular
        modules: []
```
