
0.9.0
=============
2023-08-22

* Refactoring, rewriting, and enhancements (#46) (8bb68353)
* gossip: use gossip to get cluster member information (#44) (5da8845d)
* cmd_ss: parse SocketMetric.processes (#43) (915c2c70)
* http: register pprof handlers (d3510b51)
* monitor: support collecting socket stats at el7 (#39) (ae7ee7f7)
* build: remove unused files (827bd419)
* cmd: load db from backup archive (#38) (ce37eacc)
* build(deps): bump golang.org/x/sys in /rpm (85bcf447)
* rpm: update binary file mode (bc998bb0)
* cmd: add command to generate default config (#35) (319516e4)
* log: use local timezone (#34) (91854caa)
* monitor: adds support for collecting ss details (#33) (b752ea2b)
* log: set log levels (#32) (ad6a6941)
* log: persistent log files (#31) (10cb0bf5)
* build: generate proto with docker (#30) (07cc5eb9)
* cmd: reduce reclaiming interval (#29) (61ab793b)
* Fix an issue that reclaims more keys than expected (#28) (42038dec)
* datastore: reclaim the oldest metrics (#27) (7afa1343)

0.1.0
=============
2023-08-04

* RPM packagings (#25) (e10b1522)
* cmd: write default config file (#21) (1200b6dc)
* http: add backup API (#18) (ec621db5)
* Some chores, cleaning trashes, upgrading dependencies (#16) (22455dac)
* Add new API /metrics to get them (#17) (d08a1354)
* build: fix inf loop when the tool not exists (7a0638ca)
* monitor: fix nil pointer access when collection timeout  (#14) (a928a94f)
* monitor: fix netstat metric collection in the newer Linux (#13) (be240960)
* proto: fix typo (31e6ff55)
* proto: snake case the fields (4031df9f)
* build: go mod tidy (9951aca9)
* introduce HTTP server (661fef96)
* build: add golangci-lint and the checking for toolchains (591b210c)
* monitor: introduce netstat counters (#11) (5face0a0)
* proto: follow protobuf guidelines (#10) (236324f3)
* build: fix missing modules (#9) (eceaaed4)
* monitor: get interval from config (#6) (bc8269d6)
* cmd: add flags and configuration parsing (#5) (30b28596)
* monitor: add unit tests for command parsers (#4) (9d029d8a)
* Nic monitor with protobuf (#3) (7908a301)
* Create go.yml (5b610b25)
* introduce metrics storage (c38df12f)
* refactor: rename iface records (a4ad4f0f)
* Support iface info collection (#1) (eec41b6a)
* Create project (dd6cffcd)
* Initial commit (7bfb8899)


