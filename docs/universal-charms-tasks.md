# Universal Charms — Design Review & Task List

## Implementation Delta (May 2026)

### Verified

- Hook pipeline wiring is present and tested:
   - `pebble-ready` via `PebblePoller`
   - `pebble-check-failed` and `pebble-check-recovered` via `PebbleNoticer`
   - workload event to hook mapping in the uniter workload resolver
- IAAS runner wiring is in place in unit manifolds with uniter dependency on
   `iaas-container-runner` when `ContainerNames` is non-empty.
- Runtime provisioning moved to cloud-init defaults for this POC (`containerd`,
   `nerdctl`, `pebble`) and validated by cloudconfig unit tests.

### Fixed During Debugging

- Pebble startup crash from incorrect `PEBBLE` env was fixed by using a
   directory path in the container (`PEBBLE=/charm/container`).
- Cloud-init shell compatibility issue (`set -o pipefail`) was fixed.
- Runtime provisioning script was made bounded and best-effort so download
   failures do not indefinitely block machine agent installation.

### Current Blocker

- Some integration runs were interrupted while unit agent remained
   `installing agent`. In this state, charm hooks do not run yet, so
   `pebble services` returning `Plan has no services` is expected.

### Immediate Next Tasks

1. Re-run `tests/suites/universal_charms` without early interruption and wait
    until unit agent reaches `idle`.
2. If still stuck in `installing agent`, capture machine cloud-init logs and
    unit-agent logs from introspection artifacts for root cause.
3. Once idle, validate event transitions explicitly:
    - absent `/trigger` => `pebble_check_failed`
    - create `/trigger` => `pebble_check_recovered`
    - remove `/trigger` => `pebble_check_failed` again

## Design Review: Identified Gaps

### Gap 1: `Sidecar` Flag Interaction

The uniter has a `Sidecar bool` field (used by CAAS containeragent) that affects resolver
behaviour. The plan does not mention how this field interacts with universal charms. For
IAAS units running FormatV2 charms, `Sidecar` should remain `false` — the existing IAAS
uniter path is correct — but this needs explicit verification. The plan should state:
"Sidecar remains false for IAAS universal charms; it is only true when running inside a
K8s sidecar container."

### Gap 2: OCI Image Resource Resolution

The plan's `ResourcesAPI.GetContainerImageInfo()` is described as a new interface, but
the existing `domain/containerimageresourcestore` package already has
`GetContainerImageMetadata` which returns registry path, username, and password. The
plan should leverage this existing domain service rather than creating a parallel path.
The deployer facade or a new dedicated facade needs to expose this to the machine agent.

### Gap 3: Charm Upgrade / OCI Image Update Path

The plan mentions "On charm upgrade (new OCI image resource): stop, pull new image,
restart" but doesn't specify:
- How the `iaacontainerrunner` gets notified of resource updates
- Whether this is triggered by a watcher on the resource revision or by the uniter
  signalling the worker
- What happens to in-flight pebble operations during container restart

### Gap 4: Container Networking Beyond `--network host`

`--network host` means all containers share the machine's network namespace. If multiple
units of the same application (or different applications) run containers that bind the
same port (e.g., pebble HTTP `:38813`), they will conflict. The plan should either:
- Use unique pebble HTTP ports per container (e.g., assign dynamically)
- Use `--network none` with only the unix socket for communication (pebble HTTP is
  optional and only needed for health checks)

### Gap 5: Security — Credential Handling for Private Registries

The plan passes `--username` and `--password` directly on the nerdctl command line.
These will be visible in `ps aux` output. The correct approach is:
- Use `nerdctl login` with `--password-stdin`
- Or write a temporary docker config.json and pass `--config`

### Gap 6: No Graceful Container Shutdown on Unit Removal

The plan's `stopContainer` calls `nerdctl stop` then `nerdctl rm`. It should:
- Respect a configurable stop timeout (pebble needs time to drain)
- Clean up socket directories
- Remove pulled images (optional, for disk hygiene)

### Gap 7: Multiple Units on Same Machine — Socket Path Collision

The socket path `/charm/containers/<name>/pebble.socket` is not unit-scoped. If two
units of the same application land on the same machine, they'd share the same socket
path. The path should be `/var/lib/juju/agents/unit-<name>/charm/containers/<containerName>/pebble.socket`
or similar. This requires checking what `PebblePoller` uses — it uses
`path.Join("/charm/containers", containerName, "pebble.socket")` which on CAAS is rooted
in the unit's filesystem. On IAAS, `DataDir` for the unit agent is
`/var/lib/juju/agents/unit-<app>-<num>`, so the socket path needs to be relative to the
unit's data directory.

### Gap 8: Pebble Binary Version Management

If pebble is packaged with the juju snap, upgrades to juju update pebble. Running
containers still have the old pebble binary bind-mounted. The plan needs a strategy for
pebble binary upgrades — likely requiring container restart on juju agent upgrade.

### Gap 9: `assumes: containerd` Feature Registration on IAAS

The plan says to advertise `containerd` feature when snapd is available, but the
`GetSupportedFeatures` call happens at deploy time in the controller — not on the
target machine. The controller doesn't know if snapd is available on the machines.
Options:
- Always advertise `containerd` for IAAS models (and fail at runtime if not available)
- Have the machine agent report supported features back to the controller

### Gap 10: Log Integration

Container stdout/stderr goes to nerdctl/containerd logs. The plan mentions "future
enhancement" for log integration but charm operators expect `juju debug-log` to show
workload messages. At minimum, the `iaacontainerrunner` should tail container logs and
forward to the Juju logging pipeline.

### Gap 11: Backwards Compatibility — VM Charms Still Work

The user requirement states: "Juju should still support running regular VM charms on VMs
and k8s charms on k8s." The plan correctly notes that when `ContainerNames` is empty, no
new behaviour is triggered. However, the plan needs explicit verification that:
- FormatV1 charms on IAAS continue to work (no `ContainerNames`, no container runner)
- FormatV2 charms on CAAS continue to work (existing K8s provider path unchanged)
- The new `GetContainerNames` callback in `ContextConfig` returns nil for FormatV1 charms

### Gap 12: Hook Execution Context — `JUJU_DISPATCH_PATH` Differences

On CAAS, hooks run inside the charm container via pebble exec. On IAAS universal charms,
hooks run on the host (inside the unit agent). The charm can interact with the workload
via pebble client connecting to the socket. This is the correct model but should be
explicitly documented: the dispatch path for hooks is identical to regular IAAS charms;
pebble connectivity to workload containers is via the socket path.

### Gap 13: init System Integration

On CAAS, pod restart guarantees container restart. On IAAS, if the machine reboots,
nerdctl containers won't auto-restart unless configured with `--restart=always` or a
systemd unit. The plan should add `--restart=unless-stopped` to the nerdctl run command.

---

## Task List

### Task 0: Pebble Binary Distribution

**Description:** Package the pebble binary with the juju snap so it's available at
`/snap/juju/current/bin/pebble` (and symlinked to `/usr/lib/juju/bin/pebble` or similar).
At unit deploy time, copy or symlink to the unit's charm directory.

**Files to modify:**
- `snap/snapcraft.yaml` — add pebble part that builds from canonical/pebble
- Makefile/build scripts — download or build pebble binary

**Agent prompt:**
```
Add pebble binary to the juju snap build. In snap/snapcraft.yaml, add a new part
called "pebble" that clones https://github.com/canonical/pebble and builds it with
`go build -o $SNAPCRAFT_PART_INSTALL/bin/pebble ./cmd/pebble`. Place the binary at
bin/pebble in the snap. Ensure the part uses the same Go version as the juju part.
Reference the existing parts structure in the file for the correct pattern.
```

---

### Task 1: Add `UnitContainerNames` to the Deployer Facade (Server Side)

**Description:** Add a new RPC method `UnitContainerNames` to the server-side deployer
facade that returns the container names from a unit's charm metadata. This requires:
1. Adding a method to the `ApplicationService` interface to get charm metadata containers
   for a given unit.
2. Implementing `UnitContainerNames` on `DeployerAPI`.
3. Bumping the facade version.

**Files to modify:**
- `apiserver/facades/agent/deployer/deployer.go` — add `UnitContainerNames` method
- `apiserver/facades/agent/deployer/deployer.go` — extend `ApplicationService` interface
  with `GetUnitContainerNames(context.Context, coreunit.Name) ([]string, error)`
- `domain/application/service/application.go` — implement `GetUnitContainerNames`
- `domain/application/state/application.go` — state layer query
- `apiserver/facades/agent/deployer/register.go` — bump version if needed

**Agent prompt:**
```
Add a new RPC method `UnitContainerNames` to the server-side deployer API facade in
apiserver/facades/agent/deployer/deployer.go.

1. Add to the ApplicationService interface:
   GetUnitContainerNames(ctx context.Context, unitName coreunit.Name) ([]string, error)

2. Implement on DeployerAPI:
   func (d *DeployerAPI) UnitContainerNames(ctx context.Context, args params.Entities) (params.StringsResults, error)
   - Use params.Entities (batch call pattern matching the existing Life/Remove methods)
   - For each entity, parse the unit tag, authorize via d.getAuth, call
     d.applicationService.GetUnitContainerNames, return sorted container names
   - Return empty slice (not error) for charms with no containers

3. Implement GetUnitContainerNames in the application domain service layer:
   - Look up the unit's application
   - Get the charm metadata for that application  
   - Extract and return sorted keys from meta.Containers map
   - Return nil/empty if no containers section

4. Write unit tests following the existing test patterns in deployer_test.go.
   Use tc.Assert/tc.Check as per AGENTS.md conventions.

Read the existing DeployerAPI methods (Life, Remove, WatchUnits) for the correct
authorization and error handling patterns. Use params.StringsResults (plural) for
the batch response.
```

---

### Task 2: Add `UnitContainerNames` to the Deployer Facade (Client Side)

**Description:** Add client-side method to `api/agent/deployer` package to call the
new `UnitContainerNames` RPC.

**Files to modify:**
- `api/agent/deployer/deployer.go` — add `UnitContainerNames` method to `Client`
- `api/agent/deployer/unit.go` — add `ContainerNames` method to `Unit`

**Agent prompt:**
```
Add client-side support for the UnitContainerNames RPC call to api/agent/deployer/.

1. In api/agent/deployer/unit.go, add:
   func (u *Unit) ContainerNames(ctx context.Context) ([]string, error)
   - Call u.client.facade.FacadeCall(ctx, "UnitContainerNames", args, &result)
   - Use params.Entities{Entities: []params.Entity{{Tag: u.tag.String()}}} as args
   - Use params.StringsResults as result type, extract first result
   - Return nil (not error) if result is empty

2. In api/agent/deployer/deployer.go, add a convenience method:
   func (c *Client) UnitContainerNames(ctx context.Context, tag names.UnitTag) ([]string, error)
   - Create the unit via c.Unit()
   - Call unit.ContainerNames()

Follow the existing patterns in unit.go (SetPassword, SetStatus, Remove) for
FacadeCall conventions. Write tests in deployer_test.go following the existing
test patterns.
```

---

### Task 3: Thread `ContainerNames` Through the Deployer to the Uniter

**Description:** When the machine agent deploys a unit, look up the unit's container
names via the deployer facade and pass them through the nested context → unit agent →
unit manifolds → uniter.

**Files to modify:**
- `internal/worker/deployer/nested.go` — add `GetContainerNames` to `ContextConfig`,
  call it in `newUnitAgent()`
- `internal/worker/deployer/unit_agent.go` — add `ContainerNames []string` to
  `UnitAgentConfig` and thread to `UnitManifoldsConfig`
- `internal/worker/deployer/unit_manifolds.go` — add `ContainerNames []string` to
  `UnitManifoldsConfig`, pass to uniter `ManifoldConfig`
- `internal/worker/deployer/manifold.go` — implement `GetContainerNames` callback
  using the deployer facade client

**Agent prompt:**
```
Thread container names from the deployer facade through to the uniter manifold for
IAAS units. This enables PebblePoller to start for FormatV2 charms on VMs.

1. In internal/worker/deployer/nested.go:
   - Add to ContextConfig:
     GetContainerNames func(ctx context.Context, unitTag names.UnitTag) ([]string, error)
   - Do NOT add it to Validate() — it's optional (nil means no containers)
   - In newUnitAgent(), after creating unitConfig, if c.config.GetContainerNames != nil:
     call it with the unit tag, set unitConfig.ContainerNames = result
   - Store the ContextConfig on nestedContext so newUnitAgent can access it

2. In internal/worker/deployer/unit_agent.go:
   - Add ContainerNames []string to UnitAgentConfig
   - In the start() method where UnitManifoldsConfig is constructed, set
     ContainerNames: a.containerNames (store from config in NewUnitAgent)

3. In internal/worker/deployer/unit_manifolds.go:
   - Add ContainerNames []string to UnitManifoldsConfig (after Clock field)
   - In the uniter manifold config (uniterName), add:
     ContainerNames: config.ContainerNames,

4. In internal/worker/deployer/manifold.go:
   - In newWorker(), when building contextConfig, add:
     GetContainerNames: func(ctx context.Context, unitTag names.UnitTag) ([]string, error) {
         return deployerFacade.UnitContainerNames(ctx, unitTag)
     },

The existing CAAS path is unchanged — it sets ContainerNames directly in
cmd/containeragent/unit/manifolds.go from the JUJU_CONTAINER_NAMES env var.
For IAAS, the deployer now fetches container names from the controller.

When ContainerNames is empty (FormatV1 charms), the uniter behaves exactly as before.

Write tests: test that newUnitAgent() calls GetContainerNames and passes the result.
Test that when GetContainerNames returns nil, ContainerNames is empty.
```

---

### Task 4: Create `internal/worker/iaascontainerrunner` Package

**Description:** New worker package that manages OCI workload containers on IAAS VMs
using nerdctl. This is the core new component.

**Files to create:**
- `internal/worker/iaascontainerrunner/worker.go`
- `internal/worker/iaascontainerrunner/manifold.go`
- `internal/worker/iaascontainerrunner/doc.go`
- `internal/worker/iaascontainerrunner/worker_test.go`

**Agent prompt:**
```
Create a new worker package at internal/worker/iaascontainerrunner/ that runs OCI
workload containers on IAAS machines using nerdctl/containerd.

IMPORTANT DESIGN DECISIONS:
- Use "iaascontainerrunner" (not "iaacontainerrunner") for correct spelling
- Socket path must be unit-scoped: use dataDir from agent config, not global /charm/
- Use --restart=unless-stopped for container resilience across reboots
- Use nerdctl login --password-stdin for registry auth (not CLI args)
- Use unique pebble HTTP port per container or disable HTTP (--http "" flag)

Package structure:

1. doc.go — Package documentation per AGENTS.doc-dot-go-rules.md

2. worker.go:
   - Define CommandRunner interface for testability:
     type CommandRunner interface {
         Run(ctx context.Context, name string, args ...string) ([]byte, error)
         RunStdin(ctx context.Context, stdin string, name string, args ...string) ([]byte, error)
     }
   - Define ResourcesAPI interface (to get OCI image info):
     type ResourcesAPI interface {
         GetContainerImageInfo(ctx context.Context, resourceName string) (ImageDetails, error)
     }
     type ImageDetails struct {
         RegistryPath string
         Username     string
         Password     string
     }
   - Define Config struct:
     type Config struct {
         Logger         logger.Logger
         UnitTag        names.UnitTag
         DataDir        string
         ContainerNames []string
         CharmMeta      map[string]ContainerMeta // container name -> meta
         Resources      ResourcesAPI
         CommandRunner  CommandRunner
     }
     type ContainerMeta struct {
         ResourceName string
         Mounts       []Mount
     }
     type Mount struct {
         StorageName string
         Location    string
     }
   - Worker struct with tomb.Tomb
   - loop() method:
     a. Call ensureNerdctl() — check nerdctl available, if not run snap install
     b. For each container, call ensureRunning()
     c. Block on tomb.Dying(), then stop all containers
   - ensureRunning():
     a. Check if container already running via nerdctl inspect
     b. Get image details from ResourcesAPI
     c. Login to registry if credentials present (via nerdctl login --password-stdin)
     d. Pull image
     e. Run container with:
        --name juju-<unitID>-<containerName>
        --network host
        --restart unless-stopped
        -v <dataDir>/charm/bin/pebble:/charm/bin/pebble:ro
        -v <dataDir>/charm/containers/<name>:/charm/container
        -e JUJU_CONTAINER_NAME=<name>
        -e PEBBLE_SOCKET=/charm/container/pebble.socket
        -e PEBBLE=/charm/bin/pebble
        -e PEBBLE_COPY_ONCE=/var/lib/pebble/default
        --entrypoint /charm/bin/pebble
        <image>
        run --create-dirs --hold --http "" --verbose
   - stopContainer(): nerdctl stop --time 30 <id> && nerdctl rm <id>
   - isRunning(): nerdctl inspect --format {{.State.Running}} <id>

3. manifold.go:
   - ManifoldConfig struct with AgentName, APICallerName, ContainerNames, Logger
   - Manifold() function: if len(ContainerNames) == 0, return no-op manifold
     (return dependency.ErrMissing from Start)
   - Otherwise use engine.AgentAPIManifold pattern
   - start() function: get agent config for DataDir and UnitTag, create
     ResourcesAPI client from API caller, get charm metadata for container->resource
     mapping, construct Config and call New()

4. worker_test.go:
   - Mock CommandRunner interface
   - Test ensureRunning calls nerdctl with correct args
   - Test isRunning detection
   - Test that credentials are passed via stdin, not CLI
   - Test graceful shutdown calls stop on all containers
   - Use tc.Assert/tc.Check per AGENTS.md conventions

Follow the patterns in existing workers like internal/worker/uniter/ for tomb usage
and internal/worker/trace/ for manifold structure. The CommandRunner abstraction allows
tests to verify nerdctl commands without actually executing them.
```

---

### Task 5: Wire `iaascontainerrunner` into Unit Manifolds

**Description:** Add the new container runner manifold to the IAAS unit manifolds,
gated on non-empty `ContainerNames`. The uniter should depend on the container runner
so containers are started before PebblePoller begins polling.

**Files to modify:**
- `internal/worker/deployer/unit_manifolds.go` — add the new manifold, add dependency
  from uniter to container runner

**Agent prompt:**
```
Wire the iaascontainerrunner worker into the IAAS unit dependency engine in
internal/worker/deployer/unit_manifolds.go.

1. Add import:
   "github.com/juju/juju/internal/worker/iaascontainerrunner"

2. Add constant:
   iaasContainerRunnerName = "iaas-container-runner"

3. In UnitManifolds(), add the manifold BEFORE uniterName:
   iaasContainerRunnerName: ifNotMigrating(iaascontainerrunner.Manifold(
       iaascontainerrunner.ManifoldConfig{
           AgentName:      agentName,
           APICallerName:  apiCallerName,
           ContainerNames: config.ContainerNames,
           Logger:         config.LoggerContext.GetLogger("juju.worker.iaascontainerrunner"),
       },
   )),

4. Add iaasContainerRunnerName to the uniter manifold's Inputs list so that the
   uniter waits for containers to be running before starting PebblePoller. This
   requires changing the uniter.ManifoldConfig to accept an optional dependency.

   ALTERNATIVE (simpler): Since iaascontainerrunner.Manifold returns ErrMissing when
   ContainerNames is empty, and the uniter's Start function currently doesn't depend
   on it — add a new optional input to uniter.ManifoldConfig:
   ContainerRunnerName string  // optional, only for IAAS universal charms
   
   In uniter/manifold.go Manifold(), conditionally add to Inputs if non-empty, and
   in Start(), if set, call getter.Get() to ensure it's running (just validate
   existence, don't extract a value).

The key constraint: when ContainerNames is empty (standard IAAS charm), NOTHING
changes — no new manifold is started, no new dependency is added to the uniter.
Run gci on all modified .go files after editing.
```

---

### Task 6: Register `containerd` Feature for `assumes` Validation

**Description:** Add a `containerd` feature to the assumes feature set so charm
authors can declare `assumes: [containerd]`. For IAAS models, advertise this feature
unconditionally (runtime failures will surface clearly in `juju status`).

**Files to modify:**
- `core/assumes/features.go` — add `containerd` to `UserFriendlyFeatureDescriptions`
  and `featureMissingErrs`, add `ContainerdFeature()` constructor
- `domain/application/service/` (or wherever `GetSupportedFeatures` is implemented) —
  include containerd feature for IAAS models

**Agent prompt:**
```
Add containerd feature support to the assumes framework so charms can declare
`assumes: [containerd]`.

1. In core/assumes/features.go:
   - Add to UserFriendlyFeatureDescriptions map:
     "containerd": "the containerd runtime for running OCI workload containers"
   - Add to featureMissingErrs map:
     "containerd": "charm requires containerd runtime but model does not support it"
   - Add constructor function:
     func ContainerdFeature() Feature {
         return Feature{
             Name:        "containerd",
             Description: UserFriendlyFeatureDescriptions["containerd"],
         }
     }

2. Find where GetSupportedFeatures is implemented for IAAS models (search for
   GetSupportedFeatures in domain/application/service/). Add assumes.ContainerdFeature()
   to the feature set returned for IAAS models. This makes `assumes: [containerd]`
   pass validation at deploy time for any IAAS model.

   For CAAS models, do NOT add this feature — K8s charms use the K8s runtime directly.

3. Write a test that verifies ContainerdFeature is included in the supported features
   for IAAS models.

Run gci on modified files. Follow existing patterns in features.go for the constructor.
```

---

### Task 7: Pebble Binary Setup During Unit Deployment

**Description:** When deploying a unit that has containers (FormatV2 on IAAS), ensure
the pebble binary is available at the expected path within the unit's data directory.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — add `ensurePebbleBinary()` that
  copies/symlinks pebble from the juju snap location to the unit's charm/bin/ dir

**Agent prompt:**
```
Add pebble binary setup to the iaascontainerrunner worker. Before starting containers,
the worker must ensure pebble is available at <dataDir>/charm/bin/pebble.

In internal/worker/iaascontainerrunner/worker.go, add an ensurePebbleBinary() method:

1. Define the source path: /snap/juju/current/bin/pebble (when running from snap)
   or look up via the agent tools directory
2. Create <dataDir>/charm/bin/ directory (mkdir -p)
3. Copy the pebble binary to <dataDir>/charm/bin/pebble
4. chmod 0755 the binary
5. If source not found, return a clear error that surfaces in juju status

Call this at the start of loop() before ensureNerdctl().

Use the CommandRunner interface for filesystem operations to keep it testable,
or use os package directly (simpler, test via integration tests).

Also ensure <dataDir>/charm/containers/<name>/ directories exist for each container
before starting them — this is where the pebble socket will appear.
```

---

### Task 8: Container Status Reporting

**Description:** Surface per-container status in `juju status` for IAAS universal
units. The `iaascontainerrunner` worker should report container state via the existing
unit status mechanisms.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — periodic container health check
- Status reporting via the uniter's existing status API

**Agent prompt:**
```
Add container health monitoring and status reporting to the iaascontainerrunner worker.

1. In the worker's loop(), after starting all containers, enter a monitoring loop:
   - Every 30 seconds, check each container's state via nerdctl inspect
   - If a container has crashed (State.Running == false, State.ExitCode != 0):
     a. Log the crash with exit code
     b. Attempt restart (nerdctl start <id>)
     c. If restart fails 3 times consecutively, report blocked status
   - Report container states to a status callback

2. Add a StatusReporter interface to Config:
   type StatusReporter interface {
       ReportContainerStatus(ctx context.Context, containerName string, status ContainerStatus) error
   }
   type ContainerStatus struct {
       State   string // "running", "waiting", "crashed", "unknown"
       Message string
   }

3. The manifold should wire this to the uniter's existing workload status mechanism
   or directly to the agent status API. For the initial implementation, logging the
   status is sufficient — full juju status integration can follow.

4. Use tomb.Dying() channel in select with the health check timer to ensure clean
   shutdown.

Write tests that verify:
- Crashed container triggers restart attempt
- Status is reported correctly
- Worker exits cleanly on Kill()
```

---

### Task 9: Container Restart Policy and Machine Reboot Handling

**Description:** Ensure containers survive machine reboots and that the worker handles
the case where containers are already running on startup (idempotent).

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — use `--restart=unless-stopped`,
  handle pre-existing containers on startup

**Agent prompt:**
```
Ensure the iaascontainerrunner handles machine reboots and worker restarts correctly.

1. Add --restart=unless-stopped to the nerdctl run command. This ensures containers
   auto-restart via containerd after a machine reboot, even before the juju agent
   starts.

2. In ensureRunning(), the existing isRunning() check already handles idempotency.
   Add handling for containers that exist but are stopped:
   - nerdctl inspect succeeds but State.Running is false
   - In this case, call nerdctl start <id> instead of nerdctl run

3. Handle the case where a container exists with old configuration (e.g., different
   image after charm upgrade):
   - Compare the running container's image with the expected image
   - If different: stop, rm, and re-run with new image
   - Use nerdctl inspect --format {{.Image}} to get current image

4. Add a method isCreated() that checks if a container exists (running or stopped):
   - nerdctl inspect <id> succeeds (regardless of state)

Write tests verifying:
- Already-running container is not restarted
- Stopped container is started (not re-created)
- Container with wrong image is replaced
```

---

### Task 10: Storage Mount Support

**Description:** Map charm metadata storage mounts to nerdctl `-v` flags. The existing
IAAS storage provisioner handles volume creation and attachment; this task maps those
host paths into the container.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — add storage mount resolution
- `internal/worker/iaascontainerrunner/manifold.go` — wire storage API dependency

**Agent prompt:**
```
Add storage mount support to iaascontainerrunner so that containers can access
Juju-provisioned storage.

1. Add a StorageResolver interface to Config:
   type StorageResolver interface {
       GetStorageMountPath(ctx context.Context, storageName string) (string, error)
   }

2. In ContainerMeta, the Mounts field already has StorageName and Location.
   In ensureRunning(), for each mount in the container's metadata:
   - Call StorageResolver.GetStorageMountPath(ctx, mount.StorageName) to get host path
   - Add -v <hostPath>:<mount.Location> to the nerdctl run args
   - If storage is not yet attached, return a retryable error (the worker will bounce)

3. In the manifold start() function, create a StorageResolver implementation:
   - Use the existing storage attachment API to resolve storage name → host path
   - The unit agent's data directory contains storage attachment info

4. For the initial implementation, only support filesystem storage (type: filesystem).
   Block storage and shared storage are deferred to a follow-up task.

5. If a mount's storage is not found, log a warning and skip the mount (don't fail
   the container start) — the storage may not be attached yet and a storage-attached
   hook will fire later.

Write tests that verify mount flags are correctly generated from ContainerMeta.
```

---

### Task 11: Integration Testing

**Description:** Create integration tests that deploy a FormatV2 charm on an IAAS model
and verify pebble-ready hook fires, workload status is reported, and pebble commands
work.

**Files to create:**
- `tests/suites/universal_charms/task.sh`
- Test charm with containers section (or use existing test charm)

**Agent prompt:**
```
Create integration tests for universal charms in tests/suites/universal_charms/.

1. Create tests/suites/universal_charms/task.sh following the pattern from other
   test suites (e.g., tests/suites/deploy/task.sh).

2. Test cases:
   a. test_deploy_formatv2_on_iaas:
      - Bootstrap an LXD model
      - Deploy a charm that has a containers: section (FormatV2)
      - Wait for active/idle status
      - Verify pebble-ready hook was fired (check debug-log)
      - Verify pebble services command works via juju exec
   
   b. test_formatv1_still_works:
      - Deploy a traditional FormatV1 charm on the same IAAS model
      - Verify it reaches active/idle (regression test)
   
   c. test_formatv2_on_k8s_unchanged:
      - Deploy the same FormatV2 charm on a K8s model
      - Verify it works as before (regression test)

3. You may need a minimal test charm in testcharms/ that has:
   metadata.yaml with containers: section
   A simple pebble layer that runs a basic service
   hooks/install and hooks/<container>-pebble-ready

Use the test framework from tests/includes/ (juju.sh, wait-for.sh, check.sh).
Follow the conventions in existing test suites for structure.
```

---

### Task 12: Documentation Updates

**Description:** Update user-facing documentation explaining that FormatV2 (K8s) charms
can now be deployed on IAAS models, and what the requirements are.

**Files to modify:**
- `docs/` — new user-facing doc on universal charms
- Update any deploy documentation that states K8s charms are K8s-only

**Agent prompt:**
```
Create user documentation for the universal charms feature.

1. Create docs/user/universal-charms.md explaining:
   - What universal charms are (FormatV2 charms that work on both K8s and IAAS)
   - Prerequisites: nerdctl/containerd on target machine (auto-installed via snap)
   - How to deploy: `juju deploy <formatv2-charm>` on an IAAS model (same command)
   - How storage mounts work on IAAS vs K8s
   - The `assumes: [containerd]` directive for charm authors who want to require this
   - Known limitations (initial release)
   - Troubleshooting: checking container status, logs, nerdctl commands

2. Follow the documentation rules in AGENTS.documentation.rules.md.

3. Keep it concise and practical with examples.
```

---

### Task 13: OCI Image Update on Charm Refresh

**Description:** When a charm is refreshed (`juju refresh`) with a new OCI image resource,
the `iaascontainerrunner` must detect the change and replace the running container with
one using the new image. This addresses Gap 3.

The uniter already processes resource changes during upgrade-charm. The cleanest approach
is to have the uniter bounce the `iaascontainerrunner` worker (by restarting the unit
engine or signalling the worker) when an OCI resource revision changes. On restart, the
worker compares the running container's image with the expected image (Task 9 already
handles image mismatch detection) and replaces it.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — accept a channel or watcher for
  resource change notifications; on notification, re-resolve image and replace container
- `internal/worker/iaascontainerrunner/manifold.go` — wire a resource revision watcher
  or accept a "bounce" signal from the uniter
- `internal/worker/uniter/uniter.go` — after processing upgrade-charm with new resources,
  signal the container runner to refresh

**Design options (pick one):**

A. **Worker restart on upgrade-charm** (simplest): The uniter already causes the unit
   engine to bounce during upgrade-charm. If the container runner is a sibling manifold,
   it restarts too. On restart, the image mismatch logic from Task 9 handles replacement.
   This works if the engine bounces — verify this is the case.

B. **Resource revision watcher**: The container runner watches resource revisions via
   an API watcher. On change, it re-fetches image info and replaces the container.
   More responsive but adds complexity.

**Agent prompt:**
```
Implement OCI image update handling for charm refresh in the iaascontainerrunner.

Preferred approach: rely on the unit dependency engine restart that occurs during
upgrade-charm. Verify this by reading the uniter's upgrade-charm handling code.

1. Verify: When `juju refresh` is called, does the unit dependency engine restart?
   Check internal/worker/uniter/ upgrade handling. If the engine restarts, then the
   iaascontainerrunner worker also restarts, and the image mismatch logic from Task 9
   (compare running image vs expected image, replace if different) handles it.

2. If the engine does NOT restart on upgrade-charm:
   - Add a ResourceWatcher to the iaascontainerrunner Config:
     type ResourceWatcher interface {
         Changes() <-chan struct{}
     }
   - In the worker's monitoring loop, select on the watcher channel
   - On change: for each container, re-fetch image info from ResourcesAPI
   - If image differs from running container, call replaceContainer():
     a. nerdctl stop --time 30 <id>
     b. nerdctl rm <id>
     c. Pull new image
     d. nerdctl run with new image (same config as initial start)
   - After replacement, wait for pebble socket to reappear (PebblePoller handles this)

3. Handle in-flight pebble operations: The pebble socket disappears during container
   replacement. PebblePoller will detect the socket is gone and fire
   <container>-pebble-custom-notice or stop responding until the new container creates
   the socket. The charm's pebble-ready hook fires again on the new container. This is
   the correct behaviour — document it in a code comment.

4. Write tests:
   - Test that worker detects image mismatch on restart and replaces container
   - Test that resource watcher change triggers re-evaluation (if option B)
   - Test graceful replacement: stop with timeout before rm

Use tc.Assert/tc.Check per AGENTS.md. Follow existing watcher patterns in the codebase.
```

---

### Task 14: Pebble Binary Upgrade on Agent Upgrade

**Description:** When the juju agent is upgraded, the pebble binary in the snap changes.
Running containers have the old binary bind-mounted. The `iaascontainerrunner` must
detect the mismatch and restart containers with the new pebble. This addresses Gap 8.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — add pebble binary hash comparison
  logic in the startup/monitoring path

**Design:**

On worker startup (which happens after an agent upgrade because the engine restarts):
1. Compute SHA256 of the source pebble binary (`/snap/juju/current/bin/pebble`)
2. Compute SHA256 of the deployed pebble binary (`<dataDir>/charm/bin/pebble`)
3. If different: copy new binary, then restart all containers (stop + start)
4. The `--restart=unless-stopped` policy means containerd may have restarted containers
   with the old binary after reboot; this check ensures they get the new one.

**Agent prompt:**
```
Add pebble binary upgrade detection to the iaascontainerrunner worker. This ensures
containers use the latest pebble binary after a juju agent upgrade.

In internal/worker/iaascontainerrunner/worker.go:

1. Add a method ensurePebbleCurrent() that runs AFTER ensurePebbleBinary() (Task 7):
   - Compute SHA256 hash of source binary (e.g., /snap/juju/current/bin/pebble)
   - Compute SHA256 hash of deployed binary (<dataDir>/charm/bin/pebble)
   - If hashes match, return nil (no action needed)
   - If different:
     a. Copy new binary to <dataDir>/charm/bin/pebble (overwrite)
     b. Set w.pebbleUpgraded = true (flag for container restart)

2. In ensureRunning(), after the isRunning() check returns true:
   - If w.pebbleUpgraded is true, force container replacement:
     a. Stop container (nerdctl stop --time 30)
     b. Remove container (nerdctl rm)
     c. Re-run container with the new bind-mounted binary
   - Clear pebbleUpgraded after all containers are processed

3. This handles the sequence:
   - Agent upgrades → snap updates → new pebble binary in snap
   - Agent restarts → iaascontainerrunner restarts
   - ensurePebbleCurrent() detects hash mismatch
   - All containers are restarted with new pebble

4. Write tests:
   - Test that matching hashes result in no restart
   - Test that mismatched hashes trigger binary copy + container restart
   - Mock the hash computation and file copy via CommandRunner

Use crypto/sha256 for hashing. Use io.Copy for binary copy (don't shell out for cp).
This is safe because the worker only runs once per agent start — not on every poll.
```

---

### Task 15: Container Log Forwarding

**Description:** Forward container stdout/stderr logs into the Juju logging pipeline
so workload messages appear in `juju debug-log`. Without this, operators cannot see
what's happening inside workload containers. This addresses Gap 10.

**Files to modify:**
- `internal/worker/iaascontainerrunner/worker.go` — add log tailing goroutine per
  container
- `internal/worker/iaascontainerrunner/manifold.go` — accept a log sink dependency

**Design:**

After a container starts successfully, launch a goroutine that runs
`nerdctl logs --follow --timestamps <containerID>` and parses each line into the
Juju logging system. Lines are tagged with the container name so they can be filtered.

**Agent prompt:**
```
Add container log forwarding to the iaascontainerrunner so workload output appears
in `juju debug-log`.

In internal/worker/iaascontainerrunner/worker.go:

1. Add a LogSink interface to Config:
   type LogSink interface {
       // Log writes a workload log message tagged with the container name.
       Log(containerName string, timestamp time.Time, message string)
   }
   Make it optional — if nil, log forwarding is disabled (for tests).

2. After ensureRunning() succeeds for a container, start a log-tailing goroutine:
   func (w *worker) tailLogs(containerName string) error {
       id := w.containerID(containerName)
       ctx, cancel := context.WithCancel(context.Background())
       // Store cancel func so Kill() can stop tailing
       w.logCancels[containerName] = cancel

       cmd := exec.CommandContext(ctx, "nerdctl", "logs", "--follow", "--timestamps", id)
       stdout, _ := cmd.StdoutPipe()
       stderr, _ := cmd.StderrPipe()
       cmd.Start()

       // Merge stdout and stderr into a single scanner
       // Parse timestamp prefix from nerdctl log format: "2024-01-01T00:00:00Z message"
       // Forward each line to w.config.LogSink.Log(containerName, ts, msg)
   }

3. Register the goroutine with the tomb so it's cleaned up on shutdown.
   On tomb.Dying(), cancel all log-tail contexts.

4. For the log sink implementation (in manifold.go):
   - Create a LogSink adapter that writes to the unit agent's existing log sender
     (logsender.LogRecordCh or similar)
   - Tag messages with module "unit.<unitname>.workload.<containername>"
   - Use logger.INFO level for stdout, logger.ERROR for stderr

5. Handle container restart: when a container is replaced (image update, crash
   recovery), cancel the old tail goroutine and start a new one for the new container.

6. Write tests:
   - Test that tailLogs calls nerdctl logs with correct container ID
   - Test that lines are parsed and forwarded to LogSink
   - Test that cancellation stops the tail
   - Mock exec via CommandRunner (add a RunStreaming method or use exec.Cmd directly
     in a testable wrapper)

For the initial implementation, it's acceptable to use exec.Cmd directly (not through
CommandRunner) since streaming requires pipes. Test via integration tests for the
streaming path and unit test the line parsing logic separately.
```

---

## Execution Order

The tasks should be executed in this order due to dependencies:

```
Task 0  (Pebble binary)       — can be done independently
Task 1  (Server facade)       — no dependencies
Task 2  (Client facade)       — depends on Task 1
Task 3  (Threading)           — depends on Task 2
Task 4  (Worker)              — can be done in parallel with Tasks 1-3
Task 5  (Wiring)              — depends on Tasks 3 and 4
Task 6  (Assumes)             — can be done independently
Task 7  (Pebble setup)        — depends on Task 4
Task 8  (Status)              — depends on Task 4
Task 9  (Reboot)              — depends on Task 4
Task 10 (Storage)             — depends on Task 4
Task 11 (Integration)         — depends on all above
Task 12 (Docs)                — depends on all above
Task 13 (Image update)        — depends on Tasks 4 and 9
Task 14 (Pebble upgrade)      — depends on Tasks 4 and 7
Task 15 (Log forwarding)      — depends on Task 4
```

Parallelizable groups:
- **Group A** (independent): Tasks 0, 1, 4, 6
- **Group B** (after Task 1): Task 2
- **Group C** (after Tasks 2, 4): Tasks 3, 7, 8, 9, 10, 15
- **Group D** (after Task 3, 4): Task 5
- **Group E** (after Tasks 7, 9): Tasks 13, 14
- **Group F** (after all): Tasks 11, 12

---

## Compatibility Matrix (Must Remain Working)

| Charm Format | Model Type | Behaviour | Changed? |
|---|---|---|---|
| FormatV1 | IAAS | Traditional VM charm, no containers | NO |
| FormatV1 | CAAS | Rejected at deploy time | NO |
| FormatV2 | CAAS | K8s sidecar containers via pod spec | NO |
| FormatV2 | IAAS | **NEW**: nerdctl containers on VM | YES |

---

## Gap Coverage Matrix

| Gap | Description | Task(s) | Coverage |
|-----|-------------|---------|----------|
| 1 | Sidecar flag interaction | Task 3 | Implicit (Sidecar never set true for IAAS); add unit test in Task 3 |
| 2 | OCI image resolution | Task 4 | ResourcesAPI wired to existing containerimageresourcestore |
| 3 | Charm upgrade / image update | **Task 13** | Worker detects image mismatch on restart + optional watcher |
| 4 | Port conflicts | Task 4 | `--http ""` disables pebble HTTP listener |
| 5 | Credential security | Task 4 | `nerdctl login --password-stdin` |
| 6 | Graceful shutdown | Task 4 | `nerdctl stop --time 30` |
| 7 | Socket path collision | Task 4 | Unit-scoped `<dataDir>/charm/containers/<name>` |
| 8 | Pebble binary versioning | **Task 14** | SHA256 comparison on startup; restart if changed |
| 9 | `assumes: containerd` | Task 6 | Advertised unconditionally for IAAS |
| 10 | Log integration | **Task 15** | `nerdctl logs --follow` tailed into juju log pipeline |
| 11 | Backwards compatibility | Tasks 3, 11 | Empty ContainerNames = no change; regression tests |
| 12 | Hook execution context | Task 12 | Documented in user docs |
| 13 | Reboot resilience | Task 9 | `--restart=unless-stopped` + idempotent startup |
