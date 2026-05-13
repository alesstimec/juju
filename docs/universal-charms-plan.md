# Universal Charms — Design Plan

## Concept

Allow FormatV2 (K8s/CAAS) charms to deploy on IAAS (VM) models unchanged.
Juju runs the OCI workload containers on the VM using containerd/nerdctl.
The existing `PebblePoller` fires workload hooks identically on both substrates.
No `metadata.yaml` changes needed by charm authors.

## Key Insight

The `PebblePoller` in `internal/worker/uniter/pebblepoller.go` already polls
`/charm/containers/<name>/pebble.socket`. On K8s this socket is created by
pebble running inside a K8s sidecar container. On a VM, Juju runs the same
OCI image via `nerdctl` with a bind mount; pebble creates the same socket path
and the poller works without any changes.

**Important:** On IAAS, each unit's agent runs in its own data directory
(e.g., `/var/lib/juju/agents/unit-myapp-0`). The socket path
`/charm/containers/<name>/pebble.socket` is relative to that data directory, so
multiple units on the same machine do not collide. The `iaascontainerrunner`
must use `<dataDir>/charm/containers/<name>/` as the bind-mount source.

## Current Verification Status (May 2026)

What has been verified in code and tests:

- The IAAS uniter starts Pebble workers when `ContainerNames` is non-empty.
- `PebblePoller` emits `<container>-pebble-ready` when Pebble boot ID changes.
- `PebbleNoticer` maps Pebble `perform-check` and `recover-check` changes to
    `pebble_check_failed` and `pebble_check_recovered` hooks.
- The workload hook resolver maps these events to `hooks.PebbleReady`,
    `hooks.PebbleCheckFailed`, and `hooks.PebbleCheckRecovered`.

What is still pending in integration:

- End-to-end `universal_charms` confirmation in a clean run that reaches
    unit agent `idle`.
- Runs interrupted while unit stayed in `installing agent`; in that state the
    charm has not started yet, so `pebble services` correctly shows
    "Plan has no services".

## Background: How CAAS Works Today

- The K8s provider builds a pod spec in `internal/provider/kubernetes/application/application.go`.
- For each workload container, it adds a `corev1.Container` running `/charm/bin/pebble run`
  with a bind-mounted shared volume so the charm container can reach the pebble socket at
  `/charm/containers/<name>/pebble.socket`.
- The `JUJU_CONTAINER_NAMES` env var is set in the pod spec, which is read by
  `cmd/containeragent/unit/agent.go` and passed as `ContainerNames` to the uniter manifold.
- When `len(u.containerNames) > 0`, `uniter.go` starts `PebblePoller` and `PebbleNoticer`,
  which fire `<name>-pebble-ready` and related workload hooks.

## Important Note: No Deploy Block Exists

**Note:** There is no explicit code block rejecting FormatV2 charms on IAAS models.
The existing code only rejects FormatV1 charms on CAAS models (see
`apiserver/facades/client/application/deploy.go` lines 113-119). FormatV2 charms
can technically already be deployed on IAAS — they just don't work because:

1. `ContainerNames` is never populated for IAAS units
2. No container runtime starts the OCI workloads
3. Pebble sockets never appear at the expected paths

This plan adds the infrastructure to **make** FormatV2 charms work on IAAS,
not remove a gate that blocks them.

## Compatibility Requirements

This feature MUST NOT break existing behaviour:

| Charm Format | Model Type | Behaviour | Changed? |
|---|---|---|---|
| FormatV1 | IAAS | Traditional VM charm, no containers | NO |
| FormatV1 | CAAS | Rejected at deploy time | NO |
| FormatV2 | CAAS | K8s sidecar containers via pod spec | NO |
| FormatV2 | IAAS | **NEW**: nerdctl containers on VM | YES |

**Sidecar flag:** The uniter's `Sidecar bool` field MUST remain `false` for IAAS
universal charms. `Sidecar` is only `true` when running inside a K8s sidecar
container (set by `cmd/containeragent`). The IAAS uniter path is correct for
universal charms — hooks run on the host; workload access is via pebble socket.

## Required Changes

### 1. Distribute Pebble Binary with Juju Agent

On K8s, the pebble binary is available at `/charm/bin/pebble` via a shared volume
populated by the charm container. On IAAS VMs, we must distribute pebble separately.

**Approach:** Package pebble binary with the juju snap/deb at `/usr/lib/juju/bin/pebble`.
At unit deploy time, copy to `/charm/bin/pebble` on the machine.

**Files to modify:**
- `snap/snapcraft.yaml` — include pebble binary in snap
- Packaging scripts for deb/rpm builds
- Unit deployment code to copy pebble to `/charm/bin/`

---

### 2. Populate `ContainerNames` for IAAS units

**File: `internal/worker/deployer/deployer.go`**

When starting an IAAS unit worker, read the unit's charm metadata from the API
and extract container names from `meta.Containers`:

```go
containerNames := make([]string, 0, len(meta.Containers))
for name := range meta.Containers {
    containerNames = append(containerNames, name)
}
sort.Strings(containerNames)
```

Pass these into `UnitManifoldsConfig.ContainerNames`. This causes the existing
`PebblePoller` and `PebbleNoticer` to start for IAAS units that have containers,
polling the same `/charm/containers/<name>/pebble.socket` paths as on CAAS.

---

### 3. New worker: `internal/worker/iaacontainerrunner/`

This is the only substantively new code. It runs the OCI workload containers on
the VM using `nerdctl`.

**Responsibilities:**

1. Look up the OCI image resource name from `charm.Meta().Containers[name].Resource`.
2. Fetch image details via `ResourcesAPI` (registry path + credentials).
3. `mkdir -p /charm/containers/<name>` and `/charm/bin` on host.
4. Pull the image: `nerdctl pull <image>`.
5. Run the container with proper mounts, env vars, and entrypoint (see corrected command below).
6. Map storage mounts from charm metadata to `-v` flags.
7. On unit removal: `nerdctl stop juju-<unitname>-<name> && nerdctl rm ...`
8. On charm upgrade (new OCI image resource): stop, pull new image, restart.
9. On worker restart: check if container already running via `nerdctl inspect`.
10. Watch container health via `nerdctl events` or periodic inspect; restart crashed containers.

**Corrected nerdctl command:**

```bash
nerdctl run -d \
  --name juju-<unitname>-<name> \
  --network host \
  --restart unless-stopped \
  -v <dataDir>/charm/bin/pebble:/charm/bin/pebble:ro \
  -v <dataDir>/charm/containers/<name>:/charm/container \
  -v <storage-host-path>:<mount-location> \  # for each storage mount
  -e JUJU_CONTAINER_NAME=<name> \
  -e PEBBLE_SOCKET=/charm/container/pebble.socket \
  -e PEBBLE=/charm/bin/pebble \
  -e PEBBLE_COPY_ONCE=/var/lib/pebble/default \
  --entrypoint /charm/bin/pebble \
  <image> \
  run --create-dirs --hold --http "" --verbose
```

Where `<dataDir>` is the unit agent's data directory, e.g.,
`/var/lib/juju/agents/unit-myapp-0`.

**Key differences from original plan:**
- Added pebble binary mount: `-v <dataDir>/charm/bin/pebble:/charm/bin/pebble:ro`
- Added `--entrypoint /charm/bin/pebble` to override container entrypoint
- Added pebble run arguments: `run --create-dirs --hold --http "" --verbose`
- Added `--network host` for simple networking
- Added `--restart unless-stopped` for reboot resilience (see Section 3a)
- Added storage mount placeholders
- Fixed `PEBBLE_COPY_ONCE` to use concrete path `/var/lib/pebble/default`
- **Disabled pebble HTTP listener** (`--http ""`) to avoid port conflicts when
  multiple units/containers run on the same machine. Communication is via the
  unix socket only.
- **Unit-scoped paths**: All host paths use `<dataDir>` prefix to avoid collisions
  between multiple units on the same machine

**Interface `ResourcesAPI`** (defined in the new package):

```go
type ResourcesAPI interface {
    GetContainerImageInfo(ctx context.Context, unitTag names.UnitTag, resourceName string) (ImageDetails, error)
}

type ImageDetails struct {
    RegistryPath string
    Username     string
    Password     string
}
```

**Interface `StorageAPI`** (for storage mount resolution):

```go
type StorageAPI interface {
    GetStorageAttachmentPath(ctx context.Context, unitTag names.UnitTag, storageName string) (string, error)
}
```

**Interface `CharmMetaAPI`** (for container mount definitions):

```go
type CharmMetaAPI interface {
    GetContainerMounts(ctx context.Context, unitTag names.UnitTag, containerName string) ([]Mount, error)
}

type Mount struct {
    StorageName string
    Location    string
}
```

---

### 4. Wire the new worker

**File: `internal/worker/deployer/unit_manifolds.go`**

Add `iaaContainerRunnerName` manifold gated on `len(containerNames) > 0 && modelType == IAAS`.
The uniter manifold should depend on `iaaContainerRunnerName` so the containers are running
and the pebble socket is present before the `PebblePoller` starts polling.

---

### 5. Storage Mounts for Universal Charms

Universal charms may declare storage requirements with container mounts:

```yaml
storage:
  database:
    type: filesystem
    minimum-size: 1G
containers:
  myworkload:
    resource: myworkload-image
    mounts:
      - storage: database
        location: /data/db
```

**How it works on K8s:** The K8s provider creates PersistentVolumeClaims and adds
VolumeMount entries to the pod spec.

**How it works on IAAS VMs:**

1. Existing IAAS storage provisioner handles storage creation and attachment
   (loop devices, cloud volumes, etc.) — no changes needed.
2. Storage is attached to the machine at a host path (e.g., `/var/lib/juju/storage/database/0`).
3. `iaacontainerrunner` reads charm metadata `containers[].mounts[]`.
4. Worker calls `StorageAPI.GetStorageAttachmentPath()` to resolve storage name → host path.
5. Worker adds `-v <host-path>:<mount-location>` flags to nerdctl command.

**Block storage handling:**

When a charm declares `type: block` storage, the IAAS storage provisioner attaches
a block device to the machine. For universal charms, `iaacontainerrunner`:

1. Detects block storage type from charm metadata.
2. Creates a filesystem on the block device (if not already formatted): `mkfs.ext4 /dev/sdb`.
3. Mounts the filesystem to a host path: `mount /dev/sdb /var/lib/juju/storage/myblock/0`.
4. Passes the mount path to the container as a normal `-v` bind mount.

This approach treats block storage as filesystem storage from the container's perspective,
which is the common use case. Raw block device passthrough to containers is not supported.

**Shared storage handling:**

When a charm declares `shared: true` storage, multiple units need access to the same data.
`iaacontainerrunner` uses nerdctl named volumes for shared storage:

1. For the first unit on a machine with shared storage:
   ```bash
   nerdctl volume create juju-<model>-<storage-name>
   ```
2. All containers mounting this storage use the named volume:
   ```bash
   nerdctl run -v juju-<model>-<storage-name>:<mount-location> ...
   ```
3. For cross-machine shared storage (e.g., NFS), the storage provisioner mounts
   the shared filesystem to a host path, and containers bind-mount that path.

Named volumes ensure data persists across container restarts and can be shared
between multiple containers on the same machine.

**New API method needed:**

```go
// GetContainerMounts returns the storage mounts declared for a container.
func (s *ApplicationService) GetContainerMounts(ctx context.Context, unitTag names.UnitTag, containerName string) ([]Mount, error) {
    // Get charm metadata for unit
    // Return containers[containerName].mounts
}

// StorageInfo includes type and sharing information for mount handling.
type StorageInfo struct {
    Name       string
    Type       string // "filesystem" or "block"
    Shared     bool
    DevicePath string // for block storage
    MountPath  string // host path where storage is mounted
}
```

---

### 6. Runtime Installation and Feature Detection

**How runtime tools get onto the machine:**

Runtime prerequisites are installed during machine provisioning in cloud-init,
not by the worker at runtime.

- Base packages are added (`containerd`, `curl`, `tar`).
- A provisioning script installs `nerdctl` and `pebble` to
  `/usr/lib/juju/bin/`.
- `containerd` is enabled via systemd.

**Important reliability adjustment:**

The runtime provisioning script is best-effort and bounded. Downloads have
timeouts/retries, and failures are logged as warnings instead of aborting
cloud-init. This avoids blocking machine-agent startup in `installing agent`
when external downloads are slow or unavailable.

**Worker behavior:**

The `iaascontainerrunner` now expects a runtime CLI to already exist and only
detects available `nerdctl` binaries. If none are found, it fails with a clear
error.

**File: `domain/application/service/provider.go` — `GetSupportedFeatures()`**

For IAAS models, advertise a `containerd` feature in the FeatureSet when snapd
is available on the target machine (since nerdctl can always be snap-installed).
Charm authors declare:

```yaml
assumes:
  - containerd
```

**Detection approach:** Since the worker installs nerdctl via snap automatically,
the `containerd` feature can be advertised on any machine where `snapd` is running
(i.e., Ubuntu-based machines). Machines without snapd cannot support universal charms.

**Important caveat:** `GetSupportedFeatures` runs on the controller at deploy time,
not on the target machine. The controller cannot know if snapd is available on every
machine. The recommended approach is to advertise `containerd` unconditionally for
IAAS models and fail at runtime with a clear error in `juju status` if nerdctl cannot
be installed.

**New feature in `core/assumes/features.go`:**

```go
var UserFriendlyFeatureDescriptions = map[string]string{
    // ... existing ...
    "containerd": "the containerd runtime is available for running OCI workload containers (installed via nerdctl snap)",
}

func ContainerdFeature() Feature {
    return Feature{Name: "containerd"}
}
```

---

### 7. `juju status` — Container Status for IAAS Units

Surface per-container status for IAAS universal units in `juju status`, mirroring
the CAAS workload status display. Query container state via `nerdctl inspect` or
the containerd gRPC API.

**Status states:**
- `waiting` — Container not yet started
- `active` — Container running
- `blocked` — Container crashed or in error state
- `unknown` — Cannot determine container state

---

### 8. OCI Image Update on Charm Refresh

When a charm is refreshed (`juju refresh`) with a new OCI image resource, the
`iaascontainerrunner` must replace running containers with the new image.

**Mechanism:** The uniter processes upgrade-charm and bounces the unit dependency
engine. When the `iaascontainerrunner` restarts, it compares the running container's
image (via `nerdctl inspect --format {{.Image}}`) against the expected image from
`ResourcesAPI`. If they differ, it stops the old container and starts a new one.

**Sequence:**
1. `juju refresh myapp --resource myworkload-image=newregistry/newimage:v2`
2. Uniter processes upgrade-charm hook, bounces unit engine
3. `iaascontainerrunner` restarts, calls `ensureRunning()`
4. `ensureRunning()` detects image mismatch → `nerdctl stop` + `nerdctl rm` + re-run
5. Pebble socket disappears briefly; PebblePoller detects socket gone
6. New container starts, pebble creates socket, `<container>-pebble-ready` fires again

**In-flight pebble operations:** The pebble socket disappears during container
replacement. PebblePoller handles this naturally — it polls and will detect the new
socket when the new container creates it. The charm's pebble-ready hook fires again.

---

### 9. Pebble Binary Upgrade on Agent Upgrade

When the juju agent is upgraded, the pebble binary in the snap changes. Running
containers still have the old binary bind-mounted. The `iaascontainerrunner` must
detect the mismatch and restart containers with the new pebble.

**Mechanism:** On worker startup (after agent upgrade causes engine restart):
1. Compute SHA256 of source pebble binary (`/snap/juju/current/bin/pebble`)
2. Compute SHA256 of deployed pebble binary (`<dataDir>/charm/bin/pebble`)
3. If different: copy new binary, then restart all containers (stop + start)

This is safe because the worker only runs this check once per startup, not
continuously.

---

### 10. Container Log Forwarding

Container stdout/stderr must be forwarded into the Juju logging pipeline so
workload messages appear in `juju debug-log`.

**Mechanism:** After a container starts, launch a goroutine that runs
`nerdctl logs --follow --timestamps <containerID>` and parses each line into the
Juju log sender. Lines are tagged with module
`unit.<unitname>.workload.<containername>`.

```go
type LogSink interface {
    Log(containerName string, timestamp time.Time, message string)
}
```

- stdout lines → `logger.INFO` level
- stderr lines → `logger.ERROR` level
- On container replacement (crash recovery, image update), cancel the old tail
  goroutine and start a new one

---

### 11. Machine Reboot Resilience

On CAAS, pod restart guarantees container restart. On IAAS, the nerdctl run command
MUST include `--restart=unless-stopped` so that containerd auto-restarts containers
after a machine reboot, even before the juju agent starts.

The `iaascontainerrunner` handles the case where containers are already running on
startup (idempotent `isRunning()` check). It also handles containers that exist but
are stopped (calls `nerdctl start` instead of `nerdctl run`).

---

## Files That Require Zero Changes

| File | Reason |
|---|---|
| `internal/worker/uniter/pebblepoller.go` | Already polls `/charm/containers/<name>/pebble.socket` |
| `internal/worker/uniter/pebblenotices.go` | Same |
| `internal/worker/uniter/uniter.go` (lines 832–839) | Already starts pollers when `len(containerNames) > 0` |
| `internal/worker/uniter/container/` (all workload hook dispatch) | Unchanged |
| `domain/deployment/charm/hooks/hooks.go` | `PebbleReady` etc. already exist |
| `domain/deployment/charm/meta.go` | No metadata additions needed |
| `domain/deployment/charm/resource/type.go` | No new resource types |
| `internal/worker/storageprovisioner/` | Existing IAAS storage provisioning works unchanged |
| The ops library / charm author hook code | Fully substrate-agnostic once socket exists |

## Key Reference Files

| File | Purpose |
|---|---|
| `internal/worker/uniter/pebblepoller.go` | Polls `/charm/containers/<name>/pebble.socket` |
| `internal/worker/uniter/manifold.go` | `ManifoldConfig.ContainerNames` already exists (line 58) |
| `internal/provider/kubernetes/application/application.go` | K8s pod spec builder — reference for bind mount and env vars (lines 1864-1906) |
| `cmd/containeragent/unit/agent.go` | Reads `JUJU_CONTAINER_NAMES` env var, source of `ContainerNames` on CAAS |
| `cmd/containeragent/unit/manifolds.go` | CAAS uniter manifold config with `ContainerNames` (line 129) |
| `internal/worker/deployer/unit_manifolds.go` | IAAS uniter manifold config — add `ContainerNames` here |
| `internal/worker/storageprovisioner/` | Existing IAAS storage attachment handling |
| `domain/deployment/charm/meta.go` | Container and Mount type definitions (lines 283-296) |
| `caas/application.go` | `ApplicationConfig` and `ContainerConfig` structs |

---

## Appendix: Proposed Code Changes

### A. `internal/worker/deployer/unit_manifolds.go`

**A.1 — Add `ContainerNames` to `UnitManifoldsConfig`**

```go
// UnitManifoldsConfig allows specialisation of the result of Manifolds.
type UnitManifoldsConfig struct {
    // ... existing fields ...

    // Clock supplies timekeeping services to various workers.
    Clock clock.Clock

+   // ContainerNames holds the names of workload containers declared in the
+   // charm's metadata.yaml containers: section. When non-empty, the uniter
+   // starts PebblePoller and PebbleNoticer so that workload hooks are fired.
+   ContainerNames []string
}
```

**A.2 — Pass `ContainerNames` to the uniter manifold, and add `iaaContainerRunnerName`**

```go
+   const iaaContainerRunnerName = "iaa-container-runner"

    // In UnitManifolds(), inside the returned dependency.Manifolds map:

+   // The iaaContainerRunner starts OCI workload containers on the VM using
+   // nerdctl so that pebble sockets appear at the paths PebblePoller expects.
+   // Only started on IAAS models with container-enabled charms.
+   iaaContainerRunnerName: engine.Housing{
+       Flags: []string{migrationInactiveFlagName},
+       Occupy: migrationFortressName,
+   }.Decorate(iaacontainerrunner.Manifold(iaacontainerrunner.ManifoldConfig{
+       AgentName:      agentName,
+       APICallerName:  apiCallerName,
+       ContainerNames: config.ContainerNames,
+       Logger:         config.LoggerContext.GetLogger("juju.worker.iaacontainerrunner"),
+   })),  // only installed when len(config.ContainerNames) > 0 && IAAS

    uniterName: ifNotMigrating(uniter.Manifold(uniter.ManifoldConfig{
        AgentName:             agentName,
        ModelType:             model.IAAS,
        APICallerName:         apiCallerName,
        S3CallerName:          s3CallerName,
        TraceName:             traceName,
        MachineLock:           config.MachineLock,
        Clock:                 config.Clock,
        LeadershipTrackerName: leadershipTrackerName,
        CharmDirName:          charmDirName,
        HookRetryStrategyName: hookRetryStrategyName,
        TranslateResolverErr:  uniter.TranslateFortressErrors,
        Logger:                config.LoggerContext.GetLogger("juju.worker.uniter"),
+       ContainerNames:        config.ContainerNames,
    })),
```

When `config.ContainerNames` is empty (a pure IAAS charm with no containers), neither
the new manifold nor PebblePoller starts, so existing behaviour is unchanged.

---

### B. `internal/worker/deployer/unit_agent.go`

**B.1 — Add `ContainerNames` to `UnitAgentConfig` and `UnitAgent`**

```go
 // UnitAgentConfig is a params struct with the values necessary to
 // construct a working unit agent.
 type UnitAgentConfig struct {
     Name             string
     DataDir          string
     FlightRecorder   flightrecorder.FlightRecorder
     Clock            clock.Clock
     Logger           logger.Logger
     UnitEngineConfig func() dependency.EngineConfig
     UnitManifolds    func(UnitManifoldsConfig) dependency.Manifolds
     SetupLogging     func(logger.LoggerContext, agent.Config)
+    ContainerNames   []string
 }
```

```go
 type UnitAgent struct {
     tag    names.UnitTag
     name   string
     clock  clock.Clock
     logger logger.Logger
     // ...
+    containerNames []string
 }
```

In `NewUnitAgent()`:
```go
 unit := &UnitAgent{
     // ... existing fields ...
+    containerNames: config.ContainerNames,
 }
```

**B.2 — Thread `ContainerNames` into `UnitManifoldsConfig` inside `start()`**

```go
 manifolds := a.unitManifolds(UnitManifoldsConfig{
     LoggerContext:       loggerContext,
     Agent:               a,
     LogSource:           bufferedLogger.Logs(),
     LeadershipGuarantee: 30 * time.Second,
     AgentConfigChanged:  a.configChangedVal,
     ValidateMigration:   a.validateMigration,
     UpdateLoggerConfig:  updateAgentConfLogging,
     MachineLock:         machineLock,
     Clock:               a.clock,
+    ContainerNames:      a.containerNames,
 })
```

---

### C. `internal/worker/deployer/nested.go`

The `nestedContext` must know each unit's container names before it starts the unit's
dependency engine. The cleanest approach is to add a callback to `ContextConfig` that
the deployer manifold satisfies using its existing API caller.

**C.1 — Add `GetContainerNames` callback to `ContextConfig`**

```go
 type ContextConfig struct {
     Agent                    agent.Agent
     FlightRecorder           flightrecorder.FlightRecorder
     Clock                    clock.Clock
     Logger                   logger.Logger
     UnitEngineConfig         func() dependency.EngineConfig
     SetupLogging             func(logger.LoggerContext, agent.Config)
     UnitManifolds            func(config UnitManifoldsConfig) dependency.Manifolds
     RebootMonitorStatePurger RebootMonitorStatePurger

+    // GetContainerNames returns the names of workload containers in the
+    // charm deployed for the given unit. Returns nil when the charm has
+    // no containers: section (pure IAAS charm).
+    GetContainerNames func(ctx context.Context, unitTag names.UnitTag) ([]string, error)
 }
```

Store the config on `nestedContext`:
```go
 type nestedContext struct {
     logger      logger.Logger
     agent       agent.Agent
     agentConfig agent.Config
     baseUnitConfig UnitAgentConfig
+    config      ContextConfig
     // ...
 }
```

**C.2 — Populate `ContainerNames` in `newUnitAgent()`**

```go
 func (c *nestedContext) newUnitAgent(unitName string) (*UnitAgent, error) {
     unitConfig := c.baseUnitConfig
     unitConfig.Name = unitName

+    if c.config.GetContainerNames != nil {
+        tag := names.NewUnitTag(unitName)
+        names, err := c.config.GetContainerNames(context.TODO(), tag)
+        if err != nil {
+            c.logger.Warningf(context.TODO(),
+                "could not get container names for %q: %v", unitName, err)
+        } else {
+            unitConfig.ContainerNames = names
+        }
+    }

     // ... existing engine config override code ...
     return NewUnitAgent(unitConfig)
 }
```

---

### D. `internal/worker/deployer/manifold.go`

Implement `GetContainerNames` using the deployer API facade's existing `apiCaller`.
This requires a new method on the deployer facade (see Change E below).

```go
 func (config ManifoldConfig) newWorker(_ context.Context, a agent.Agent, apiCaller base.APICaller) (worker.Worker, error) {
     // ...
     deployerFacade := apideployer.NewClient(apiCaller)
     contextConfig := ContextConfig{
         Agent:            a,
         FlightRecorder:   config.FlightRecorder,
         Clock:            config.Clock,
         Logger:           config.Logger,
         UnitEngineConfig: config.UnitEngineConfig,
         SetupLogging:     config.SetupLogging,
         UnitManifolds:    UnitManifolds,
+        GetContainerNames: func(ctx context.Context, unitTag names.UnitTag) ([]string, error) {
+            return deployerFacade.UnitContainerNames(ctx, unitTag)
+        },
     }
     // ...
 }
```

---

### E. `api/agent/deployer/` and `apiserver/facades/agent/deployer/`

A new `UnitContainerNames` RPC call is needed on both sides of the deployer facade so
that the machine agent can ask the controller "what containers does this unit's charm
declare?".

**E.1 — Client side: `api/agent/deployer/unit.go` (new method)**

```go
// ContainerNames returns the names of workload containers declared in the
// charm's metadata for this unit. Returns nil for FormatV1 (pure IAAS) charms.
func (u *Unit) ContainerNames(ctx context.Context) ([]string, error) {
    var result params.StringsResult
    args := params.Entity{Tag: u.tag.String()}
    if err := u.client.facade.FacadeCall(ctx, "UnitContainerNames", args, &result); err != nil {
        return nil, err
    }
    if result.Error != nil {
        return nil, result.Error
    }
    return result.Result, nil
}
```

**Add to the `Client`:**
```go
// UnitContainerNames returns the workload container names for the given unit.
func (c *Client) UnitContainerNames(ctx context.Context, tag names.UnitTag) ([]string, error) {
    unit, err := c.Unit(ctx, tag)
    if err != nil {
        return nil, err
    }
    return unit.ContainerNames(ctx)
}
```

**E.2 — Server side: add `UnitContainerNames` to the deployer facade**

In `apiserver/facades/agent/deployer/deployer.go`, add a new method that looks up the
unit's assigned charm via the application service and reads `meta.Containers`:

```go
// UnitContainerNames returns the workload container names for the given unit.
func (d *DeployerAPI) UnitContainerNames(ctx context.Context, args params.Entity) (params.StringsResult, error) {
    // authorise caller
    tag, err := names.ParseUnitTag(args.Tag)
    if err != nil {
        return params.StringsResult{Error: apiservererrors.ServerError(err)}, nil
    }
    charmMeta, err := d.applicationService.GetUnitCharmMeta(ctx, tag)
    if err != nil {
        return params.StringsResult{Error: apiservererrors.ServerError(err)}, nil
    }
    names := make([]string, 0, len(charmMeta.Containers))
    for name := range charmMeta.Containers {
        names = append(names, name)
    }
    sort.Strings(names)
    return params.StringsResult{Result: names}, nil
}
```

---

### F. New package: `internal/worker/iaacontainerrunner/`

**F.1 — `worker.go` skeleton**

```go
// Package iaacontainerrunner runs OCI workload containers on IAAS machines
// using nerdctl/containerd, so that pebble sockets appear at the same paths
// that PebblePoller already polls.
package iaacontainerrunner

import (
    "context"
    "fmt"
    "os/exec"
    "path"

    "github.com/juju/errors"
    "github.com/juju/names/v6"
    "gopkg.in/tomb.v2"

    "github.com/juju/juju/core/logger"
)

// ResourcesAPI provides OCI image details for a named charm resource.
type ResourcesAPI interface {
    GetContainerImageInfo(ctx context.Context, unitTag names.UnitTag, resourceName string) (ImageDetails, error)
}

// ImageDetails holds the information needed to pull and run an OCI image.
type ImageDetails struct {
    RegistryPath string
    Username     string
    Password     string
}

// ContainerConfig describes a single workload container.
type ContainerConfig struct {
    Name         string // charm container name, e.g. "myapp"
    ResourceName string // charm resource that holds the OCI image ref
}

type worker struct {
    tomb       tomb.Tomb
    logger     logger.Logger
    unitTag    names.UnitTag
    containers []ContainerConfig
    resources  ResourcesAPI
}

// New returns a Worker that ensures each named container is running.
func New(unitTag names.UnitTag, containers []ContainerConfig, resources ResourcesAPI, log logger.Logger) (*worker, error) {
    w := &worker{
        logger:     log,
        unitTag:    unitTag,
        containers: containers,
        resources:  resources,
    }
    w.tomb.Go(w.loop)
    return w, nil
}

func (w *worker) Kill()            { w.tomb.Kill(nil) }
func (w *worker) Wait() error      { return w.tomb.Wait() }

func (w *worker) loop() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for _, c := range w.containers {
        if err := w.ensureRunning(ctx, c); err != nil {
            return errors.Annotatef(err, "starting container %q", c.Name)
        }
    }
    <-w.tomb.Dying()
    for _, c := range w.containers {
        w.stopContainer(ctx, c.Name)
    }
    return tomb.ErrDying
}

// containerID returns the nerdctl container name for a given unit+container.
func (w *worker) containerID(containerName string) string {
    return fmt.Sprintf("juju-%s-%s", w.unitTag.Id(), containerName)
}

// socketDir returns the host-side directory that is bind-mounted into the
// container as /charm/container. Pebble creates its socket here.
// Uses unit-scoped dataDir to avoid collisions between units on same machine.
func (w *worker) socketDir(containerName string) string {
    return path.Join(w.dataDir, "charm/containers", containerName)
}

func (w *worker) ensureRunning(ctx context.Context, c ContainerConfig) error {
    id := w.containerID(c.Name)

    // Idempotent: skip if already running.
    if isRunning(ctx, id) {
        w.logger.Debugf(ctx, "container %q already running", id)
        return nil
    }

    img, err := w.resources.GetContainerImageInfo(ctx, w.unitTag, c.ResourceName)
    if err != nil {
        return errors.Annotatef(err, "fetching image for resource %q", c.ResourceName)
    }

    dir := socketDir(c.Name)
    if err := exec.CommandContext(ctx, "mkdir", "-p", dir).Run(); err != nil {
        return errors.Annotatef(err, "creating socket directory %q", dir)
    }

    // Registry authentication — use --password-stdin to avoid credentials
    // being visible in `ps aux` output (OWASP credential exposure).
    if img.Username != "" {
        loginCmd := exec.CommandContext(ctx, "nerdctl", "login",
            "--username", img.Username, "--password-stdin",
            registryHost(img.RegistryPath))
        loginCmd.Stdin = strings.NewReader(img.Password)
        if out, err := loginCmd.CombinedOutput(); err != nil {
            return errors.Annotatef(err, "nerdctl login: %s", out)
        }
    }

    // nerdctl pull
    pullArgs := []string{"pull", img.RegistryPath}
    if out, err := exec.CommandContext(ctx, "nerdctl", pullArgs...).CombinedOutput(); err != nil {
        return errors.Annotatef(err, "nerdctl pull %q: %s", img.RegistryPath, out)
    }

    // nerdctl run — with pebble binary mount and entrypoint override.
    // Uses unit-scoped paths to avoid collisions between units on same machine.
    // Uses --restart=unless-stopped for reboot resilience.
    // Disables pebble HTTP (--http "") to avoid port conflicts.
    runArgs := []string{
        "run", "-d",
        "--name", id,
        "--network", "host",
        "--restart", "unless-stopped",
        "-v", fmt.Sprintf("%s/charm/bin/pebble:/charm/bin/pebble:ro", w.dataDir),
        "-v", fmt.Sprintf("%s:/charm/container", dir),
        "-e", fmt.Sprintf("JUJU_CONTAINER_NAME=%s", c.Name),
        "-e", "PEBBLE_SOCKET=/charm/container/pebble.socket",
        "-e", "PEBBLE=/charm/bin/pebble",
        "-e", "PEBBLE_COPY_ONCE=/var/lib/pebble/default",
        "--entrypoint", "/charm/bin/pebble",
        img.RegistryPath,
        "run", "--create-dirs", "--hold", "--http", "", "--verbose",
    }
    // TODO: Add storage mount flags via storageMountFlags()
    if out, err := exec.CommandContext(ctx, "nerdctl", runArgs...).CombinedOutput(); err != nil {
        return errors.Annotatef(err, "nerdctl run %q: %s", id, out)
    }

    w.logger.Infof(ctx, "started container %q for unit %q", id, w.unitTag)
    return nil
}

func (w *worker) stopContainer(ctx context.Context, containerName string) {
    id := w.containerID(containerName)
    // Graceful stop with 30s timeout to allow pebble to drain.
    _ = exec.CommandContext(ctx, "nerdctl", "stop", "--time", "30", id).Run()
    _ = exec.CommandContext(ctx, "nerdctl", "rm", id).Run()
    // Clean up socket directory.
    _ = os.RemoveAll(w.socketDir(containerName))
}

func isRunning(ctx context.Context, id string) bool {
    out, err := exec.CommandContext(ctx,
        "nerdctl", "inspect", "--format", "{{.State.Running}}", id,
    ).Output()
    return err == nil && string(out) == "true\n"
}
```

**F.2 — `manifold.go` skeleton**

```go
package iaacontainerrunner

import (
    "context"

    "github.com/juju/worker/v5"
    "github.com/juju/worker/v5/dependency"

    "github.com/juju/juju/agent"
    "github.com/juju/juju/api/base"
    "github.com/juju/juju/core/logger"
    "github.com/juju/juju/agent/engine"
)

// ManifoldConfig holds the manifold dependencies.
type ManifoldConfig struct {
    AgentName      string
    APICallerName  string
    ContainerNames []string
    Logger         logger.Logger
    NewWorker      func(Config) (worker.Worker, error)
}

// Manifold returns a dependency.Manifold that starts the worker only when
// ContainerNames is non-empty.
func Manifold(config ManifoldConfig) dependency.Manifold {
    if len(config.ContainerNames) == 0 {
        return dependency.Manifold{Start: func(_ context.Context, _ dependency.Getter) (worker.Worker, error) {
            return nil, dependency.ErrMissing
        }}
    }
    return engine.AgentAPIManifold(
        engine.AgentAPIManifoldConfig{
            AgentName:     config.AgentName,
            APICallerName: config.APICallerName,
        },
        config.start,
    )
}

func (config ManifoldConfig) start(_ context.Context, a agent.Agent, apiCaller base.APICaller) (worker.Worker, error) {
    agentCfg := a.CurrentConfig()
    unitTag, err := names.ParseUnitTag(agentCfg.Tag().String())
    if err != nil {
        return nil, err
    }

    resourcesClient := newResourcesClient(apiCaller) // thin wrapper around resources facade
    containers := make([]ContainerConfig, len(config.ContainerNames))
    // ContainerConfig.ResourceName will be looked up inside the worker from
    // charm metadata. For the initial skeleton it can be passed as the same
    // string (convention: resource named same as container).
    for i, name := range config.ContainerNames {
        containers[i] = ContainerConfig{Name: name, ResourceName: name}
    }
    return New(unitTag, containers, resourcesClient, config.Logger)
}
```

---

### G. `internal/worker/deployer/unit_manifolds.go` — wire the new manifold

```go
 import (
     // ... existing imports ...
+    "github.com/juju/juju/internal/worker/iaacontainerrunner"
 )

 const (
     // ... existing constants ...
+    iaaContainerRunnerName = "iaa-container-runner"
 )

 // In UnitManifolds():
 manifolds := dependency.Manifolds{
     // ... existing manifolds ...

+    iaaContainerRunnerName: ifNotMigrating(iaacontainerrunner.Manifold(iaacontainerrunner.ManifoldConfig{
+        AgentName:      agentName,
+        APICallerName:  apiCallerName,
+        ContainerNames: config.ContainerNames,
+        Logger:         config.LoggerContext.GetLogger("juju.worker.iaacontainerrunner"),
+    })),

     uniterName: ifNotMigrating(uniter.Manifold(uniter.ManifoldConfig{
         AgentName:             agentName,
         ModelType:             model.IAAS,
         APICallerName:         apiCallerName,
         S3CallerName:          s3CallerName,
         TraceName:             traceName,
         MachineLock:           config.MachineLock,
         Clock:                 config.Clock,
         LeadershipTrackerName: leadershipTrackerName,
         CharmDirName:          charmDirName,
         HookRetryStrategyName: hookRetryStrategyName,
         TranslateResolverErr:  uniter.TranslateFortressErrors,
         Logger:                config.LoggerContext.GetLogger("juju.worker.uniter"),
+        ContainerNames:        config.ContainerNames,
     })),
 }
```

The `iaacontainerrunner.Manifold()` returns a no-op manifold when `ContainerNames`
is empty, so pure IAAS charms (no `containers:` section) have zero overhead.

---

## Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| Pebble binary distribution | High | Package with juju snap/deb; validate during build |
| containerd not installed on VMs | Medium | `assumes: containerd` validation; clear error messages |
| OCI image credentials | Medium | Use `nerdctl login --password-stdin` (not CLI args) |
| Container networking complexity | Low | Use `--network host` initially; disable pebble HTTP to avoid port conflicts |
| Charm upgrade races | Medium | Worker holds off uniter start until containers ready |
| Storage mount failures | Medium | Fail unit clearly if storage not attached before container start |
| Block device filesystem creation | Low | Use ext4 by default; detect existing filesystem to avoid data loss |
| Shared storage coordination | Medium | Use nerdctl named volumes for same-machine sharing; rely on storage provisioner for cross-machine |
| Socket path collision (multi-unit) | High | Unit-scoped paths under `<dataDir>/charm/containers/` |
| Pebble binary version skew | Medium | SHA256 hash check on worker startup; restart containers if changed |
| Container logs invisible | Medium | Forward via `nerdctl logs --follow` to juju log pipeline |
| Machine reboot loses containers | Medium | `--restart=unless-stopped` + idempotent startup |

---

## Implementation Decisions

| Decision | Rationale |
|----------|-----------|
| Use nerdctl CLI (not containerd gRPC API) | Simpler implementation; optimize later if needed |
| Use `--network host` | Avoids container networking complexity; matches typical VM deployment patterns |
| Disable pebble HTTP (`--http ""`) | Avoids port conflicts when multiple containers run on same machine; unix socket is sufficient |
| Package pebble with juju agent | Most reliable distribution; no download at runtime |
| Advertise `containerd` unconditionally for IAAS | Controller can't probe machines at deploy time; runtime failure is clearer |
| Reuse existing storage provisioner | IAAS storage attachment already works; just map to container mounts |
| Block storage → filesystem mount | Create filesystem on block device, mount to host path, bind-mount into container; covers common use cases without raw device complexity |
| Shared storage via nerdctl named volumes | Named volumes persist across container restarts and can be shared between containers on same machine; cross-machine sharing uses underlying storage provisioner (NFS, etc.) |
| Unit-scoped paths (`<dataDir>/charm/...`) | Prevents socket/binary collisions when multiple units on same machine |
| `--restart=unless-stopped` | Containers survive machine reboots via containerd; worker handles idempotent startup |
| `nerdctl login --password-stdin` | Credentials not exposed in process list (OWASP compliance) |
| `Sidecar` stays false for IAAS | IAAS universal charms run hooks on host, not in container; existing IAAS uniter path is correct |
| Image mismatch detection on restart | Leverages existing engine bounce on upgrade-charm; no additional watcher needed |

---

## Corrections from Original Plan

| Original Plan | Corrected |
|---------------|-----------|
| "Remove FormatV2 deploy block on IAAS" | No such block exists — infrastructure must be *added*, not a gate removed |
| Missing pebble binary distribution | Added Section 1: distribute pebble binary with juju agent |
| Incomplete nerdctl command (missing entrypoint) | Added `--entrypoint /charm/bin/pebble` and pebble run arguments |
| Missing pebble binary mount | Added `-v <dataDir>/charm/bin/pebble:/charm/bin/pebble:ro` |
| `PEBBLE_COPY_ONCE=<pebble-dir>` | Fixed to concrete path `/var/lib/pebble/default` |
| Storage mounts not addressed | Added Section 5: Storage Mounts for Universal Charms |
| containerd detection unclear | Clarified: advertise unconditionally for IAAS; fail at runtime |
| Missing health monitoring | Added: worker watches container state, restarts crashed containers |
| Global socket paths `/charm/containers/` | Fixed: unit-scoped under `<dataDir>/charm/containers/` to avoid multi-unit collisions |
| Credentials on CLI (`--username`/`--password`) | Fixed: use `nerdctl login --password-stdin` to avoid exposure in `ps aux` |
| Missing `--restart` policy | Added: `--restart=unless-stopped` for machine reboot resilience |
| Pebble HTTP port conflict (`:38813`) | Fixed: disabled pebble HTTP (`--http ""`) — unix socket only |
| No charm upgrade / image update mechanism | Added Section 8: image mismatch detection on worker restart |
| No pebble binary upgrade strategy | Added Section 9: SHA256 hash comparison on worker startup |
| Container logs invisible | Added Section 10: `nerdctl logs --follow` forwarding to juju log pipeline |
| Sidecar flag interaction unclear | Clarified: Sidecar stays `false` for IAAS universal charms |

---

## Verification Steps

1. **Unit tests:** `go test -race ./internal/worker/iaacontainerrunner/...`
2. **Deploy FormatV2 charm on LXD:**
   ```bash
   juju bootstrap lxd test
   juju deploy cockroachdb  # FormatV2 charm with containers
   ```
3. **Verify pebble-ready hook fires:** Check `juju debug-log` for `cockroachdb-pebble-ready`
4. **Verify pebble commands work:**
   ```bash
   juju exec --unit cockroachdb/0 -- pebble services
   ```
5. **Verify filesystem storage mounts:**
   ```bash
   juju deploy myapp --storage data=1G
   juju exec --unit myapp/0 -- ls -la /data  # check mount is accessible
   ```
6. **Verify block storage (treated as filesystem):**
   ```bash
   juju deploy myapp --storage blockdata=1G,block
   juju exec --unit myapp/0 -- df -h /blockdata  # should show mounted filesystem
   ```
7. **Verify shared storage between units:**
   ```bash
   juju deploy myapp -n 2 --storage shared=1G
   juju exec --unit myapp/0 -- touch /shared/testfile
   juju exec --unit myapp/1 -- ls /shared/testfile  # should exist
   ```
8. **Run static analysis:** `make pre-check`
9. **Integration tests:** `cd tests && ./main.sh universal_charms`

---

## Future Enhancements

- **containerd gRPC API:** Replace nerdctl CLI with direct API for better performance and error handling
- **Container resource limits:** Map charm metadata resource constraints to nerdctl `--memory`, `--cpus` flags
- **Raw block device passthrough:** Direct `/dev/` device access for containers (not commonly needed)
- **Cross-machine shared storage improvements:** Better coordination for distributed storage backends
- **Container networking modes:** Support bridge networking or CNI for multi-container service discovery
