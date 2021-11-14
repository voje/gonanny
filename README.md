# gonanny
Parental control system.
Os-agnostic software - currently developed for Linux, should be easy to add Windows handlers (we need to implement a shutdown or logout syscall).

Features:
* Allowed playtime interval (e.g. from 8:00 to 22:00)
* Accumulating playtime - each day fund account with N seconds. Spend seconds by being online. When playtime reaches zero, shutdown.
