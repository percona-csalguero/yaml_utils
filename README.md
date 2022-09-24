# yaml_utils

This project is only a POC about how to interact with yaml files from the command line.
Currently lets you read a base yaml file and then:
- merge other yaml files into it.
- set new values for specific paths.


## Flags

|Flag name|Short|Description|
|-----|-----|
|--base-yaml||Base yaml file (required)|
|--merge||Yaml files to merge onto the base. This argument can be repeated to specify more than one file|
|--set-path|Path and value to modify from the base yaml. Format `--set-path='path.to.key=value'`|
|--output-format|-o|Sets output format: json, yaml|

## Usage example

Suppose you have these yaml files:

```yaml
#base.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: perconaservermongodbs.psmdb.percona.com
spec:
  group: psmdb.percona.com
  names:
    kind: PerconaServerMongoDB
    listKind: PerconaServerMongoDBList
    plural: perconaservermongodbs
    singular: perconaservermongodb
    shortNames:
      - psmdb
  scope: Namespaced
```

and

```yaml
# to-merge.yaml
spec:
  names:
    shortNames:
    - psmdb2
    singular: zapp brannigan
```
and you want to:
1. Merge the `to-merge.yaml` file into `base.yaml`
2. Change the value under `spec.scope` key to `cluster-scoped`.

then:
```
./yaml_utils --base-yaml base.yaml --merge=to-merge.yaml --set-path='spec.scope=cluster-scoped'
```

The result will be:
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: perconaservermongodbs.psmdb.percona.com
spec:
  group: psmdb.percona.com
  names:
    kind: PerconaServerMongoDB
    listKind: PerconaServerMongoDBList
    plural: perconaservermongodbs
    shortNames:
    - psmdb2
    singular: zapp brannigan
  scope: cluster-scoped
```

## Features that might be added.
- Enable `--set-path` to add new fields.
- Add the `--unset-path` flag.
