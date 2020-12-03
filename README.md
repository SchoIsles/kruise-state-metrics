## Introduction

kruise-state-metrics is a prometheus-exporter for workflow metrics of [Openkruise](https://github.com/openkruise/kruise)

Based on [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)

### Installation

#### Docker

```bash
docker pull okletswin/kruise-state-metrics
docker run --rm -p 8080:8080 -p 8081:8081 okletswin/kruise-state-metrics
```

### Limits
Currently only `cloneset` workflow is provided